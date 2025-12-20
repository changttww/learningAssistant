<template>
  <div class="study-room">
    <header class="room-header">
      <div class="brand">
        <div class="brand-icon">
          <iconify-icon icon="mdi:book-open-page-variant" width="22"></iconify-icon>
        </div>
        <div>
          <div class="brand-title">{{ roomInfo.name }}</div>
          <div class="brand-sub">在线自习室 · {{ roomInfo.status }}</div>
        </div>
      </div>
      <div class="room-date">{{ formattedDate }}</div>
      <div class="room-meta">
        <div class="meta-item">
          <iconify-icon icon="mdi:account-group" width="18"></iconify-icon>
          <span>{{ roomInfo.currentUsers || renderedMembers.length }}/{{ roomCapacityLabel }}人在线</span>
        </div>
        <div class="meta-item">
          <iconify-icon icon="mdi:clock-outline" width="18"></iconify-icon>
          <span>已专注{{ formattedStudyDuration }}</span>
        </div>
        <button class="ghost-btn" @click="goBack">
          <iconify-icon icon="bi:arrow-left" width="16"></iconify-icon>
          返回
        </button>
      </div>
    </header>

    <main class="room-main">
      <section class="focus-area">
        <div class="card">
          <div class="card-header">
            <div>
              <h2>专注计时器</h2>
              <p class="muted">设定目标，开启番茄钟专注一段时间。</p>
            </div>
            <div class="chip">番茄钟</div>
          </div>

          <div class="timer-wrap">
            <div class="timer-circle">
              <div class="timer-text">{{ formattedTimer }}</div>
            </div>
            <div class="timer-actions">
              <button class="circle-btn play" @click="startTimer" :disabled="timerRunning">
                <iconify-icon icon="mdi:play" width="22"></iconify-icon>
              </button>
              <button class="circle-btn pause" @click="pauseTimer" :disabled="!timerRunning">
                <iconify-icon icon="mdi:pause" width="22"></iconify-icon>
              </button>
              <button class="circle-btn reset" @click="resetTimer">
                <iconify-icon icon="mdi:refresh" width="22"></iconify-icon>
              </button>
            </div>
          </div>

          <div class="goal-forms">
            <div class="goal-form">
              <label>学习目标</label>
              <input v-model="focusGoal" type="text" placeholder="输入当前学习任务" />
            </div>
            <div class="goal-form">
              <label>目标时长（分钟）</label>
              <input v-model.number="focusMinutes" type="number" min="10" step="5" />
            </div>
          </div>

          <div class="stats-grid">
            <div class="stat-card blue">
              <div class="stat-label">今日学习</div>
              <div class="stat-value">{{ todayStudy }}</div>
            </div>
            <div class="stat-card green">
              <div class="stat-label">完成任务</div>
              <div class="stat-value">{{ completedTasks }}</div>
            </div>
          </div>
        </div>
      </section>

      <aside class="sidebar">
        <div class="card members-card">
          <div class="card-header">
            <h2>自习同伴</h2>
            <div class="online-pill"><span>{{ onlineCount }}</span>人在线</div>
          </div>
          <div class="member-filters">
            <button
              :class="['filter-btn', { active: memberFilter === 'all' }]"
              @click="memberFilter = 'all'"
            >
              全部状态
            </button>
            <button
              :class="['filter-btn', { active: memberFilter === 'focus' }]"
              @click="memberFilter = 'focus'"
            >
              学习中
            </button>
            <button
              :class="['filter-btn', { active: memberFilter === 'rest' }]"
              @click="memberFilter = 'rest'"
            >
              休息中
            </button>
          </div>
          <div class="member-list">
            <div class="member-row" v-for="member in filteredMembers" :key="member.id">
              <button
                class="avatar"
                :class="`avatar-${member.avatarType}`"
                :title="`和${member.name}私聊`"
                @click.stop="openPrivateChat(member)"
              ></button>
              <div class="member-meta">
                <div class="member-name">{{ member.name }}</div>
                <div class="member-time">
                  <iconify-icon icon="mdi:clock-outline" width="14"></iconify-icon>
                  <span>{{ member.focusTime || defaultMemberTime }}</span>
                </div>
              </div>
              <span class="status-tag" :class="memberStatusClass(member)">
                {{ memberStatusText(member) }}
              </span>
            </div>
            <div class="empty-state" v-if="!filteredMembers.length">等待小伙伴加入...</div>
          </div>
        </div>

        <div class="card chat-card">
          <div class="card-header">
            <h2>自习室聊天</h2>
          </div>
          <div class="chat-messages" ref="messagesContainer">
            <div class="chat-group" v-for="(group, index) in groupedMessages" :key="index">
              <div class="chat-time">{{ group.time }}</div>
              <div
                class="chat-message"
                v-for="message in group.messages"
                :key="message.id"
                :class="{ self: message.isSelf }"
              >
                <div class="avatar small" :class="`avatar-${message.avatarType}`"></div>
                <div class="bubble">
                  <div class="bubble-header">
                    <span class="bubble-name">{{ message.senderName }}</span>
                    <span class="bubble-time">{{ message.time }}</span>
                  </div>
                  <div class="bubble-text">{{ message.content }}</div>
                </div>
              </div>
            </div>
          </div>
          <div class="chat-input">
            <input
              v-model="newMessage"
              type="text"
              placeholder="输入消息..."
              @keydown.enter="sendMessage"
            />
            <button :disabled="!newMessage.trim()" @click="sendMessage">发送</button>
          </div>
        </div>
      </aside>
    </main>

    <div v-if="directChatVisible" class="direct-chat-overlay" @click.self="closePrivateChat">
      <div class="direct-chat-panel">
        <div class="direct-chat-header">
          <div>
            <div class="direct-chat-title">与 {{ activeDirectChatName }} 私聊</div>
            <div class="direct-chat-sub">快速私聊 · 单人聊天室</div>
          </div>
          <button class="ghost-btn small" @click="closePrivateChat">
            <iconify-icon icon="mdi:close" width="16"></iconify-icon>
            关闭
          </button>
        </div>
        <div class="direct-chat-body" ref="directMessagesContainer">
          <div v-if="!activeDirectMessages.length" class="direct-chat-empty">
            现在可以和对方私聊啦～
          </div>
          <div
            v-for="message in activeDirectMessages"
            :key="message.id"
            class="direct-chat-message"
            :class="{ self: message.isSelf }"
          >
            <div class="direct-chat-bubble">
              <div class="direct-chat-meta">
                <span class="direct-chat-name">{{ message.senderName }}</span>
                <span class="direct-chat-time">{{ message.time }}</span>
              </div>
              <div class="direct-chat-text">{{ message.content }}</div>
            </div>
          </div>
        </div>
        <div class="direct-chat-input">
          <input
            v-model="directChatInput"
            type="text"
            placeholder="输入私聊消息..."
            @keydown.enter="sendDirectMessage"
          />
          <button :disabled="!directChatInput.trim()" @click="sendDirectMessage">发送</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { useCurrentUser } from "@/composables/useCurrentUser";
import {
  getStudyRoomDetail,
  getRoomChatHistory,
  sendRoomChatMessage,
} from "@/api/modules/study";
import { apiConfig } from "@/config";

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
    const currentUserId = computed(() => profile.value?.id || 1);

    return {
      currentUserName,
      currentUserId,
    };
  },
  data() {
    return {
      roomInfo: {
        id: null,
        name: "学习空间",
        status: "进行中",
        currentUsers: 0,
        maxUsers: 0,
        studyTime: "0h",
      },
      members: [],
      messages: [],
      newMessage: "",
      focusGoal: "",
      focusMinutes: 60,
      elapsedSeconds: 0,
      timerInterval: null,
      timerRunning: false,
      memberFilter: "all",
      ws: null,
      wsConnected: false,
      completedTasks: "4/7",
      directChatVisible: false,
      directChatInput: "",
      directChats: {},
      activeDirectChatId: null,
      activeDirectChatName: "",
      pendingDirectChat: null,
    };
  },
  computed: {
    renderedMembers() {
      return this.members;
    },
    filteredMembers() {
      return this.renderedMembers.filter((member) => {
        if (this.memberFilter === "focus") {
          return !member.isResting;
        }
        if (this.memberFilter === "rest") {
          return member.isResting;
        }
        return true;
      });
    },
    groupedMessages() {
      const groups = new Map();
      this.messages.forEach((message) => {
        const timeGroup = message.timeGroup || message.sentAt || "最新消息";
        if (!groups.has(timeGroup)) {
          groups.set(timeGroup, { time: timeGroup, messages: [] });
        }
        groups.get(timeGroup).messages.push(message);
      });
      return Array.from(groups.values());
    },
    formattedTimer() {
      const { hours, minutes, seconds } = this.formatTimerParts(this.elapsedSeconds);
      return `${hours}:${minutes}:${seconds}`;
    },
    formattedDate() {
      return new Date().toLocaleDateString("zh-CN", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        weekday: "long",
      });
    },
    formattedStudyDuration() {
      return this.roomInfo.studyTime || this.formatHoursMinutes(this.elapsedSeconds);
    },
    todayStudy() {
      return this.formatHoursMinutes(this.elapsedSeconds);
    },
    onlineCount() {
      return this.filteredMembers.length || this.renderedMembers.length;
    },
    roomCapacityLabel() {
      if (!this.roomInfo.maxUsers || this.roomInfo.maxUsers <= 0) {
        return "不限";
      }
      return this.roomInfo.maxUsers;
    },
    defaultMemberTime() {
      return "00:00:00";
    },
    activeDirectMessages() {
      if (!this.activeDirectChatId) return [];
      const chat = this.directChats[this.activeDirectChatId];
      return chat?.messages || [];
    },
  },
  mounted() {
    const roomId = this.$route.params.roomId;
    if (roomId) {
      this.loadRoomInfo(roomId);
      this.loadChatHistory();
      try {
        localStorage.setItem("study:lastRoomId", roomId);
      } catch (error) {
        console.warn("无法保存最近房间ID", error);
      }
    }
    this.pendingDirectChat = this.readChatQuery();
    this.connectWebSocket();
  },
  beforeUnmount() {
    this.stopTimerInterval();
    if (this.ws) {
      this.ws.close();
    }
  },
  methods: {
    goBack() {
      this.$router.push("/study-room");
    },
    startTimer() {
      if (this.timerRunning) return;
      this.timerRunning = true;
      this.timerInterval = setInterval(() => {
        this.elapsedSeconds += 1;
      }, 1000);
    },
    pauseTimer() {
      if (!this.timerRunning) return;
      this.timerRunning = false;
      this.stopTimerInterval();
    },
    resetTimer() {
      this.stopTimerInterval();
      this.elapsedSeconds = 0;
      this.timerRunning = false;
    },
    stopTimerInterval() {
      if (this.timerInterval) {
        clearInterval(this.timerInterval);
        this.timerInterval = null;
      }
    },
    formatTimerParts(totalSeconds) {
      const hours = Math.floor(totalSeconds / 3600);
      const minutes = Math.floor((totalSeconds % 3600) / 60);
      const seconds = totalSeconds % 60;
      return {
        hours: String(hours).padStart(2, "0"),
        minutes: String(minutes).padStart(2, "0"),
        seconds: String(seconds).padStart(2, "0"),
      };
    },
    formatHoursMinutes(totalSeconds) {
      const hours = Math.floor(totalSeconds / 3600);
      const minutes = Math.floor((totalSeconds % 3600) / 60);
      if (!hours && !minutes) return "0分钟";
      const parts = [];
      if (hours) parts.push(`${hours}小时`);
      if (minutes) parts.push(`${minutes}分`);
      return parts.join("");
    },
    memberStatusText(member) {
      return member.isResting ? "休息中" : "学习中";
    },
    memberStatusClass(member) {
      return member.isResting ? "tag-rest" : "tag-focus";
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
        this.sendWs("state_request", {});
      };
      this.ws.onclose = () => {
        this.wsConnected = false;
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
        case "chat":
          this.handleIncomingChat(data);
          break;
        case "direct_chat":
          this.handleIncomingDirectChat(data);
          break;
        default:
          break;
      }
    },
    applyState(data) {
      if (data?.room) {
        this.roomInfo = {
          ...this.roomInfo,
          name: data.room.name || this.roomInfo.name,
          status: data.room.status || this.roomInfo.status,
          currentUsers: data.room.current_users ?? data.room.currentUsers ?? this.roomInfo.currentUsers,
          maxUsers: data.room.max_users ?? data.room.maxUsers ?? this.roomInfo.maxUsers,
          studyTime: data.room.study_time || data.room.studyTime || this.roomInfo.studyTime,
        };
      }
      if (!data?.members) return;
      this.members = data.members.map((m, index) => ({
        id: m.user_id,
        user_id: m.user_id,
        name: m.display_name,
        role: m.role || "同伴",
        online: true,
        avatarType: (index % 6) + 1,
        isResting: !!m.is_resting,
        focusTime: m.focus_time || m.study_time || "",
      }));
      this.tryOpenPendingChat();
    },
    applyMemberJoined(data) {
      if (!data) return;
      const exists = this.members.find((m) => m.user_id === data.user_id);
      if (!exists) {
        this.members.push({
          id: data.user_id,
          user_id: data.user_id,
          name: data.display_name || "新同学",
          role: data.role || "同伴",
          online: true,
          avatarType: (this.members.length % 6) + 1,
          isResting: !!data.is_resting,
          focusTime: data.focus_time || "",
        });
      }
      this.tryOpenPendingChat();
    },
    async sendMessage() {
      const content = this.newMessage.trim();
      if (!content) return;
      this.sendWs("chat", { content });
      const roomId = this.$route.params.roomId;
      if (roomId) {
        try {
          await sendRoomChatMessage(roomId, {
            user_id: this.currentUserId,
            content,
          });
        } catch (error) {
          console.error("发送消息失败", error);
          ElMessage.error("发送消息失败");
        }
      }
      this.newMessage = "";
    },
    async loadChatHistory() {
      const roomId = this.$route.params.roomId;
      if (!roomId) return;
      try {
        const res = await getRoomChatHistory(roomId, { limit: 100 });
        const items = res?.data?.messages || [];
        this.messages = items
          .slice()
          .reverse()
          .map((item) => this.normalizeMessage(item));
      } catch (error) {
        console.error("加载聊天记录失败", error);
      }
    },
    handleIncomingChat(data) {
      const content = data?.content || "";
      if (!content) return;
      const sentAt = data?.sent_at ? new Date(data.sent_at) : new Date();
      const timeStr = sentAt.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
      const dateLabel = sentAt.toLocaleDateString();
      this.messages.push({
        id: data.id || Date.now(),
        senderName: data.display_name || "成员",
        senderRole: "",
        content,
        time: timeStr,
        timeGroup: `${dateLabel} ${timeStr}`,
        avatarType: (data.user_id % 6) + 1,
        isSelf: data.user_id === this.currentUserId,
      });
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container) container.scrollTop = container.scrollHeight;
      });
    },
    handleIncomingDirectChat(data) {
      const content = data?.content || "";
      if (!content) return;
      const fromId = Number(data?.from_id || 0);
      const toId = Number(data?.to_id || 0);
      if (!fromId || !toId) return;
      const isSelf = fromId === this.currentUserId;
      const peerId = isSelf ? toId : fromId;
      const senderName = isSelf ? this.currentUserName : data?.display_name || "同学";
      const peerName = this.resolveMemberName(peerId);
      const sentAt = data?.sent_at ? new Date(data.sent_at) : new Date();
      const timeStr = sentAt.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
      this.ensureDirectChat(peerId, peerName);
      this.directChats[peerId].messages.push({
        id: `${fromId}-${toId}-${sentAt.getTime()}`,
        senderName,
        content,
        time: timeStr,
        isSelf,
      });
      if (this.activeDirectChatId === peerId && this.directChatVisible) {
        this.$nextTick(() => {
          const container = this.$refs.directMessagesContainer;
          if (container) container.scrollTop = container.scrollHeight;
        });
      }
    },
    normalizeMessage(item) {
      const sentAt = item.sent_at ? new Date(item.sent_at) : new Date();
      const timeStr = sentAt.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
      const dateLabel = sentAt.toLocaleDateString();
      return {
        id: item.id || `${item.user_id}-${item.sent_at}`,
        senderName: item.display_name || "成员",
        senderRole: "",
        content: item.content || "",
        time: timeStr,
        timeGroup: `${dateLabel} ${timeStr}`,
        avatarType: (item.user_id % 6) + 1,
        isSelf: item.user_id === this.currentUserId,
      };
    },
    openPrivateChat(member) {
      if (!member) return;
      if (member.user_id === this.currentUserId) {
        ElMessage.warning("不能和自己私聊");
        return;
      }
      const peerId = member.user_id;
      const peerName = member.name || "同学";
      this.ensureDirectChat(peerId, peerName);
      this.activeDirectChatId = peerId;
      this.activeDirectChatName = peerName;
      this.directChatVisible = true;
      this.$nextTick(() => {
        const input = this.$el.querySelector(".direct-chat-input input");
        if (input) input.focus();
      });
      this.$nextTick(() => {
        const container = this.$refs.directMessagesContainer;
        if (container) container.scrollTop = container.scrollHeight;
      });
    },
    closePrivateChat() {
      this.directChatVisible = false;
      this.directChatInput = "";
    },
    sendDirectMessage() {
      const content = this.directChatInput.trim();
      if (!content || !this.activeDirectChatId) return;
      if (!this.wsConnected) {
        ElMessage.warning("聊天连接未建立，请稍后重试");
        return;
      }
      this.sendWs("direct_chat", { target_id: this.activeDirectChatId, content });
      this.directChatInput = "";
    },
    ensureDirectChat(peerId, peerName) {
      if (!peerId) return;
      if (!this.directChats[peerId]) {
        this.directChats[peerId] = {
          name: peerName || "同学",
          messages: [],
        };
      } else if (peerName && !this.directChats[peerId].name) {
        this.directChats[peerId].name = peerName;
      }
      if (this.activeDirectChatId === peerId || !this.activeDirectChatId) {
        this.activeDirectChatName = this.directChats[peerId].name || peerName || "同学";
      }
    },
    resolveMemberName(userId) {
      const match = this.members.find((m) => m.user_id === userId);
      return match?.name || this.directChats[userId]?.name || "同学";
    },
    readChatQuery() {
      const { chatUserId, chatUserName } = this.$route.query || {};
      if (chatUserId || chatUserName) {
        return {
          id: chatUserId ? Number(chatUserId) : 0,
          name: chatUserName || "",
        };
      }
      return null;
    },
    tryOpenPendingChat() {
      if (!this.pendingDirectChat) return;
      const { id, name } = this.pendingDirectChat;
      let member = null;
      if (id) {
        member = this.members.find((m) => m.user_id === id);
      }
      if (!member && name) {
        member = this.members.find((m) => m.name === name);
      }
      if (!member) return;
      this.pendingDirectChat = null;
      this.openPrivateChat(member);
    },
  },
};
</script>

<style scoped>
.study-room {
  min-height: 100vh;
  background: #f5f7fb;
  color: #1f2937;
  display: flex;
  flex-direction: column;
}

.room-header {
  height: 80px;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 28px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.04);
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-icon {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.brand-title {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
}

.brand-sub {
  font-size: 13px;
  color: #6b7280;
  margin-top: 4px;
}

.room-date {
  font-weight: 600;
  color: #374151;
}

.room-meta {
  display: flex;
  align-items: center;
  gap: 14px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #4b5563;
  font-size: 14px;
}

.ghost-btn {
  border: 1px solid #d1d5db;
  background: #fff;
  border-radius: 10px;
  padding: 8px 12px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: all 0.2s ease;
  color: #1f2937;
}

.ghost-btn:hover {
  border-color: #93c5fd;
  color: #2563eb;
}

.ghost-btn.small {
  padding: 6px 10px;
  font-size: 12px;
}

.room-main {
  display: flex;
  gap: 20px;
  padding: 20px 24px 28px;
  flex: 1;
}

.focus-area {
  flex: 7;
}

.sidebar {
  flex: 5;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.04);
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.card-header h2 {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
}

.muted {
  color: #6b7280;
  margin-top: 6px;
  font-size: 13px;
}

.chip {
  background: #eff6ff;
  color: #2563eb;
  border-radius: 999px;
  padding: 6px 12px;
  font-weight: 600;
  font-size: 13px;
}

.timer-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.timer-circle {
  width: 260px;
  height: 260px;
  border-radius: 50%;
  background: radial-gradient(circle at 30% 30%, #60a5fa, #2563eb 60%, #1d4ed8);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 20px 50px rgba(37, 99, 235, 0.35);
}

.timer-text {
  font-size: 46px;
  color: #fff;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
  letter-spacing: 2px;
  font-weight: 700;
}

.timer-actions {
  display: flex;
  gap: 14px;
}

.circle-btn {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: none;
  color: #fff;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.12);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.circle-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  box-shadow: none;
}

.circle-btn:hover:not(:disabled) {
  transform: translateY(-2px);
}

.play {
  background: #10b981;
}

.pause {
  background: #6b7280;
}

.reset {
  background: #2563eb;
}

.goal-forms {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
  margin-top: 12px;
}

.goal-form {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.goal-form label {
  font-size: 14px;
  color: #4b5563;
}

.goal-form input {
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 12px;
  font-size: 14px;
  outline: none;
  transition: border 0.2s ease, box-shadow 0.2s ease;
}

.goal-form input:focus {
  border-color: #93c5fd;
  box-shadow: 0 0 0 4px rgba(147, 197, 253, 0.3);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  margin-top: 18px;
}

.stat-card {
  border-radius: 14px;
  padding: 14px;
  color: #0f172a;
}

.stat-card.blue {
  background: #eff6ff;
  color: #1d4ed8;
}

.stat-card.green {
  background: #ecfdf3;
  color: #15803d;
}

.stat-label {
  font-size: 13px;
  color: #6b7280;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  margin-top: 4px;
}

.members-card .member-filters {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.filter-btn {
  border-radius: 999px;
  border: 1px solid #d1d5db;
  padding: 8px 12px;
  background: #fff;
  color: #4b5563;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 13px;
}

.filter-btn.active {
  background: #2563eb;
  color: #fff;
  border-color: #2563eb;
  box-shadow: 0 10px 20px rgba(37, 99, 235, 0.2);
}

.member-list {
  max-height: 320px;
  overflow-y: auto;
  padding-right: 4px;
}

.member-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border: 1px solid #f1f5f9;
  border-left: 4px solid #10b981;
  border-radius: 12px;
  margin-bottom: 10px;
  background: #fff;
  box-shadow: 0 6px 14px rgba(0, 0, 0, 0.03);
}

.avatar {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  flex-shrink: 0;
  border: none;
  padding: 0;
  cursor: pointer;
}

.avatar-1 {
  background: linear-gradient(135deg, #3b82f6, #6366f1);
}

.avatar-2 {
  background: linear-gradient(135deg, #10b981, #22d3ee);
}

.avatar-3 {
  background: linear-gradient(135deg, #f59e0b, #ef4444);
}

.avatar-4 {
  background: linear-gradient(135deg, #8b5cf6, #3b82f6);
}

.avatar-5 {
  background: linear-gradient(135deg, #ec4899, #f97316);
}

.avatar-6 {
  background: linear-gradient(135deg, #2dd4bf, #0ea5e9);
}

.avatar.small {
  width: 34px;
  height: 34px;
}

.member-meta {
  flex: 1;
}

.member-name {
  font-weight: 600;
  color: #111827;
}

.member-time {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #6b7280;
  margin-top: 4px;
}

.status-tag {
  border-radius: 12px;
  padding: 6px 10px;
  font-size: 12px;
  font-weight: 700;
}

.tag-focus {
  background: #ecfdf3;
  color: #15803d;
}

.tag-rest {
  background: #f3f4f6;
  color: #4b5563;
}

.online-pill {
  background: #eff6ff;
  color: #2563eb;
  border-radius: 999px;
  padding: 6px 12px;
  font-weight: 600;
}

.empty-state {
  text-align: center;
  color: #9ca3af;
  padding: 16px 0;
  font-size: 14px;
}

.chat-card {
  display: flex;
  flex-direction: column;
  height: 420px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding-right: 4px;
}

.chat-group + .chat-group {
  margin-top: 12px;
}

.chat-time {
  text-align: center;
  font-size: 12px;
  color: #9ca3af;
  margin: 8px 0;
}

.chat-message {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  align-items: flex-start;
}

.chat-message.self {
  flex-direction: row-reverse;
}

.bubble {
  background: #f3f4f6;
  border-radius: 12px;
  padding: 10px 12px;
  max-width: 420px;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.04);
}

.chat-message.self .bubble {
  background: #e0f2fe;
}

.bubble-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.bubble-name {
  font-weight: 700;
  color: #111827;
}

.bubble-time {
  font-size: 12px;
  color: #6b7280;
}

.bubble-text {
  color: #1f2937;
  line-height: 1.5;
  word-break: break-word;
}

.chat-input {
  display: flex;
  gap: 10px;
  margin-top: 12px;
}

.chat-input input {
  flex: 1;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 12px;
  font-size: 14px;
  outline: none;
  transition: border 0.2s ease, box-shadow 0.2s ease;
}

.chat-input input:focus {
  border-color: #93c5fd;
  box-shadow: 0 0 0 4px rgba(147, 197, 253, 0.3);
}

.chat-input button {
  width: 88px;
  border: none;
  background: #2563eb;
  color: #fff;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.chat-input button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 10px 20px rgba(37, 99, 235, 0.18);
}

.chat-input button:disabled {
  background: #cbd5e1;
  cursor: not-allowed;
  box-shadow: none;
}

.direct-chat-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 60;
}

.direct-chat-panel {
  width: min(520px, 92vw);
  max-height: 80vh;
  background: #fff;
  border-radius: 18px;
  box-shadow: 0 20px 50px rgba(15, 23, 42, 0.2);
  display: flex;
  flex-direction: column;
  padding: 18px 18px 16px;
  border: 1px solid #e5e7eb;
}

.direct-chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.direct-chat-title {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
}

.direct-chat-sub {
  font-size: 12px;
  color: #6b7280;
  margin-top: 4px;
}

.direct-chat-body {
  flex: 1;
  overflow-y: auto;
  background: #f8fafc;
  border-radius: 12px;
  padding: 12px;
  margin-bottom: 12px;
}

.direct-chat-empty {
  text-align: center;
  color: #9ca3af;
  font-size: 13px;
  padding: 16px 0;
}

.direct-chat-message {
  display: flex;
  margin-bottom: 10px;
}

.direct-chat-message.self {
  justify-content: flex-end;
}

.direct-chat-bubble {
  background: #fff;
  border-radius: 12px;
  padding: 10px 12px;
  box-shadow: 0 8px 16px rgba(15, 23, 42, 0.06);
  max-width: 80%;
}

.direct-chat-message.self .direct-chat-bubble {
  background: #e0f2fe;
}

.direct-chat-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  font-size: 12px;
  margin-bottom: 6px;
  color: #6b7280;
}

.direct-chat-name {
  font-weight: 700;
  color: #111827;
}

.direct-chat-text {
  color: #111827;
  line-height: 1.5;
  word-break: break-word;
}

.direct-chat-input {
  display: flex;
  gap: 10px;
}

.direct-chat-input input {
  flex: 1;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 10px 12px;
  font-size: 14px;
  outline: none;
  transition: border 0.2s ease, box-shadow 0.2s ease;
}

.direct-chat-input input:focus {
  border-color: #93c5fd;
  box-shadow: 0 0 0 4px rgba(147, 197, 253, 0.3);
}

.direct-chat-input button {
  width: 72px;
  border: none;
  background: #2563eb;
  color: #fff;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.direct-chat-input button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 10px 20px rgba(37, 99, 235, 0.18);
}

.direct-chat-input button:disabled {
  background: #cbd5e1;
  cursor: not-allowed;
  box-shadow: none;
}

.member-list,
.chat-messages {
  scrollbar-width: thin;
  scrollbar-color: #d1d5db transparent;
}

.member-list::-webkit-scrollbar,
.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.member-list::-webkit-scrollbar-thumb,
.chat-messages::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 999px;
}
</style>
