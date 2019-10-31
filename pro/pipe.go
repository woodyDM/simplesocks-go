package pro

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type Pipeline struct {
	s *SClient
	c *Client
	factory encrypterFactory
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
	p.close("Pipe read error",err)
}

func (p *Pipeline) close(msg string, err error) {
	log.Printf("%sr: %v", msg, err)
	if p.s != nil {
		p.s.close()
	}
	if p.c != nil {
		p.c.close()
	}
}

func onReadBuffer(p *Pipeline, buf *buffer) {
	fmt.Println(buf)
	c, err := parseCommand(buf)
	if err != nil {
		p.close("ParseCommand ", err)
	}
	switch cmd := c.(type) {
	case *cmdConnect:
		p.connectToTarget(cmd)
	case *cmdProxy:
		fmt.Println(cmd)


	default:
		panic(errors.New("Unreachable. "))
	}
}

func (p *Pipeline) connectToTarget(cmd *cmdConnect) {
	proxyClient, e := newProxyClient(cmd)
	if e != nil {
		p.close("Error when connect to target. ",e)
	}else{
		p.configProxyClient(proxyClient)
		//TODO validate
		resp := newCmdConnectResp(true, cmd)
		p.c.Start()
		p.s.sendToClient(resp)

	}

}

func (p *Pipeline) configProxyClient(c *Client)  {
	p.c = c
}
