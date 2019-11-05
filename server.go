package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"simplesocks-go/pro"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:9999", nil))
	}()

	var port = 12021
	local := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", local)
	if err != nil {
		log.Fatalf("\nFailed to start server %v", err)
		return
	} else {
		log.Printf("SimpleSocks Server start at %d, waiting connection.", port)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept %v", err)
		} else {
			go handleNewConn(conn)
		}
	}

}

func handleNewConn(conn net.Conn) {
	pipe := pro.NewPipeline(conn)
	pipe.Start()
}
