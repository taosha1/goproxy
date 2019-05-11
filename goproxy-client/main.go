package main

import (
	"flag"
	"github.com/taosha1/goproxy/client"
)

func main() {
	var (
		LocalAddr  string
		RemoteAddr string
		//Secure     bool
		//Ipv6       bool
	)
	flag.StringVar(&LocalAddr, "l", "127.0.0.1:2233", "socks5 local server address")
	flag.StringVar(&RemoteAddr, "r", "", "socks5 remote server address")
	//flag.BoolVar(&Secure, "s", false, "Secure flag for enable https")
	//flag.BoolVar(&Ipv6, "6", false, "Flag for enable ipv6.")
	flag.Parse()
	cfg := client.Config{
		LocalAddr:LocalAddr,
		RemoteAddr:RemoteAddr,
		//Secure:Secure,
		//Ipv6:Ipv6,
	}
	client.New(&cfg).ListenAndServe()
}
