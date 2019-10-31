package pro

import (
	"encoding/hex"
	"testing"
)

func Test_should_resp_ok(t *testing.T){
	bytes, _ := hex.DecodeString(CONN)
	//skip headers and cmdtype
	c, _ := parseConnectCmd(bytes[6:])
	connect,_ := c.(*cmdConnect)
	result := newCmdConnectResp(true, connect)
	if !result.ok{
		t.Fatal("should ok")
	}
	if result.encType!="aes-cbc"{
		t.Fatal("type error")
	}



}


func Test_should_resp_fail(t *testing.T){
	bytes, _ := hex.DecodeString(CONN)
	//skip headers and cmdtype
	c, _ := parseConnectCmd(bytes[6:])
	connect,_ := c.(*cmdConnect)
	result := newCmdConnectResp(false, connect)
	if result.ok{
		t.Fatal("should fail")
	}
	if result.encType!="aes-cbc"{
		t.Fatal("type error")
	}



}
