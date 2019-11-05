package pro

import (
	"errors"
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
	c, err := parseCommand(buf, p.s.encFactory)
	if err != nil {
		p.close("ParseCommand ", err)
	}
	switch cmd := c.(type) {
	case *cmdConnect:
		p.connectToTarget(cmd)
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

func (p *Pipeline) connectToTarget(cmd *cmdConnect) {
	proxyClient, e := newProxyClient(cmd)
	if e != nil {
		p.close("Error when connect to target server .", e)
	} else {
		p.configProxyClient(proxyClient)
		//TODO validate

		resp := newCmdConnectResp(true, cmd.enctype)
		p.createEncFactory(resp)
		p.c.Start()
		p.s.sendToClient(resp)
	}

}

func (p *Pipeline) createEncFactory(cmd *cmdConnectResp) {
	//TODO add aes encFactory
	p.s.encFactory = caesarFactory{offset: cmd.iv[0]}

}

func (p *Pipeline) configProxyClient(c *Client) {
	p.c = c
	c.onError = func(cn net.Conn, err error) {
		p.close("Error in proxy client.", err)
	}
	p.c.onRead = func(l int, data []byte) {
		id := time.Now().Format("20060102150405001")
		resp := newServerProxyData(id, data, p.s.encFactory)
		p.s.sendToClient(resp)
	}
}
