package pro

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

type Pipeline struct {
	s *SClient
	c *Client
}

func NewPipeline(cn net.Conn) *Pipeline {
	s := NewSClient(cn)
	pipe := &Pipeline{
		s: s,
	}
	pipe.s.onReadBuff = func(buf *buffer) {
		onReadBuffer(pipe, buf)
	}
	pipe.s.onError = func(cn net.Conn, err error) {
		onReadError(pipe, err)
	}

	return pipe

}

func (p *Pipeline) Start() {
	//start listen from client
	p.s.Start()
}

func onReadError(p *Pipeline, err error) {
	p.close("Pipe read error", err)
}

func (p *Pipeline) close(msg string, err error) {
	//log.Printf("%s : %v", msg, err)
	if p.s != nil {
		p.s.close()
	}
	if p.c != nil {
		p.c.close()
	}
}

func onReadBuffer(p *Pipeline, buf *buffer) {
	c, err := parseCommand(buf, p.s.encrypter)
	if err != nil {
		p.close("ParseCommand ", err)
	}
	switch cmd := c.(type) {
	case *cmdConnect:
		p.connectToTargetIfAuthOk(cmd)
	case *cmdProxy:
		p.sendClientDataToProxyTargetAndResponse(cmd)
	default:
		panic(errors.New("Unreachable. "))
	}
}

func (p *Pipeline) sendClientDataToProxyTargetAndResponse(cmd *cmdProxy) {
	_, err := p.c.cn.Write(cmd.data)
	if err != nil {
		p.close("send to proxy client err. ", err)
	} else {
		resp := newProxyCmdResp(true, cmd.id)
		p.s.sendToClient(resp)
	}
}

func (p *Pipeline) connectToTargetIfAuthOk(cmd *cmdConnect) {
	if cmd.auth == Config.Auth {
		p.tryConnectingTarget(cmd)
	} else {
		p.sendFailResponse(cmd)
	}

}

func (p *Pipeline) sendFailResponse(cmd *cmdConnect) {
	resp := newCmdConnectFailResp(cmd.enctype)
	p.s.sendToClient(resp)
	msg := fmt.Sprintf("Invalid auth from client. Actual :%s , desire:%s, target host is %s\n", cmd.auth, Config.Auth, cmd.getHost())
	log.Printf(msg)
	p.close(msg, nil)
}

func (p *Pipeline) tryConnectingTarget(cmd *cmdConnect) {
	proxyClient, e := newProxyClient(cmd)
	if e != nil {
		p.close("Error when connect to target server .", e)
	} else {
		p.configProxyClient(proxyClient)
		resp := newCmdConnectResp(true, cmd.enctype)
		p.createSClientEncrypter(resp)
		p.c.Start()
		p.s.sendToClient(resp)
	}
}

func (p *Pipeline) createSClientEncrypter(cmd *cmdConnectResp) {
	//TODO add aes encFactory
	iv := generateIV(cmd.encType)
	cmd.configIv(iv)
	switch cmd.encType {
	case ENC_CAESAR:
		p.s.encrypter = &caesarEncrypter{offset: iv[0]}
	case ENC_AES_CBC:
		p.s.encrypter = &aesCBCEncrypter{
			iv:  cmd.iv,
			key: paddingEncKey(Config.Auth),
		}
	case ENC_AES_CFB:
		p.s.encrypter = &aesCFBEncrypter{
			iv:  cmd.iv,
			key: paddingEncKey(Config.Auth),
		}
	}
}

func (p *Pipeline) configProxyClient(c *Client) {
	p.c = c
	c.onError = func(cn net.Conn, err error) {
		p.close("Error in proxy client.", err)
	}
	p.c.onRead = func(l int, data []byte) {
		rd := rand.Intn(10000)
		t := time.Now().Format("2006-01-02 15:04:05")
		id := fmt.Sprintf("%s:%d", t, rd)
		resp := newServerProxyData(id, data, p.s.encrypter)
		p.s.sendToClient(resp)
	}
}
