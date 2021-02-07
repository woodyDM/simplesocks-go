package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"simplesocks-go/pkg"
	"time"

)

func main() {

	initConfig()

	local := fmt.Sprintf(":%d", pkg.Config.Port)
	listener, err := net.Listen("tcp", local)
	if err != nil {
		log.Fatalf("\nFailed to start server %v", err)
		return
	} else {
		log.Printf("\nSimpleSocks Server start at port %d with auth [%s] .\n", pkg.Config.Port, pkg.Config.Auth)
	}
	go func() {
		// 启动一个 http server，注意 pprof 相关的 handler 已经自动注册过了
		if err := http.ListenAndServe(":8023", nil); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept %v\n", err)
			continue
		} else {
			go handleNewConn(conn)
		}
	}
}

func handleNewConn(conn net.Conn) {
	//todo
	//log.Printf("new connection from %v\n",conn.RemoteAddr())
	pipe := pkg.NewPipe(conn)
	pipe.Start()
}

/**
read config from config.json
 */
func initConfig() {
	configPath := "./config.json"
	pkg.LoadFile(configPath)
	rand.Seed(time.Now().Unix())
}


