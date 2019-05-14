package socks5

import (
	"errors"
	"log"
	"net"
	"strconv"
)

//解析socks5协议

type Target struct {
	//从socks5通信中获取的被转发（代理）的信息
	Host string
	Port string
}

const (
	ipv4   = 0x01
	domain = 0x03
	ipv6   = 0x04
)

var (
	firstRep = []byte{0x05, 0x00}
	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	infoRep = []byte{
		0x05, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00,
		0x10, 0x10,
	}
)

func ListenAndServe(localAddr string, f func(conn net.Conn, target *Target)) {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalln("Socks5 server start fail")
	}
	log.Println("Socks5 server listening at", localAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleSocks5(conn, f)
	}
}

func handleSocks5(conn net.Conn, f func(conn net.Conn, target *Target)) {
	target, err := socks5Handshake(conn)
	//握手失败
	if err != nil {
		conn.Close()
		log.Println("Socks5 handshake error:", err)
		return
	}
	//数据转发,客户端在收到来自服务器成功的响应后，就会开始发送数据了，服务端在收到来自客户端的数据后，会转发到目标服务。
	go f(conn, target)
}

func socks5Handshake(conn net.Conn) (*Target, error) {
	bytes := make([]byte, 1500)
	n, err := conn.Read(bytes)
	if err != nil {
		return nil, err
	}
	if bytes[0] != 0x05 {
		return nil, errors.New("不是socks5 协议")
	}
	_, err = conn.Write(firstRep)
	if err != nil {
		return nil, err
	}

	// get target info.
	n, err = conn.Read(bytes)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(infoRep)
	if err != nil {
		return nil, err
	}

	//跳过前三个字节，转换
	return parseTargetInfo(bytes[3:n])
}

func parseTargetInfo(buf []byte) (*Target, error) {
	errInvalid := errors.New("Invalid target info")
	target := &Target{}
	var len = len(buf)
	target.Port = strconv.Itoa(int(buf[len-2])<<8 | int(buf[len-1]))
	switch buf[0] {
	case ipv4:
		if len != 7 {
			return nil, errInvalid
		}
		target.Host = net.IP(buf[1:5]).String()
	case domain:
		domainlength := int(buf[1])
		if len-domainlength != 4 {
			return nil, errInvalid
		}
		target.Host = string(buf[2 : domainlength+2])
		//if host := string(buf[2 : domainlength+2]); !util.IsDomain(host) {
		//	return nil, errInvalid
		//} else {
		//	target.Host = host
		//}
	case ipv6:
		if len != 19 {
			return nil, errInvalid
		}
		target.Host = net.IP(buf[1:17]).String()
	default:
		return nil, errInvalid
	}
	//log.Println(target)
	return target, nil
}
