<template>
  <div class="studyroom-page">
    <div class="cyber-bg"></div>
    <div class="cyber-shell">
      <header class="hero">
        <div>
          <div class="hero-kicker">Cyber Study Dock</div>
          <h1 class="hero-title">在线自习室</h1>
          <p class="hero-sub">进入学习舱，选择你的专注频道，开启高能模式。</p>
        </div>
        <div class="hero-actions">
          <button class="ghost-btn" type="button" @click="enterRandomRoom" :disabled="roomsLoading || !rooms.length">
            <iconify-icon icon="mdi:shuffle-variant" width="16" height="16"></iconify-icon>
            随机进入一个自习室
          </button>
          <button
            @click="showCreateRoom = true"
            class="primary-btn"
          >
            <iconify-icon icon="mdi:plus" width="16" height="16"></iconify-icon>
            创建房间
          </button>
        </div>
      </header>


      <section class="insight-grid">
        <div class="panel-card meteor-card">
          <div class="panel-header">
            <div>
              <h2>流星进度条</h2>
              <p class="panel-sub">现实时间进度（08:00 - 22:00）</p>
            </div>
            <div class="goal-pill">14h 时段</div>
          </div>
          <div class="meteor-track">
            <div class="meteor-fill" :style="{ width: `${realTimeProgressPercent}%` }">
              <div class="meteor-head"></div>
            </div>
          </div>
          <div class="meteor-meta">
            <span>当前时间 {{ currentTimeLabel }}</span>
            <span>{{ realTimeProgressPercent }}%</span>
          </div>
        </div>
      </section>

      <section class="recommend-section">
        <div class="section-header">
          <div>
            <h2>推荐房间</h2>
            <p class="panel-sub">优先展示当前最热的自习舱</p>
          </div>
          <button class="ghost-btn small" type="button" @click="fetchRooms">刷新</button>
        </div>
        <div v-if="roomsLoading" class="panel-empty">正在扫描自习舱...</div>
        <div v-else-if="!recommendedRooms.length" class="panel-empty">暂时没有可推荐的房间</div>
        <div v-else class="room-grid recommend-grid">
          <div
            v-for="room in recommendedRooms"
            :key="room.id"
            class="room-card"
            @click="enterRoom(room.id)"
          >
            <div class="room-header">
              <div class="room-title">
                <h3>{{ room.name }}</h3>
                <span v-if="room.isPrivate" class="room-tag">私密</span>
              </div>
              <span class="room-status" :class="statusTone(room.status)">{{ room.status }}</span>
            </div>
            <p class="room-desc">{{ room.description }}</p>
            <div class="room-tags" v-if="room.tags?.length">
              <span v-for="tag in room.tags" :key="tag">{{ tag }}</span>
            </div>
          </div>
        </div>
      </section>

      <section class="room-section">
        <div class="section-header">
          <div>
            <h2>全部房间</h2>
            <p class="panel-sub">选择一个频道，开始你的专注航行。</p>
          </div>
        </div>
        <div v-if="roomsLoading" class="panel-empty">正在加载房间数据...</div>
        <div v-else-if="roomsError" class="panel-empty error">{{ roomsError }}</div>
        <div v-else-if="rooms.length === 0" class="panel-empty">暂无房间，点击右上角创建一个吧～</div>
        <div v-else class="room-grid">
          <div
            v-for="room in rooms"
            :key="room.id"
            class="room-card"
            @click="enterRoom(room.id)"
          >
            <div class="room-header">
              <div class="room-title">
                <h3>{{ room.name }}</h3>
                <span v-if="room.isPrivate" class="room-tag">私密</span>
              </div>
              <span class="room-status" :class="statusTone(room.status)">{{ room.status }}</span>
            </div>
            <p class="room-desc">{{ room.description }}</p>
            <div class="room-tags" v-if="room.tags?.length">
              <span v-for="tag in room.tags" :key="tag">{{ tag }}</span>
            </div>
          </div>
        </div>
      </section>
    </div>
    <aside class="motto-float">
      <div class="motto-title">学习广播</div>
      <div class="motto-text">{{ currentMotto }}</div>
      <button class="motto-btn" type="button" @click="nextMotto">换一个</button>
    </aside>
    <!-- 创建房间弹窗 -->
    <CreateRoom 
      :visible="showCreateRoom"
      @close="hideCreateRoomModal"
      @created="onRoomCreated"
    />

    <!-- 私密房间密码弹窗 -->
    <div v-if="passwordDialogVisible" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-30">
      <div class="bg-white rounded-lg shadow-lg w-80 p-5">
        <h3 class="text-lg font-semibold mb-3">输入房间密码</h3>
        <div v-if="passwordError" class="text-red-500 text-xs mb-2">{{ passwordError }}</div>
        <input
          v-model="passwordInput"
          type="password"
          class="w-full border border-gray-200 rounded px-3 py-2 mb-4 outline-none focus:border-blue-500"
          placeholder="请输入房间密码"
          @keydown.enter="confirmPassword"
        />
        <div class="flex justify-end gap-3 text-sm">
          <button class="px-3 py-1.5 rounded border border-gray-200" @click="cancelPassword">取消</button>
          <button class="px-3 py-1.5 rounded bg-blue-600 text-white" @click="confirmPassword">确认</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ElMessage } from 'element-plus'
import CreateRoom from '@/components/StudyRoom/CreateRoom.vue'
import { getStudyRooms, joinStudyRoom } from '@/api/modules/study'
import { useCurrentUser } from '@/composables/useCurrentUser'

export default {
  name: 'StudyRoom',
  components: {
    CreateRoom
  },
  setup() {
    const { profile } = useCurrentUser()
    return {
      profile
    }
  },
  data() {
    return {
      showCreateRoom: false,
      rooms: [],
      roomsLoading: false,
      roomsError: '',
      currentTime: new Date(),
      timeTicker: null,
      mottoTicker: null,
      mottos: [
        '稳一点，你会更强。',
        '今天的专注，明天的光。',
        '慢慢推进，就是最快的路。',
        '再坚持 10 分钟，胜利就会靠近。',
        '把心收回到当下，一切都会更清晰。',
        '哪怕是一小步，也是在前进。',
        '自律不是束缚，是对自己的温柔。',
        '你不需要完美，只需要开始。',
      ],
      mottoIndex: 0,
      passwordDialogVisible: false,
      passwordInput: '',
      pendingRoomId: null,
      passwordError: ''
    }
  },
  computed: {
    recommendedRooms() {
      return [...this.rooms]
        .filter((room) => !room.isPrivate)
        .sort((a, b) => b.currentUsers - a.currentUsers)
        .slice(0, 3)
    },
    currentTimeLabel() {
      return this.currentTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    },
    realTimeProgressPercent() {
      const startHour = 8
      const endHour = 22
      const totalMinutes = (endHour - startHour) * 60
      const minutes = this.currentTime.getHours() * 60 + this.currentTime.getMinutes()
      const currentMinutes = Math.min(
        Math.max(minutes - startHour * 60, 0),
        totalMinutes
      )
      return Math.round((currentMinutes / totalMinutes) * 100)
    },
    currentMotto() {
      return this.mottos[this.mottoIndex] || ''
    }
  },
  created() {
    this.fetchRooms()
    this.mottoIndex = Math.floor(Math.random() * this.mottos.length)
  },
  mounted() {
    this.timeTicker = setInterval(() => {
      this.currentTime = new Date()
    }, 30000)
    this.mottoTicker = setInterval(() => {
      this.nextMotto()
    }, 12000)
  },
  beforeUnmount() {
    if (this.timeTicker) {
      clearInterval(this.timeTicker)
    }
    if (this.mottoTicker) {
      clearInterval(this.mottoTicker)
    }
  },
  methods: {
    hideCreateRoomModal() {
      this.showCreateRoom = false
    },
    async fetchRooms() {
      this.roomsLoading = true
      this.roomsError = ''
      try {
        const res = await getStudyRooms()
        const rooms = res?.data?.rooms || res?.data || []
        this.rooms = rooms.map((room) => this.transformRoom(room)).filter(Boolean)
      } catch (error) {
        console.error('加载学习房间失败:', error)
        this.roomsError = error?.message || '加载学习房间失败'
        ElMessage.error(this.roomsError)
      } finally {
        this.roomsLoading = false
      }
    },
    transformRoom(room) {
      if (!room) return null
      return {
        id: room.id,
        name: room.name,
        status: room.status || '进行中',
        statusClass: room.status_class || 'bg-green-100 text-green-800',
        description: room.is_private ? '私密房间，加入后可见' : (room.description || '这个房间暂未填写介绍'),
        currentUsers: room.current_users ?? room.currentUsers ?? 0,
        maxUsers: room.max_users ?? room.maxUsers ?? 0,
        studyTime: room.study_time || room.studyTime || '0h',
        tags: room.tags || [],
        isPrivate: room.is_private ?? room.isPrivate ?? false
      }
    },
    nextMotto() {
      if (!this.mottos.length) return
      let next = Math.floor(Math.random() * this.mottos.length)
      if (this.mottos.length > 1) {
        while (next === this.mottoIndex) {
          next = Math.floor(Math.random() * this.mottos.length)
        }
      }
      this.mottoIndex = next
    },
    enterRoom(roomId) {
      const room = this.rooms.find(r => r.id === roomId)
      if (room?.isPrivate) {
        this.pendingRoomId = roomId
        this.passwordInput = ''
        this.passwordError = ''
        this.passwordDialogVisible = true
        return
      }
      this.tryJoinRoom(roomId, '')
    },
    enterRandomRoom() {
      const candidates = this.rooms.filter((room) => !room.isPrivate)
      if (!candidates.length) {
        ElMessage.warning('暂无可随机进入的房间')
        return
      }
      const pick = candidates[Math.floor(Math.random() * candidates.length)]
      this.enterRoom(pick.id)
    },
    statusTone(status = '') {
      const value = String(status)
      if (value.includes('进行') || value.toLowerCase().includes('active')) return 'active'
      if (value.includes('休息') || value.toLowerCase().includes('rest')) return 'rest'
      return 'idle'
    },
    onRoomCreated(roomData) {
      const normalized = this.transformRoom(roomData)
      if (normalized) {
        this.rooms = [normalized, ...this.rooms]
        this.enterRoom(normalized.id)
      }
      this.hideCreateRoomModal()
    },
    async tryJoinRoom(roomId, password) {
      try {
        await joinStudyRoom(roomId, { password, user_id: this.profile?.value?.id || 1 })
        this.passwordDialogVisible = false
        this.passwordError = ''
        this.$router.push({ name: 'VideoRoom', params: { roomId } })
      } catch (error) {
        const msg = error?.response?.data?.message || error?.message || '加入房间失败'
        this.passwordError = msg
      }
    },
    confirmPassword() {
      if (!this.passwordInput.trim()) {
        ElMessage.warning('请输入房间密码')
        return
      }
      const roomId = this.pendingRoomId
      this.tryJoinRoom(roomId, this.passwordInput.trim())
    },
    cancelPassword() {
      this.passwordDialogVisible = false
      this.passwordInput = ''
      this.pendingRoomId = null
      this.passwordError = ''
    }
  }
}
</script>

<style scoped>
.studyroom-page {
  position: relative;
  min-height: 100vh;
  padding: 28px 20px 60px;
  background: #f7f9fc;
  color: #0f172a;
  overflow: hidden;
  font-family: "Chakra Petch", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.cyber-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 20% 20%, rgba(59, 130, 246, 0.08), transparent 45%),
    radial-gradient(circle at 80% 10%, rgba(14, 165, 233, 0.08), transparent 40%),
    radial-gradient(circle at 60% 80%, rgba(249, 115, 22, 0.08), transparent 45%);
  opacity: 0.8;
}

.cyber-shell {
  position: relative;
  z-index: 1;
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.hero {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 24px;
  border-radius: 20px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 24px 50px rgba(15, 23, 42, 0.08);
}

.hero-kicker {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.2em;
  color: #2563eb;
  margin-bottom: 8px;
}

.hero-title {
  font-size: 32px;
  font-weight: 700;
  color: #0f172a;
}

.hero-sub {
  margin-top: 8px;
  color: #64748b;
  font-size: 14px;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.ghost-btn,
.primary-btn {
  border-radius: 999px;
  padding: 10px 18px;
  font-weight: 600;
  font-size: 13px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.ghost-btn {
  border: 1px solid #cbd5f5;
  background: #f8fafc;
  color: #1e293b;
}

.ghost-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 24px rgba(59, 130, 246, 0.2);
}

.primary-btn {
  border: none;
  background: linear-gradient(135deg, #2d5bff, #22d3ee);
  color: #0f172a;
}

.primary-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 16px 32px rgba(34, 211, 238, 0.35);
}

.poster-banner {
  padding: 22px 24px;
  border-radius: 18px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 18px 32px rgba(15, 23, 42, 0.08);
  position: relative;
  overflow: hidden;
}

.poster-banner::after {
  content: "";
  position: absolute;
  inset: 0;
  background: linear-gradient(120deg, rgba(59, 130, 246, 0.08), transparent 45%);
}

.poster-main {
  position: relative;
  z-index: 1;
  font-size: 26px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #0f172a;
}

.poster-sub {
  position: relative;
  z-index: 1;
  margin-top: 6px;
  font-size: 14px;
  color: #64748b;
}

.insight-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 20px;
}

.panel-card {
  padding: 20px;
  border-radius: 18px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.05);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.panel-header h2 {
  font-size: 18px;
  font-weight: 700;
  color: #0f172a;
}

.panel-sub {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}

.panel-empty {
  color: #94a3b8;
  font-size: 13px;
  padding: 16px 0;
}

.panel-empty.error {
  color: #fca5a5;
}

.timeline-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.timeline-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.timeline-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #2563eb;
  box-shadow: 0 0 8px rgba(37, 99, 235, 0.25);
  margin-top: 6px;
}

.timeline-title {
  font-weight: 600;
  color: #0f172a;
}

.timeline-meta {
  font-size: 12px;
  color: #64748b;
  margin-top: 2px;
}

.meteor-card {
  display: flex;
  flex-direction: column;
}

.goal-pill {
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(59, 130, 246, 0.15);
  color: #2563eb;
  font-size: 12px;
}

.meteor-track {
  position: relative;
  height: 14px;
  border-radius: 999px;
  background: rgba(148, 163, 184, 0.2);
  overflow: hidden;
}

.meteor-fill {
  height: 100%;
  background: linear-gradient(90deg, rgba(59, 130, 246, 0.2), #3b82f6, #60a5fa);
  position: relative;
  border-radius: 999px;
  transition: width 0.4s ease;
}

.meteor-head {
  position: absolute;
  right: -8px;
  top: 50%;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #ffffff;
  transform: translateY(-50%);
  box-shadow: 0 0 12px rgba(59, 130, 246, 0.4);
}

.meteor-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #64748b;
  margin-top: 8px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.section-header h2 {
  font-size: 20px;
  font-weight: 700;
  color: #0f172a;
}

.room-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.room-card {
  padding: 16px;
  border-radius: 16px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.room-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 20px 40px rgba(15, 23, 42, 0.4);
}

.room-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.room-title h3 {
  font-size: 16px;
  font-weight: 600;
  color: #0f172a;
}

.room-tag {
  margin-left: 8px;
  font-size: 10px;
  color: #f59e0b;
  padding: 2px 8px;
  border-radius: 999px;
  border: 1px solid rgba(245, 158, 11, 0.4);
}

.room-status {
  font-size: 11px;
  padding: 4px 10px;
  border-radius: 999px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.room-status.active {
  background: rgba(34, 197, 94, 0.2);
  color: #0f172a;
}

.room-status.rest {
  background: rgba(248, 113, 113, 0.2);
  color: #fecaca;
}

.room-status.idle {
  background: rgba(59, 130, 246, 0.2);
  color: #93c5fd;
}

.room-desc {
  color: #64748b;
  font-size: 13px;
  margin-bottom: 12px;
}

.room-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.room-tags span {
  font-size: 11px;
  padding: 3px 8px;
  border-radius: 999px;
  background: rgba(59, 130, 246, 0.12);
  color: #2563eb;
}

.room-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  color: #475569;
}

.room-meta div {
  display: flex;
  align-items: center;
  gap: 6px;
}

.motto-float {
  position: fixed;
  right: 24px;
  top: 30%;
  width: 220px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(15, 23, 42, 0.12);
  padding: 16px;
  z-index: 5;
}

.motto-title {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.18em;
  color: #2563eb;
  font-weight: 700;
}

.motto-text {
  margin-top: 10px;
  font-size: 15px;
  color: #0f172a;
  line-height: 1.6;
  font-weight: 600;
}

.motto-btn {
  margin-top: 12px;
  width: 100%;
  border: none;
  border-radius: 999px;
  padding: 8px 12px;
  background: #2563eb;
  color: #fff;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.motto-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 24px rgba(37, 99, 235, 0.25);
}

@media (max-width: 1100px) {
  .motto-float {
    position: static;
    width: 100%;
    margin-top: 8px;
  }
}

.record-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.record-item {
  padding: 14px 16px;
  border-radius: 14px;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(148, 163, 184, 0.2);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.record-title {
  font-weight: 600;
  color: #f8fafc;
}

.record-meta {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 4px;
}

.record-time {
  font-size: 12px;
  color: #94a3b8;
}

@media (max-width: 768px) {
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .hero-actions {
    width: 100%;
  }

  .ghost-btn,
  .primary-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
