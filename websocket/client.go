package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer. 读取pong消息60秒
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait. 发送ping消息54秒一次
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub *Hub
	// 保存客户端的websocket连接
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte

	// 保存该客户端的唯一标识
	id string
}

// 从客户端发来的json数据
type ClientMessage struct {
	To      string `json:"to"`      // 目标客户端ID
	Content string `json:"content"` // 消息内容（也可以是 []byte 或其他类型）
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c // 通知hub客户端下线
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	// 设置从客户端读取pong消息的时长
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// 接受客户端的pong消息  并增加时长到pongwait
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// 从连接中读取消息
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		var msg ClientMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("发生错误 我真的服了啊: %v", err)
		}

		privateMsg := PrivateMessage{
			From:    c.id,
			To:      msg.To,
			Content: []byte(msg.Content),
		}

		// 从客户端读取消息并私发给指定客户端
		c.hub.private <- privateMsg

		// TODO 读取客户端消息并广播
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		// 如果读取到来自其他客户端发送的消息
		case message, ok := <-c.send: // 要么通道有消息  要么通道被关闭触发return
			// 写入数据的截至时间
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok { // 没有读取成功  就发送closemessage
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
			//w, err := c.conn.NextWriter(websocket.TextMessage)
			//if err != nil {
			//	return
			//}
			//w.Write(message)

			// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline)
			//	w.Write(<-c.send)
			//}

			//if err := w.Close(); err != nil {
			//	return
			//}
		case <-ticker.C: // 每经过54秒发送一个ping
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return // 发送错误并返回 以清理
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// 连接提升为socket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 从URL参数或Header中获取客户端ID
	clientId := r.URL.Query().Get("client_id")

	// 建立一个客户端 传入连接 和创建消息缓冲 并注册client
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), id: clientId}
	client.hub.register <- client

	// 启动协程
	go client.writePump()
	go client.readPump()
}
