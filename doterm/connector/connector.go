package connector

import (
	"context"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
)

type Websocket struct {
	Request  *http.Request
	Response http.ResponseWriter

	websocket *websocket.Conn
	ctx       context.Context
	cancel    context.CancelFunc
}

// 升级get请求为webSocket协议
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *Websocket) initWebSocket(manager Manager) {
	//创建websocket
	conn, err := upGrader.Upgrade(ws.Response, ws.Request, nil)
	if err != nil {
		panic(err)
	}
	conn.SetCloseHandler(func(code int, text string) error {
		manager.Close()
		return nil
	})
	ws.websocket = conn
	ws.ctx, ws.cancel = context.WithCancel(context.Background())
}

func (ws *Websocket) outputError(err error) {
	if err := ws.websocket.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
		panic(err)
	}
}

func (ws *Websocket) readerToWebsocket(reader io.Reader) {
	go func() {
		message := make([]byte, 1*1024)
		for ws.ctx.Err() == nil {
			n, err := reader.Read(message)
			if err != nil {
				return
			} else if n > 0 {
				if err := ws.websocket.WriteMessage(websocket.TextMessage, message[0:n]); err != nil {
					return
				}
			}
		}
	}()
}

func (ws *Websocket) websocketToWriter(write io.Writer) {
	go func() {
		for ws.ctx.Err() == nil {
			_, message, err := ws.websocket.ReadMessage()
			if err != nil {
				return
			}
			if _, err := write.Write(message); err != nil {
				return
			}
		}
	}()
}

type Manager interface {
	Connect()

	Close()
}
