package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"simplesocks-go/pro"
	"time"
)

func main() {
	//TODO delete debug server
	go func() {
		log.Println(http.ListenAndServe("localhost:9999", nil))
	}()

	configPath := "./config.json"
	pro.LoadFile(configPath)
	rand.Seed(time.Now().Unix())

	local := fmt.Sprintf(":%d", pro.Config.Port)
	listener, err := net.Listen("tcp", local)
	if err != nil {
		log.Fatalf("\nFailed to start server %v", err)
		return
	} else {
		log.Printf("\nSimpleSocks Server start at port %d with auth [%s] .\n", pro.Config.Port, pro.Config.Auth)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept %v\n", err)
		} else {
			go handleNewConn(conn)
		}
	}

}

func handleNewConn(conn net.Conn) {
	pipe := pro.NewPipeline(conn)
	pipe.Start()
}
