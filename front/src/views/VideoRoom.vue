<template>
  <div class="container">
    <!-- 顶部信息栏 -->
    <div class="header">
      <div class="room-title-area">
        <div class="back-btn" @click="goBack">
          <iconify-icon icon="bi:arrow-left" width="18"></iconify-icon>
          <span style="margin-left: 6px">返回房间列表</span>
        </div>
        <div class="room-name">
          {{ roomInfo.name }}
          <span class="status-tag">{{ roomInfo.status }}</span>
        </div>
      </div>

      <div class="action-area">
        <div class="action-item">
          <iconify-icon icon="mdi:users" width="16"></iconify-icon>
          <span>{{ roomInfo.currentUsers }}/{{ roomInfo.maxUsers }}人在线</span>
        </div>
        <div class="action-item">
          <iconify-icon icon="mdi:clock-outline" width="16"></iconify-icon>
          <span>已专注{{ roomInfo.studyTime }}</span>
        </div>
        <div class="action-item" title="房间设置" @click="showSettings">
          <iconify-icon icon="mdi:cog-outline" width="16"></iconify-icon>
        </div>
      </div>
    </div>

    <!-- 中间视频通话区 -->
    <div class="main-content">
      <!-- 左侧视频容器 -->
      <div class="video-container">
        <div 
          v-for="(video, index) in activeVideos" 
          :key="index"
          class="video-box"
          @dblclick="switchVideoLayout(index)"
        >
          <div class="video-user-info">
            <span>{{ video.userName }} - {{ video.role }}</span>
          </div>
          <div class="video-controls">
            <iconify-icon
              class="control-icon"
              :icon="video.micOn ? 'mdi:microphone' : 'mdi:microphone-off'"
              width="20"
              :style="{ color: video.micOn ? 'white' : '#ff4d4f' }"
            ></iconify-icon>
            <iconify-icon
              class="control-icon"
              :icon="video.videoOn ? 'mdi:video' : 'mdi:video-off'"
              width="20"
              :style="{ color: video.videoOn ? 'white' : '#d9d9d9' }"
            ></iconify-icon>
          </div>
        </div>
      </div>

      <!-- 右侧成员列表 -->
      <div class="members-container">
        <div class="members-header">
          房间成员
          <span class="online-count">({{ members.length }}人)</span>
        </div>

        <div class="members-list">
          <div 
            v-for="member in members" 
            :key="member.id"
            :class="['member-item', { active: member.id === activeMemberId }]"
            @click="selectMember(member.id)"
          >
            <div :class="['avatar', `avatar-${member.avatarType}`]"></div>
            <div class="member-info">
              <div class="member-name">{{ member.name }}</div>
              <div class="member-title">{{ member.role }}</div>
            </div>
            <div :class="['status-indicator', member.online ? 'status-online' : 'status-offline']"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部聊天区域 -->
    <div class="chat-container">
      <div class="messages-container" ref="messagesContainer">
        <div v-for="(group, index) in groupedMessages" :key="index">
          <div class="time-divider">{{ group.time }}</div>
          <div 
            v-for="message in group.messages" 
            :key="message.id"
            :class="['message-item', { self: message.isSelf }]"
          >
            <div :class="['message-avatar', `avatar-${message.avatarType}`]"></div>
            <div class="message-content">
              <div class="sender-info">
                <div class="sender-name">{{ message.senderName }}</div>
                <div class="sender-title">{{ message.senderRole }}</div>
              </div>
              <div class="message-bubble">{{ message.content }}</div>
              <div class="message-time">{{ message.time }}</div>
            </div>
          </div>
        </div>
      </div>

      <div class="input-area">
        <input 
          v-model="newMessage"
          type="text" 
          class="message-input" 
          placeholder="输入消息..." 
          @keydown.enter="sendMessage"
        />
        <button 
          class="send-btn" 
          :disabled="!newMessage.trim()"
          @click="sendMessage"
        >
          发送
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'VideoRoom',
  data() {
    return {
      roomInfo: {
        name: '前端学习小组',
        status: '进行中',
        currentUsers: 12,
        maxUsers: 20,
        studyTime: '2.5h'
      },
      activeVideos: [
        {
          userName: '李明',
          role: '组长',
          micOn: false,
          videoOn: false
        },
        {
          userName: '张小雨',
          role: '组员',
          micOn: true,
          videoOn: true
        }
      ],
      activeMemberId: 1,
      members: [
        { id: 1, name: '李明', role: '组长', online: true, avatarType: 1 },
        { id: 2, name: '张小雨', role: '副组长', online: true, avatarType: 2 },
        { id: 3, name: '王浩然', role: '核心成员', online: true, avatarType: 3 },
        { id: 4, name: '赵舒雅', role: '新人', online: true, avatarType: 4 },
        { id: 5, name: '孙哲', role: '组员', online: true, avatarType: 5 },
        { id: 6, name: '刘佳琪', role: '新人', online: true, avatarType: 6 },
        { id: 7, name: '陈宇', role: '组员', online: false, avatarType: 1 }
      ],
      messages: [
        {
          id: 1,
          senderName: '李明',
          senderRole: '组长',
          content: '大家早上好！今天我们的目标是完成React组件库的设计文档，建议先进行1小时的代码实践',
          time: '09:30',
          timeGroup: '今天 09:30',
          avatarType: 1,
          isSelf: false
        },
        {
          id: 2,
          senderName: '孙哲',
          senderRole: '组员',
          content: '收到，我已经开始研究Form组件的设计模式了，稍后和大家讨论设计方案',
          time: '09:32',
          timeGroup: '今天 09:30',
          avatarType: 5,
          isSelf: false
        },
        {
          id: 3,
          senderName: '赵舒雅',
          senderRole: '新人',
          content: '遇到一个问题，React Hooks中用useReducer管理复杂状态时怎么避免重复渲染？',
          time: '10:15',
          timeGroup: '今天 10:15',
          avatarType: 4,
          isSelf: true
        },
        {
          id: 4,
          senderName: '张小雨',
          senderRole: '副组长',
          content: '舒雅，建议使用上下文+useMemo进行优化。另外建议阅读我们共享文档中的性能优化部分，里面专门分析了这个情况',
          time: '10:18',
          timeGroup: '今天 10:15',
          avatarType: 2,
          isSelf: false
        },
        {
          id: 5,
          senderName: '王浩然',
          senderRole: '核心成员',
          content: '我刚更新了代码库，添加了新的DatePicker组件基础结构，大家可以拉取最新代码',
          time: '10:20',
          timeGroup: '今天 10:15',
          avatarType: 3,
          isSelf: false
        }
      ],
      newMessage: ''
    }
  },
  computed: {
    groupedMessages() {
      const groups = {}
      this.messages.forEach(message => {
        if (!groups[message.timeGroup]) {
          groups[message.timeGroup] = {
            time: message.timeGroup,
            messages: []
          }
        }
        groups[message.timeGroup].messages.push(message)
      })
      return Object.values(groups)
    }
  },
  mounted() {
    // 获取房间ID并加载房间信息
    const roomId = this.$route.params.roomId
    if (roomId) {
      this.loadRoomInfo(roomId)
    }
  },
  methods: {
    goBack() {
      this.$router.push('/study-room')
    },
    
    showSettings() {
      alert('房间设置功能开发中...')
    },
    
    switchVideoLayout(index) {
      // 视频布局切换逻辑
      console.log('切换视频布局:', index)
    },
    
    selectMember(memberId) {
      if (this.activeMemberId === memberId) return
      
      const member = this.members.find(m => m.id === memberId)
      if (member && confirm(`将与${member.name}开始视频通话，当前通话将结束`)) {
        this.activeMemberId = memberId
      }
    },
    
    sendMessage() {
      if (!this.newMessage.trim()) return
      
      const now = new Date()
      const hours = now.getHours().toString().padStart(2, '0')
      const minutes = now.getMinutes().toString().padStart(2, '0')
      const timeStr = `${hours}:${minutes}`
      
      const newMsg = {
        id: this.messages.length + 1,
        senderName: '赵舒雅',
        senderRole: '新人',
        content: this.newMessage,
        time: timeStr,
        timeGroup: `今天 ${timeStr}`,
        avatarType: 4,
        isSelf: true
      }
      
      this.messages.push(newMsg)
      this.newMessage = ''
      
      // 滚动到底部
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer
        container.scrollTop = container.scrollHeight
      })
    },
    
    loadRoomInfo(roomId) {
      // 根据房间ID加载房间信息
      console.log('加载房间信息:', roomId)
      // 这里可以调用API获取房间详情
    }
  }
}
</script>

<style scoped>
* {
  font-family: "Microsoft YaHei", sans-serif;
  box-sizing: border-box;
}

.container {
  width: 100%;
  height: 100vh;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.05);
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background: white;
}

/* 顶部信息栏 */
.header {
  height: 80px;
  background: #fff;
  border-bottom: 1px solid #e5e9f2;
  display: flex;
  align-items: center;
  padding: 0 24px;
}

.room-title-area {
  width: 60%;
  display: flex;
  align-items: center;
  position: relative;
}

.back-btn {
  cursor: pointer;
  display: flex;
  align-items: center;
  color: #8c8c8c;
  transition: color 0.2s;
}
.back-btn:hover {
  color: #1890ff;
}

.room-name {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  font-size: 20px;
  font-weight: bold;
  color: #333;
}

.status-tag {
  background: #f0fff4;
  color: #00b42a;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 20px;
  margin-left: 8px;
}

.action-area {
  width: 40%;
  display: flex;
  justify-content: flex-end;
  gap: 24px;
}

.action-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #333;
  font-size: 14px;
  cursor: pointer;
}

/* 中间区域 */
.main-content {
  height: 550px;
  display: flex;
  padding: 16px;
  background: #f5f7fa;
}

.video-container {
  width: 70%;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  padding: 16px;
  display: flex;
  gap: 4%;
}

.video-box {
  width: 48%;
  height: 100%;
  background: black;
  border-radius: 8px;
  border: 1px solid #e5e9f2;
  overflow: hidden;
  position: relative;
  cursor: pointer;
}

.video-user-info {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  font-size: 16px;
  padding: 4px 8px;
  display: flex;
  justify-content: space-between;
}

.video-controls {
  position: absolute;
  bottom: 8px;
  right: 8px;
  display: flex;
  gap: 8px;
}

.control-icon {
  cursor: pointer;
  transition: transform 0.2s;
}
.control-icon:hover {
  transform: scale(1.1);
}

.members-container {
  width: 30%;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
}

.members-header {
  padding: 12px 16px;
  font-size: 16px;
  font-weight: bold;
  color: #333;
  border-bottom: 1px solid #e5e9f2;
}

.online-count {
  font-size: 14px;
  color: #8c8c8c;
  font-weight: normal;
  margin-left: 6px;
}

.members-list {
  overflow-y: auto;
  height: 470px;
  padding: 8px 0;
}

.member-item {
  display: flex;
  padding: 8px 16px;
  align-items: center;
  transition: background 0.2s;
  cursor: pointer;
}
.member-item:hover {
  background: #f2f3f5;
}
.member-item.active {
  background: #e6f7ff;
  border-left: 3px solid #1890ff;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  flex-shrink: 0;
}

.avatar-1 {
  background: linear-gradient(135deg, #1890ff, #722ed1);
}
.avatar-2 {
  background: linear-gradient(135deg, #00b42a, #8bbb11);
}
.avatar-3 {
  background: linear-gradient(135deg, #fa8c16, #fa541c);
}
.avatar-4 {
  background: linear-gradient(135deg, #722ed1, #1890ff);
}
.avatar-5 {
  background: linear-gradient(135deg, #f5222d, #fa541c);
}
.avatar-6 {
  background: linear-gradient(135deg, #52c41a, #13c2c2);
}

.member-info {
  flex-grow: 1;
  margin-left: 12px;
  overflow: hidden;
}

.member-name {
  font-size: 14px;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.member-title {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 2px;
}

.status-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-right: 10px;
}
.status-online {
  background: #00b42a;
}
.status-offline {
  background: #8c8c8c;
}

/* 聊天区域 */
.chat-container {
  height: 350px;
  margin-top: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e5e9f2;
  display: flex;
  flex-direction: column;
  margin: 20px 16px 16px;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.time-divider {
  text-align: center;
  font-size: 12px;
  color: #bbb;
  margin: 12px 0;
  position: relative;
}
.time-divider::before,
.time-divider::after {
  content: "";
  position: absolute;
  top: 50%;
  height: 1px;
  background: #e5e9f2;
  width: 30%;
}
.time-divider::before {
  left: 0;
}
.time-divider::after {
  right: 0;
}

.message-item {
  display: flex;
  margin-bottom: 16px;
}

.message-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  flex-shrink: 0;
}

.message-content {
  margin-left: 12px;
  max-width: 500px;
}

.sender-info {
  display: flex;
  align-items: center;
  margin-bottom: 4px;
}

.sender-name {
  font-size: 14px;
  color: #333;
  font-weight: bold;
}

.sender-title {
  font-size: 12px;
  color: #8c8c8c;
  margin-left: 8px;
}

.message-bubble {
  padding: 8px 12px;
  border-radius: 8px;
  background: #f2f3f5;
  line-height: 1.5;
  font-size: 14px;
  color: #333;
}
.self .message-bubble {
  background: #e6f7ff;
  border-top-right-radius: 8px;
  border-top-left-radius: 8px;
  border-bottom-right-radius: 0;
}

.message-time {
  font-size: 12px;
  color: #8c8c8c;
  text-align: right;
  margin-top: 4px;
}

.self {
  flex-direction: row-reverse;
}
.self .message-content {
  margin-left: 0;
  margin-right: 12px;
  align-items: flex-end;
}
.self .sender-info {
  justify-content: flex-end;
}

.input-area {
  height: 60px;
  padding: 10px;
  border-top: 1px solid #e5e9f2;
  display: flex;
}

.message-input {
  flex: 1;
  border: 1px solid #e5e9f2;
  border-radius: 4px;
  padding: 0 12px;
  font-size: 14px;
  outline: none;
  transition: border 0.2s;
}
.message-input:focus {
  border-color: #1890ff;
}
.message-input::placeholder {
  color: #8c8c8c;
}

.send-btn {
  width: 80px;
  margin-left: 10px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;
}
.send-btn:hover {
  background: #096dd9;
}
.send-btn:disabled {
  background: #bae0ff;
  cursor: not-allowed;
}

/* 隐藏滚动条 */
.members-list,
.messages-container {
  scrollbar-width: none;
}
.members-list::-webkit-scrollbar,
.messages-container::-webkit-scrollbar {
  display: none;
}
</style>