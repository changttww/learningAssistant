<template>
  <div class="container w-full h-full overflow-auto px-4">
    <div class="card">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl font-bold">在线自习室</h1>
        <button
          @click="showCreateRoom = true"
          class="bg-[#2D5BFF] text-white font-medium py-2 px-4 rounded-lg text-sm hover:bg-opacity-90 transition-colors flex items-center gap-2"
        >
          <iconify-icon icon="mdi:plus" width="16" height="16"></iconify-icon>
          创建房间
        </button>
      </div>

      <!-- 自习室统计 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div class="stat-card bg-blue-50">
          <div class="text-2xl font-bold text-blue-600">{{ summary.active_rooms }}</div>
          <div class="text-gray-600 mt-1">活跃房间</div>
        </div>
        <div class="stat-card bg-green-50">
          <div class="text-2xl font-bold text-green-600">{{ summary.online_users }}</div>
          <div class="text-gray-600 mt-1">在线用户</div>
        </div>
        <div class="stat-card bg-orange-50">
          <div class="text-2xl font-bold text-orange-600">{{ summary.today_study_hours }}h</div>
          <div class="text-gray-600 mt-1">今日学习</div>
        </div>
        <div class="stat-card bg-purple-50">
          <div class="text-2xl font-bold text-purple-600">{{ summary.streak_days }}</div>
          <div class="text-gray-600 mt-1">连续天数</div>
        </div>
      </div>

      <!-- 房间列表 -->
      <div>
        <div v-if="roomsLoading" class="py-12 text-center text-gray-500">
          正在加载房间数据...
        </div>
        <div v-else-if="roomsError" class="py-12 text-center text-red-500">
          {{ roomsError }}
        </div>
        <div v-else-if="rooms.length === 0" class="py-12 text-center text-gray-500">
          暂无房间，点击右上角创建一个吧～
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg-grid-cols-3 gap-4">
          <div
            v-for="room in rooms"
            :key="room.id"
            class="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer"
            @click="enterRoom(room.id)"
          >
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center gap-2">
                <h3 class="font-medium">{{ room.name }}</h3>
                <span
                  v-if="room.isPrivate"
                  class="text-xs px-2 py-0.5 rounded-full bg-gray-100 text-gray-600"
                >
                  私密
                </span>
              </div>
              <span :class="['px-2 py-1 rounded text-xs', room.statusClass]">
                {{ room.status }}
              </span>
            </div>
            <p class="text-sm text-gray-600 mb-3">{{ room.description }}</p>
            <div v-if="room.tags?.length" class="flex flex-wrap gap-2 mb-3">
              <span
                v-for="tag in room.tags"
                :key="tag"
                class="text-xs bg-blue-50 text-blue-600 px-2 py-0.5 rounded-full"
              >
                {{ tag }}
              </span>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <iconify-icon
                  icon="mdi:account-group"
                  width="16"
                  height="16"
                  class="text-gray-500"
                ></iconify-icon>
                <span class="text-sm">
                  {{ room.currentUsers }}/{{ formatCapacity(room.maxUsers) }}
                </span>
              </div>
              <div class="flex items-center gap-2">
                <iconify-icon
                  icon="mdi:clock"
                  width="16"
                  height="16"
                  class="text-gray-500"
                ></iconify-icon>
                <span class="text-sm">{{ room.studyTime }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 我的学习记录 -->
      <div class="mt-8">
        <h2 class="font-bold text-xl mb-4">我的学习记录</h2>
        <div class="space-y-3" v-if="records.length">
          <div
            v-for="record in records"
            :key="record.id"
            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
          >
            <div class="flex items-center gap-3">
              <div class="w-3 h-3 rounded-full bg-green-500"></div>
              <div>
                <span class="font-medium">{{ record.title || '学习记录' }}</span>
                <span class="text-sm text-gray-600 ml-2">学习了 {{ (record.duration / 60).toFixed(1) }} 小时</span>
              </div>
            </div>
            <span class="text-sm text-gray-500">{{ formatRecordTime(record.recorded_at) }}</span>
          </div>
        </div>
        <div v-else class="text-gray-500 text-sm">暂无学习记录</div>
      </div>
    </div>

    <!-- 创建房间弹窗 -->
    <CreateRoom
      :visible="showCreateRoom"
      @close="hideCreateRoomModal"
      @created="onRoomCreated"
    />

    <!-- 私密房间密码弹窗 -->
    <div v-if="passwordDialogVisible" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-30">
      <div class="bg-white rounded-lg shadow-lg w-80 p-5">
        <h3 class="text-lg font-semibold mb-1">输入房间密码</h3>
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
import { getStudyRooms, getStudySummary, getStudyRecords, joinStudyRoom } from '@/api/modules/study'
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
      summary: {
        active_rooms: 0,
        online_users: 0,
        today_study_hours: 0,
        streak_days: 0
      },
      records: [],
      recordsLoading: false,
      passwordDialogVisible: false,
      passwordInput: '',
      pendingRoomId: null,
      passwordError: ''
    }
  },
  created() {
    this.fetchRooms()
    this.fetchSummary()
    this.fetchRecords()
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
    formatCapacity(maxUsers) {
      if (!maxUsers || maxUsers <= 0) {
        return '不限'
      }
      return maxUsers
    },
    formatRecordTime(timeStr) {
      if (!timeStr) return ''
      const date = new Date(timeStr)
      const today = new Date()
      const isToday = date.toDateString() === today.toDateString()
      const label = isToday ? '今天' : date.toLocaleDateString()
      const time = date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
      return `${label} ${time}`
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
    async fetchSummary() {
      try {
        const res = await getStudySummary({ user_id: this.profile?.value?.id || 1 })
        this.summary = res?.data || this.summary
      } catch (error) {
        console.error('获取学习汇总失败:', error)
      }
    },
    async fetchRecords() {
      this.recordsLoading = true
      try {
        const res = await getStudyRecords({ user_id: this.profile?.value?.id || 1 })
        this.records = res?.data?.records || []
      } catch (error) {
        console.error('获取学习记录失败:', error)
      } finally {
        this.recordsLoading = false
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
.container {
  max-width: 1440px;
  margin: 0 auto;
  padding: 20px;
}
</style>
