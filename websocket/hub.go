package websocket

import "log"

// 连接管理中心
type Hub struct {
	// 保存在线的客户端
	clients map[*Client]bool

	// 通过id映射的客户端
	clientsById map[string]*Client

	// 从客户端泵入的消息（用于广播）
	broadcast chan []byte

	private chan PrivateMessage

	// 来自客户端的注册（在线）请求
	register chan *Client

	// 来自客户端的离线请求
	unregister chan *Client
}

type PrivateMessage struct {
	To      string // 目标用户id
	From    string // 发送者id
	Content []byte // 消息内容
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		private:     make(chan PrivateMessage),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		clientsById: make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register: // 读取到需要在线的客户端  并设置为true
			h.clients[client] = true
			h.clientsById[client.id] = client
		case client := <-h.unregister: // 读取到需要离线的客户端 并删除
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientsById, client.id) // 删除id索引
				close(client.send)
			}
		case message := <-h.private:
			targetClient := h.clientsById[message.To] // 获取目标客户端
			select {
			case targetClient.send <- message.Content:
				// 发送成功
			default:
				// 目标客户端的发送通道满了，关闭连接
				log.Printf("客户端 %s 的发送通道已满，关闭连接", message.To)
				close(targetClient.send)
				delete(h.clients, targetClient)
				delete(h.clientsById, targetClient.id)
			}

		case message := <-h.broadcast: // 获取广播的消息
			for client := range h.clients { // 遍历key 即具体client
				select {
				case client.send <- message: // 发送给客户端消息
				default: // 如果客户端不在线就关闭发送通道  并删除client
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
