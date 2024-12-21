package handler

import (
	"net/http"
	"time"
	"sync-video/internal/model"
	"sync-video/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Index 首页处理
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "同步视频观看",
	})
}

// RoomPage 房间页面处理
func RoomPage(c *gin.Context) {
	roomID := c.Param("id")
	roomService := service.GetRoomService()
	
	room, err := roomService.GetRoom(roomID)
	if err != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"error": "房间不存在",
		})
		return
	}

	c.HTML(http.StatusOK, "room.html", gin.H{
		"roomID": room.ID,
		"video":  room.VideoURL,
	})
}

// CreateRoom 创建房间
func CreateRoom(c *gin.Context) {
	videoURL := c.PostForm("video_url")
	if videoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "视频URL不能为空",
		})
		return
	}

	roomService := service.GetRoomService()
	room, err := roomService.CreateRoom()
	// 如果房间创建失败 则返回错误
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	room.VideoURL = videoURL
	c.JSON(http.StatusOK, gin.H{
		"room_id": room.ID,
	})
}

// JoinRoom 加入房间
func JoinRoom(c *gin.Context) {
	roomID := c.Param("id")
	roomService := service.GetRoomService()
	
	room, err := roomService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "房间不存在",
		})
		return
	}

	user := &model.User{
		ID:       uuid.New().String(),
		JoinedAt: time.Now(),
	}

	if !room.AddUser(user) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "房间已满",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": user.ID,
		"room_id": room.ID,
	})
} 