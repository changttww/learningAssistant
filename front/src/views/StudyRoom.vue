<template>
  <div class="w-full h-full overflow-auto px-4">
    <div class="card">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl font-bold">在线自习室</h1>
        <button class="bg-[#2D5BFF] text-white font-medium py-2 px-4 rounded-lg text-sm hover:bg-opacity-90 transition-colors flex items-center gap-2">
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
        <div class="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-medium">前端学习小组</h3>
            <span class="bg-green-100 text-green-800 px-2 py-1 rounded text-xs">进行中</span>
          </div>
          <p class="text-sm text-gray-600 mb-3">专注于前端技术学习和交流</p>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:account-group" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">12/20</span>
            </div>
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:clock" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">2.5h</span>
            </div>
          </div>
        </div>

        <div class="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-medium">算法刷题房</h3>
            <span class="bg-blue-100 text-blue-800 px-2 py-1 rounded text-xs">热门</span>
          </div>
          <p class="text-sm text-gray-600 mb-3">一起刷算法题，提升编程能力</p>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:account-group" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">8/15</span>
            </div>
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:clock" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">1.8h</span>
            </div>
          </div>
        </div>

        <div class="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-medium">设计师工作室</h3>
            <span class="bg-purple-100 text-purple-800 px-2 py-1 rounded text-xs">创意</span>
          </div>
          <p class="text-sm text-gray-600 mb-3">UI/UX设计学习和作品分享</p>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:account-group" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">6/10</span>
            </div>
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:clock" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">3.2h</span>
            </div>
          </div>
        </div>

        <div class="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-medium">考研冲刺班</h3>
            <span class="bg-red-100 text-red-800 px-2 py-1 rounded text-xs">紧急</span>
          </div>
          <p class="text-sm text-gray-600 mb-3">考研最后冲刺，互相监督学习</p>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:account-group" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">25/30</span>
            </div>
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:clock" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">4.1h</span>
            </div>
          </div>
        </div>

        <div class="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-medium">英语角</h3>
            <span class="bg-yellow-100 text-yellow-800 px-2 py-1 rounded text-xs">语言</span>
          </div>
          <p class="text-sm text-gray-600 mb-3">英语口语练习和交流</p>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:account-group" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">9/12</span>
            </div>
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:clock" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">1.5h</span>
            </div>
          </div>
        </div>

        <div class="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-medium">安静学习室</h3>
            <span class="bg-gray-100 text-gray-800 px-2 py-1 rounded text-xs">静音</span>
          </div>
          <p class="text-sm text-gray-600 mb-3">专注学习，禁止聊天</p>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:account-group" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">18/25</span>
            </div>
            <div class="flex items-center gap-2">
              <iconify-icon icon="mdi:clock" width="16" height="16" class="text-gray-500"></iconify-icon>
              <span class="text-sm">2.8h</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 我的学习记录 -->
      <div class="mt-8">
        <h2 class="font-bold text-xl mb-4">我的学习记录</h2>
        <div class="space-y-3">
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center gap-3">
              <div class="w-3 h-3 rounded-full bg-green-500"></div>
              <div>
                <span class="font-medium">前端学习小组</span>
                <span class="text-sm text-gray-600 ml-2">学习了 2.5 小时</span>
              </div>
            </div>
            <span class="text-sm text-gray-500">今天 14:30 - 17:00</span>
          </div>
          
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
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
  </div>
</template>

<script>
export default {
  name: 'StudyRoom'
}
</script>

<style scoped>
.container {
  max-width: 1440px;
  margin: 0 auto;
  padding: 20px;
}
</style>
