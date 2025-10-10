<template>
  <div id="app">
    <!-- 顶部导航 -->
    <div class="header">
      <div class="flex items-center space-x-10">
        <div class="text-xl font-bold text-[#2D5BFF] flex items-center">
          <iconify-icon icon="mdi:brain" width="26" height="26"></iconify-icon>
          <span class="ml-2">学辅助手</span>
        </div>
        <div class="flex space-x-6 font-medium text-gray-700">
          <div class="py-1 hover:text-[#2D5BFF]">
            <router-link to="/">首页</router-link>
          </div>
          <div class="py-1 hover:text-[#2D5BFF]">
            <router-link to="/personal-tasks">个人任务</router-link>
          </div>
          <div class="py-1 hover:text-[#2D5BFF]">
            <router-link to="/team-tasks">团队任务</router-link>
          </div>
          <div class="py-1 hover:text-[#2D5BFF]">
            <router-link to="/study-room">在线自习室</router-link>
          </div>
          <div class="py-1 hover:text-[#2D5BFF]">
            <router-link to="/task-manager">任务进度管理</router-link>
          </div>
        </div>
      </div>
      <div class="flex items-center space-x-5">
        <div class="relative">
          <iconify-icon icon="mdi:bell-outline" width="24" height="24" class="text-gray-600"></iconify-icon>
          <div class="notification-dot"></div>
        </div>
        <div class="flex items-center">
          <router-link to="/profile" class="flex items-center">
            <div class="w-8 h-8 rounded-full bg-gray-300"></div>
            <span class="ml-2 font-medium">张同学</span>
          </router-link>
        </div>
      </div>
    </div>

    <!-- 路由视图 -->
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
.router-link-active {
  color: #2D5BFF !important;
}
</style>