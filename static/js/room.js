class VideoSync {
    constructor(roomId, videoElement) {
        this.roomId = roomId;
        this.video = videoElement;
        this.userId = localStorage.getItem('user_id');
        this.ws = null;
        this.syncThreshold = 1; // 同步阈值（秒）
        this.lastUpdate = 0;
        this.connecting = false;
        this.isRemoteUpdate = false; // 添加标志位
        this.updateInterval = 1000; // 更新间隔（毫秒）
    }

    // 初始化
    init() {
        this.connectWebSocket();
        this.bindVideoEvents();
        this.bindUIEvents();
    }

    // 连接WebSocket
    connectWebSocket() {
        if (this.connecting) return;
        this.connecting = true;

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws/${this.roomId}?user_id=${this.userId}`;

        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
            console.log('WebSocket连接成功');
            this.connecting = false;
        };

        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            this.handleMessage(message);
        };

        this.ws.onclose = () => {
            console.log('WebSocket连接断开，尝试重连...');
            this.connecting = false;
            setTimeout(() => this.connectWebSocket(), 3000);
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket错误:', error);
            this.connecting = false;
        };
    }

    // 绑定视频事件
    bindVideoEvents() {
        // 播放/暂停
        this.video.onplay = () => {
            if (!this.isRemoteUpdate) {
                this.sendMessage('play');
            }
        };

        this.video.onpause = () => {
            if (!this.isRemoteUpdate) {
                this.sendMessage('pause');
            }
        };

        // 进度更新 - 优化更新频率
        this.video.ontimeupdate = () => {
            const now = Date.now();
            if (!this.isRemoteUpdate && now - this.lastUpdate > this.updateInterval) {
                this.sendMessage('seek', this.video.currentTime);
                this.lastUpdate = now;
            }
        };
    }

    // 绑定UI事件
    bindUIEvents() {
        const copyButton = document.getElementById('copyRoomId');
        copyButton.onclick = () => {
            navigator.clipboard.writeText(this.roomId)
                .then(() => alert('房间号已复制到剪贴板'))
                .catch(err => console.error('复制失败:', err));
        };
    }

    // 发送消息
    sendMessage(type, data = null) {
        if (this.isRemoteUpdate) return; // 如果是远程更新触发的事件,不发送消息

        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({
                type,
                data,
                userId: this.userId
            }));
        }
    }

    // 处理接收到的消息
    handleMessage(message) {
        this.isRemoteUpdate = true; // 设置远程更新标志
        switch (message.type) {
            case 'play':
                if (this.video.paused) {
                    this.video.play();
                }
                break;
            case 'pause':
                if (!this.video.paused) {
                    this.video.pause();
                }
                break;
            case 'seek':
                const currentTime = parseFloat(message.data);
                if (Math.abs(this.video.currentTime - currentTime) > this.syncThreshold) {
                    this.video.currentTime = currentTime;
                }
                break;
        }
        setTimeout(() => {
            this.isRemoteUpdate = false;
        }, 100); // 短暂延迟后重置标志
    }
}

// 初始化视频同步
document.addEventListener('DOMContentLoaded', function () {
    const videoPlayer = document.getElementById('videoPlayer');
    const videoSync = new VideoSync(window.ROOM_ID, videoPlayer);
    videoSync.init();
}); 