package pro

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

func TestCallBack(t *testing.T) {
	p ,_ := newProxyClient(&cmdConnect{
		htype: DOMAIN,
		requestPort:  10099,
		requestHost:  "localhost",
	})

	p.onError = func(e net.Conn, err error) {
		defer e.Close()
		t.Log(err)
	}
	p.onRead = func(l int, data []byte) {
		t.Logf("Receive %d as string %s",l, string(data))
	}
	p.tryConnect()
	defer p.close()
	if p.cn == nil {
		t.Fatal("not connection")
	}
	p.Start()
	p.cn.Write([]byte("哈哈hH1Aa\r\n"))
	time.Sleep(5 * time.Second)
}

func TestServer(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:10099")
	if err != nil {
		log.Fatalln("Failed to start server")
	}
	for {
		conn, _ := listener.Accept()
		go func(cn net.Conn) {
			scanner := bufio.NewScanner(cn)
			for scanner.Scan() {
				fmt.Printf("Receive: %s", scanner.Text())
				echo(cn, scanner.Text(), 1*time.Second)
			}
			e := cn.Close()
			if e != nil {
				t.Log(e)
			}
		}(conn)
	}
}

func echo(cn net.Conn, t string, delay time.Duration) {

	fmt.Fprintln(cn, "\t", strings.ToUpper(t))
	time.Sleep(delay)
	fmt.Fprintln(cn, "\t", strings.ToLower(t))

}
