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
            <div class="video-status" v-if="callStatus === 'idle' && !localStream">
              摄像头暂未开启，通话建立时自动开启
            </div>
          </div>
          <div class="video-box">
            <video ref="remoteVideo" autoplay playsinline class="video-element"></video>
            <div class="video-overlay">
              <span>{{ activePartnerName || '等待成员' }}</span>
              <span>远端成员</span>
            </div>
            <div class="video-status" v-if="!remoteStream">
              {{ callStatusText }}
            </div>
          </div>
        </div>
        <div class="call-panel">
          <div class="panel-row">
            <label class="panel-label">通话状态</label>
            <div class="peer-id-box">
              <span>{{ callStatusText }}</span>
              <button
                class="panel-btn danger"
                :disabled="!isCalling"
                @click="endCall"
              >
                结束通话
              </button>
            </div>
          </div>
          <p class="panel-hint">
            点击右侧成员列表发起或接听通话，系统会自动匹配并交换通话ID。
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
          <span class="online-count">({{ renderedMembers.length }}人)</span>
        </div>

        <div class="members-list">
          <div
            v-for="member in renderedMembers"
            :key="member.id"
            :class="['member-item', { active: member.id === activeMemberId, busy: member.isBusy }]"
            @click="selectMember(member.id)"
          >
            <div :class="['avatar', `avatar-${member.avatarType}`]"></div>
            <div class="member-info">
              <div class="member-name">{{ member.name }}</div>
              <div class="member-title">{{ member.role }}</div>
            </div>
            <div :class="['status-indicator', member.isBusy ? 'status-busy' : 'status-online']"></div>
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
import { computed } from "vue";
import Peer from "peerjs";
import { ElMessage } from "element-plus";
import { useCurrentUser } from "@/composables/useCurrentUser";
import { getStudyRoomDetail } from "@/api/modules/study";
import { apiConfig } from "@/config";

export default {
  name: "VideoRoom",
  setup() {
    const { profile, loadCurrentUser } = useCurrentUser();

    const currentUserName = computed(
      () => profile.value?.display_name || "学习者"
    );
    const currentUserRole = computed(() => profile.value?.role || "组长");
    const currentUserId = computed(() => profile.value?.id || 1);

    return {
      currentUserName,
      currentUserRole,
      currentUserId,
      loadCurrentUser,
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
      members: [],
      messages: [],
      newMessage: "",
      peer: null,
      peerId: "",
      callInstance: null,
      callStatus: "idle",
      callError: "",
      mediaReady: false,
      localStream: null,
      remoteStream: null,
      ws: null,
      wsConnected: false,
      activePartnerId: null,
      activePartnerName: "",
      activeMemberId: null,
      wsReconnectTimer: null,
      mediaErrorDetail: "",
    };
  },
  computed: {
    renderedMembers() {
      return this.members;
    },
    groupedMessages() {
      const groups = {};
      this.messages.forEach((message) => {
        const timeGroup = message.timeGroup || message.sentAt || "最新消息";
        if (!groups[timeGroup]) {
          groups[timeGroup] = {
            time: timeGroup,
            messages: [],
          };
        }
        groups[timeGroup].messages.push(message);
      });
      return Object.values(groups);
    },
    callStatusText() {
      switch (this.callStatus) {
        case "dialing":
          return "正在呼叫对方...";
        case "incoming":
          return "收到通话请求";
        case "connected":
          return "通话中";
        default:
          return "等待通话";
      }
    },
    isCalling() {
      return this.callStatus === "dialing" || this.callStatus === "connected";
    },
    roomCapacityLabel() {
      return this.formatCapacity(this.roomInfo.maxUsers);
    },
  },
  mounted() {
    this.initVideoRoom();
  },
  beforeUnmount() {
    this.cleanupPeer();
    if (this.ws) {
      this.ws.close();
    }
    if (this.wsReconnectTimer) {
      clearTimeout(this.wsReconnectTimer);
      this.wsReconnectTimer = null;
    }
  },
  methods: {
    async initVideoRoom() {
      const roomId = this.$route.params.roomId;
      if (roomId) {
        await this.loadRoomInfo(roomId);
      }
      try {
        await this.loadCurrentUser();
      } catch (error) {
        console.error("加载用户信息失败:", error);
      }
      this.initPeerInstance();
      this.connectWebSocket();
    },
    goBack() {
      this.$router.push("/study-room");
    },
    showSettings() {
      ElMessage.info("房间设置功能开发中");
    },
    formatCapacity(maxUsers) {
      if (!maxUsers || maxUsers <= 0) {
        return "不限";
      }
      return maxUsers;
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
    initPeerInstance() {
      if (this.peer) {
        this.peer.destroy();
      }
      this.peer = new Peer();
      this.peer.on("open", (id) => {
        this.peerId = id;
        this.registerPeerId();
      });
      this.peer.on("call", async (call) => {
        const ok = await this.ensureLocalStream();
        if (!ok || !this.localStream) {
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
    async ensureLocalStream() {
      if (this.localStream) return true;
      const isSecureContext =
        window.location.protocol === "https:" || window.location.hostname === "localhost";
      try {
        const stream = await navigator.mediaDevices.getUserMedia({
          video: true,
          audio: true,
        });
        this.localStream = stream;
        this.mediaReady = true;
        this.mediaErrorDetail = "";
        this.attachLocalStream();
        return true;
      } catch (error) {
        this.mediaErrorDetail =
          !isSecureContext && error?.name === "NotAllowedError"
            ? "请在 HTTPS 或 localhost 环境下使用摄像头"
            : error?.message || "无法获取摄像头或麦克风权限";
        this.callError = this.mediaErrorDetail;
        ElMessage.error(this.callError);
        return false;
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
    handleCallLifecycle(call) {
      this.callInstance = call;
      call.on("stream", (remoteStream) => {
        this.callStatus = "connected";
        this.attachRemoteStream(remoteStream);
      });
      call.on("close", () => {
        this.notifyCallEnd();
        this.resetCallState();
      });
      call.on("error", (error) => {
        this.callError = error?.message || "通话过程中出现错误";
        ElMessage.error(this.callError);
        this.notifyCallEnd();
        this.resetCallState();
      });
    },
    async beginCallWithPeer(peerId, isInitiator) {
      this.callError = "";
      this.callStatus = isInitiator ? "dialing" : "incoming";
      this.remoteStream = null;
      const mediaOk = await this.ensureLocalStream();
      if (!mediaOk || !this.localStream || !this.peer) {
        this.callError = this.mediaErrorDetail || "无法建立通话";
        this.callStatus = "idle";
        return;
      }
      if (isInitiator) {
        try {
          const call = this.peer.call(peerId, this.localStream);
          if (!call) {
            throw new Error("发起通话失败");
          }
          this.handleCallLifecycle(call);
        } catch (error) {
          this.callError = error?.message || "发起通话失败";
          this.resetCallState(true);
        }
      }
    },
    endCall() {
      this.notifyCallEnd();
      if (this.callInstance) {
        this.callInstance.close();
      }
      this.resetCallState(true);
    },
    resetCallState(manual = false) {
      if (this.$refs.remoteVideo) {
        this.$refs.remoteVideo.srcObject = null;
      }
      if (this.remoteStream) {
        const tracks = this.remoteStream.getTracks?.() || [];
        tracks.forEach((track) => track.stop());
      }
      this.remoteStream = null;
      this.callInstance = null;
      this.callStatus = "idle";
      this.activePartnerId = null;
      this.activePartnerName = "";
      if (manual) {
        this.callError = "";
      }
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
      if (this.wsReconnectTimer) {
        clearTimeout(this.wsReconnectTimer);
        this.wsReconnectTimer = null;
      }
    },
    notifyCallEnd() {
      if (this.activePartnerId) {
        this.sendWs("call_end", { partner_id: this.activePartnerId });
      }
    },
    wsUrl() {
      const roomId = this.$route.params.roomId;
      const base = new URL(apiConfig.baseURL);
      const protocol = base.protocol === "https:" ? "wss:" : "ws:";
      const host = base.host;
      const params = new URLSearchParams({
        user_id: this.currentUserId,
        display_name: this.currentUserName,
      });
      return `${protocol}//${host}/api/study/rooms/${roomId}/ws?${params.toString()}`;
    },
    connectWebSocket() {
      try {
        this.ws = new WebSocket(this.wsUrl());
      } catch (error) {
        console.error("WS 创建失败", error);
        return;
      }
      this.ws.onopen = () => {
        this.wsConnected = true;
        this.registerPeerId();
        this.sendWs("state_request", {});
        if (this.wsReconnectTimer) {
          clearTimeout(this.wsReconnectTimer);
          this.wsReconnectTimer = null;
        }
      };
      this.ws.onclose = () => {
        this.wsConnected = false;
        if (!this.wsReconnectTimer) {
          this.wsReconnectTimer = setTimeout(() => {
            this.wsReconnectTimer = null;
            this.connectWebSocket();
          }, 2000);
        }
      };
      this.ws.onerror = (err) => {
        console.error("WS error", err);
      };
      this.ws.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data);
          this.handleWsMessage(payload);
        } catch (error) {
          console.error("WS parse error", error);
        }
      };
    },
    sendWs(type, data = {}) {
      if (!this.wsConnected || !this.ws) return;
      this.ws.send(JSON.stringify({ type, data }));
    },
    registerPeerId() {
      if (this.wsConnected && this.peerId) {
        this.sendWs("register_peer", { peer_id: this.peerId });
      }
    },
    handleWsMessage(message) {
      const { type, data } = message;
      switch (type) {
        case "state":
          this.applyState(data);
          break;
        case "member_joined":
          this.applyMemberJoined(data);
          break;
        case "member_left":
          this.members = this.members.filter((m) => m.user_id !== data.user_id);
          break;
        case "incoming_call":
          this.handleIncomingCall(data);
          break;
        case "call_start":
          this.handleCallStart(data);
          break;
        case "call_denied":
          this.callError = data?.reason || "对方忙线";
          ElMessage.warning(this.callError);
          this.resetCallState(true);
          break;
        case "call_ended":
          this.resetCallState(true);
          break;
        case "chat":
          this.handleIncomingChat(data);
          break;
      }
    },
    applyState(data) {
      if (!data?.members) return;
      this.members = data.members.map((m, index) => ({
        id: m.user_id,
        user_id: m.user_id,
        name: m.display_name,
        role: m.is_busy ? "通话中" : "空闲",
        online: true,
        avatarType: (index % 6) + 1,
        isBusy: !!m.is_busy,
        partnerId: m.partner_id || null,
        peerId: m.peer_id || "",
      }));
    },
    applyMemberJoined(data) {
      if (!data) return;
      const exists = this.members.find((m) => m.user_id === data.user_id);
      if (!exists) {
        this.members.push({
          id: data.user_id,
          user_id: data.user_id,
          name: data.display_name,
          role: "空闲",
          online: true,
          avatarType: (this.members.length % 6) + 1,
          isBusy: false,
          partnerId: null,
          peerId: data.peer_id || "",
        });
      }
    },
    handleIncomingCall(data) {
      if (!data?.from_id) return;
      if (!this.peerId) {
        ElMessage.warning("通话ID未就绪，暂无法接听");
        this.sendWs("call_reject", { from_id: data.from_id, reason: "peer_not_ready" });
        return;
      }
      const accept = window.confirm(`是否接听 ${data.from_name || "对方"} 的通话请求？`);
      if (accept) {
        // 被叫方在同意时提前准备本地媒体，触发权限弹窗
        this.ensureLocalStream();
        this.sendWs("call_accept", { from_id: data.from_id });
      } else {
        this.sendWs("call_reject", { from_id: data.from_id, reason: "rejected" });
      }
    },
    handleCallStart(data) {
      const { caller_id, callee_id, caller_peer_id, callee_peer_id, caller_name, callee_name } =
        data || {};
      if (!caller_id || !callee_id) return;
      const isInitiator = caller_id === this.currentUserId;
      const partnerId = isInitiator ? callee_id : caller_id;
      this.activePartnerId = partnerId;
      this.activeMemberId = partnerId;
      this.activePartnerName = isInitiator ? callee_name : caller_name;
      const partnerPeerId = isInitiator ? callee_peer_id : caller_peer_id;
      this.beginCallWithPeer(partnerPeerId, isInitiator);
    },
    selectMember(memberId) {
      if (!memberId || memberId === this.currentUserId) return;
      const member = this.members.find((m) => m.user_id === memberId);
      if (!member) return;
      if (this.isCalling) {
        ElMessage.warning("您已在通话中");
        return;
      }
      if (member.isBusy) {
        ElMessage.warning("对方正在通话中");
        return;
      }
      if (!this.peerId) {
        ElMessage.warning("通话ID生成中，请稍候");
        return;
      }
      // 主动呼叫方先准备好本地媒体，触发浏览器权限弹窗
      this.ensureLocalStream();
      this.sendWs("call_request", { target_id: memberId });
      this.callStatus = "dialing";
      this.activePartnerId = memberId;
      this.activeMemberId = memberId;
      this.activePartnerName = member.name;
    },
    handleIncomingChat(data) {
      const senderName = data?.display_name || "成员";
      const content = data?.content || "";
      if (!content) return;
      const timeLabel = data?.sent_at
        ? new Date(data.sent_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
        : "";
      this.messages.push({
        id: data.id || Date.now(),
        senderName,
        senderRole: "",
        content,
        time: timeLabel,
        timeGroup: data?.sent_at || "实时消息",
        avatarType: 2,
        isSelf: data.user_id === this.currentUserId,
      });
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container) container.scrollTop = container.scrollHeight;
      });
    },
    sendMessage() {
      const content = this.newMessage.trim();
      if (!content) return;
      this.sendWs("chat", { content });
      this.newMessage = "";
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
.member-item.busy {
  opacity: 0.7;
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
.status-busy {
  background: #faad14;
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
