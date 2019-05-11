package tunnel

import "github.com/gorilla/websocket"

type Tunnel struct {
	websocket.Conn
}

func (t *Tunnel) Read(p []byte) (int,error) {
	_, data, err := t.ReadMessage()
	if err != nil {
		return 0, err
	}
	return copy(p,data),nil
	//bug：
	//copy(p, data)
	//return len(p), nil
}

func (t *Tunnel) Write(p []byte) (n int, err error) {
	err = t.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0,err
	}
	return len(p), nil
}
