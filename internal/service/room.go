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

// RoomService 房间服务 通过锁来保证线程安全
type RoomService struct {
	rooms map[string]*model.Room
	mu    sync.RWMutex
}

var (
	roomService *RoomService
	once        sync.Once // 使用单例模式
)

// GetRoomService 获取房间服务单例
func GetRoomService() *RoomService {
	once.Do(func() { // Do 方法通过原子操作和双重验证锁来确保初始化是线程安全的 这里确保了roomService只初始化一次
		roomService = &RoomService{
			rooms: make(map[string]*model.Room),
		}
		go roomService.startCleanup()
	})
	return roomService
}

// CreateRoom 创建新房间
func (s *RoomService) CreateRoom() (*model.Room, error) {
	// 加互斥锁 阻塞其他goroutine的读写操作
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否达到最大房间数限制
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
	// 加读锁 阻塞其他goroutine的写操作
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
		// 生成6位随机数 这里id的可枚举数量要大于config.MaxRooms 否则程序会一直卡在生成id上
		id := strconv.Itoa(100000 + rand.Intn(900000))
		// 检查房间id是否已存在 如果存在则继续生成 否则返回
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