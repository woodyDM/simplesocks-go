package pro

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	NowVersion     = byte(0x01)
	LenHeader      = 4
	LenVersion     = 1
	LenHeaderTotal = LenHeader + LenVersion
	LenType =1
)

type buffer struct {
	readHeader bool
	readBody   bool
	header     *lenHolder
	body       *lenHolder
}

func readWithRemainingBuffer(data []byte, b *buffer) []*buffer {
	//side conditions
	if b == nil {
		b = new(buffer)
	}
	var result []*buffer
	var leftData = data
	for ; len(leftData) > 0; {
		if b.isNew() {
			checkProtocolVersion(leftData[0])
			b.readHeader = true
			b.header = NewLenHolder(LenHeader)
			leftData = data[1:]
		} else {
			if !b.readBody {
				ok, left := b.header.read(leftData)
				leftData = left
				if ok {
					b.readBody = true
					b.body = NewLenHolder(b.header.parseLength() - LenHeaderTotal)
				}
			} else {
				ok, left := b.body.read(leftData)
				leftData = left
				if ok {
					if !b.full(){
						panic(errors.New("Should be full when reach here "))
					}
					result = append(result, b)
					b = new(buffer)
				}
			}
		}
	}
	if !b.isNew(){
		result = append(result, b)
	}
	return result
}

func checkProtocolVersion(b byte) {
	if b != NowVersion {
		panic(errors.New("Unsupported version "))
	}
}

type lenHolder struct {
	totalLength int
	content     []byte
	pos         int
	ok          bool
}

func (buf *buffer) isNew() bool{
	return !buf.readHeader
}

func (buf *buffer) full() bool {
	return buf.header != nil && buf.body != nil && buf.header.ok && buf.body.ok
}

func NewLenHolder(l int) *lenHolder {
	if l <= 0 {
		panic(errors.New(fmt.Sprintf("Invalid content-length %d ", l)))
	}
	return &lenHolder{
		totalLength: l,
		content:     make([]byte, l),
	}
}

/**
	试图读入 lenHolder 长度的字节数
    返回true 表示字节已读满，[]byte为剩下的有效的slice
*/
func (h *lenHolder) read(data []byte) (bool, []byte) {
	h.assertNotFull()
	canRead := len(data)
	shouldRead := h.totalLength - h.pos
	readLen := min(canRead, shouldRead)
	//TODO change to use Copy
	for i := 0; i < readLen; i++ {
		h.content[h.pos+i] = data[i]
	}
	h.pos += readLen
	if canRead >= shouldRead {
		h.ok = true
		return true, data[readLen:]
	} else {
		return false, nil
	}

}

//big endian
func (h *lenHolder) parseLength() int {
	var result int
	for i := 0; i < h.totalLength; i++ {
		result += int(h.content[i]) << (8 * (3 - i))
	}
	return result
}

func (h *lenHolder) parseLength0() int {
	u := binary.BigEndian.Uint32(h.content)
	return int(u)
}

func min(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func (h *lenHolder) assertNotFull() {
	if h.ok {
		panic(errors.New("The holder is already full. "))
	}
}
