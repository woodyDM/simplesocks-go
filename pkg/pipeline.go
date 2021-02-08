package pkg

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

const defaultChanSize = 10
const timeOut = 30
const connectTimeOut = 10
const defaultBufferSize = 4096

/**
client and target server data pipeline
*/
func NewPipe(cn net.Conn) *pipe {
	return &pipe{
		client:   cn,
		done:     make(chan interface{}),
		from:     make(chan *buffer, defaultChanSize),
		toClient: make(chan serverCommand, defaultChanSize),
		out:      make(chan []byte, defaultChanSize),
		buf:      new(buffer),
	}
}

type pipe struct {
	client net.Conn
	target net.Conn
	/**
	goroutine listen on this channel, do not send data to this channel, only select on it.
	*/
	done chan interface{}
	/**
	current client command buffer
	*/
	buf *buffer
	/**
	client connect meta
	*/
	meta *cmdConnect
	/**
	encrypter for this connection
	*/
	enc      encrypter
	from     chan *buffer       //goroutine listen on this channel
	toClient chan serverCommand //goroutine listen on this channel
	out      chan []byte        //goroutine listen on this channel

	once sync.Once
}

func (p *pipe) Start() {
	go p.readClient()     //cmd from client
	go p.sendTarget()     //cmd send to target
	go p.listenToClient() //cmd send to client
}

func (p *pipe) isDone() bool {
	select {
	case <-p.done:
		return true
	default:
		return false
	}
}

func (p *pipe) close(msg string, err error, needLog bool) {
	if needLog && !skipLog(err) {
		log.Printf("Close pipe \n[message]: %v\n[ error ]: %v\n", msg, err)
	}
	//close signal
	p.once.Do(func() {
		close(p.done)
	})
}

func skipLog_(err error) bool {
	return false
}

func skipLog(err error) bool {
	if err == nil {
		return false
	}
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return true
	}
	msg := err.Error()
	b := strings.Contains(msg, "use of closed network connection")
	if b {
		return true
	}
	b = strings.Contains(msg, "i/o timeout")
	if b {
		return true
	}
	b = strings.Contains(msg, "read: connection reset by peer")
	return b
}

func (p *pipe) sendTarget() {
	for !p.isDone() {
		select {
		case buf := <-p.from:
			c, err := parseCommand(buf, p.enc)
			if err != nil {
				p.close("ParseCommand error", err, true)
				return
			}
			//todo
			//log.Printf("send to target %v\n", buf)
			switch cmd := c.(type) {
			case *cmdConnect:
				p.connect(cmd)
			case *cmdProxy:
				p.out <- cmd.data
			default:
				panic(errors.New("Unreachable. "))
			}
		case <-p.done:
			return
		}
	}
}

func (p *pipe) listenToClient() {
	closed := false
	//no for condition, for graceful shutdown
	for {
		select {
		case cmd := <-p.toClient:
			for _, b := range cmd.body() {
				for len(b) > 0 {
					errS := p.client.SetWriteDeadline(time.Now().Add(time.Second * timeOut))
					if errS != nil {
						p.close("Failed to set write timeout to client", errS, true)
						break
					}
					n, err := p.client.Write(b)
					if err != nil {
						p.close("Failed to send back to client", err, true)
						break
					}
					b = b[n:]
				}
			}
		case <-p.done:
			//graceful shutdown: should close client until all data in chan[toClient] sent.
			if closed {
				_ = p.client.Close()
				return
			} else {
				closed = true
			}
		}
	}
}

func (p *pipe) connect(cmd *cmdConnect) {
	p.meta = cmd
	if cmd.auth != Config.Auth {
		resp := newCmdConnectFailResp(cmd.enctype)
		p.toClient <- resp
		log.Printf("Invalid auth from client. Actual :%s , desire:%s, target host is %s\n", cmd.auth, Config.Auth, cmd.getHost())
		p.close("Invalid auth", nil, false)
		return
	}
	host := cmd.getHost()
	conn, e := net.DialTimeout("tcp", host, time.Second*connectTimeOut)
	if e != nil {
		p.close("Failed to connect to target", e, true)
		return
	}

	p.target = conn
	go p.readIn()
	go p.sendOut()
	resp := newCmdConnectResp(true, cmd.enctype)
	err := p.createEnc(resp)
	if err != nil {
		p.close("Failed create enc for client ", e, true)
		return
	}
	p.toClient <- resp
}

func (p *pipe) createEnc(cmd *cmdConnectResp) error {
	iv, err := generateIV(cmd.encType)
	if err != nil {
		return err
	}
	cmd.configIv(iv)
	switch cmd.encType {
	case ENC_CAESAR:
		p.enc = &caesarEncrypter{offset: iv[0]}
	case ENC_AES_CBC:
		p.enc = &aesCBCEncrypter{
			iv:  cmd.iv,
			key: paddingEncKey(Config.Auth),
		}
	default:
		panic("unreachable")

	}
	return nil
}
func (p *pipe) readClient() {
	for !p.isDone() {
		leftData := make([]byte, defaultBufferSize)
		er := p.client.SetReadDeadline(time.Now().Add(time.Second * timeOut))
		if er != nil {
			p.close("Error when set client read timeout", er, true)
			return
		}
		n, err := p.client.Read(leftData)
		if err != nil {
			p.close("Error when client read", err, true)
			return
		}
		leftData = leftData[0:n]
		//todo
		//log.Printf("read local len %d\n", len(leftData))
		for len(leftData) > 0 {
			if p.buf.isNew() {
				checkProtocolVersion(leftData[0])
				p.buf.readHeader = true
				p.buf.header = NewLenHolder(LenHeader)
				leftData = leftData[1:]
			} else {
				if !p.buf.readBody {
					ok, left := p.buf.header.read(leftData)
					leftData = left
					if ok {
						p.buf.readBody = true
						p.buf.body = NewLenHolder(p.buf.header.parseLength0() - LenHeaderTotal)
					}
				} else {
					ok, left := p.buf.body.read(leftData)
					leftData = left
					if ok {
						if !p.buf.full() {
							p.close("Should be full when reach here ", nil, true)
							return
						}
						p.from <- p.buf
						p.buf = new(buffer)
					}
				}
			}
		}
	}
}

func (p *pipe) readIn() {
	if p.meta != nil {
		log.Printf("[Goroutine]===> Create for [%s:%d]\n", p.meta.requestHost, p.meta.requestPort)
	} else {
		panic(errors.New("No meta but read in started\n"))
	}
	for !p.isDone() {
		data := make([]byte, defaultBufferSize)
		er := p.target.SetReadDeadline(time.Now().Add(time.Second * timeOut))
		if er != nil {
			p.close("Failed to set read timeout for target server", er, true)
			break
		}
		n, err := p.target.Read(data)
		if err != nil {
			p.close(fmt.Sprintf("Failed to read from target server %v ", p.target.RemoteAddr()), err, true)
			break
		} else {
			if n == 0 {
				panic(errors.New("The read bytes length should >0 "))
			}
			rd := rand.Intn(10000)
			t := time.Now().Format("2006-01-02 15:04:05")
			id := fmt.Sprintf("%s:%d", t, rd)
			data = data[0:n]
			//todo
			//log.Printf("read from target server %d\n", len(data))
			resp := newServerProxyData(id, data, p.enc)
			p.toClient <- resp
		}
	}
	log.Printf("[Goroutine]           <=== Close [%s:%d]\n", p.meta.requestHost, p.meta.requestPort)
}

func (p *pipe) sendOut() {
	closed := false
	for {
		select {
		case d := <-p.out:
			for len(d) > 0 {
				er := p.target.SetWriteDeadline(time.Now().Add(time.Second * timeOut))
				if er != nil {
					p.close("Error when set write timeout for target", er, true)
					return
				}
				n, err := p.target.Write(d)
				if err != nil {
					p.close("Error when send data to target", err, true)
					return
				}
				d = d[n:]
			}
		case <-p.done:
			if closed {
				_ = p.target.Close()
				return
			} else {
				closed = true
			}
		}
	}
}
