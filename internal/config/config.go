package config

const (
    // 房间ID长度
    RoomIDLength = 6
    
    // 最大房间数
    MaxRooms = 1000
    
    // 房间最大人数
    MaxUsersPerRoom = 2
    
    // 房间过期时间（小时）
    RoomExpireHours = 2
)

func Init() {
    // 未来可添加更多配置初始化逻辑
} 