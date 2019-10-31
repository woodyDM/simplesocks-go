package pro

import (
	"errors"
	"net"
)

type Client struct {
	target  *cmdConnect
	cn      net.Conn
	onError func(cn net.Conn, err error)
	onRead  func(l int, data []byte)
}

/**
used in SClient
*/
func NewClient(cn net.Conn) *Client {
	p := &Client{
		cn: cn,
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
		return p, nil
	}
}

func (p *Client) close() {
	if p.cn != nil {
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
		for {
			data := make([]byte, 4096)
			n, err := p.cn.Read(data)
			if err != nil {
				p.onError(p.cn, err)
				break
			} else {
				p.onRead(n, data[0:n])
			}
		}
	}()
}

func (p *Client) assertRunning() {
	if p.cn == nil {
		panic(errors.New("should running"))
	}
}
