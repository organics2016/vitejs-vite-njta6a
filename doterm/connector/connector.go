package connector

import (
	"context"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
)

type Websocket struct {
	wsConn *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

// 升级get请求为webSocket协议
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitWebSocket(w http.ResponseWriter, r *http.Request, ctx context.Context) *Websocket {
	//创建websocket
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	ws := &Websocket{
		wsConn: conn,
	}
	ws.ctx, ws.cancel = context.WithCancel(ctx)
	ws.wsConn.SetCloseHandler(func(code int, text string) error {
		ws.cancel()
		return nil
	})
	return ws
}

func (ws *Websocket) OutputError(err error) {
	ws.wsConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	ws.cancel()
	ws.wsConn.Close()
}

func (ws *Websocket) readerToWebsocket(reader io.Reader) {
	go func() {
		message := make([]byte, 100*1024)
		for ws.ctx.Err() == nil {
			n, err := reader.Read(message)
			if err != nil {
				ws.cancel()
				return
			} else if n > 0 {
				if err := ws.wsConn.WriteMessage(websocket.TextMessage, message[0:n]); err != nil {
					return
				}
			}
		}
	}()
}

func (ws *Websocket) websocketToWriter(write io.Writer) {
	go func() {
		for ws.ctx.Err() == nil {
			_, message, err := ws.wsConn.ReadMessage()
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

type TTY interface {
	Connect() error

	// Close 可以为阻塞方法，收听 ctx 的取消事件。
	// 也可以为非阻塞方法，用户收听到来自 ctx 的取消事件时调用 Close
	// 任何导致终端失败的异常都应该发送取消事件，然后由 Close 统一处理
	Close()
}
