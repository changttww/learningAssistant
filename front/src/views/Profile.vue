<template>
  <div class="container">
    <div class="card">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl font-bold">个人资料</h1>
        <button class="bg-[#2D5BFF] text-white font-medium py-2 px-4 rounded-lg text-sm hover:bg-opacity-90 transition-colors flex items-center gap-2">
          <iconify-icon icon="mdi:pencil" width="16" height="16"></iconify-icon>
          编辑资料
        </button>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 左侧个人信息 -->
        <div class="lg:col-span-2">
          <div class="flex items-start gap-6 mb-6">
            <div class="w-24 h-24 rounded-full bg-gray-300"></div>
            <div class="flex-1">
              <h2 class="text-xl font-bold mb-2">张同学</h2>
              <p class="text-gray-600 mb-3">专注于教育科技领域，热衷于学习新技术，致力于打造高效的学习工具。</p>
              <div class="flex items-center gap-4">
                <div class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">技术项目经理</div>
                <div class="bg-green-100 text-green-800 px-3 py-1 rounded-full text-sm font-medium flex items-center gap-1">
                  <span class="w-2 h-2 rounded-full bg-green-500"></span>
                  在线
                </div>
              </div>
            </div>
          </div>

          <!-- 基本信息 -->
          <div class="space-y-4">
            <h3 class="font-bold text-lg">基本信息</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">姓名</label>
                <input type="text" value="张同学" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500" readonly>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">邮箱</label>
                <input type="email" value="zhang@example.com" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500" readonly>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">学校</label>
                <input type="text" value="北京大学" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500" readonly>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">专业</label>
                <input type="text" value="计算机科学与技术" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500" readonly>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">地区</label>
                <input type="text" value="北京市" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500" readonly>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">加入时间</label>
                <input type="text" value="2023年9月" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500" readonly>
              </div>
            </div>
          </div>

          <!-- 技能标签 -->
          <div class="mt-6">
            <h3 class="font-bold text-lg mb-3">技能标签</h3>
            <div class="flex flex-wrap gap-2">
              <span class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm">Vue.js</span>
              <span class="bg-green-100 text-green-800 px-3 py-1 rounded-full text-sm">JavaScript</span>
              <span class="bg-purple-100 text-purple-800 px-3 py-1 rounded-full text-sm">Node.js</span>
              <span class="bg-orange-100 text-orange-800 px-3 py-1 rounded-full text-sm">项目管理</span>
              <span class="bg-pink-100 text-pink-800 px-3 py-1 rounded-full text-sm">UI设计</span>
              <span class="bg-yellow-100 text-yellow-800 px-3 py-1 rounded-full text-sm">数据分析</span>
            </div>
          </div>
        </div>

        <!-- 右侧统计信息 -->
        <div class="space-y-4">
          <div class="bg-blue-50 p-4 rounded-lg">
            <div class="text-center">
              <div class="text-2xl font-bold text-blue-600">学霸 Lv.4</div>
              <div class="text-sm text-gray-600 mt-1">当前等级</div>
              <div class="mt-3">
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 75%"></div>
                </div>
                <div class="text-xs text-gray-600 mt-1">距离下一级还需 250 积分</div>
              </div>
            </div>
          </div>

          <div class="space-y-3">
            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <div class="flex items-center gap-3">
                <iconify-icon icon="mdi:clock" width="20" height="20" class="text-blue-600"></iconify-icon>
                <span class="font-medium">总学习时长</span>
              </div>
              <span class="font-bold text-blue-600">87.5h</span>
            </div>

            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <div class="flex items-center gap-3">
                <iconify-icon icon="mdi:check-circle" width="20" height="20" class="text-green-600"></iconify-icon>
                <span class="font-medium">完成任务</span>
              </div>
              <span class="font-bold text-green-600">156</span>
            </div>

            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <div class="flex items-center gap-3">
                <iconify-icon icon="mdi:certificate" width="20" height="20" class="text-purple-600"></iconify-icon>
                <span class="font-medium">获得证书</span>
              </div>
              <span class="font-bold text-purple-600">24</span>
            </div>

            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <div class="flex items-center gap-3">
                <iconify-icon icon="mdi:account-group" width="20" height="20" class="text-orange-600"></iconify-icon>
                <span class="font-medium">学习小组</span>
              </div>
              <span class="font-bold text-orange-600">3</span>
            </div>
          </div>

          <!-- 最近成就 -->
          <div class="mt-6">
            <h3 class="font-bold text-lg mb-3">最近成就</h3>
            <div class="space-y-2">
              <div class="flex items-center gap-3 p-2 bg-yellow-50 rounded-lg">
                <iconify-icon icon="mdi:trophy" width="20" height="20" class="text-yellow-600"></iconify-icon>
                <div>
                  <div class="font-medium text-sm">连续学习7天</div>
                  <div class="text-xs text-gray-600">2024-01-15</div>
                </div>
              </div>
              
              <div class="flex items-center gap-3 p-2 bg-blue-50 rounded-lg">
                <iconify-icon icon="mdi:star" width="20" height="20" class="text-blue-600"></iconify-icon>
                <div>
                  <div class="font-medium text-sm">完成Vue项目</div>
                  <div class="text-xs text-gray-600">2024-01-12</div>
                </div>
              </div>
              
              <div class="flex items-center gap-3 p-2 bg-green-50 rounded-lg">
                <iconify-icon icon="mdi:medal" width="20" height="20" class="text-green-600"></iconify-icon>
                <div>
                  <div class="font-medium text-sm">团队协作达人</div>
                  <div class="text-xs text-gray-600">2024-01-10</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Profile'
}
</script>

<style scoped>
.container {
  max-width: 1440px;
  margin: 0 auto;
  padding: 20px;
}
</style>