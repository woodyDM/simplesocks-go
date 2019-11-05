package pro

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type hostType byte
type cmdType byte

type cmdConnect struct {
	auth        string
	enctype     string
	htype       hostType
	requestPort int
	offset      byte
	requestHost string
}

type cmdProxy struct {
	idLength int
	id       string
	data     []byte
}

type cmdServerProxyData struct {
	header
	cmdProxy
	idBytes []byte
}

type header struct {
	version       byte
	contentLength int
	cType         cmdType
}

type cmdConnectResp struct {
	header
	ok           bool
	encType      string
	encTypeBytes []byte
	iv           []byte
}

type cmdProxyDataResp struct {
	header
	ok       bool
	idLength int
	idBtyes  []byte
}

type serverCommand interface {
	body() [][]byte
}

type command interface {
}

//********************** TYPES **********************
const (
	IPV4   hostType = 0x01
	DOMAIN hostType = 0x03
	IPV6   hostType = 0x04

	CONNECT          cmdType = 0x01
	CONNECT_RESPONSE cmdType = 0x11
	PROXY            cmdType = 0x02
	PROXY_RESPONSE   cmdType = 0x12
	SUCESS           byte    = 0x01
	FAIL             byte    = 0x02
)

func (c *cmdConnectResp) body() [][]byte {
	var result = make([][]byte, 6)
	fillHeader(&c.header, result)
	var isOk byte
	if c.ok {
		isOk = SUCESS
	} else {
		isOk = FAIL
	}
	result[3] = []byte{isOk, byte(len(c.encTypeBytes)), byte(len(c.iv))}
	result[4] = c.encTypeBytes
	result[5] = c.iv
	return result
}

func (c *cmdProxyDataResp) body() [][]byte {
	var result = make([][]byte, 5)
	fillHeader(&c.header, result)
	var isOk byte
	if c.ok {
		isOk = SUCESS
	} else {
		isOk = FAIL
	}
	result[3] = []byte{isOk, byte(len(c.idBtyes))}
	result[4] = c.idBtyes
	return result
}

func (c *cmdServerProxyData) body() [][]byte {
	result := make([][]byte, 6)
	fillHeader(&c.header, result)
	result[3] = []byte{byte(c.idLength)}
	result[4] = c.idBytes
	result[5] = c.data
	return result
}

func fillHeader(h *header, result [][]byte) {
	result[0] = []byte{h.version}
	var lengthBytes = make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(h.contentLength))
	result[1] = lengthBytes
	result[2] = []byte{byte(h.cType)}
}

/**
when fail, still give enctype back but iv is []byte len=1 value=0
*/
func newCmdConnectResp(ok bool, enctype string) *cmdConnectResp {

	var iv []byte
	if ok {
		iv = []byte{100}
	} else {
		iv = []byte{0}
	}
	//should fill contentLength
	result := &cmdConnectResp{
		header: header{
			version: NowVersion,
			cType:   CONNECT_RESPONSE,
		},
		ok:           ok,
		encType:      enctype,
		encTypeBytes: []byte(enctype),
		iv:           iv, //TODO change
	}
	result.fillContentLength()
	return result
}

func newProxyCmdResp(ok bool, id string) *cmdProxyDataResp {
	idBytes := []byte(id)
	l := len(idBytes)
	result := &cmdProxyDataResp{
		header: header{
			version:       NowVersion,
			cType:         PROXY_RESPONSE,
			contentLength: LenHeaderTotal + LenType + 2 + l,
		},
		ok:       ok,
		idBtyes:  idBytes,
		idLength: l,
	}
	return result
}

// the data is encrypted!
func newServerProxyData(id string, data []byte, factory encFactory) *cmdServerProxyData {
	idBytes := []byte(id)
	l := len(idBytes)
	e := factory.newEncrypter()
	enc := e.enc(data)

	result := &cmdServerProxyData{
		header: header{
			version:       NowVersion,
			cType:         PROXY,
			contentLength: LenHeaderTotal + LenType + 1 + l + len(enc),
		},
		cmdProxy: cmdProxy{
			idLength: l,
			id:       id,
			data:     enc,
		},
		idBytes: idBytes,
	}
	return result
}

func (c *cmdConnectResp) fillContentLength() {
	var result = LenHeaderTotal + LenType + 3 + len(c.iv) + len(c.encTypeBytes)
	c.contentLength = result
}

func parseCommand(buf *buffer, factory encFactory) (command, error) {
	data := buf.body.content
	body := data[1:]
	switch cmdType(data[0]) {
	case CONNECT:
		return parseConnectCmd(body)
	case PROXY:
		return parseProxyCmd(body, factory)
	default:
		return nil, errors.New(fmt.Sprintf("Type %d unsupported. ", data[0]))
	}

}

func (cmd *cmdConnect) getHost() string {
	return fmt.Sprintf("%s:%d", cmd.requestHost, cmd.requestPort)
}

//TODO PANIC
func parseConnectCmd(data []byte) (command, error) {
	authLen := int(data[0])
	encTypeLen := int(data[1])
	//offset skip the two field
	pos := 2 + authLen + encTypeLen

	hType := hostType(data[pos])
	port := parsePort(data[pos+1 : pos+3])
	offset := data[pos+3]
	it := caesarEncrypter{offset: offset}

	auth := string(it.dec(data[2 : 2+authLen]))
	encType := string(it.dec(data[2+authLen : pos]))
	host := string(it.dec(data[pos+4:]))
	result := &cmdConnect{
		auth:        auth,
		enctype:     encType,
		htype:       hType,
		requestPort: port,
		offset:      offset,
		requestHost: host,
	}
	return result, nil
}

func parseProxyCmd(data []byte, factory encFactory) (*cmdProxy, error) {
	l := int(data[0])
	id := string(data[1 : l+1])
	enc := data[l+1:]
	dec := factory.newEncrypter().dec(enc)
	return &cmdProxy{
		idLength: l,
		id:       id,
		data:     dec,
	}, nil

}

func parsePort(p []byte) int {
	port := int(p[0])<<8 + int(p[1])
	return port
}
