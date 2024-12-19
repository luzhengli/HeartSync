package main

import (
	"log"
	"sync-video/internal/config"
	"sync-video/internal/handler"
	
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()
	
	// 创建gin实例
	r := gin.Default()
	
	// 加载静态文件
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	
	// 注册路由
	r.GET("/", handler.Index)
	r.GET("/room/:id", handler.RoomPage)
	r.POST("/room/create", handler.CreateRoom)
	r.GET("/room/join/:id", handler.JoinRoom)
	r.GET("/ws/:room_id", handler.HandleWebSocket)
	
	// 启动服务器
	log.Fatal(r.Run(":8080"))
} 