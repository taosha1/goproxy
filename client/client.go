package client

import (
	"github.com/gorilla/websocket"
	"github.com/taosha1/goproxy/socks5"
	"github.com/taosha1/goproxy/tunnel"
	"github.com/taosha1/goproxy/util"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

type Config struct {
	// host:port
	LocalAddr  string
	RemoteAddr string

	//Secure bool
	//Ipv6   bool
	//Host string
}

type Client struct {
	config *Config

	//remote host port
	host string
	port string
	//schema + remoteAddr
	wsAddr string

	header http.Header

	dialer websocket.Dialer
}

func New(config *Config) *Client {
	client := &Client{}
	client.config = config
	client.dialer = websocket.Dialer{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		WriteBufferPool:  new(sync.Pool),
		//HandshakeTimeout: time.Second * 10,
	}
	client.completeHostAndPort()
	client.completeWSAddr()
	//client.completeSecureSetting()
	client.completeHeader()
	return client
}

// 增强conn 实现ReadWriteCloser接口 复用io包的copy函数
// 使用websocket 建立隧道
func (client *Client) CreateTunnel(targethost, targetport string) (io.ReadWriteCloser, error) {
	//base64编码避免明文传输
	url := client.wsAddr + "/free?h=" + util.Encode(targethost) + "&p=" + util.Encode(targetport)
	conn, _, err := client.dialer.Dial(url, client.header)
	if err != nil {
		log.Println("conn fail")
		return nil, err
	}
	log.Println("conn success")
	t := &tunnel.Tunnel{*conn}
	return t, nil
}

// create a local socks5 server.
// broswer - socks5 server | local - remote
// 解析浏览器转发的socks5报文，对每次请求都打通一个隧道
func (client *Client) ListenAndServe() {
	log.Println("server:", client.wsAddr)
	//形参的conn 是socks5的隧道
	handleConn := func(conn net.Conn, target *socks5.Target) {
		log.Println("client f()")
		//t是local - remote之间的隧道
		t, err := client.CreateTunnel(target.Host, target.Port)
		if err != nil {
			log.Fatalln(err)
		}

		//请求转发
		go io.Copy(t,conn)
		//响应转发
		go io.Copy(conn,t)

	}
	socks5.ListenAndServe(client.config.LocalAddr, handleConn)
}

// 根据config 截取host和port
func (client *Client) completeHostAndPort() {
	host, port, err := net.SplitHostPort(client.config.RemoteAddr)
	if err != nil {
		log.Fatal("error remote address")
	}
	client.host = host
	client.port = port
}

// 根据host 是否为域名 拼接为websocket地址
func (client *Client) completeWSAddr() {
	var schema string = "ws://"
	// ip
	if !util.IsDomain(client.host) {
		client.wsAddr = schema + client.host + ":" + client.port
		return
	}
	// domain to ip
	log.Println("lookup for server")
	ip, err := util.Lookup(client.host)
	if err != nil {
		log.Fatalln("error lookup ip")
		return
	}
	client.wsAddr = schema + ip.String() + ":" + client.port

}

// 设置dialer的secure init tls setting.
//func (client *Client) completeSecureSetting() {
//	if !client.config.Secure {
//		return
//	}
//	cfg := tls.Config{
//		ServerName: client.host,
//	}
//	client.dialer.TLSClientConfig = cfg
//}

//设置请求头
func (client *Client) completeHeader() {
	header := http.Header{}
	header.Add("host", "www.bilibili.com")
	client.header = header
}
