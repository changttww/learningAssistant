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
            <div class="header-actions">
              <div class="chip">{{ timerMode === "rest" ? "休息中" : "番茄钟" }}</div>
              <button class="icon-btn" type="button" @click="toggleTimerSettings">
                <iconify-icon icon="mdi:tune-variant" width="18"></iconify-icon>
              </button>
            </div>
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
              <button class="circle-btn end" @click="endStudy">
                <iconify-icon icon="mdi:check-circle" width="22"></iconify-icon>
              </button>
              <button class="circle-btn reset" @click="resetTimer">
                <iconify-icon icon="mdi:refresh" width="22"></iconify-icon>
              </button>
            </div>
            <div class="timer-hint">
              {{ timerMode === "rest" ? "正在休息" : "正在专注" }} ·
              目标 {{ timerMode === "rest" ? restMinutes : focusMinutes }} 分钟
            </div>
          </div>

          <div v-if="timerSettingsVisible" class="timer-settings">
            <div class="settings-title">番茄钟设置</div>
            <div class="settings-grid">
              <label>
                默认目标时长（分钟）
                <input v-model.number="focusMinutes" type="number" min="10" step="5" />
              </label>
              <label>
                默认休息时间（分钟）
                <input v-model.number="restMinutes" type="number" min="5" step="5" />
              </label>
              <label>
                每日专注时间（分钟）
                <input v-model.number="dailyFocusTargetMinutes" type="number" min="10" step="10" />
              </label>
              <label>
                每日任务目标（个）
                <input v-model.number="dailyTaskTarget" type="number" min="1" step="1" />
              </label>
            </div>
          </div>

          <div class="goal-forms">
            <div class="goal-form">
              <label>学习目标</label>
              <div class="goal-input-wrap">
                <input
                  v-model="focusGoal"
                  type="text"
                  placeholder="输入当前学习任务"
                  @focus="openTaskDropdown"
                  @blur="closeTaskDropdown"
                />
                <button
                  v-if="focusGoal"
                  class="clear-btn"
                  type="button"
                  @mousedown.prevent
                  @click="clearFocusGoal"
                >
                  <iconify-icon icon="mdi:close-circle" width="16"></iconify-icon>
                </button>
              </div>
              <div
                v-if="taskDropdownOpen && filteredWeekTasks.length"
                class="task-dropdown"
                @mousedown.prevent
              >
                <div class="task-dropdown-header">
                  <div>
                    <div class="task-dropdown-title">未来一周任务</div>
                    <div class="task-dropdown-sub">
                      {{ focusGoal ? "匹配到" : "共" }} {{ filteredWeekTasks.length }} 项任务
                    </div>
                  </div>
                  <button class="ghost-btn small" type="button" @click="refreshPersonalTasks">
                    <iconify-icon icon="mdi:refresh" width="14"></iconify-icon>
                    刷新
                  </button>
                </div>
                <div class="task-dropdown-list">
                  <button
                    v-for="task in filteredWeekTasks"
                    :key="task.id"
                    class="task-dropdown-item"
                    type="button"
                    @mousedown.prevent="selectFocusTask(task)"
                  >
                    <div class="task-dropdown-main">
                      <div class="task-title">{{ task.title }}</div>
                      <div class="task-meta">
                        <span>{{ task.dateLabel }}</span>
                        <span class="task-dot">•</span>
                        <span>{{ task.statusLabel }}</span>
                      </div>
                    </div>
                    <div class="task-pill" :class="task.statusClass">{{ task.category }}</div>
                  </button>
                </div>
              </div>
              <div v-else-if="taskDropdownOpen" class="task-dropdown empty">
                <div class="task-empty">
                  {{ tasksLoading ? "正在加载任务..." : "未来一周暂无可用任务" }}
                </div>
              </div>
              <div v-if="selectedFocusTask" class="selected-task">
                <div class="selected-task-info">
                  <div class="selected-task-title">{{ selectedFocusTask.title }}</div>
                  <div class="selected-task-meta">
                    <span>{{ selectedFocusTask.dateLabel }}</span>
                    <span class="task-dot">•</span>
                    <span>{{ selectedFocusTask.statusLabel }}</span>
                  </div>
                </div>
                <button
                  class="task-complete-btn"
                  type="button"
                  :disabled="selectedFocusTask.status === 'completed'"
                  @click="openTaskCompleteModal(selectedFocusTask)"
                >
                  {{ selectedFocusTask.status === "completed" ? "已完成" : "完成任务" }}
                </button>
              </div>
            </div>
            <div class="goal-form">
              <label>目标时长（分钟）</label>
              <input v-model.number="focusMinutes" type="number" min="10" step="5" />
            </div>
          </div>

          <div class="atmosphere-card">
            <div class="atmosphere-header">
              <div>
                <div class="atmosphere-title">今日氛围</div>
                <div class="atmosphere-sub">把注意力交给当下的小进步</div>
              </div>
              <button class="ghost-btn small" type="button" @click="nextEncouragement">
                <iconify-icon icon="mdi:refresh" width="14"></iconify-icon>
                换一句
              </button>
            </div>
            <div class="atmosphere-quote">
              {{ currentEncouragement }}
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
            <div>
              <h2>自习室讨论</h2>
              <p class="muted">把讨论整理成要点，方便继续学习。</p>
            </div>
            <button class="ghost-btn small" type="button" @click="summarizeRoomChat" :disabled="aiLoading">
              <iconify-icon icon="mdi:card-text" width="14"></iconify-icon>
              总结当前讨论
            </button>
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

    <button class="ai-floating-btn" type="button" @click="toggleAiPanel">
      <iconify-icon icon="mdi:robot-happy-outline" width="44"></iconify-icon>
      <span>AI 助理</span>
    </button>
    <div v-if="aiPanelVisible" class="ai-panel">
      <div class="ai-panel-header">
        <div>
          <div class="ai-panel-title">AI 学习助理</div>
          <div class="ai-panel-sub">计划/问答一键搞定</div>
        </div>
        <div class="ai-panel-actions">
          <button
            class="ghost-btn small"
            type="button"
            :disabled="aiPlanLoading"
            @click="generatePlanWithAI"
          >
            <iconify-icon icon="mdi:calendar-star" width="14"></iconify-icon>
            {{ aiPlanLoading ? "生成中..." : "帮我生成学习计划" }}
          </button>
          <button class="ghost-btn small" type="button" @click="toggleAiPanel">
            <iconify-icon icon="mdi:close" width="14"></iconify-icon>
            关闭
          </button>
        </div>
      </div>
      <div class="ai-messages" ref="aiMessagesContainer">
        <div v-if="!aiMessages.length" class="ai-empty">
          试试问我：今天适合先学什么？
        </div>
        <div
          v-for="message in aiMessages"
          :key="message.id"
          class="ai-message"
          :class="message.role"
        >
          <div class="ai-avatar">
            {{ message.role === "user" ? "👤" : "🤖" }}
          </div>
          <div class="ai-bubble">
            <div class="ai-text">{{ message.content }}</div>
            <div v-if="message.plan" class="ai-plan">
              <div class="ai-plan-title">今日学习计划</div>
              <div class="ai-plan-summary">{{ message.plan.summary }}</div>
              <div class="ai-plan-recommendation">
                {{ message.plan.recommendation }}
              </div>
              <div class="ai-plan-list">
                <div
                  v-for="(slot, index) in message.plan.schedule"
                  :key="`${message.id}-${index}`"
                  class="ai-plan-item"
                  :class="`plan-${slot.type}`"
                >
                  <div class="ai-plan-time">{{ slot.start }} - {{ slot.end }}</div>
                  <div class="ai-plan-body">
                    <div class="ai-plan-name">{{ slot.title }}</div>
                    <div class="ai-plan-note">{{ slot.notes || slot.taskTitle }}</div>
                  </div>
                </div>
              </div>
              <div v-if="message.plan.tips?.length" class="ai-plan-tips">
                <span>小贴士：</span>
                <span>{{ message.plan.tips.join(" / ") }}</span>
              </div>
            </div>
          </div>
        </div>
        <div v-if="aiLoading" class="ai-message assistant">
          <div class="ai-avatar">🤖</div>
          <div class="ai-bubble">
            <div class="ai-text">思考中...</div>
          </div>
        </div>
      </div>
      <div class="ai-input">
        <input
          v-model="aiInput"
          type="text"
          placeholder="输入你的问题..."
          @keydown.enter="sendAiMessage"
        />
        <button :disabled="aiLoading || !aiInput.trim()" @click="sendAiMessage">
          发送
        </button>
      </div>
    </div>

    <div v-if="taskCompleteModalVisible" class="task-complete-overlay" @click.self="closeTaskCompleteModal">
      <div class="task-complete-panel">
        <div class="task-complete-header">
          <div>
            <div class="task-complete-title">🎉 任务完成提示</div>
            <div class="task-complete-sub">确认完成后将同步到个人任务</div>
          </div>
          <button class="ghost-btn small" type="button" @click="closeTaskCompleteModal">
            <iconify-icon icon="mdi:close" width="16"></iconify-icon>
            关闭
          </button>
        </div>
        <div class="task-complete-body">
          <div class="task-complete-card">
            <div class="task-complete-name">{{ completingTask?.title }}</div>
            <div class="task-complete-desc">
              {{ completingTask?.description || "暂无任务描述" }}
            </div>
            <div class="task-complete-grid">
              <div>
                <div class="task-complete-label">学习时间</div>
                <div class="task-complete-value">{{ formattedStudyDuration }}</div>
              </div>
              <div>
                <div class="task-complete-label">任务时间</div>
                <div class="task-complete-value">{{ completingTask?.dateLabel || "-" }}</div>
              </div>
              <div>
                <div class="task-complete-label">剩余任务</div>
                <div class="task-complete-value">{{ remainingWeekTasks }} 项</div>
              </div>
            </div>
          </div>
        </div>
        <div class="task-complete-actions">
          <button class="ghost-btn" type="button" @click="closeTaskCompleteModal">稍后再说</button>
          <button
            class="primary-btn"
            type="button"
            :disabled="completingTaskSubmitting"
            @click="confirmTaskCompletion"
          >
            {{ completingTaskSubmitting ? "提交中..." : "确认完成并同步" }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="dailyLimitModalVisible" class="task-complete-overlay" @click.self="closeDailyLimitModal">
      <div class="task-complete-panel">
        <div class="task-complete-header">
          <div>
            <div class="task-complete-title">🌙 今日学习完成啦</div>
            <div class="task-complete-sub">已经达到今日小目标，记得好好休息哦～</div>
          </div>
          <button class="ghost-btn small" type="button" @click="closeDailyLimitModal">
            <iconify-icon icon="mdi:close" width="16"></iconify-icon>
            关闭
          </button>
        </div>
        <div class="task-complete-body">
          <div class="task-complete-card">
            <div class="task-complete-name">今日专注：{{ todayStudyLabel }}</div>
            <div class="task-complete-desc">
              已完成 {{ tasksCompletedToday }} / {{ dailyTaskTarget }} 个任务
            </div>
          </div>
        </div>
        <div class="task-complete-actions">
          <button class="ghost-btn" type="button" @click="closeDailyLimitModal">今天先收工</button>
          <button class="primary-btn" type="button" @click="continueAfterLimit">继续学习</button>
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
import { getPersonalTasks, completeTask } from "@/api/modules/task";
import { chatWithAI, generateStudyPlan } from "@/api/modules/ai";
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
      restMinutes: 10,
      dailyFocusMinutes: 0,
      dailyFocusTargetMinutes: 120,
      tasksCompletedToday: 0,
      dailyTaskTarget: 3,
      dailyStatsDate: "",
      timerMode: "focus",
      elapsedSeconds: 0,
      timerInterval: null,
      timerRunning: false,
      memberFilter: "all",
      ws: null,
      wsConnected: false,
      directChatVisible: false,
      directChatInput: "",
      directChats: {},
      activeDirectChatId: null,
      activeDirectChatName: "",
      pendingDirectChat: null,
      timerSettingsVisible: false,
      personalTasks: [],
      tasksLoading: false,
      tasksError: "",
      taskDropdownOpen: false,
      selectedFocusTask: null,
      taskCompleteModalVisible: false,
      completingTask: null,
      completingTaskSubmitting: false,
      dailyLimitModalVisible: false,
      overrideDailyLimit: false,
      encouragementIndex: 0,
      aiMessages: [],
      aiInput: "",
      aiLoading: false,
      aiPlanLoading: false,
      aiPanelVisible: false,
      encouragements: [
        "你正在慢慢把今天变得更亮一点。就像一盏小灯，稳稳地亮着，照见每一步努力的细节。",
        "别急着追赶，稳定的节奏才是你的优势。每一个专注的呼吸，都在给未来的你积攒底气。",
        "你已经开始了，这就是今天最重要的胜利。把注意力放回任务上，其他的交给时间。",
        "今天的你只需要向前走一点点。把目标拆成下一步，然后完成它，就这么简单。",
        "你在做的事情，正在变成你。哪怕很小，也是在向更好的自己靠近。",
        "深呼吸一下，手边的任务会因为你的耐心而变得更清晰。继续，慢慢来。",
        "别担心速度，真正的进步是可持续的。你每一次专注，都是给自己的一份礼物。",
        "如果觉得累，就把注意力缩小到这一分钟。做完这一分钟，你就赢了。",
        "今天的努力不会喧哗，但会沉淀成你最可靠的力量。",
        "你正在打造属于自己的学习节奏。稳一点、轻一点、也更长久一点。",
        "欢迎来到专注结界，结界稳定运行中。把世界的喧闹调成静音，留下你和目标的对话。",
        "你是自己的主角卡，经验值正在悄悄上涨。每一次点击“开始”，都是剧情推进。",
        "现在进入副本：专注训练。怪物叫分心，但你有耐心护盾，稳稳输出。",
        "今天的主线任务是：把一件事做完。支线任务是：别忘了给自己一个微笑。",
        "把注意力当成魔力，持续施法 25 分钟，你会看到不一样的结局。",
        "你的自习伙伴们都在加载中，你的专注就是最亮的特效。",
        "这是你与未来自己的连线时间。每一次专注，都是给未来发出的“我做到了”。",
        "学习不是冷冰冰的指令，而是你在构建自己的世界观。每个知识点都是新的地图。",
        "专注模式启动！背景音乐只剩心跳与键盘声，进度条也在安静上涨。",
        "你正在攒能量条，满格之后就能释放“成就感”大招。",
      ],
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
      const targetSeconds = this.timerTargetSeconds;
      const remainingSeconds = Math.max(0, targetSeconds - this.elapsedSeconds);
      const { hours, minutes, seconds } = this.formatTimerParts(remainingSeconds);
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
    todayStudyLabel() {
      return this.formatHoursMinutes(this.dailyFocusMinutes * 60);
    },
    completedTasksLabel() {
      return `${this.tasksCompletedToday}/${this.dailyTaskTarget}`;
    },
    timerTargetSeconds() {
      const minutes = this.timerMode === "rest" ? this.restMinutes : this.focusMinutes;
      return Math.max(1, Number(minutes || 0)) * 60;
    },
    currentEncouragement() {
      return this.encouragements[this.encouragementIndex] || "";
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
    weekTasks() {
      const today = this.getLocalDateString(new Date());
      const endDate = new Date();
      endDate.setDate(endDate.getDate() + 6);
      const endStr = this.getLocalDateString(endDate);
      return this.personalTasks.filter((task) => {
        const start = task.startDate || task.date;
        const end = task.endDate || task.date;
        return start <= endStr && end >= today && task.status !== "completed";
      });
    },
    filteredWeekTasks() {
      const keyword = this.normalizeKeyword(this.focusGoal);
      if (!keyword) return this.weekTasks;
      return this.weekTasks.filter((task) => this.isTaskMatched(task, keyword));
    },
    remainingWeekTasks() {
      return this.weekTasks.length;
    },
  },
  mounted() {
    const roomId = this.$route.params.roomId;
    this.ensureDailyStatsDate();
    this.encouragementIndex = this.randomEncouragementIndex();
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
    this.refreshPersonalTasks();
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
      this.ensureDailyStatsDate();
      if (this.timerMode === "focus" && this.isDailyLimitReached() && !this.overrideDailyLimit) {
        this.dailyLimitModalVisible = true;
        return;
      }
      if (this.timerMode === "focus") {
        this.setCurrentUserResting(false);
      }
      this.timerRunning = true;
      this.timerInterval = setInterval(() => {
        this.elapsedSeconds += 1;
        if (this.elapsedSeconds >= this.timerTargetSeconds) {
          this.handleTimerComplete();
        }
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
      this.timerMode = "focus";
      this.setCurrentUserResting(false);
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
    toggleTimerSettings() {
      this.timerSettingsVisible = !this.timerSettingsVisible;
    },
    ensureDailyStatsDate() {
      const today = this.getLocalDateString(new Date());
      if (!this.dailyStatsDate || this.dailyStatsDate !== today) {
        this.dailyStatsDate = today;
        this.dailyFocusMinutes = 0;
        this.tasksCompletedToday = 0;
        this.overrideDailyLimit = false;
      }
    },
    isDailyLimitReached() {
      if (this.dailyFocusTargetMinutes > 0 && this.dailyFocusMinutes >= this.dailyFocusTargetMinutes) {
        return true;
      }
      if (this.dailyTaskTarget > 0 && this.tasksCompletedToday >= this.dailyTaskTarget) {
        return true;
      }
      return false;
    },
    handleTimerComplete() {
      this.stopTimerInterval();
      this.timerRunning = false;
      this.elapsedSeconds = this.timerTargetSeconds;
      if (this.timerMode === "focus") {
        this.handleFocusCompleted();
      } else {
        this.handleRestCompleted();
      }
    },
    handleFocusCompleted() {
      this.dailyFocusMinutes += Number(this.focusMinutes || 0);
      this.tasksCompletedToday += 1;
      this.overrideDailyLimit = false;
      ElMessage.success("专注结束啦～现在进入休息时间");
      this.startRestMode();
    },
    handleRestCompleted() {
      ElMessage.success("休息完成！准备开启下一轮专注吧");
      this.timerMode = "focus";
      this.elapsedSeconds = 0;
      this.setCurrentUserResting(false);
    },
    startRestMode() {
      if (this.timerRunning) {
        this.stopTimerInterval();
        this.timerRunning = false;
      }
      this.timerMode = "rest";
      this.elapsedSeconds = 0;
      this.setCurrentUserResting(true);
      this.startTimer();
    },
    setCurrentUserResting(isResting) {
      const index = this.members.findIndex((m) => m.user_id === this.currentUserId);
      if (index !== -1) {
        this.members[index].isResting = isResting;
        return;
      }
      this.members.push({
        id: this.currentUserId,
        user_id: this.currentUserId,
        name: this.currentUserName,
        role: "我",
        online: true,
        avatarType: (this.members.length % 6) + 1,
        isResting,
        focusTime: this.defaultMemberTime,
      });
    },
    randomEncouragementIndex() {
      if (!this.encouragements.length) return 0;
      return Math.floor(Math.random() * this.encouragements.length);
    },
    nextEncouragement() {
      if (!this.encouragements.length) return;
      let next = this.randomEncouragementIndex();
      if (this.encouragements.length > 1) {
        while (next === this.encouragementIndex) {
          next = this.randomEncouragementIndex();
        }
      }
      this.encouragementIndex = next;
    },
    async sendAiMessage() {
      const content = (this.aiInput || "").trim();
      if (!content || this.aiLoading) return;
      const messageId = `${Date.now()}-${Math.random().toString(16).slice(2)}`;
      this.aiMessages.push({ id: messageId, role: "user", content });
      this.aiInput = "";
      this.aiLoading = true;
      this.scrollAiToBottom();
      try {
        const response = await chatWithAI({ message: content });
        const reply = response?.data?.reply || response?.reply || "我还在整理思路，请再问一次～";
        this.aiMessages.push({
          id: `${messageId}-reply`,
          role: "assistant",
          content: reply,
        });
      } catch (error) {
        console.error("AI 对话失败:", error);
        this.aiMessages.push({
          id: `${messageId}-error`,
          role: "assistant",
          content: "我刚才卡壳了，再问一次我一定努力回答！",
        });
      } finally {
        this.aiLoading = false;
        this.scrollAiToBottom();
      }
    },
    async summarizeRoomChat() {
      if (this.aiLoading) return;
      const summaryPrompt = this.buildChatSummaryPrompt();
      if (!summaryPrompt) {
        ElMessage.warning("当前没有可总结的讨论内容");
        return;
      }
      const messageId = `${Date.now()}-summary`;
      this.aiMessages.push({
        id: messageId,
        role: "assistant",
        content: "复杂任务处理中,请稍后",
      });
      this.aiLoading = true;
      this.aiPanelVisible = true;
      this.scrollAiToBottom();
      try {
        const response = await chatWithAI({ message: summaryPrompt });
        const reply = response?.data?.reply || response?.reply || "总结结果暂时不可用，请稍后再试。";
        this.aiMessages.push({
          id: `${messageId}-result`,
          role: "assistant",
          content: reply,
        });
        this.aiMessages = this.aiMessages.filter((item) => item.id !== messageId);
      } catch (error) {
        console.error("总结讨论失败:", error);
        this.aiMessages.push({
          id: `${messageId}-error`,
          role: "assistant",
          content: "总结失败了，再试一次吧～",
        });
        this.aiMessages = this.aiMessages.filter((item) => item.id !== messageId);
      } finally {
        this.aiLoading = false;
        this.scrollAiToBottom();
      }
    },
    buildChatSummaryPrompt() {
      const recent = this.messages.slice(-80);
      if (!recent.length) return "";
      const lines = recent.map((item) => {
        const name = item.senderName || "成员";
        const time = item.time || "";
        return `[${time}] ${name}: ${item.content}`;
      });
      return [
        "你是一个学习助理，请根据以下自习室讨论内容，整理一份简洁的总结，输出要点、待解决问题和下一步建议。",
        "要求：",
        "1. 输出不超过 6 条要点。",
        "2. 提到每个要点时尽量对应到讨论内容。",
        "3. 语言友好、条理清晰。",
        "",
        "讨论内容：",
        lines.join("\n"),
      ].join("\n");
    },
    async generatePlanWithAI() {
      if (this.aiPlanLoading) return;
      this.aiPlanLoading = true;
      const now = new Date();
      const payload = {
        current_time: now.toLocaleString("zh-CN"),
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || "Asia/Shanghai",
        focus_minutes: this.focusMinutes,
        rest_minutes: this.restMinutes,
        meal_minutes: 45,
        post_meal_rest_minutes: 20,
        tasks: this.weekTasks.map((task) => ({
          title: task.title,
          description: task.description,
          startDate: task.startDate,
          endDate: task.endDate,
        })),
      };
      const messageId = `${Date.now()}-plan`;
      this.aiMessages.push({
        id: messageId,
        role: "user",
        content: "帮我生成学习计划",
      });
      this.aiMessages.push({
        id: `${messageId}-processing`,
        role: "assistant",
        content: "复杂任务处理中,请稍后",
      });
      this.aiPanelVisible = true;
      this.scrollAiToBottom();
      try {
        const response = await generateStudyPlan(payload);
        const plan = response?.data || response;
        const summary = plan?.summary || "学习计划已生成，请查看安排";
        this.aiMessages.push({
          id: `${messageId}-result`,
          role: "assistant",
          content: summary,
          plan,
        });
        this.aiMessages = this.aiMessages.filter((item) => item.id !== `${messageId}-processing`);
      } catch (error) {
        console.error("生成学习计划失败:", error);
        this.aiMessages.push({
          id: `${messageId}-error`,
          role: "assistant",
          content: "计划生成失败了，稍后再试试吧～",
        });
        this.aiMessages = this.aiMessages.filter((item) => item.id !== `${messageId}-processing`);
      } finally {
        this.aiPlanLoading = false;
        this.scrollAiToBottom();
      }
    },
    scrollAiToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.aiMessagesContainer;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },
    toggleAiPanel() {
      this.aiPanelVisible = !this.aiPanelVisible;
      if (this.aiPanelVisible) {
        this.scrollAiToBottom();
      }
    },
    endStudy() {
      if (this.timerRunning) {
        this.stopTimerInterval();
        this.timerRunning = false;
      }
      this.ensureDailyStatsDate();
      if (this.timerMode === "focus") {
        this.handleFocusCompleted();
        return;
      }
      this.handleRestCompleted();
    },
    getLocalDateString(date) {
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      return `${year}-${month}-${day}`;
    },
    normalizeKeyword(value) {
      return (value || "").trim().toLowerCase();
    },
    isTaskMatched(task, keyword) {
      const haystack = `${task.title || ""} ${task.description || ""}`.toLowerCase();
      return haystack.includes(keyword);
    },
    openTaskDropdown() {
      this.taskDropdownOpen = true;
    },
    closeTaskDropdown() {
      setTimeout(() => {
        this.taskDropdownOpen = false;
      }, 150);
    },
    clearFocusGoal() {
      this.focusGoal = "";
      this.selectedFocusTask = null;
    },
    selectFocusTask(task) {
      this.selectedFocusTask = task;
      this.focusGoal = task.title;
      this.taskDropdownOpen = false;
    },
    async refreshPersonalTasks() {
      this.tasksLoading = true;
      this.tasksError = "";
      try {
        const response = await getPersonalTasks();
        if (response?.code === 0) {
          const items = response.data || [];
          this.personalTasks = items.map((task) => this.normalizePersonalTask(task));
        } else {
          this.tasksError = response?.msg || "加载任务失败";
        }
      } catch (error) {
        console.error("加载个人任务失败:", error);
        this.tasksError = error?.message || "加载任务失败";
      } finally {
        this.tasksLoading = false;
      }
    },
    normalizePersonalTask(task) {
      const startDate = task.start_at
        ? new Date(task.start_at).toISOString().split("T")[0]
        : this.getLocalDateString(new Date());
      const endDate = task.due_at
        ? new Date(task.due_at).toISOString().split("T")[0]
        : startDate;
      const status = task.status === 2 ? "completed" : task.status === 1 ? "in-progress" : "pending";
      const dateLabel =
        startDate === endDate ? startDate : `${startDate} - ${endDate}`;
      return {
        id: task.id,
        title: task.title || "未命名任务",
        description: task.description || "",
        category: task.category?.name || "个人任务",
        date: startDate,
        startDate,
        endDate,
        status,
        statusLabel:
          status === "completed" ? "已完成" : status === "in-progress" ? "进行中" : "待开始",
        statusClass:
          status === "completed"
            ? "pill-complete"
            : status === "in-progress"
              ? "pill-progress"
              : "pill-pending",
        dateLabel,
      };
    },
    openTaskCompleteModal(task) {
      if (!task) return;
      this.completingTask = task;
      this.taskCompleteModalVisible = true;
    },
    closeTaskCompleteModal() {
      this.taskCompleteModalVisible = false;
      this.completingTask = null;
      this.completingTaskSubmitting = false;
    },
    async confirmTaskCompletion() {
      if (!this.completingTask || this.completingTaskSubmitting) return;
      this.completingTaskSubmitting = true;
      try {
        const response = await completeTask(this.completingTask.id);
        if (response?.code === 0 || response?.code === 200) {
          const task = this.personalTasks.find((t) => t.id === this.completingTask.id);
          if (task) {
            task.status = "completed";
            task.statusLabel = "已完成";
            task.statusClass = "pill-complete";
          }
          if (this.selectedFocusTask?.id === this.completingTask.id) {
            this.selectedFocusTask.status = "completed";
            this.selectedFocusTask.statusLabel = "已完成";
            this.selectedFocusTask.statusClass = "pill-complete";
          }
          ElMessage.success("任务已完成并同步");
          this.closeTaskCompleteModal();
          this.startRestMode();
        } else {
          throw new Error(response?.msg || "任务完成失败");
        }
      } catch (error) {
        console.error("完成任务失败:", error);
        ElMessage.error(error?.message || "任务完成失败，请稍后重试");
      } finally {
        this.completingTaskSubmitting = false;
      }
    },
    closeDailyLimitModal() {
      this.dailyLimitModalVisible = false;
    },
    continueAfterLimit() {
      this.overrideDailyLimit = true;
      this.dailyLimitModalVisible = false;
      this.startTimer();
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
  --ai-scale: 1.6;
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.icon-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  border: 1px solid #dbeafe;
  background: #eff6ff;
  color: #2563eb;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.icon-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 16px rgba(37, 99, 235, 0.15);
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

.timer-hint {
  font-size: 13px;
  color: #64748b;
}

.timer-settings {
  margin: 16px 0 6px;
  padding: 14px;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
}

.settings-title {
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 10px;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  font-size: 13px;
  color: #475569;
}

.settings-grid label {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.settings-grid input {
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 10px;
  font-size: 13px;
  outline: none;
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

.end {
  background: #f97316;
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

.goal-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.goal-input-wrap input {
  width: 100%;
  padding-right: 36px;
}

.clear-btn {
  position: absolute;
  right: 10px;
  border: none;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  padding: 0;
}

.clear-btn:hover {
  color: #2563eb;
}

.task-dropdown {
  margin-top: 10px;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  background: #fff;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
  overflow: hidden;
}

.task-dropdown.empty {
  padding: 14px;
}

.task-dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 14px;
  border-bottom: 1px solid #eef2f7;
  background: #f8fafc;
}

.task-dropdown-title {
  font-weight: 700;
  color: #111827;
  font-size: 14px;
}

.task-dropdown-sub {
  font-size: 12px;
  color: #6b7280;
  margin-top: 2px;
}

.task-dropdown-list {
  max-height: 220px;
  overflow-y: auto;
}

.task-dropdown-item {
  width: 100%;
  border: none;
  background: #fff;
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 14px;
  cursor: pointer;
  transition: background 0.2s ease;
  text-align: left;
}

.task-dropdown-item:hover {
  background: #f1f5ff;
}

.task-dropdown-main {
  flex: 1;
}

.task-title {
  font-weight: 600;
  color: #111827;
  margin-bottom: 4px;
}

.task-meta {
  font-size: 12px;
  color: #6b7280;
  display: flex;
  align-items: center;
  gap: 6px;
}

.task-dot {
  color: #cbd5f5;
}

.task-pill {
  align-self: center;
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 999px;
  font-weight: 600;
}

.pill-complete {
  background: #dcfce7;
  color: #15803d;
}

.pill-progress {
  background: #ffedd5;
  color: #c2410c;
}

.pill-pending {
  background: #e5e7eb;
  color: #4b5563;
}

.task-empty {
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
}

.selected-task {
  margin-top: 10px;
  padding: 12px 14px;
  border-radius: 14px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.selected-task-title {
  font-weight: 700;
  color: #111827;
}

.selected-task-meta {
  font-size: 12px;
  color: #6b7280;
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.task-complete-btn {
  border: none;
  background: #10b981;
  color: #fff;
  font-weight: 700;
  border-radius: 10px;
  padding: 8px 12px;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
  min-width: 84px;
}

.task-complete-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 8px 18px rgba(16, 185, 129, 0.3);
}

.task-complete-btn:disabled {
  background: #cbd5e1;
  cursor: not-allowed;
  box-shadow: none;
}

.task-complete-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 70;
}

.task-complete-panel {
  width: min(560px, 92vw);
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 30px 60px rgba(15, 23, 42, 0.25);
  border: 1px solid #e5e7eb;
  padding: 18px 20px 16px;
}

.task-complete-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.task-complete-title {
  font-size: 18px;
  font-weight: 700;
  color: #0f172a;
}

.task-complete-sub {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}

.task-complete-card {
  background: #f8fafc;
  border-radius: 16px;
  padding: 16px;
  border: 1px solid #e2e8f0;
}

.task-complete-name {
  font-size: 16px;
  font-weight: 700;
  color: #111827;
}

.task-complete-desc {
  color: #64748b;
  font-size: 13px;
  margin-top: 6px;
  line-height: 1.5;
}

.task-complete-grid {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}

.task-complete-label {
  font-size: 12px;
  color: #94a3b8;
}

.task-complete-value {
  font-weight: 700;
  color: #0f172a;
  margin-top: 4px;
}

.task-complete-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
}

.primary-btn {
  border: none;
  background: #2563eb;
  color: #fff;
  padding: 10px 16px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.primary-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 12px 24px rgba(37, 99, 235, 0.25);
}

.primary-btn:disabled {
  background: #cbd5e1;
  cursor: not-allowed;
  box-shadow: none;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  margin-top: 18px;
}

.atmosphere-card {
  margin-top: 18px;
  padding: 16px;
  border-radius: 16px;
  background: linear-gradient(135deg, #eff6ff, #ecfeff);
  border: 1px solid #dbeafe;
  box-shadow: 0 12px 30px rgba(59, 130, 246, 0.12);
}

.atmosphere-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.atmosphere-title {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
}

.atmosphere-sub {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}

.atmosphere-quote {
  font-size: 15px;
  color: #1e293b;
  line-height: 1.7;
  background: rgba(255, 255, 255, 0.85);
  border-radius: 14px;
  padding: 14px 16px;
  border: 1px solid rgba(148, 163, 184, 0.3);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.9);
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

.ai-panel {
  position: fixed;
  right: calc(24px * var(--ai-scale));
  bottom: calc(92px * var(--ai-scale));
  width: min(calc(420px * var(--ai-scale)), 96vw);
  height: min(calc(420px * var(--ai-scale)), 80vh);
  max-height: min(calc(420px * var(--ai-scale)), 80vh);
  background: #fff;
  border-radius: calc(18px * var(--ai-scale));
  border: 1px solid #e5e7eb;
  box-shadow: 0 24px 50px rgba(15, 23, 42, 0.2);
  display: flex;
  flex-direction: column;
  z-index: 80;
}

.ai-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: calc(12px * var(--ai-scale));
  padding: calc(14px * var(--ai-scale)) calc(16px * var(--ai-scale))
    calc(10px * var(--ai-scale));
  border-bottom: 1px solid #eef2f7;
}

.ai-panel-title {
  font-weight: 700;
  color: #0f172a;
  font-size: calc(16px * var(--ai-scale));
}

.ai-panel-sub {
  font-size: calc(12px * var(--ai-scale));
  color: #64748b;
  margin-top: calc(4px * var(--ai-scale));
}

.ai-panel-actions {
  display: flex;
  align-items: center;
  gap: calc(8px * var(--ai-scale));
}

.ai-messages {
  flex: 1;
  overflow-y: auto;
  padding: calc(12px * var(--ai-scale)) calc(14px * var(--ai-scale)) 0;
}

.ai-empty {
  text-align: center;
  color: #9ca3af;
  font-size: calc(13px * var(--ai-scale));
  padding: calc(16px * var(--ai-scale)) 0;
}

.ai-message {
  display: flex;
  gap: calc(10px * var(--ai-scale));
  margin-bottom: calc(12px * var(--ai-scale));
  align-items: flex-start;
}

.ai-message.user {
  flex-direction: row-reverse;
}

.ai-avatar {
  width: calc(32px * var(--ai-scale));
  height: calc(32px * var(--ai-scale));
  border-radius: 50%;
  background: #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: calc(16px * var(--ai-scale));
}

.ai-bubble {
  background: #f1f5f9;
  border-radius: calc(12px * var(--ai-scale));
  padding: calc(10px * var(--ai-scale)) calc(12px * var(--ai-scale));
  max-width: 100%;
  box-shadow: 0 8px 16px rgba(15, 23, 42, 0.06);
}

.ai-message.user .ai-bubble {
  background: #e0f2fe;
}

.ai-text {
  color: #1f2937;
  line-height: 1.6;
  font-size: calc(13px * var(--ai-scale));
}

.ai-plan {
  margin-top: calc(10px * var(--ai-scale));
  border-radius: calc(12px * var(--ai-scale));
  padding: calc(12px * var(--ai-scale));
  background: #fff;
  border: 1px solid #e2e8f0;
}

.ai-plan-title {
  font-weight: 700;
  color: #0f172a;
  font-size: calc(14px * var(--ai-scale));
  margin-bottom: calc(6px * var(--ai-scale));
}

.ai-plan-summary {
  font-size: calc(13px * var(--ai-scale));
  color: #475569;
  margin-bottom: calc(6px * var(--ai-scale));
}

.ai-plan-recommendation {
  font-size: calc(12px * var(--ai-scale));
  color: #2563eb;
  margin-bottom: calc(10px * var(--ai-scale));
}

.ai-plan-list {
  display: flex;
  flex-direction: column;
  gap: calc(8px * var(--ai-scale));
}

.ai-plan-item {
  display: flex;
  gap: calc(12px * var(--ai-scale));
  padding: calc(8px * var(--ai-scale)) calc(10px * var(--ai-scale));
  border-radius: calc(10px * var(--ai-scale));
  background: #f8fafc;
  border: 1px solid #e2e8f0;
}

.ai-plan-time {
  font-size: calc(12px * var(--ai-scale));
  color: #64748b;
  min-width: calc(90px * var(--ai-scale));
}

.ai-plan-name {
  font-weight: 600;
  color: #0f172a;
  font-size: calc(13px * var(--ai-scale));
}

.ai-plan-note {
  font-size: calc(12px * var(--ai-scale));
  color: #64748b;
  margin-top: calc(2px * var(--ai-scale));
}

.plan-study {
  border-left: 4px solid #2563eb;
}

.plan-break {
  border-left: 4px solid #10b981;
}

.plan-meal {
  border-left: 4px solid #f97316;
}

.plan-rest {
  border-left: 4px solid #8b5cf6;
}

.plan-buffer {
  border-left: 4px solid #94a3b8;
}

.ai-plan-tips {
  margin-top: calc(10px * var(--ai-scale));
  font-size: calc(12px * var(--ai-scale));
  color: #475569;
}

.ai-input {
  display: flex;
  gap: calc(10px * var(--ai-scale));
  margin: calc(12px * var(--ai-scale)) calc(14px * var(--ai-scale))
    calc(14px * var(--ai-scale));
}

.ai-input input {
  flex: 1;
  border: 1px solid #e5e7eb;
  border-radius: calc(10px * var(--ai-scale));
  padding: calc(10px * var(--ai-scale)) calc(12px * var(--ai-scale));
  font-size: calc(13px * var(--ai-scale));
  outline: none;
  transition: border 0.2s ease, box-shadow 0.2s ease;
}

.ai-input input:focus {
  border-color: #93c5fd;
  box-shadow: 0 0 0 4px rgba(147, 197, 253, 0.3);
}

.ai-input button {
  width: calc(72px * var(--ai-scale));
  border: none;
  background: #2563eb;
  color: #fff;
  border-radius: calc(10px * var(--ai-scale));
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.ai-input button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 10px 20px rgba(37, 99, 235, 0.18);
}

.ai-input button:disabled {
  background: #cbd5e1;
  cursor: not-allowed;
  box-shadow: none;
}

.ai-floating-btn {
  position: fixed;
  right: calc(24px * var(--ai-scale));
  bottom: calc(24px * var(--ai-scale));
  border: none;
  background: #2563eb;
  color: #fff;
  border-radius: 999px;
  padding: calc(12px * var(--ai-scale)) calc(18px * var(--ai-scale));
  display: inline-flex;
  align-items: center;
  gap: calc(8px * var(--ai-scale));
  font-weight: 700;
  font-size: calc(14px * var(--ai-scale));
  cursor: pointer;
  box-shadow: 0 16px 32px rgba(37, 99, 235, 0.3);
  z-index: 80;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.ai-floating-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 20px 36px rgba(37, 99, 235, 0.35);
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
.chat-messages,
.ai-messages {
  scrollbar-width: thin;
  scrollbar-color: #d1d5db transparent;
}

.member-list::-webkit-scrollbar,
.chat-messages::-webkit-scrollbar,
.ai-messages::-webkit-scrollbar {
  width: 6px;
}

.member-list::-webkit-scrollbar-thumb,
.chat-messages::-webkit-scrollbar-thumb,
.ai-messages::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 999px;
}
</style>
