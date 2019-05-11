package server

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/taosha1/goproxy/tunnel"
	"github.com/taosha1/goproxy/util"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

//请求路由 注册
var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	WriteBufferPool:  new(sync.Pool),
}

func ListenAndServe(port string) {
	http.HandleFunc("/free", func(responseWriter http.ResponseWriter, request *http.Request) {
		addr, err := getTarget(request)
		if err != nil {
			return
		}
		channle := make(chan net.Conn, 1)
		go createConn(addr, channle)

		wsconn, err := upgrader.Upgrade(responseWriter, request, nil)
		if err != nil {
			log.Println("upgrade websocket fail")
			return
		}

		tunnel := &tunnel.Tunnel{*wsconn}

		google := <-channle
		if google == nil {
			tunnel.Close()
			return
		}

		//请求转发
		go io.Copy(google, tunnel)
		//响应转发
		go io.Copy(tunnel, google)
	})
	log.Fatalln(http.ListenAndServe(":"+port, nil))
	// responseHeader包含在对客户端升级请求的响应中。
	// 使用responseHeader指定cookie（Set-Cookie）和应用程序协商的子协议（Sec-WebSocket-Protocol）。
	// 如果升级失败，则升级将使用HTTP错误响应回复客户端
	// 返回一个 Conn 指针，拿到后，可使用 Conn 读写数据与客户端通信。
	//var handler = func handler(w http.ResponseWriter, r *http.Request) {
	//	conn, err := upgrader.Upgrade(w, r, nil)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//	//... Use conn to send and receive messages.
	//}
}

func getTarget(r *http.Request) (string, error) {
	map1 := r.URL.Query()
	host := util.Decode(map1.Get("h"))
	port := util.Decode(map1.Get("p"))
	if host == "" || port == "" {
		return "", errors.New("parameter nil")
	}
	addr := net.JoinHostPort(host, port)
	return addr, nil
}

func createConn(addr string, c chan net.Conn) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		//失败传入nil 避免阻塞
		c <- nil
		return
	}
	c <- conn
}
