<template>
  <div class="meeting-room">
    <header class="meeting-header">
      <div class="header-left">
        <div class="meeting-title">
          <div class="meeting-icon">
            <iconify-icon icon="mdi:video" width="20"></iconify-icon>
          </div>
          <div>
            <div class="title-text">快速会议室</div>
            <div class="title-sub">团队私有 · 仅限当前团队成员进入</div>
          </div>
        </div>
      </div>
      <div class="meeting-meta">
        <div class="meta-item">
          <iconify-icon icon="mdi:account-group" width="16"></iconify-icon>
          <span>{{ onlineCount }}人在线</span>
        </div>
        <div class="meta-item">
          <iconify-icon icon="mdi:calendar" width="16"></iconify-icon>
          <span>{{ formattedDate }}</span>
        </div>
      </div>
    </header>

    <main class="meeting-main">
      <section class="meeting-board">
        <div class="meeting-card">
          <div class="card-header">
            <div>
              <h2>会议概览</h2>
              <p class="muted">团队专属的临时会议空间，支持实时聊天与在线成员。</p>
            </div>
            <div class="pill">私有会议</div>
          </div>
          <div class="overview-grid">
            <div class="overview-item">
              <div class="overview-label">团队ID</div>
              <div class="overview-value">{{ teamIdLabel }}</div>
            </div>
            <div class="overview-item">
              <div class="overview-label">会议类型</div>
              <div class="overview-value">团队内部</div>
            </div>
            <div class="overview-item">
              <div class="overview-label">当前在线</div>
              <div class="overview-value">{{ onlineCount }}人</div>
            </div>
          </div>
          <div class="notice-card">
            <h3>会议提醒</h3>
            <ul>
              <li>本会议室不会出现在公共自习室列表。</li>
              <li>仅支持团队成员进入与聊天。</li>
              <li>需要讨论时可直接在右侧发消息。</li>
            </ul>
          </div>
        </div>
      </section>

      <aside class="meeting-sidebar">
        <div class="meeting-card members-card">
          <div class="card-header">
            <h2>在线成员</h2>
            <div class="pill ghost">{{ onlineCount }}人在线</div>
          </div>
          <div class="member-list">
            <div class="member-row" v-for="member in members" :key="member.id">
              <div class="avatar" :class="`avatar-${member.avatarType}`"></div>
              <div class="member-meta">
                <div class="member-name">{{ member.name }}</div>
                <div class="member-status">{{ member.statusLabel }}</div>
              </div>
              <span class="status-tag" :class="member.statusClass">{{ member.statusText }}</span>
            </div>
            <div class="empty-state" v-if="!members.length">等待成员进入会议室...</div>
          </div>
        </div>

        <div class="meeting-card chat-card">
          <div class="card-header">
            <h2>会议聊天</h2>
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
              placeholder="输入会议消息..."
              @keydown.enter="sendMessage"
            />
            <button :disabled="!newMessage.trim()" @click="sendMessage">发送</button>
          </div>
        </div>
      </aside>
    </main>
  </div>
</template>

<script>
import { computed, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { useCurrentUser } from "@/composables/useCurrentUser";
import { getRoomChatHistory, sendRoomChatMessage } from "@/api/modules/study";
import { apiConfig } from "@/config";

export default {
  name: "TeamMeetingRoom",
  setup() {
    const { profile, loadCurrentUser } = useCurrentUser();

    onMounted(() => {
      loadCurrentUser().catch((error) => {
        console.error("加载用户信息失败:", error);
      });
    });

    const currentUserName = computed(
      () => profile.value?.display_name || "会议成员"
    );
    const currentUserId = computed(() => profile.value?.id || 1);

    return {
      currentUserName,
      currentUserId,
    };
  },
  data() {
    return {
      members: [],
      messages: [],
      newMessage: "",
      ws: null,
      wsConnected: false,
    };
  },
  computed: {
    teamIdLabel() {
      return this.$route.params.teamId || "-";
    },
    teamIdValue() {
      const raw = this.$route.params.teamId;
      const numeric = Number(raw);
      return Number.isFinite(numeric) && numeric > 0 ? numeric : 0;
    },
    onlineCount() {
      return this.members.length;
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
    formattedDate() {
      return new Date().toLocaleDateString("zh-CN", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        weekday: "long",
      });
    },
  },
  mounted() {
    if (!this.ensureEntry()) {
      return;
    }
    if (this.teamIdValue) {
      this.loadChatHistory();
    }
    this.connectWebSocket();
  },
  beforeUnmount() {
    if (this.ws) {
      this.ws.close();
    }
    this.clearEntryFlag();
  },
  methods: {
    ensureEntry() {
      const currentId = String(this.$route.params.teamId || "");
      let allowed = "";
      try {
        allowed = sessionStorage.getItem("team:quickMeeting") || "";
      } catch (error) {
        console.warn("无法读取会议室访问标记", error);
      }
      if (!currentId || !allowed || allowed !== currentId) {
        alert("该会议室仅可从团队任务入口进入");
        this.$router.push("/team-tasks");
        return false;
      }
      return true;
    },
    clearEntryFlag() {
      try {
        sessionStorage.removeItem("team:quickMeeting");
      } catch (error) {
        console.warn("无法清理会议室访问标记", error);
      }
    },
    wsUrl() {
      const base = new URL(apiConfig.baseURL);
      const protocol = base.protocol === "https:" ? "wss:" : "ws:";
      const host = base.host;
      const params = new URLSearchParams({
        user_id: this.currentUserId,
        display_name: this.currentUserName,
      });
      return `${protocol}//${host}/api/study/rooms/${this.teamIdValue}/ws?${params.toString()}`;
    },
    connectWebSocket() {
      if (!this.teamIdValue) return;
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
        default:
          break;
      }
    },
    applyState(data) {
      if (!data?.members) return;
      this.members = data.members.map((m, index) => ({
        id: m.user_id,
        user_id: m.user_id,
        name: m.display_name || "成员",
        avatarType: (index % 6) + 1,
        statusText: "在线",
        statusLabel: "会议中",
        statusClass: "tag-focus",
      }));
    },
    applyMemberJoined(data) {
      if (!data) return;
      const exists = this.members.find((m) => m.user_id === data.user_id);
      if (!exists) {
        this.members.push({
          id: data.user_id,
          user_id: data.user_id,
          name: data.display_name || "成员",
          avatarType: (this.members.length % 6) + 1,
          statusText: "在线",
          statusLabel: "会议中",
          statusClass: "tag-focus",
        });
      }
    },
    async sendMessage() {
      const content = this.newMessage.trim();
      if (!content) return;
      if (!this.wsConnected) {
        ElMessage.warning("聊天连接未建立，请稍后重试");
        return;
      }
      this.sendWs("chat", { content });
      if (this.teamIdValue) {
        try {
          await sendRoomChatMessage(this.teamIdValue, {
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
      if (!this.teamIdValue) return;
      try {
        const res = await getRoomChatHistory(this.teamIdValue, { limit: 100 });
        const items = res?.data?.messages || [];
        this.messages = items
          .slice()
          .reverse()
          .map((item) => this.normalizeMessage(item));
        this.$nextTick(() => {
          const container = this.$refs.messagesContainer;
          if (container) container.scrollTop = container.scrollHeight;
        });
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
    normalizeMessage(item) {
      const sentAt = item.sent_at ? new Date(item.sent_at) : new Date();
      const timeStr = sentAt.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
      const dateLabel = sentAt.toLocaleDateString();
      return {
        id: item.id || `${item.user_id}-${item.sent_at}`,
        senderName: item.display_name || "成员",
        content: item.content || "",
        time: timeStr,
        timeGroup: `${dateLabel} ${timeStr}`,
        avatarType: (item.user_id % 6) + 1,
        isSelf: item.user_id === this.currentUserId,
      };
    },
  },
};
</script>

<style scoped>
.meeting-room {
  min-height: 100vh;
  background: linear-gradient(180deg, #f8fafc, #eef2ff);
  color: #111827;
  display: flex;
  flex-direction: column;
}

.meeting-header {
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  padding: 20px 28px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.04);
  flex-wrap: wrap;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.meeting-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.meeting-icon {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  background: linear-gradient(135deg, #2563eb, #1d4ed8);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.title-text {
  font-size: 18px;
  font-weight: 700;
}

.title-sub {
  font-size: 13px;
  color: #6b7280;
  margin-top: 4px;
}

.meeting-meta {
  display: flex;
  align-items: center;
  gap: 14px;
  color: #4b5563;
  font-size: 14px;
}

.meeting-meta .meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.meeting-main {
  display: flex;
  gap: 20px;
  padding: 20px 24px 28px;
  flex: 1;
}

.meeting-board {
  flex: 7;
}

.meeting-sidebar {
  flex: 5;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.meeting-card {
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

.pill {
  background: #eff6ff;
  color: #2563eb;
  border-radius: 999px;
  padding: 6px 12px;
  font-weight: 600;
  font-size: 12px;
}

.pill.ghost {
  background: #f3f4f6;
  color: #4b5563;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.overview-item {
  background: #f8fafc;
  border-radius: 12px;
  padding: 12px;
}

.overview-label {
  font-size: 12px;
  color: #6b7280;
}

.overview-value {
  font-size: 16px;
  font-weight: 700;
  margin-top: 6px;
}

.notice-card {
  background: #f8fafc;
  border-radius: 12px;
  padding: 14px;
}

.notice-card h3 {
  font-size: 15px;
  font-weight: 700;
  margin-bottom: 8px;
}

.notice-card ul {
  padding-left: 18px;
  color: #4b5563;
  line-height: 1.8;
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

.member-list {
  max-height: 260px;
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

.member-status {
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

@media (max-width: 1024px) {
  .meeting-main {
    flex-direction: column;
  }
}
</style>
