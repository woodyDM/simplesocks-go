package pro

import (
	"encoding/hex"
	"testing"
)

const  (
	CONN="010000002c010a07cccdcecfd0d1d280544ffc000ec8fefdfe0301bb9bfeff09c9fd0a0a0ffe0e0ec9fe0a08"
)


func Test_buff_read_from_bytes(t *testing.T){
	bytes, _ := hex.DecodeString(CONN)
	buffers := readWithRemainingBuffer(bytes, nil)
	if len(buffers)!=1{
		t.Fatal("Len should be one")
	}
	buf :=buffers[0]
	if !buf.full(){
		t.Fatal("should full")
	}
	if buf.body.totalLength!=44 - LenHeaderTotal{
		t.Fatal("length should be 39")
	}


}

func Test_cmd_connect(t *testing.T){
	bytes, _ := hex.DecodeString(CONN)
	//skip headers and cmdtype
	c, _ := parseConnectCmd(bytes[6:])
	connect,_ := c.(*cmdConnect)
	if connect.requestPort!=443{
		t.Fatal("port 443")
	}
	if connect.auth!="1234567年"{
		t.Fatal("auth fail")
	}
	if connect.htype!=DOMAIN{
		t.Fatal("type fail")
	}
	if connect.requestHost!="cdn.bootcss.com"{
		t.Fatal("host fail")
	}
	if connect.offset!=155{
		t.Fatal("offset 155")
	}
	if connect.enctype!="aes-cbc"{
		t.Fatal("encType fail")
	}
}

func Test_lenHolder_when_not_full(t *testing.T) {
	holder := NewLenHolder(4)
	data := make([]byte, 3)
	for i := 0; i < 3; i++ {
		data[i] = byte(i + 1)
	}
	b, bytes := holder.read(data)
	if b || holder.ok {
		t.Fatal("should false ")
	}
	if bytes != nil {
		t.Fatal("should nil")
	}
	if holder.pos != 3 {
		t.Fatal("length should match ")
	}

}

func Test_lenHolder_when_exactly_full(t *testing.T) {
	holder := NewLenHolder(4)
	l := 4
	data := make([]byte, l)
	for i := 0; i < l; i++ {
		data[i] = byte(i + 1)
	}
	b, bytes := holder.read(data)
	if !b || !holder.ok {
		t.Fatal("should true ")
	}
	if len(bytes) != 0 {
		t.Fatal("should empty")
	}
	if holder.pos != holder.totalLength {
		t.Fatal("length should match ")
	}
}

func Test_lenHolder_when_full_and_exceed_need(t *testing.T) {
	holder := NewLenHolder(4)
	l := 5
	data := make([]byte, l)
	for i := 0; i < l; i++ {
		data[i] = byte(i + 1)
	}
	b, bytes := holder.read(data)
	if !b || !holder.ok {
		t.Fatal("should true ")
	}
	if len(bytes) != l-4 {
		t.Fatal("should empty")
	}
	if holder.pos != holder.totalLength {
		t.Fatal("length should match ")
	}
}

func Test_buffer_when_read_empty(t *testing.T) {


}
