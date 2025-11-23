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
          <span>{{ roomInfo.currentUsers }}/{{ roomCapacityLabel }}人在线</span>
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
        <div class="video-grid">
          <div class="video-box">
            <video ref="localVideo" autoplay playsinline muted class="video-element"></video>
            <div class="video-overlay">
              <span>{{ currentUserName }}</span>
              <span>{{ currentUserRole }}</span>
            </div>
            <div class="video-status" v-if="!mediaReady">
              正在初始化摄像头...
            </div>
          </div>
          <div class="video-box">
            <video ref="remoteVideo" autoplay playsinline class="video-element"></video>
            <div class="video-overlay">
              <span>{{ remotePeerId || '等待成员' }}</span>
              <span>远端成员</span>
            </div>
            <div class="video-status" v-if="!remoteStream">
              {{ callStatusText }}
            </div>
          </div>
        </div>
        <div class="call-panel">
          <div class="panel-row">
            <label class="panel-label">我的通话ID</label>
            <div class="peer-id-box">
              <span>{{ peerId || '生成中...' }}</span>
              <button class="panel-btn" :disabled="!peerId" @click.stop="copyPeerId">
                复制
              </button>
            </div>
          </div>
          <div class="panel-row">
            <label class="panel-label">对方通话ID</label>
            <input
              v-model="remotePeerId"
              type="text"
              class="peer-input"
              placeholder="请输入对方的 Peer ID"
            />
            <div class="call-actions">
              <button
                class="panel-btn primary"
                :disabled="!mediaReady || !peerId"
                @click="startCall"
              >
                发起通话
              </button>
              <button
                class="panel-btn danger"
                :disabled="!isCalling"
                @click="endCall"
              >
                结束
              </button>
            </div>
          </div>
          <p class="panel-hint">
            分享上方ID给同学，对方输入后即可基于 WebRTC 建立连接。
          </p>
        </div>
        <div v-if="callError" class="error-banner">
          {{ callError }}
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
            v-for="member in renderedMembers" 
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
import { computed, onMounted } from "vue";
import Peer from "peerjs";
import { ElMessage } from "element-plus";
import { useCurrentUser } from "@/composables/useCurrentUser";
import { getStudyRoomDetail } from "@/api/modules/study";

export default {
  name: "VideoRoom",
  setup() {
    const { profile, loadCurrentUser } = useCurrentUser();

    onMounted(() => {
      loadCurrentUser().catch((error) => {
        console.error("加载用户信息失败:", error);
      });
    });

    const currentUserName = computed(
      () => profile.value?.display_name || "学习者"
    );
    const currentUserRole = computed(() => profile.value?.role || "组长");

    return {
      currentUserName,
      currentUserRole,
    };
  },
  data() {
    return {
      roomInfo: {
        id: null,
        name: "在线学习房间",
        status: "进行中",
        currentUsers: 0,
        maxUsers: 0,
        studyTime: "0h",
      },
      activeMemberId: 1,
      members: [
        { id: 1, name: "", role: "组长", online: true, avatarType: 1 },
        { id: 2, name: "张小雨", role: "副组长", online: true, avatarType: 2 },
        { id: 3, name: "王浩然", role: "核心成员", online: true, avatarType: 3 },
        { id: 4, name: "赵舒雅", role: "新人", online: true, avatarType: 4 },
        { id: 5, name: "孙哲", role: "组员", online: true, avatarType: 5 },
        { id: 6, name: "刘佳琪", role: "新人", online: true, avatarType: 6 },
        { id: 7, name: "陈宇", role: "组员", online: false, avatarType: 1 },
      ],
      messages: [
        {
          id: 1,
          senderName: "",
          senderRole: "组长",
          content:
            "大家早上好！今天我们的目标是完成React组件库的设计文档，建议先进行1小时的代码实践",
          time: "09:30",
          timeGroup: "今天 09:30",
          avatarType: 1,
          isSelf: false,
        },
        {
          id: 2,
          senderName: "孙哲",
          senderRole: "组员",
          content:
            "收到，我已经开始研究Form组件的设计模式了，稍后和大家讨论设计方案",
          time: "09:32",
          timeGroup: "今天 09:30",
          avatarType: 5,
          isSelf: false,
        },
        {
          id: 3,
          senderName: "赵舒雅",
          senderRole: "新人",
          content:
            "遇到一个问题，React Hooks中用useReducer管理复杂状态时怎么避免重复渲染？",
          time: "10:15",
          timeGroup: "今天 10:15",
          avatarType: 4,
          isSelf: true,
        },
        {
          id: 4,
          senderName: "张小雨",
          senderRole: "副组长",
          content:
            "舒雅，建议使用上下文+useMemo进行优化。另外建议阅读我们共享文档中的性能优化部分，里面专门分析了这个情况",
          time: "10:18",
          timeGroup: "今天 10:15",
          avatarType: 2,
          isSelf: false,
        },
        {
          id: 5,
          senderName: "王浩然",
          senderRole: "核心成员",
          content:
            "我刚更新了代码库，添加了新的DatePicker组件基础结构，大家可以拉取最新代码",
          time: "10:20",
          timeGroup: "今天 10:15",
          avatarType: 3,
          isSelf: false,
        },
      ],
      newMessage: "",
      peer: null,
      peerId: "",
      remotePeerId: "",
      callInstance: null,
      callStatus: "idle",
      localStream: null,
      remoteStream: null,
      callError: "",
      mediaReady: false,
      isCalling: false,
    };
  },
  computed: {
    currentMemberId() {
      return this.members.length ? this.members[0].id : null;
    },
    renderedMembers() {
      return this.members.map((member) => {
        if (member.id === this.currentMemberId) {
          return {
            ...member,
            name: this.currentUserName,
            role: this.currentUserRole || member.role || "组员",
          };
        }
        return member;
      });
    },
    groupedMessages() {
      const groups = {};
      this.messages.forEach((message) => {
        const timeGroup = message.timeGroup;
        if (!groups[timeGroup]) {
          groups[timeGroup] = {
            time: timeGroup,
            messages: [],
          };
        }
        groups[timeGroup].messages.push({
          ...message,
          senderName: message.isSelf
            ? this.currentUserName
            : message.senderName,
          senderRole: message.isSelf
            ? this.currentUserRole
            : message.senderRole,
        });
      });
      return Object.values(groups);
    },
    callStatusText() {
      if (!this.mediaReady) {
        return "等待摄像头和麦克风授权...";
      }
      if (this.callStatus === "dialing") {
        return "正在呼叫对方...";
      }
      if (this.callStatus === "connected") {
        return "通话中";
      }
      return "输入对方ID后即可发起通话";
    },
    roomCapacityLabel() {
      return this.formatCapacity(this.roomInfo.maxUsers);
    },
  },
  mounted() {
    const roomId = this.$route.params.roomId;
    if (roomId) {
      this.loadRoomInfo(roomId);
    }
    this.initMedia();
  },
  beforeUnmount() {
    this.cleanupPeer();
  },
  methods: {
    goBack() {
      this.$router.push("/study-room");
    },
    showSettings() {
      ElMessage.info("房间设置功能开发中");
    },
    selectMember(memberId) {
      if (this.activeMemberId === memberId) return;
      const member = this.renderedMembers.find((m) => m.id === memberId);
      if (member && confirm(`将与${member.name}开始视频通话，当前通话将结束`)) {
        this.activeMemberId = memberId;
      }
    },
    sendMessage() {
      if (!this.newMessage.trim()) return;
      const now = new Date();
      const hours = now.getHours().toString().padStart(2, "0");
      const minutes = now.getMinutes().toString().padStart(2, "0");
      const timeStr = `${hours}:${minutes}`;
      const newMsg = {
        id: this.messages.length + 1,
        senderName: this.currentUserName,
        senderRole: this.currentUserRole,
        content: this.newMessage,
        time: timeStr,
        timeGroup: `今天 ${timeStr}`,
        avatarType: 4,
        isSelf: true,
      };
      this.messages.push(newMsg);
      this.newMessage = "";
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },
    async loadRoomInfo(roomId) {
      try {
        const res = await getStudyRoomDetail(roomId);
        const room = res?.data?.room || res?.data;
        if (room) {
          this.roomInfo = {
            id: room.id,
            name: room.name || this.roomInfo.name,
            status: room.status || "进行中",
            currentUsers: room.current_users ?? room.currentUsers ?? 0,
            maxUsers: room.max_users ?? room.maxUsers ?? 0,
            studyTime: room.study_time || room.studyTime || "0h",
          };
        }
      } catch (error) {
        console.error("加载房间信息失败:", error);
      }
    },
    formatCapacity(maxUsers) {
      if (!maxUsers || maxUsers <= 0) {
        return "不限";
      }
      return maxUsers;
    },
    async initMedia() {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({
          video: true,
          audio: true,
        });
        this.localStream = stream;
        this.mediaReady = true;
        this.attachLocalStream();
        this.initPeerInstance();
      } catch (error) {
        this.callError = "无法获取摄像头或麦克风权限";
        ElMessage.error(this.callError);
      }
    },
    attachLocalStream() {
      const videoEl = this.$refs.localVideo;
      if (videoEl && this.localStream) {
        videoEl.srcObject = this.localStream;
        videoEl.muted = true;
        const playPromise = videoEl.play();
        if (playPromise?.catch) {
          playPromise.catch(() => {});
        }
      }
    },
    initPeerInstance() {
      if (this.peer) {
        this.peer.destroy();
      }
      this.peer = new Peer();
      this.peer.on("open", (id) => {
        this.peerId = id;
      });
      this.peer.on("call", (call) => {
        if (!this.localStream) {
          call.close();
          return;
        }
        call.answer(this.localStream);
        this.handleCallLifecycle(call);
      });
      this.peer.on("error", (err) => {
        const message = err?.message || "通话服务出现问题";
        this.callError = message;
        ElMessage.error(message);
      });
    },
    handleCallLifecycle(call) {
      this.callInstance = call;
      this.isCalling = true;
      call.on("stream", (remoteStream) => {
        this.callStatus = "connected";
        this.attachRemoteStream(remoteStream);
      });
      call.on("close", () => {
        this.resetCallState();
      });
      call.on("error", (error) => {
        this.callError = error?.message || "通话过程中出现错误";
        ElMessage.error(this.callError);
        this.resetCallState();
      });
    },
    attachRemoteStream(stream) {
      this.remoteStream = stream;
      const videoEl = this.$refs.remoteVideo;
      if (videoEl) {
        videoEl.srcObject = stream;
        const playPromise = videoEl.play();
        if (playPromise?.catch) {
          playPromise.catch(() => {});
        }
      }
    },
    async startCall() {
      if (!this.peer || !this.mediaReady) {
        ElMessage.warning("请先允许访问摄像头和麦克风");
        return;
      }
      const targetId = this.remotePeerId.trim();
      if (!targetId) {
        ElMessage.warning("请输入对方的通话ID");
        return;
      }
      this.callError = "";
      this.callStatus = "dialing";
      try {
        const call = this.peer.call(targetId, this.localStream);
        if (!call) {
          throw new Error("发起通话失败");
        }
        this.handleCallLifecycle(call);
      } catch (error) {
        this.callError = error?.message || "发起通话失败";
        ElMessage.error(this.callError);
        this.resetCallState(true);
      }
    },
    endCall() {
      if (this.callInstance) {
        this.callInstance.close();
      }
      this.resetCallState(true);
    },
    resetCallState(manual = false) {
      this.isCalling = false;
      this.callStatus = "idle";
      if (this.$refs.remoteVideo) {
        this.$refs.remoteVideo.srcObject = null;
      }
      if (this.remoteStream) {
        const tracks = this.remoteStream.getTracks?.() || [];
        tracks.forEach((track) => track.stop());
      }
      this.remoteStream = null;
      this.callInstance = null;
    },
    copyPeerId() {
      if (!this.peerId) return;
      if (navigator.clipboard?.writeText) {
        navigator.clipboard
          .writeText(this.peerId)
          .then(() => {
            ElMessage.success("已复制通话ID");
          })
          .catch(() => {
            this.fallbackCopyPeerId();
          });
      } else {
        this.fallbackCopyPeerId();
      }
    },
    fallbackCopyPeerId() {
      const input = document.createElement("input");
      input.value = this.peerId;
      document.body.appendChild(input);
      input.select();
      document.execCommand("copy");
      document.body.removeChild(input);
      ElMessage.success("已复制通话ID");
    },
    cleanupPeer() {
      if (this.callInstance) {
        this.callInstance.close();
      }
      if (this.peer) {
        this.peer.destroy();
        this.peer = null;
      }
      if (this.localStream) {
        this.localStream.getTracks().forEach((track) => track.stop());
        this.localStream = null;
      }
      this.resetCallState(true);
    },
  },
};
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
  flex-direction: column;
  gap: 16px;
}

.video-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 16px;
  min-height: 320px;
}

.video-box {
  position: relative;
  border: 1px solid #e5e9f2;
  border-radius: 8px;
  overflow: hidden;
  background: #000;
  min-height: 240px;
}

.video-box .video-element {
  width: 100%;
  height: 100%;
  object-fit: cover;
  background: #000;
}

.video-overlay {
  position: absolute;
  left: 12px;
  bottom: 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  color: #fff;
  font-size: 14px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.6);
}

.video-status {
  position: absolute;
  top: 12px;
  left: 12px;
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.45);
  color: #fff;
  border-radius: 4px;
  font-size: 12px;
}

.call-panel {
  background: #f7f9fc;
  border: 1px dashed #d9e2ec;
  border-radius: 8px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.panel-row {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.panel-label {
  font-size: 13px;
  color: #555;
}

.peer-id-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 1px solid #e5e9f2;
  border-radius: 6px;
  padding: 8px 12px;
  background: #fff;
  font-size: 14px;
  color: #333;
  gap: 12px;
}

.peer-id-box span {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.peer-input {
  border: 1px solid #e5e9f2;
  border-radius: 6px;
  padding: 8px 12px;
  font-size: 14px;
  outline: none;
}
.peer-input:focus {
  border-color: #2d5bff;
}

.call-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.panel-btn {
  border: none;
  border-radius: 6px;
  padding: 8px 14px;
  font-size: 14px;
  cursor: pointer;
  background: #e5e9f2;
  color: #333;
  transition: background 0.2s;
}
.panel-btn.primary {
  background: #2d5bff;
  color: #fff;
}
.panel-btn.primary:hover:not(:disabled) {
  background: #1d47e0;
}
.panel-btn.danger {
  background: #ff4d4f;
  color: #fff;
}
.panel-btn.danger:hover:not(:disabled) {
  background: #d9363e;
}
.panel-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.panel-hint {
  font-size: 12px;
  color: #8c8c8c;
}

.error-banner {
  background: #fff2f0;
  border: 1px solid #ffccc7;
  color: #d4380d;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 13px;
  margin-top: 4px;
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
