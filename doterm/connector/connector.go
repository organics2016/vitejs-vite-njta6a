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

func (ws *Websocket) initWebSocket() {
	//创建websocket
	conn, err := upGrader.Upgrade(ws.Response, ws.Request, nil)
	if err != nil {
		panic(err)
	}
	ws.websocket = conn
	ws.ctx, ws.cancel = context.WithCancel(context.Background())

	ws.websocket.SetCloseHandler(func(code int, text string) error {
		ws.cancel()
		return nil
	})
}

func (ws *Websocket) outputError(err error) {
	if err := ws.websocket.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
		ws.cancel()
		panic(err)
	}
	ws.cancel()
}

func (ws *Websocket) readerToWebsocket(reader io.Reader) {
	go func() {
		message := make([]byte, 1*1024)
		for ws.ctx.Err() == nil {
			n, err := reader.Read(message)
			if err != nil {
				ws.cancel()
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
				ws.cancel()
				return
			}
			if _, err := write.Write(message); err != nil {
				ws.cancel()
				return
			}
		}
	}()
}

type Manager interface {
	Connect()

	// Close 可以为阻塞方法，收听 ctx 的取消事件。
	// 也可以为非阻塞方法，用户收听到来自 ctx 的取消事件时调用 Close
	// 任何导致终端失败的异常都应该发送取消事件，然后由 Close 统一处理
	Close()
}
