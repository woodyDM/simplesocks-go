package pro

import (
	"testing"
)

func Test_should_cn_resp_ok(t *testing.T) {

	enctype := "caesar"
	result := newCmdConnectResp(true, enctype)
	if !result.ok {
		t.Fatal("should ok")
	}
	if result.encType != "caesar" {
		t.Fatal("type error")
	}
	if result.contentLength != 16 {
		t.Fatal("Lenght 16")
	}

}

func Test_should_cn_resp_fail(t *testing.T) {
	enctype := "caesar"
	result := newCmdConnectResp(false, enctype)
	if result.ok {
		t.Fatal("should fail")
	}
	if result.encType != "caesar" {
		t.Fatal("type error")
	}
	if result.contentLength != 16 {
		t.Fatal("Length 16 ")
	}

}

func Test_cmp_proxy_resp(t *testing.T) {
	resp := newProxyCmdResp(true, "123456")
	if !resp.ok {
		t.Fatal("should ok")
	}
	if resp.contentLength != 14 {
		t.Fatal("should 14")
	}
	if resp.idLength != 6 {
		t.Fatal("should len 6")
	}
}
