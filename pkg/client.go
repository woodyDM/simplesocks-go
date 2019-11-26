package pkg

import (
	"errors"
	"log"
	"net"
)

type Client struct {
	target  *cmdConnect
	cn      net.Conn
	running bool
	onError func(cn net.Conn, err error)
	onRead  func(l int, data []byte)
}

/**
used in SClient
*/
func NewClient(cn net.Conn) *Client {
	p := &Client{
		cn:      cn,
		running: true,
	}
	return p
}

/**
used to target server
*/
func newProxyClient(target *cmdConnect) (*Client, error) {
	p := &Client{
		target: target,
	}
	err := p.tryConnect()
	if err != nil {
		return nil, err
	} else {
		p.running = true
		return p, nil
	}
}

func (p *Client) close() {
	if p.cn != nil {
		p.running = false
		p.cn.Close()
	}
}

func (p *Client) tryConnect() error {
	host := p.target.getHost()
	conn, e := net.Dial("tcp", host)
	if e != nil {
		return e
	} else {
		p.cn = conn
		return nil
	}
}

func (p *Client) Start() {
	p.assertRunning()
	go func() {
		if p.target != nil {
			log.Printf("[Goroutine]===> Create for [%s:%d]\n", p.target.requestHost, p.target.requestPort)
		}
		for p.running {
			data := make([]byte, 4096)
			n, err := p.cn.Read(data)
			if err != nil {
				p.running = false
				p.onError(p.cn, err)
				break
			} else {
				if n == 0 {
					panic(errors.New("The read bytes length should >0 "))
				}
				p.onRead(n, data[0:n])
			}
		}
		if p.target != nil {
			log.Printf("[Goroutine]           <=== Close [%s:%d]\n", p.target.requestHost, p.target.requestPort)
		}
	}()
}

func (p *Client) assertRunning() {
	if p.cn == nil {
		panic(errors.New("should running"))
	}
}
