document.addEventListener('DOMContentLoaded', function() {
    const createRoomForm = document.getElementById('createRoomForm');
    const joinRoomForm = document.getElementById('joinRoomForm');

    // 创建房间
    createRoomForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        const videoURL = new FormData(this).get('video_url');

        try {
            const response = await fetch('/room/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `video_url=${encodeURIComponent(videoURL)}`
            });

            const data = await response.json();
            if (response.ok) {
                localStorage.setItem('user_id', data.user_id);
                // 创建成功，加入房间
                window.location.href = `/room/${data.room_id}`;
            } else {
                alert(data.error || '创建房间失败');
            }
        } catch (err) {
            alert('创建房间失败: ' + err.message);
        }
    });

    // 加入房间
    joinRoomForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        const roomId = new FormData(this).get('room_id');

        try {
            const response = await fetch(`/room/join/${roomId}`, {
                method: 'GET'
            });

            const data = await response.json();
            if (response.ok) {
                // 加入成功，进入房间
                localStorage.setItem('user_id', data.user_id); 
                window.location.href = `/room/${data.room_id}`;
            } else {
                alert(data.error || '加入房间失败');
            }
        } catch (err) {
            alert('加入房间失败: ' + err.message);
        }
    });
}); 