<template>
  <div class="w-full h-full overflow-auto px-4">
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
          <div class="text-2xl font-bold text-blue-600">15</div>
          <div class="text-gray-600 mt-1">活跃房间</div>
        </div>
        <div class="stat-card bg-green-50">
          <div class="text-2xl font-bold text-green-600">128</div>
          <div class="text-gray-600 mt-1">在线用户</div>
        </div>
        <div class="stat-card bg-orange-50">
          <div class="text-2xl font-bold text-orange-600">3.5h</div>
          <div class="text-gray-600 mt-1">今日学习</div>
        </div>
        <div class="stat-card bg-purple-50">
          <div class="text-2xl font-bold text-purple-600">7</div>
          <div class="text-gray-600 mt-1">连续天数</div>
        </div>
      </div>

      <!-- 房间列表 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="room in rooms"
          :key="room.id"
          class="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer"
          @click="enterRoom(room.id)"
        >
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-medium">{{ room.name }}</h3>
            <span :class="['px-2 py-1 rounded text-xs', room.statusClass]">
              {{ room.status }}
            </span>
          </div>
          <p class="text-sm text-gray-600 mb-3">{{ room.description }}</p>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <iconify-icon
                icon="mdi:account-group"
                width="16"
                height="16"
                class="text-gray-500"
              ></iconify-icon>
              <span class="text-sm">{{ room.currentUsers }}/{{ room.maxUsers }}</span>
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

      <!-- 我的学习记录 -->
      <div class="mt-8">
        <h2 class="font-bold text-xl mb-4">我的学习记录</h2>
        <div class="space-y-3">
          <div
            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
          >
            <div class="flex items-center gap-3">
              <div class="w-3 h-3 rounded-full bg-green-500"></div>
              <div>
                <span class="font-medium">前端学习小组</span>
                <span class="text-sm text-gray-600 ml-2">学习了 2.5 小时</span>
              </div>
            </div>
            <span class="text-sm text-gray-500">今天 14:30 - 17:00</span>
          </div>

          <div
            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
          >
            <div class="flex items-center gap-3">
              <div class="w-3 h-3 rounded-full bg-blue-500"></div>
              <div>
                <span class="font-medium">算法刷题房</span>
                <span class="text-sm text-gray-600 ml-2">学习了 1.8 小时</span>
              </div>
            </div>
            <span class="text-sm text-gray-500">昨天 19:00 - 20:48</span>
          </div>
        </div>
      </div>
    </div>
    <!-- 创建房间弹窗 -->
    <CreateRoom 
      :visible="showCreateRoom"
      @close="hideCreateRoomModal"
      @created="onRoomCreated"
    />
  </div>
</template>

<script>
import CreateRoom from '@/components/StudyRoom/CreateRoom.vue'

export default {
  name: 'StudyRoom',
  components: {
    CreateRoom
  },
  data() {
    return {
      showCreateRoom: false,
      rooms: [
        {
          id: 1,
          name: '前端学习小组',
          status: '进行中',
          statusClass: 'bg-green-100 text-green-800',
          description: '专注于前端技术学习和交流',
          currentUsers: 12,
          maxUsers: 20,
          studyTime: '2.5h'
        },
        {
          id: 2,
          name: '算法刷题房',
          status: '热门',
          statusClass: 'bg-blue-100 text-blue-800',
          description: '一起刷算法题，提升编程能力',
          currentUsers: 8,
          maxUsers: 15,
          studyTime: '1.8h'
        },
        {
          id: 3,
          name: '设计师工作室',
          status: '创意',
          statusClass: 'bg-purple-100 text-purple-800',
          description: 'UI/UX设计学习和作品分享',
          currentUsers: 6,
          maxUsers: 10,
          studyTime: '3.2h'
        },
        {
          id: 4,
          name: '考研冲刺班',
          status: '紧急',
          statusClass: 'bg-red-100 text-red-800',
          description: '考研最后冲刺，互相监督学习',
          currentUsers: 25,
          maxUsers: 30,
          studyTime: '4.1h'
        },
        {
          id: 5,
          name: '英语角',
          status: '语言',
          statusClass: 'bg-yellow-100 text-yellow-800',
          description: '英语口语练习和交流',
          currentUsers: 9,
          maxUsers: 12,
          studyTime: '1.5h'
        },
        {
          id: 6,
          name: '安静学习室',
          status: '静音',
          statusClass: 'bg-gray-100 text-gray-800',
          description: '专注学习，禁止聊天',
          currentUsers: 18,
          maxUsers: 25,
          studyTime: '2.8h'
        }
      ]
    }
  },
  methods: {
    showCreateRoomModal() {
      this.showCreateRoom = true
    },
    hideCreateRoomModal() {
      this.showCreateRoom = false
    },
    enterRoom(roomId) {
      // 跳转到视频会议室
      this.$router.push({ name: 'VideoRoom', params: { roomId } })
    },
    onRoomCreated(roomData) {
      // 处理房间创建成功后的逻辑
      console.log('房间创建成功:', roomData)
      this.hideCreateRoomModal()
      // 可以添加新房间到列表或跳转到新房间
      this.enterRoom(roomData.id)
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
