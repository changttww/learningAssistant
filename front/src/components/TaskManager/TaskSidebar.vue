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

      <div class="grid grid-cols-2 gap-3 mb-6">
        <div class="bg-white rounded-lg p-3">
          <p class="text-gray-500 text-xs mb-1">累计签到</p>
          <p class="text-lg font-bold text-blue-600">
            {{ sidebarProfile.streakDays }}<span class="text-xs">天</span>
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
        <button
          class="w-full flex items-center justify-center gap-2 py-3 rounded-xl font-semibold transition disabled:cursor-not-allowed"
          :class="hasCheckedInToday ? 'bg-gray-200 text-gray-500' : 'bg-gradient-to-r from-blue-500 to-indigo-500 text-white hover:from-blue-600 hover:to-indigo-600'"
          :disabled="hasCheckedInToday || isCheckingIn"
          @click="handleCheckIn"
        >
          <iconify-icon icon="mdi:calendar-check" class="text-xl" />
          <span>
            {{ hasCheckedInToday ? "已签到" : isCheckingIn ? "签到中..." : "今日签到" }}
          </span>
        </button>
        <p class="text-xs text-gray-500 mt-2">
          {{
            hasCheckedInToday
              ? `已累计签到${sidebarProfile.streakDays}天，次日0点可再次签到`
              : "签到可增加累计签到天数"
          }}
        </p>
      </div>

      <div class="mb-6">
        <div class="flex justify-between items-center mb-3">
          <h3 class="font-bold text-gray-700 flex items-center">
            <iconify-icon icon="mdi:crystal-ball" class="mr-2 text-blue-600" />
            我的积分
          </h3>
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
    </div>
  </aside>
</template>

<script>
  import { computed, onBeforeUnmount, onMounted, ref } from "vue";
  import { ElMessage } from "element-plus";
  import {
    useCurrentUser,
    DEFAULT_USER_ID,
  } from "@/composables/useCurrentUser";
  import {
    checkInUser,
    getUserPointsLedger,
  } from "@/api/modules/user";

  const CHECK_IN_STORAGE_KEY = "task_sidebar_daily_check_in";

  function createDefaultPointsSummary() {
    return {
      currentPoints: 0,
      pointsToNextLevel: 300,
      levelLabel: "学霸 Lv.1",
      levelProgress: 0,
    };
  }

  function mapPointsSummary(summary = {}, items = [], stats = {}) {
    const latestLedger = Array.isArray(items) && items.length ? items[0] : null;
    const currentPoints =
      Number(summary.current_points ?? stats.current_points ?? latestLedger?.balance_after ?? 0) || 0;
    const nextLevelPoints =
      Number(summary.next_level_points ?? stats.next_level_points ?? 200) || 200;
    const rawDistance =
      summary.distance_to_next ??
      stats.distance_to_next ??
      nextLevelPoints - currentPoints;
    const pointsToNextLevel = Math.max(Number(rawDistance) || 0, 0);
    const rawProgress =
      summary.progress_percent ??
      stats.progress_percent ??
      (nextLevelPoints > 0 ? (currentPoints / nextLevelPoints) * 100 : 0);
    const levelProgress = Math.min(Math.max(Math.round(rawProgress), 0), 100);
    const currentLevel = summary.current_level ?? stats.current_level ?? 1;
    const levelLabel = summary.level_label || stats.level_label || `学霸 Lv.${currentLevel}`;

    return {
      currentPoints,
      pointsToNextLevel,
      levelLabel,
      levelProgress,
    };
  }

  function getTodayKey() {
    const now = new Date();
    const month = `${now.getMonth() + 1}`.padStart(2, "0");
    const day = `${now.getDate()}`.padStart(2, "0");
    return `${now.getFullYear()}-${month}-${day}`;
  }

  export default {
    name: "TaskSidebar",
    setup() {
      const { profile, studyStats, loadCurrentUser, loadStudyStats } =
        useCurrentUser();

      const hasCheckedInToday = ref(false);
      const isCheckingIn = ref(false);
      let midnightResetTimer = null;

      const pointsSummary = ref(createDefaultPointsSummary());
      const pointsLoading = ref(false);
      const pointsLoadError = ref("");

      const restoreCheckInState = () => {
        if (typeof window === "undefined") {
          hasCheckedInToday.value = false;
          return;
        }
        try {
          const storedDate = window.localStorage.getItem(CHECK_IN_STORAGE_KEY);
          hasCheckedInToday.value = storedDate === getTodayKey();
        } catch (error) {
          console.warn("读取签到状态失败:", error);
          hasCheckedInToday.value = false;
        }
      };

      const persistCheckInState = () => {
        if (typeof window === "undefined") {
          return;
        }
        try {
          if (hasCheckedInToday.value) {
            window.localStorage.setItem(CHECK_IN_STORAGE_KEY, getTodayKey());
          } else {
            window.localStorage.removeItem(CHECK_IN_STORAGE_KEY);
          }
        } catch (error) {
          console.warn("保存签到状态失败:", error);
        }
      };

      const scheduleMidnightReset = () => {
        if (typeof window === "undefined") {
          return;
        }
        if (midnightResetTimer) {
          window.clearTimeout(midnightResetTimer);
        }
        const now = new Date();
        const nextMidnight = new Date(
          now.getFullYear(),
          now.getMonth(),
          now.getDate() + 1,
          0,
          0,
          0,
          0
        );
        const timeout = nextMidnight.getTime() - now.getTime();
        midnightResetTimer = window.setTimeout(() => {
          hasCheckedInToday.value = false;
          persistCheckInState();
          scheduleMidnightReset();
        }, timeout);
      };

      const loadPointsData = async (userId) => {
        const targetUserId = userId ?? profile.value?.id ?? DEFAULT_USER_ID;
        if (pointsLoading.value) return;
        pointsLoading.value = true;
        pointsLoadError.value = "";
        try {
          const res = await getUserPointsLedger(targetUserId, { limit: 1 });
          const data = res?.data || {};
          pointsSummary.value = mapPointsSummary(
            data.summary,
            data.items,
            studyStats.value || {}
          );
        } catch (error) {
          console.error("获取积分数据失败:", error);
          pointsLoadError.value = error?.message || "获取积分数据失败";
          pointsSummary.value = mapPointsSummary({}, [], studyStats.value || {});
        } finally {
          pointsLoading.value = false;
        }
      };

      onMounted(async () => {
        restoreCheckInState();
        scheduleMidnightReset();
        try {
          const loadedProfile = await loadCurrentUser();
          const stats = await loadStudyStats(loadedProfile?.id ?? DEFAULT_USER_ID);
          pointsSummary.value = mapPointsSummary({}, [], stats || {});
          await loadPointsData(loadedProfile?.id ?? DEFAULT_USER_ID);
        } catch (error) {
          console.error("加载侧边栏用户信息失败:", error);
        }
      });

      onBeforeUnmount(() => {
        if (midnightResetTimer) {
          window.clearTimeout(midnightResetTimer);
          midnightResetTimer = null;
        }
      });

      const handleCheckIn = async () => {
        if (hasCheckedInToday.value || isCheckingIn.value) {
          return;
        }
        isCheckingIn.value = true;

        const targetUserId = profile.value?.id ?? DEFAULT_USER_ID;
        try {
          await checkInUser(targetUserId);
          hasCheckedInToday.value = true;
          persistCheckInState();
          ElMessage.success("签到成功，连续天数+1");
          const stats = await loadStudyStats(targetUserId, { force: true });
          pointsSummary.value = mapPointsSummary({}, [], stats || {});
          loadPointsData(targetUserId);
        } catch (error) {
          console.error("签到失败:", error);
          ElMessage.error(error?.message || "签到失败，请稍后重试");
        } finally {
          isCheckingIn.value = false;
        }
      };

      const sidebarProfile = computed(() => {
        const user = profile.value || {};
        const stats = studyStats.value || {};
        const points = pointsSummary.value || createDefaultPointsSummary();

        return {
          displayName: user.display_name || "学习者",
          userAvatar: user.avatar_url || "",
          userRoleLabel: user.role || "学习者",
          streakDays: stats.streak_days ?? 28,
          totalTasks: stats.total_tasks ?? 1,
          currentPoints: points.currentPoints,
          levelLabel: points.levelLabel,
          levelProgress: points.levelProgress,
          pointsToNextLevel: points.pointsToNextLevel,
        };
      });

      const displayInitial = computed(() => {
        const name = sidebarProfile.value.displayName || "学习者";
        return name.trim().slice(0, 1);
      });

      return {
        sidebarProfile,
        displayInitial,
        hasCheckedInToday,
        isCheckingIn,
        handleCheckIn,
      };
    },
  };
</script>

<style scoped>
  .sidebar-shadow {
    box-shadow: 0px 0px 20px rgba(0, 0, 0, 0.05);
  }
</style>
