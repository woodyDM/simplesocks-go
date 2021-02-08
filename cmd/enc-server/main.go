package main

import (
	"fmt"
	"log"
	"net"
	"simplesocks-go/pkg"
	"time"
)

func main() {
	port := ":18087"

	listen, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	for {
		cn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		} else {
			log.Println("connect from ", cn.RemoteAddr())
			go handle(cn)

		}
	}
}

var key = []byte{1, 2, 3, 4, 5, 6, 7, 8,
	255, 254, 253,  251, 250}

var iv = []byte{111, 22, 33, 43, 55, 66, 77, 88,
	245, 235, 225, 215, 205, 195, 185, 175}
func handle(cn net.Conn) {
	data := make([]byte, 28192)
	read, err := cn.Read(data)
	defer cn.Close()
	if err != nil {
		log.Println(err)
		_ = cn.Close()
		return
	}
	data = data[0:read]

	key = pkg.PaddingKeyUsingPkcs5(key, 16)
	for _, b := range key {
		fmt.Printf("%d ", b)
	}
	bt := pkg.DecryptAsCBC(data, key, iv)
	cn.Write(bt)
	fmt.Println("")
	fmt.Println("GOT", string(bt))
	time.Sleep(2 * time.Second)
}
