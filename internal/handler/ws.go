package handler

import (
	"log"
	"net/http"
	"sync-video/internal/model"
	"sync-video/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境需要更严格的检查
	},
}

// HandleWebSocket WebSocket连接处理
func HandleWebSocket(c *gin.Context) {
	roomID := c.Param("room_id")
	userID := c.Query("user_id")

	roomService := service.GetRoomService()
	room, err := roomService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间不存在"})
		return
	}

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}
	defer conn.Close()

	// 更新用户的WebSocket连接
	for _, user := range room.GetUsers() {
		if user.ID == userID {
			user.WSConn = conn
			break
		}
	}

	// 处理WebSocket消息
	for {
		var msg model.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("读取消息失败: %v", err)
			break
		}

		// 广播消息给房间内其他用户
		for _, user := range room.GetUsers() {
			if user.ID != userID && user.WSConn != nil { // 如果用户id不等于当前用户id 并且用户有websocket连接 则发送消息
				wsConn := user.WSConn.(*websocket.Conn)
				if err := wsConn.WriteJSON(msg); err != nil {
					log.Printf("发送消息失败: %v", err)
				}
			}
		}
	}

	// 用户断开连接时清理
	room.RemoveUser(userID)
} 