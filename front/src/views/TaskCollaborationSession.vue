<template>
  <div class="collab-page">
    <header class="collab-header">
      <button class="ghost-btn" type="button" @click="goBackToTeamTasks">
        返回团队任务
      </button>
      <div>
        <div class="eyebrow">任务协作会话</div>
        <h1>{{ taskTitle }}</h1>
      </div>
      <div class="status-pill" :class="{ ended: isEnded }">
        {{ session?.status_label || (isEnded ? "已结束" : "进行中") }}
      </div>
    </header>

    <main class="collab-layout">
      <section class="chat-panel">
        <div class="panel-title">
          <div>
            <h2>会话讨论</h2>
            <p>成员可以同步进展，也可以把我的笔记发成卡片。</p>
          </div>
          <button class="danger-btn" type="button" :disabled="isEnded || dismissing" @click="dismissSession">
            {{ dismissing ? "解散中..." : "静默解散" }}
          </button>
        </div>

        <div class="messages" ref="messagesEl">
          <div v-if="messagesLoading" class="empty">正在加载消息...</div>
          <div v-else-if="!messages.length" class="empty">还没有讨论内容</div>
          <article
            v-for="message in messages"
            :key="message.id"
            class="message"
            :class="{ self: Number(message.user_id) === Number(currentUserId) }"
          >
            <div class="message-meta">
              <strong>{{ message.display_name || `用户 ${message.user_id}` }}</strong>
              <span>{{ formatTime(message.sent_at) }}</span>
            </div>
            <div v-if="message.message_type === 'system'" class="system-bubble">{{ message.content }}</div>
            <button
              v-else-if="message.message_type === 'knowledge_card'"
              class="knowledge-card"
              type="button"
              @click="openSharedKnowledge(message)"
            >
              <div class="knowledge-card-title">{{ message.content?.title || "我的笔记" }}</div>
              <p>{{ message.content?.excerpt || message.content?.summary || "暂无摘要" }}</p>
              <div class="knowledge-card-foot">
                <span>{{ message.content?.category || "未分类" }}</span>
                <span>来源用户 {{ message.content?.source_user_id }}</span>
              </div>
            </button>
            <div v-else class="bubble">{{ message.content }}</div>
          </article>
        </div>

        <div class="composer" :class="{ disabled: isEnded }">
          <button type="button" class="ghost-btn" :disabled="isEnded" @click="showKnowledgePicker = true">
            我的笔记
          </button>
          <input
            v-model="draft"
            :disabled="isEnded"
            placeholder="输入协作消息..."
            @keyup.enter="sendMessage"
          />
          <button type="button" class="primary-btn" :disabled="isEnded || !draft.trim()" @click="sendMessage">
            发送
          </button>
        </div>
      </section>

      <aside class="minutes-panel">
        <div class="members-card">
          <h2>协作成员</h2>
          <div v-if="participants.length" class="participant-list">
            <span v-for="member in participants" :key="member.user_id" class="participant-pill">
              {{ member.display_name || `用户 ${member.user_id}` }}
            </span>
          </div>
          <p v-else class="empty compact-empty">暂无成员信息</p>
        </div>
        <div class="panel-title compact">
          <div>
            <h2>AI 纪要</h2>
            <p>总结同步内容、行动项和下一步。</p>
          </div>
        </div>
        <button class="primary-btn wide" type="button" :disabled="minutesLoading" @click="generateMinutes">
          {{ minutesLoading ? "生成中..." : "生成纪要" }}
        </button>

        <div v-if="minutesError" class="notice error">{{ minutesError }}</div>
        <div v-if="minutes" class="minutes-card">
          <h3>讨论摘要</h3>
          <p>{{ minutes.summary }}</p>
          <h3>同步知识</h3>
          <ul>
            <li v-for="item in normalizedList(minutes.synchronized_knowledge)" :key="item">{{ item }}</li>
          </ul>
          <h3>行动项</h3>
          <ul>
            <li v-for="item in normalizedActions" :key="`${item.owner}-${item.action}`">
              <strong>{{ item.owner }}</strong>：{{ item.action }}
            </li>
          </ul>
          <h3>阻塞点</h3>
          <ul>
            <li v-for="item in normalizedList(minutes.blockers, '暂无明确阻塞')" :key="item">{{ item }}</li>
          </ul>
          <h3>下一步</h3>
          <ul>
            <li v-for="item in normalizedList(minutes.next_steps)" :key="item">{{ item }}</li>
          </ul>
          <div class="minutes-actions">
            <button class="ghost-btn" type="button" :disabled="savingMinutes" @click="saveMinutes">
              {{ savingMinutes ? "保存中..." : "保存到任务" }}
            </button>
            <button class="ghost-btn" type="button" :disabled="savingKnowledge" @click="saveTeamKnowledge">
              {{ savingKnowledge ? "沉淀中..." : "沉淀团队知识" }}
            </button>
          </div>
        </div>
      </aside>
    </main>

    <div v-if="showKnowledgePicker" class="modal" @click.self="showKnowledgePicker = false">
      <div class="modal-card">
        <div class="panel-title compact">
          <div>
            <h2>选择我的笔记</h2>
            <p>发送后其他成员会在本次协作中看到笔记卡片。</p>
          </div>
          <button class="ghost-btn" type="button" @click="showKnowledgePicker = false">关闭</button>
        </div>
        <div v-if="knowledgeLoading" class="empty">正在加载我的笔记...</div>
        <div v-else-if="!knowledgeEntries.length" class="empty">暂无可分享笔记</div>
        <button
          v-for="entry in knowledgeEntries"
          :key="entry.id"
          class="knowledge-option"
          type="button"
          @click="shareKnowledge(entry)"
        >
          <strong>{{ entry.title }}</strong>
          <span>{{ entry.summary || entry.content || "暂无摘要" }}</span>
        </button>
      </div>
    </div>

    <div v-if="selectedKnowledgeDetail" class="modal" @click.self="selectedKnowledgeDetail = null">
      <div class="modal-card">
        <div class="panel-title compact">
          <div>
            <h2>{{ selectedKnowledgeDetail.title || "我的笔记" }}</h2>
            <p>来自用户 {{ selectedKnowledgeDetail.source_user_id }}</p>
          </div>
          <button class="ghost-btn" type="button" @click="selectedKnowledgeDetail = null">关闭</button>
        </div>
        <div class="knowledge-detail">
          <div v-if="selectedKnowledgeDetail.category" class="detail-tag">{{ selectedKnowledgeDetail.category }}</div>
          <p v-if="selectedKnowledgeDetail.summary" class="detail-summary">{{ selectedKnowledgeDetail.summary }}</p>
          <pre>{{ selectedKnowledgeDetail.content || selectedKnowledgeDetail.excerpt || "暂无内容" }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {
  dismissTaskCollaborationSession,
  generateTaskCollaborationMinutes,
  getCollaborationSharedKnowledge,
  getTaskCollaborationSession,
  getTaskDetail,
  saveTaskCollaborationMinutes,
  saveTaskCollaborationTeamKnowledge,
  shareKnowledgeToCollaborationSession,
} from "@/api/modules/task";
import { getRoomChatHistory, sendRoomChatMessage } from "@/api/modules/study";
import { listUserKnowledge } from "@/api/modules/knowledge";
import { useCurrentUser } from "@/composables/useCurrentUser";

export default {
  name: "TaskCollaborationSession",
  setup() {
    const { profile, loadCurrentUser } = useCurrentUser();
    if (!profile.value) {
      loadCurrentUser().catch((error) => console.warn("加载当前用户失败", error));
    }
    return { currentUserProfile: profile };
  },
  data() {
    return {
      session: null,
      task: null,
      messages: [],
      draft: "",
      messagesLoading: false,
      minutesLoading: false,
      minutesError: "",
      minutes: null,
      savingMinutes: false,
      savingKnowledge: false,
      dismissing: false,
      showKnowledgePicker: false,
      knowledgeLoading: false,
      knowledgeEntries: [],
      selectedKnowledgeDetail: null,
    };
  },
  computed: {
    sessionId() {
      return this.$route.params.sessionId;
    },
    currentUserId() {
      return this.currentUserProfile?.id || this.currentUserProfile?.basic_info?.id || null;
    },
    taskTitle() {
      return this.task?.title || `协作会话 #${this.sessionId}`;
    },
    isEnded() {
      return Number(this.session?.status) === 2;
    },
    normalizedActions() {
      const items = Array.isArray(this.minutes?.action_items) ? this.minutes.action_items : [];
      return items.length ? items : [{ owner: "待确认", action: "根据讨论补充下一步" }];
    },
    participants() {
      return Array.isArray(this.session?.participants) ? this.session.participants : [];
    },
  },
  watch: {
    showKnowledgePicker(value) {
      if (value) this.loadKnowledgeEntries();
    },
  },
  async mounted() {
    await this.loadSession();
  },
  methods: {
    async loadSession() {
      const res = await getTaskCollaborationSession(this.sessionId);
      this.session = res?.data?.session || res?.data?.data?.session || res?.session;
      if (this.session?.minutes) this.minutes = this.session.minutes;
      if (this.session?.task_id) {
        const taskRes = await getTaskDetail(this.session.task_id);
        this.task = taskRes?.data || taskRes?.data?.data || taskRes;
      }
      await this.loadMessages();
    },
    async loadMessages() {
      if (!this.session?.room_id) return;
      this.messagesLoading = true;
      try {
        const res = await getRoomChatHistory(this.session.room_id, { limit: 100 });
        const list = res?.data?.messages || res?.data?.data?.messages || [];
        this.messages = Array.isArray(list) ? [...list].reverse() : [];
        this.$nextTick(this.scrollToBottom);
      } finally {
        this.messagesLoading = false;
      }
    },
    async sendMessage() {
      const content = this.draft.trim();
      if (!content || this.isEnded) return;
      this.draft = "";
      await sendRoomChatMessage(this.session.room_id, {
        user_id: this.currentUserId || undefined,
        content,
      });
      await this.loadMessages();
    },
    async loadKnowledgeEntries() {
      this.knowledgeLoading = true;
      try {
        const res = await listUserKnowledge(1, 30);
        const data = res?.data || res?.data?.data || res;
        this.knowledgeEntries = data?.items || data?.list || data?.entries || [];
      } finally {
        this.knowledgeLoading = false;
      }
    },
    async shareKnowledge(entry) {
      if (this.isEnded) return;
      await shareKnowledgeToCollaborationSession(this.session.id, entry.id);
      this.showKnowledgePicker = false;
      await this.loadMessages();
    },
    async openSharedKnowledge(message) {
      const entryId = message?.content?.knowledge_entry_id;
      if (!entryId) return;
      const fallback = message.content;
      this.selectedKnowledgeDetail = {
        ...fallback,
        content: fallback?.excerpt || fallback?.summary || "",
      };
      try {
        const res = await getCollaborationSharedKnowledge(this.session.id, entryId);
        const entry = res?.data?.entry || res?.data?.data?.entry || res?.entry;
        if (entry) this.selectedKnowledgeDetail = entry;
      } catch (error) {
        console.warn("加载共享笔记详情失败", error);
      }
    },
    async generateMinutes() {
      this.minutesLoading = true;
      this.minutesError = "";
      try {
        const res = await generateTaskCollaborationMinutes(this.session.id);
        this.minutes = res?.data?.minutes || res?.data?.data?.minutes || res?.minutes;
      } catch (error) {
        this.minutesError = error.message || "纪要生成失败";
      } finally {
        this.minutesLoading = false;
      }
    },
    async saveMinutes() {
      if (!this.minutes) return;
      this.savingMinutes = true;
      try {
        const res = await saveTaskCollaborationMinutes(this.session.id, this.minutes);
        this.session = res?.data?.session || this.session;
      } finally {
        this.savingMinutes = false;
      }
    },
    async saveTeamKnowledge() {
      if (!this.minutes) return;
      this.savingKnowledge = true;
      try {
        await saveTaskCollaborationTeamKnowledge(this.session.id, {
          title: `${this.taskTitle} 协作纪要`,
          content: this.formatMinutesForSave(this.minutes),
        });
      } finally {
        this.savingKnowledge = false;
      }
    },
    goBackToTeamTasks() {
      const query = {};
      if (this.$route.query.teamId) query.teamId = this.$route.query.teamId;
      if (this.$route.query.taskId) query.taskId = this.$route.query.taskId;
      this.$router.push({ name: "TeamTasks", query });
    },
    async dismissSession() {
      if (this.isEnded) return;
      this.dismissing = true;
      try {
        const res = await dismissTaskCollaborationSession(this.session.id);
        this.session = res?.data?.session || this.session;
      } finally {
        this.dismissing = false;
      }
    },
    normalizedList(value, fallback = "暂无") {
      if (Array.isArray(value) && value.length) return value;
      return [fallback];
    },
    formatTime(value) {
      if (!value) return "";
      return new Date(value).toLocaleString("zh-CN", { hour12: false });
    },
    scrollToBottom() {
      const el = this.$refs.messagesEl;
      if (el) el.scrollTop = el.scrollHeight;
    },
    formatMinutesForSave(minutes) {
      const lines = [];
      if (minutes?.summary) lines.push("## 讨论摘要", minutes.summary);
      const appendList = (title, items) => {
        if (!Array.isArray(items) || !items.length) return;
        lines.push(`## ${title}`);
        items.forEach((item) => lines.push(`- ${item}`));
      };
      appendList("同步内容", minutes?.synchronized_knowledge);
      if (Array.isArray(minutes?.action_items) && minutes.action_items.length) {
        lines.push("## 行动项");
        minutes.action_items.forEach((item) => lines.push(`- ${item.owner || "待确认"}：${item.action || ""}`));
      }
      appendList("阻塞点", minutes?.blockers);
      appendList("下一步", minutes?.next_steps);
      return lines.join("\n");
    },
  },
};
</script>

<style scoped>
.collab-page {
  min-height: 100vh;
  background: #f5f7fb;
  color: #1f2937;
  padding: 24px;
}

.collab-header,
.collab-layout,
.panel-title,
.composer,
.message-meta,
.knowledge-card-foot,
.minutes-actions {
  display: flex;
  align-items: center;
}

.collab-header {
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}

.eyebrow {
  font-size: 12px;
  color: #64748b;
  font-weight: 700;
}

h1,
h2,
h3,
p {
  margin: 0;
}

h1 {
  font-size: 24px;
}

.collab-layout {
  align-items: stretch;
  gap: 20px;
}

.chat-panel,
.minutes-panel,
.modal-card {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
}

.chat-panel {
  flex: 1;
  min-width: 0;
  padding: 18px;
}

.minutes-panel {
  width: 360px;
  padding: 18px;
}

.panel-title {
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.panel-title p {
  color: #64748b;
  font-size: 13px;
  margin-top: 4px;
}

.panel-title.compact {
  align-items: flex-start;
}

.status-pill {
  background: #e0f2fe;
  color: #0369a1;
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 700;
}

.status-pill.ended {
  background: #f1f5f9;
  color: #64748b;
}

.messages {
  height: 560px;
  overflow: auto;
  padding: 12px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.empty {
  color: #94a3b8;
  text-align: center;
  padding: 28px;
}

.message {
  max-width: 76%;
  margin-bottom: 12px;
}

.message.self {
  margin-left: auto;
}

.message-meta {
  justify-content: space-between;
  gap: 10px;
  color: #64748b;
  font-size: 12px;
  margin-bottom: 4px;
}

.bubble,
.knowledge-card,
.system-bubble {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 10px 12px;
  line-height: 1.6;
}

.message.self .bubble {
  background: #2563eb;
  color: #ffffff;
  border-color: #2563eb;
}

.knowledge-card {
  border-color: #bfdbfe;
  background: #eff6ff;
  width: 100%;
  cursor: pointer;
  text-align: left;
}

.knowledge-card:hover {
  border-color: #60a5fa;
  background: #dbeafe;
}

.system-bubble {
  background: #f1f5f9;
  border-color: #cbd5e1;
  color: #475569;
  font-size: 13px;
}

.knowledge-card-title {
  font-weight: 800;
  color: #1d4ed8;
  margin-bottom: 6px;
}

.knowledge-card-foot {
  justify-content: space-between;
  color: #64748b;
  font-size: 12px;
  margin-top: 8px;
}

.composer {
  gap: 10px;
  margin-top: 14px;
}

.composer input {
  flex: 1;
  min-width: 0;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  padding: 10px 12px;
}

.composer.disabled {
  opacity: 0.72;
}

.primary-btn,
.ghost-btn,
.danger-btn {
  border: 0;
  border-radius: 8px;
  padding: 9px 14px;
  font-weight: 700;
  cursor: pointer;
}

.primary-btn {
  background: #2563eb;
  color: #ffffff;
}

.primary-btn.wide {
  width: 100%;
  margin-bottom: 14px;
}

.ghost-btn {
  background: #eef2ff;
  color: #3730a3;
}

.danger-btn {
  background: #fee2e2;
  color: #b91c1c;
}

button:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.notice.error {
  background: #fef2f2;
  color: #b91c1c;
  padding: 10px;
  border-radius: 8px;
  margin-bottom: 12px;
}

.minutes-card {
  border-top: 1px solid #e5e7eb;
  padding-top: 12px;
}

.members-card {
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 14px;
  margin-bottom: 16px;
}

.members-card h2 {
  font-size: 16px;
  margin-bottom: 10px;
}

.participant-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.participant-pill {
  background: #ecfdf5;
  color: #047857;
  border: 1px solid #a7f3d0;
  border-radius: 999px;
  padding: 5px 10px;
  font-size: 12px;
  font-weight: 700;
}

.compact-empty {
  padding: 8px;
  text-align: left;
}

.minutes-card h3 {
  margin-top: 14px;
  font-size: 14px;
}

.minutes-card ul {
  margin: 8px 0 0;
  padding-left: 18px;
  color: #475569;
}

.minutes-actions {
  gap: 10px;
  margin-top: 16px;
  flex-wrap: wrap;
}

.modal {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.38);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  z-index: 50;
}

.modal-card {
  width: min(620px, 96vw);
  max-height: 80vh;
  overflow: auto;
  padding: 18px;
}

.knowledge-option {
  display: block;
  width: 100%;
  text-align: left;
  border: 1px solid #e5e7eb;
  background: #ffffff;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 10px;
}

.knowledge-option strong,
.knowledge-option span {
  display: block;
}

.knowledge-detail {
  color: #334155;
}

.detail-tag {
  display: inline-block;
  background: #eef2ff;
  color: #3730a3;
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 700;
  margin-bottom: 10px;
}

.detail-summary {
  background: #f8fafc;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 12px;
}

.knowledge-detail pre {
  white-space: pre-wrap;
  font-family: inherit;
  line-height: 1.7;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 12px;
}

.knowledge-option span {
  color: #64748b;
  margin-top: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 980px) {
  .collab-layout {
    flex-direction: column;
  }

  .minutes-panel {
    width: auto;
  }
}
</style>
