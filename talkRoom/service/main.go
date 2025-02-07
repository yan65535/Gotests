package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// 定义 WebSocket 协议升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 存储所有连接的客户端
var clients = make(map[*websocket.Conn]bool)

// 发送消息到所有客户端
func broadcastMessage(message []byte) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error sending message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// 处理每个客户端的连接
func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// 添加客户端到连接池
	clients[conn] = true
	fmt.Println("New client connected")

	// 不断读取消息并广播
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			delete(clients, conn)
			return
		}

		// 广播收到的消息到所有客户端
		if messageType == websocket.TextMessage {
			broadcastMessage(p)
		}
	}
}

func main() {
	http.HandleFunc("/chat", handleConnection)

	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
