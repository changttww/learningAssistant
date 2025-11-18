<template>
  <aside
    class="w-80 flex-shrink-0 sidebar-shadow rounded-2xl bg-white p-6 flex flex-col"
  >
    <div
      class="card mb-6 p-6 text-center bg-gradient-to-br from-blue-50 to-indigo-50"
    >
      <div class="flex justify-center mb-4">
        <div
          class="w-24 h-24 rounded-full bg-blue-100 flex items-center justify-center overflow-hidden text-2xl font-semibold text-[#2D5BFF]"
        >
          <img
            v-if="sidebarProfile.userAvatar"
            :src="sidebarProfile.userAvatar"
            :alt="sidebarProfile.displayName"
            class="w-full h-full object-cover"
          />
          <span v-else>{{ displayInitial }}</span>
        </div>
      </div>
      <h2 class="text-xl font-bold mb-2">{{ sidebarProfile.displayName }}</h2>
      <div class="flex justify-center items-center mb-4">
        <iconify-icon icon="mdi:medal" class="text-orange-500 text-xl" />
        <span class="ml-1 text-orange-500 font-medium">
          {{ sidebarProfile.userRoleLabel }}
        </span>
      </div>

      <div class="grid grid-cols-3 gap-3 mb-6">
        <div class="bg-white rounded-lg p-3">
          <p class="text-gray-500 text-xs mb-1">连续签到</p>
          <p class="text-lg font-bold text-blue-600">
            {{ sidebarProfile.streakDays }}<span class="text-xs">天</span>
          </p>
        </div>
        <div class="bg-white rounded-lg p-3">
          <p class="text-gray-500 text-xs mb-1">正在学习</p>
          <p class="text-lg font-bold text-blue-600">
            {{ sidebarProfile.coursesInProgress }}<span class="text-xs">门课</span>
          </p>
        </div>
        <div class="bg-white rounded-lg p-3">
          <p class="text-gray-500 text-xs mb-1">总任务</p>
          <p class="text-lg font-bold text-blue-600">
            {{ sidebarProfile.totalTasks }}<span class="text-xs">项</span>
          </p>
        </div>
      </div>

      <div class="mb-6">
        <div class="flex justify-between items-center mb-3">
          <h3 class="font-bold text-gray-700 flex items-center">
            <iconify-icon icon="mdi:crystal-ball" class="mr-2 text-blue-600" />
            我的积分
          </h3>
          <span
            class="badge inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
          >
            {{ sidebarProfile.pointsRankLabel }}
          </span>
        </div>

        <div class="flex justify-center mb-4">
          <div class="relative">
            <div
              class="w-20 h-20 rounded-full relative flex items-center justify-center"
            >
              <div
                class="absolute inset-0 bg-blue-600 bg-opacity-10 rounded-full"
              ></div>
              <span class="text-2xl font-bold text-blue-600 relative z-10">
                {{ sidebarProfile.currentPoints }}
              </span>
            </div>
            <div class="absolute bottom-0 right-0 transform translate-y-1">
              <iconify-icon
                icon="mdi:star-four-points"
                class="text-yellow-400 text-lg"
              />
            </div>
          </div>
        </div>
      </div>

      <div class="mb-4">
        <div class="flex justify-between items-center mb-2">
          <span class="text-sm font-medium text-gray-700">
            {{ sidebarProfile.levelLabel }}
          </span>
          <span class="text-xs text-gray-500">距离下一级</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-3 mb-2">
          <div
            class="bg-gradient-to-r from-blue-500 to-purple-500 h-3 rounded-full"
            :style="{ width: `${sidebarProfile.levelProgress}%` }"
          ></div>
        </div>
        <div class="text-xs text-gray-600">
          还需 {{ sidebarProfile.pointsToNextLevel }} 积分升级到下一级
        </div>
      </div>

    </div>

    <div class="space-y-3">
      <router-link
        to="/personal-tasks"
        class="flex items-center p-3 text-gray-700 hover:bg-gray-50 rounded-lg transition"
      >
        <iconify-icon
          icon="mdi:calendar-check"
          class="text-xl text-green-500 mr-3"
        />
        <span>今日任务</span>
      </router-link>
      <button
        @click="$emit('show-achievements')"
        class="w-full flex items-center p-3 text-gray-700 hover:bg-gray-50 rounded-lg transition"
      >
        <iconify-icon
          icon="mdi:trophy"
          class="text-xl text-yellow-500 mr-3"
        />
        <span>我的成就</span>
      </button>
      <button
        @click="$emit('show-settings')"
        class="w-full flex items-center p-3 text-gray-700 hover:bg-gray-50 rounded-lg transition"
      >
        <iconify-icon icon="mdi:cog" class="text-xl text-gray-500 mr-3" />
        <span>系统设置</span>
      </button>
    </div>
  </aside>
</template>

<script>
  import { computed, onMounted } from "vue";
  import {
    useCurrentUser,
    DEFAULT_USER_ID,
  } from "@/composables/useCurrentUser";

  export default {
    name: "TaskSidebar",
    emits: ["show-achievements", "show-settings"],
    setup() {
      const { profile, studyStats, loadCurrentUser, loadStudyStats } =
        useCurrentUser();

      onMounted(async () => {
        try {
          const loadedProfile = await loadCurrentUser();
          await loadStudyStats(loadedProfile?.id ?? DEFAULT_USER_ID);
        } catch (error) {
          console.error("加载侧边栏用户信息失败:", error);
        }
      });

      const sidebarProfile = computed(() => {
        const user = profile.value || {};
        const stats = studyStats.value || {};
        const progress = Math.min(
          Math.max(stats.progress_percent ?? 75, 0),
          100
        );

        return {
          displayName: user.display_name || "学习者",
          userAvatar: user.avatar_url || "",
          userRoleLabel: user.role || "学习者",
          streakDays: stats.streak_days ?? 28,
          coursesInProgress: stats.courses_in_progress ?? 3,
          totalTasks: stats.total_tasks ?? 1,
          pointsRankLabel: stats.rank_label || "TOP 15%",
          currentPoints: stats.current_points ?? 0,
          levelLabel: stats.level_label || "成长中学员",
          levelProgress: progress,
          pointsToNextLevel: stats.distance_to_next ?? 250,
        };
      });

      const displayInitial = computed(() => {
        const name = sidebarProfile.value.displayName || "学习者";
        return name.trim().slice(0, 1);
      });

      return {
        sidebarProfile,
        displayInitial,
      };
    },
  };
</script>

<style scoped>
  .sidebar-shadow {
    box-shadow: 0px 0px 20px rgba(0, 0, 0, 0.05);
  }

  .badge {
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0% {
      box-shadow: 0 0 0 0 rgba(45, 91, 255, 0.4);
    }
    70% {
      box-shadow: 0 0 0 10px rgba(45, 91, 255, 0);
    }
    100% {
      box-shadow: 0 0 0 0 rgba(45, 91, 255, 0);
    }
  }
</style>
