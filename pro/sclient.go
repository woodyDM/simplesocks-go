package pro

import (
	"net"
)

type SClient struct {
	*Client
	encrypter
	onReadBuff func(buf *buffer)
	buf        *buffer
}

func NewSClient(cn net.Conn) *SClient {
	c := NewClient(cn)
	s := &SClient{
		Client: c,
	}
	c.onRead = func(l int, data []byte) {
		onRead(s, l, data)
	}
	return s
}

func (s *SClient) sendToClient(cmd serverCommand) {
	for _, b := range cmd.body() {
		_, err := s.cn.Write(b)
		if err != nil {
			s.onError(s.cn, err)
			break
		}
	}
}

func onRead(s *SClient, l int, data []byte) {
	buffers := readWithRemainingBuffer(data, s.buf)
	for _, buf := range buffers {
		if buf.full() {
			s.onReadBuff(buf)
			s.buf = nil
		} else {
			//only last of buffers may be not full buffer, hold the data when next package of bytes come.
			s.buf = buf
		}
	}
}

func (s *SClient) Start() {
	s.Client.Start()
}
