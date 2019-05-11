package main

import (
	"flag"
	"github.com/taosha1/goproxy/server"
	"log"
)

func main() {
	var port string
	flag.StringVar(&port, "p", "80", "server listen port")
	flag.Parse()
	log.Println("Server listening at", port)
	server.ListenAndServe(port)
}
