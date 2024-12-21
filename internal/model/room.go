package model

import (
	"sync"
	"time"
)

type RoomStatus string

const (
	RoomStatusEmpty   RoomStatus = "empty"
	RoomStatusWaiting RoomStatus = "waiting"
	RoomStatusFull    RoomStatus = "full"
)

// Room 房间模型
type Room struct {
	ID        string
	VideoURL  string
	Users     map[string]*User
	Status    RoomStatus
	CreatedAt time.Time
	mu        sync.RWMutex
}

// User 用户模型
type User struct {
	ID       string
	WSConn   interface{}
	JoinedAt time.Time
}

// Message WebSocket消息结构
type Message struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	UserID string      `json:"userId"`
}

// NewRoom 创建新房间
func NewRoom(id string) *Room {
	return &Room{
		ID:        id,
		Users:     make(map[string]*User),
		Status:    RoomStatusEmpty,
		CreatedAt: time.Now(),
	}
}

// AddUser 添加用户到房间
func (r *Room) AddUser(user *User) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Users) >= 2 {
		return false
	}

	r.Users[user.ID] = user
	r.updateStatus()
	return true
}

// RemoveUser 从房间移除用户
func (r *Room) RemoveUser(userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.Users, userID)
	r.updateStatus()
}

// updateStatus 更新房间状态
func (r *Room) updateStatus() {
	switch len(r.Users) {
	case 0:
		r.Status = RoomStatusEmpty
	case 1:
		r.Status = RoomStatusWaiting
	case 2:
		r.Status = RoomStatusFull
	}
}

// GetUsers 获取房间内的所有用户
func (r *Room) GetUsers() []*User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, 0, len(r.Users))
	for _, user := range r.Users {
		users = append(users, user)
	}
	return users
}
