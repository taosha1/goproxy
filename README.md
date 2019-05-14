# goproxy
基于socks5转发http请求，并使用websocket与代理服务器通信。

- websocket是一种浏览器与服务器进行全双工通讯的网络技术，属于应用层协议。它基于TCP传输协议，并复用HTTP的握手通道。

## 如何建立连接

1、客户端：申请协议升级

首先，客户端发起协议升级请求。可以看到，采用的是标准的HTTP报文格式，且只支持GET方法。
```
GET / HTTP/1.1
Host: localhost:8080
Origin: http://127.0.0.1:3000
Connection: Upgrade
Upgrade: websocket
Sec-WebSocket-Version: 13
Sec-WebSocket-Key: w4v7O6xFTi36lq3RNcgctw==
```
重点请求首部意义如下：

- Connection: Upgrade：表示要升级协议
- Upgrade: websocket：表示要升级到websocket协议。
- Sec-WebSocket-Version: 13：表示websocket的版本。如果服务端不支持该版本，需要返回一个Sec-WebSocket-Version header，里面包含服务端支持的版本号。
- Sec-WebSocket-Key：与后面服务端响应首部的Sec-WebSocket-Accept是配套的，提供基本的防护，比如恶意的连接，或者无意的连接。

2、服务端：响应协议升级

服务端返回内容如下，状态代码101表示协议切换。到此完成协议升级，后续的数据交互都按照新的协议来。

```
HTTP/1.1 101 Switching Protocols
Connection:Upgrade
Upgrade: websocket
Sec-WebSocket-Accept: Oy4NRAQ13jhfONC7bP8dTKb4PTU=
```
> 备注：每个header都以\r\n结尾，并且最后一行加上一个额外的空行\r\n。此外，服务端回应的HTTP状态码只能在握手阶段使用。过了握手阶段后，就只能采用特定的错误码。
## Usage
```
build:
go build -o client.exe goproxy-client/main.go
go build -o server.exe goproxy-server/main.go
help:
./client.exe -h
./server.exe -h
```

