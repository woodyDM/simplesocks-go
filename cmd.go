package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	_ "net/http/pprof"
	"simplesocks-go/pkg"
	"time"
)

func main() {
	configPath := "./config.json"
	pkg.LoadFile(configPath)
	rand.Seed(time.Now().Unix())

	local := fmt.Sprintf(":%d", pkg.Config.Port)
	listener, err := net.Listen("tcp", local)
	if err != nil {
		log.Fatalf("\nFailed to start server %v", err)
		return
	} else {
		log.Printf("\nSimpleSocks Server start at port %d with auth [%s] .\n", pkg.Config.Port, pkg.Config.Auth)
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

	pipe := pkg.NewPipeline(conn)
	pipe.Start()

}
