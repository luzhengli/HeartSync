package service

import (
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"sync-video/internal/config"
	"sync-video/internal/model"
	"time"
)

var (
	ErrRoomNotFound = errors.New("房间不存在")
	ErrRoomFull     = errors.New("房间已满")
	ErrInvalidRoom  = errors.New("无效的房间")
)

type RoomService struct {
	rooms map[string]*model.Room
	mu    sync.RWMutex
}

var (
	roomService *RoomService
	once        sync.Once
)

// GetRoomService 获取房间服务单例
func GetRoomService() *RoomService {
	once.Do(func() {
		roomService = &RoomService{
			rooms: make(map[string]*model.Room),
		}
		go roomService.startCleanup()
	})
	return roomService
}

// CreateRoom 创建新房间
func (s *RoomService) CreateRoom() (*model.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.rooms) >= config.MaxRooms {
		return nil, errors.New("达到最大房间数限制")
	}

	roomID := s.generateRoomID()
	room := model.NewRoom(roomID)
	s.rooms[roomID] = room

	return room, nil
}

// GetRoom 获取房间
func (s *RoomService) GetRoom(roomID string) (*model.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.rooms[roomID]
	if !exists {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

// generateRoomID 生成房间ID
func (s *RoomService) generateRoomID() string {
	for {
		// 生成6位随机数
		id := strconv.Itoa(100000 + rand.Intn(900000))
		if _, exists := s.rooms[id]; !exists {
			return id
		}
	}
}

// startCleanup 启动清理过期房间的goroutine
func (s *RoomService) startCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		s.cleanupExpiredRooms()
	}
}

// cleanupExpiredRooms 清理过期房间
func (s *RoomService) cleanupExpiredRooms() {
	s.mu.Lock()
	defer s.mu.Unlock()

	expireTime := time.Now().Add(-time.Duration(config.RoomExpireHours) * time.Hour)
	for id, room := range s.rooms {
		if room.CreatedAt.Before(expireTime) && len(room.GetUsers()) == 0 {
			delete(s.rooms, id)
		}
	}
} 