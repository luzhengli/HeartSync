# Sync Video - 双人同步观看视频系统

[![Go Version](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

一个简单的双人同步观看视频系统，支持实时同步播放、暂停和进度控制。

## 功能特点

- 🎥 支持多种视频格式 (mp4, webm, ogg)
- 🔄 实时同步播放状态
- 🔗 简单的房间系统
- 📱 响应式设计
- ⚡ 低延迟同步
- 🛡️ 自动清理过期房间

## 技术栈

- 后端：Go + Gin + WebSocket
- 前端：HTML5 + JavaScript + CSS3
- 通信：WebSocket

## 快速开始

### 前置要求

- Go 1.20+
- 支持 HTML5 的现代浏览器

### 安装

1. 克隆仓库
```bash
git clone https://github.com/yourusername/sync-video.git
cd sync-video
```

2. 安装依赖
```bash
go mod tidy
```

3. 运行服务
```bash
go run cmd/main.go
```

4. 访问 `http://localhost:8080` 开始使用

## 使用说明

1. 创建房间
   - 在首页输入视频URL
   - 点击"创建房间"
   - 系统会生成一个6位数的房间号

2. 加入房间
   - 输入房间号
   - 点击"加入房间"
   - 开始同步观看视频

## 系统限制

- 每个房间最多2人
- 支持最大房间数：1000
- 空闲房间2小时后自动清理
- 仅支持直接视频URL

## 项目结构
```
├── cmd/
│ └── main.go # 程序入口
├── internal/
│ ├── config/ # 配置文件
│ ├── handler/ # HTTP和WebSocket处理器
│ ├── model/ # 数据模型
│ └── service/ # 业务逻辑
├── static/
│ ├── css/ # 样式文件
│ └── js/ # 前端脚本
└── templates/ # HTML模板
```

## 开发计划

- [ ] 更稳定的同步算法
- [ ] 支持多用户 
- [ ] 添加用户昵称系统
- [ ] 支持更多视频源
- [ ] 添加房间内聊天功能
- [ ] 优化同步算法
- [ ] 添加房间密码功能

## 贡献指南

欢迎提交 Pull Request 或 Issue！

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件