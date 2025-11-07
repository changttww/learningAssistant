<template>
  <div id="app" class="h-screen w-full flex flex-col overflow-hidden">
    <!-- 顶部导航 -->
    <div class="header w-full">
      <div class="flex items-center space-x-10">
        <div class="text-xl font-bold text-[#2D5BFF] flex items-center">
          <iconify-icon icon="mdi:brain" width="26" height="26"></iconify-icon>
          <span class="ml-2">智学空间</span>
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
          <iconify-icon
            icon="mdi:bell-outline"
            width="24"
            height="24"
            class="text-gray-600"
          ></iconify-icon>
          <div class="notification-dot"></div>
        </div>
        <div class="flex items-center">
          <router-link to="/profile" class="flex items-center">
            <div
              class="w-8 h-8 rounded-full bg-gray-200 overflow-hidden flex items-center justify-center text-sm font-semibold text-[#2D5BFF]"
            >
              <img
                v-if="avatarUrl"
                :src="avatarUrl"
                :alt="displayName"
                class="w-full h-full object-cover"
              />
              <span v-else>
                {{ displayName.slice(0, 1) }}
              </span>
            </div>
            <span class="ml-2 font-medium">{{ displayName }}</span>
          </router-link>
        </div>
      </div>
    </div>

    <!-- 路由视图 -->
    <main class="flex-1 min-h-0 w-full overflow-auto px-[5%]">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted } from "vue";
import { useCurrentUser } from "@/composables/useCurrentUser";
import { getToken } from "@/utils/auth";

const {
  profile,
  profileLoaded,
  loadCurrentUser,
} = useCurrentUser();

const displayName = computed(() => profile.value?.display_name || "访客");
const avatarUrl = computed(() => profile.value?.avatar_url || "");

onMounted(() => {
  if (getToken() && !profileLoaded.value) {
    loadCurrentUser(profile.value?.id);
  }
});

defineOptions({
  name: "App",
});
</script>

<style>
  .router-link-active {
    color: #2d5bff !important;
  }
</style>
