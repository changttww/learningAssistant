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
      <button
        @click="openAchievements"
        class="w-full flex items-center p-3 text-gray-700 hover:bg-gray-50 rounded-lg transition"
      >
        <iconify-icon
          icon="mdi:trophy"
          class="text-xl text-yellow-500 mr-3"
        />
        <span>我的成就</span>
      </button>
      <button
        @click="openSettings"
        class="w-full flex items-center p-3 text-gray-700 hover:bg-gray-50 rounded-lg transition"
      >
        <iconify-icon icon="mdi:cog" class="text-xl text-gray-500 mr-3" />
        <span>系统设置</span>
      </button>
    </div>
  </aside>

  <!-- 我的成就弹窗 -->
  <div
    v-if="showAchievements"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="showAchievements = false"
  >
    <div
      class="bg-white rounded-2xl p-6 w-[620px] max-h-[80vh] overflow-y-auto"
      @click.stop
    >
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-800 flex items-center">
          <iconify-icon
            icon="mdi:trophy"
            class="text-yellow-500 mr-2"
          ></iconify-icon>
          我的成就
        </h2>
        <button
          @click="showAchievements = false"
          class="text-gray-400 hover:text-gray-600"
        >
          <iconify-icon icon="mdi:close" class="text-2xl"></iconify-icon>
        </button>
      </div>

      <div
        v-if="achievementsLoading"
        class="flex items-center justify-center py-10 text-gray-500"
      >
        正在加载成就...
      </div>
      <div
        v-else-if="achievementsError"
        class="bg-red-50 border border-red-200 rounded-xl p-4 mb-4 text-red-600"
      >
        {{ achievementsError }}
      </div>
      <template v-else>
        <div class="grid grid-cols-3 gap-4 mb-6">
          <div
            class="bg-gradient-to-br from-yellow-50 to-orange-50 p-4 rounded-xl text-center"
          >
            <iconify-icon
              icon="mdi:medal"
              class="text-3xl text-yellow-500 mb-2"
            ></iconify-icon>
            <div class="text-2xl font-bold text-gray-800">{{ achievementStats.unlocked }}</div>
            <div class="text-sm text-gray-600">已获得成就</div>
          </div>
          <div
            class="bg-gradient-to-br from-blue-50 to-indigo-50 p-4 rounded-xl text-center"
          >
            <iconify-icon
              icon="mdi:map-marker-distance"
              class="text-3xl text-blue-500 mb-2"
            ></iconify-icon>
            <div class="text-2xl font-bold text-gray-800">{{ achievementStats.upcoming }}</div>
            <div class="text-sm text-gray-600">待完成成就</div>
          </div>
          <div
            class="bg-gradient-to-br from-green-50 to-emerald-50 p-4 rounded-xl text-center"
          >
            <iconify-icon
              icon="mdi:target"
              class="text-3xl text-green-500 mb-2"
            ></iconify-icon>
            <div class="text-2xl font-bold text-gray-800">{{ achievementStats.completionRate }}%</div>
            <div class="text-sm text-gray-600">完成度</div>
          </div>
        </div>

        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-bold text-gray-700 mb-3">已获得</h3>
            <span class="text-xs text-gray-500">最新优先</span>
          </div>
          <div
            v-if="achievementsData.unlocked.length === 0"
            class="p-4 bg-gray-50 border border-gray-200 rounded-xl text-gray-500"
          >
            还没有解锁任何成就，继续加油~
          </div>
          <div
            v-for="item in achievementsData.unlocked"
            :key="`unlocked-${item.id}`"
            class="flex flex-col p-4 bg-yellow-50 border border-yellow-200 rounded-xl"
          >
            <div class="flex items-center">
              <div
                class="w-12 h-12 rounded-full flex items-center justify-center mr-4"
                :class="item.completed ? 'bg-yellow-500' : 'bg-gray-300'"
              >
                <iconify-icon
                  :icon="item.icon || 'mdi:trophy'"
                  class="text-white text-xl"
                ></iconify-icon>
              </div>
              <div class="flex-1">
                <h4 class="font-bold text-gray-800">{{ item.name }}</h4>
                <p class="text-sm text-gray-600">{{ item.description }}</p>
                <p class="text-xs text-gray-500 mt-1">达成于 {{ item.awarded_at || "刚刚" }}</p>
              </div>
              <span class="text-green-600 text-sm font-semibold">已完成</span>
            </div>

            <div
              v-if="item.history && item.history.length > 0"
              class="mt-3"
            >
              <button
                class="text-xs text-blue-600 hover:underline flex items-center gap-1"
                @click="toggleHistory(item.code)"
              >
                <iconify-icon
                  :icon="isHistoryExpanded(item.code) ? 'mdi:chevron-up' : 'mdi:chevron-down'"
                  class="text-base"
                ></iconify-icon>
                {{ isHistoryExpanded(item.code) ? "收起早期成就" : `展开早期成就 (${item.history.length})` }}
              </button>
              <div
                v-if="isHistoryExpanded(item.code)"
                class="mt-2 space-y-2"
              >
                <div
                  v-for="hist in item.history"
                  :key="`hist-${hist.id}`"
                  class="flex items-center p-3 bg-white rounded-lg border border-yellow-100"
                >
                  <iconify-icon
                    :icon="hist.icon || 'mdi:trophy'"
                    class="text-yellow-500 text-lg mr-3"
                  ></iconify-icon>
                  <div class="flex-1">
                    <div class="font-medium text-gray-800">{{ hist.name }}</div>
                    <div class="text-xs text-gray-500">
                      {{ hist.description }} · {{ hist.awarded_at || "已获得" }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <h3 class="text-lg font-bold text-gray-700 mb-3 mt-6">待完成</h3>
          <div
            v-if="achievementsData.upcoming.length === 0"
            class="p-4 bg-gray-50 border border-gray-200 rounded-xl text-gray-500"
          >
            没有待完成的成就，太强了！
          </div>
          <div
            v-for="item in achievementsData.upcoming"
            :key="`upcoming-${item.id}`"
            class="flex items-center p-4 bg-gray-50 border border-gray-200 rounded-xl"
          >
            <div
              class="w-12 h-12 rounded-full flex items-center justify-center mr-4 bg-gray-200"
            >
              <iconify-icon
                :icon="item.icon || 'mdi:trophy-outline'"
                class="text-gray-700 text-xl"
              ></iconify-icon>
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-gray-800">{{ item.name }}</h4>
              <p class="text-sm text-gray-600">{{ item.description }}</p>
              <div class="w-full bg-gray-200 rounded-full h-2 mt-2">
                <div
                  class="bg-blue-500 h-2 rounded-full"
                  :style="{ width: `${Math.min(100, Math.round(progressPercent(item)))}%` }"
                ></div>
              </div>
              <div class="text-xs text-gray-500 mt-1">
                进度: {{ Math.floor(item.current_value) }} / {{ Math.floor(item.target_value) }}
              </div>
            </div>
            <span class="text-gray-500 text-sm font-semibold">未完成</span>
          </div>
        </div>
      </template>
    </div>
  </div>

  <!-- 系统设置弹窗 -->
  <div
    v-if="showSettings"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="showSettings = false"
  >
    <div
      class="bg-white rounded-2xl p-6 w-[500px] max-h-[80vh] overflow-y-auto"
      @click.stop
    >
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-800 flex items-center">
          <iconify-icon
            icon="mdi:cog"
            class="text-gray-500 mr-2"
          ></iconify-icon>
          系统设置
        </h2>
        <button
          @click="showSettings = false"
          class="text-gray-400 hover:text-gray-600"
        >
          <iconify-icon icon="mdi:close" class="text-2xl"></iconify-icon>
        </button>
      </div>

      <div v-if="settingsLoading" class="py-12 text-center text-gray-500">
        正在加载设置...
      </div>
      <div
        v-else-if="settingsLoadError"
        class="py-10 text-center flex flex-col items-center gap-4 text-sm text-gray-500"
      >
        <p>{{ settingsLoadError }}</p>
        <button
          class="px-4 py-2 bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 font-medium"
          @click="fetchUserSettings"
        >
          重试加载
        </button>
      </div>

      <!-- 设置选项 -->
      <div v-else class="space-y-6">
        <!-- 通知设置 -->
        <div>
          <h3 class="text-lg font-bold text-gray-700 mb-3 flex items-center">
            <iconify-icon
              icon="mdi:bell"
              class="text-blue-500 mr-2"
            ></iconify-icon>
            通知设置
          </h3>
          <div class="space-y-3">
            <div
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div class="font-medium text-gray-800">学习提醒</div>
                <div class="text-sm text-gray-600">
                  接收学习进度和每日提醒邮件
                </div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  class="sr-only peer"
                  v-model="settingsForm.notifications.email"
                />
                <div
                  class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
                ></div>
              </label>
            </div>

            <div
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div class="font-medium text-gray-800">任务截止提醒</div>
                <div class="text-sm text-gray-600">在应用内推送待办提醒</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  class="sr-only peer"
                  v-model="settingsForm.notifications.inApp"
                />
                <div
                  class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
                ></div>
              </label>
            </div>

            <div
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div class="font-medium text-gray-800">社交互动</div>
                <div class="text-sm text-gray-600">好友动态短信提醒</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  class="sr-only peer"
                  v-model="settingsForm.notifications.sms"
                />
                <div
                  class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
                ></div>
              </label>
            </div>

            <div
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div class="font-medium text-gray-800">每周学习周报</div>
                <div class="text-sm text-gray-600">
                  每周发送总结与建议至邮箱
                </div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  class="sr-only peer"
                  v-model="settingsForm.notifications.summary"
                />
                <div
                  class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
                ></div>
              </label>
            </div>
          </div>
        </div>

        <!-- 学习偏好 -->
        <div>
          <h3 class="text-lg font-bold text-gray-700 mb-3 flex items-center">
            <iconify-icon
              icon="mdi:account-cog"
              class="text-green-500 mr-2"
            ></iconify-icon>
            学习偏好
          </h3>
          <div class="space-y-3">
            <div class="p-3 bg-gray-50 rounded-lg">
              <label class="block text-sm font-medium text-gray-700 mb-2"
                >每日学习目标</label
              >
              <select
                class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                v-model.number="settingsForm.studyHabits.dailyGoalMinutes"
              >
                <option
                  v-for="option in studyGoalOptions"
                  :key="option.value"
                  :value="option.value"
                >
                  {{ option.label }}
                </option>
              </select>
            </div>

            <div class="p-3 bg-gray-50 rounded-lg">
              <label class="block text-sm font-medium text-gray-700 mb-2"
                >首选学习时段</label
              >
              <div class="space-y-2">
                <label
                  v-for="option in preferredPeriodOptions"
                  :key="option.value"
                  class="flex items-center"
                >
                  <input
                    type="radio"
                    class="text-blue-600"
                    :value="option.value"
                    v-model="settingsForm.studyHabits.preferredPeriod"
                  />
                  <span class="ml-2 text-sm text-gray-700">
                    {{ option.label }}
                  </span>
                </label>
              </div>
            </div>

            <div
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div class="font-medium text-gray-800">专注模式</div>
                <div class="text-sm text-gray-600">
                  开启后将屏蔽多余提醒，集中学习
                </div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  class="sr-only peer"
                  v-model="settingsForm.studyHabits.focusMode"
                />
                <div
                  class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
                ></div>
              </label>
            </div>
          </div>
        </div>

        <!-- 隐私设置 -->
        <div>
          <h3 class="text-lg font-bold text-gray-700 mb-3 flex items-center">
            <iconify-icon
              icon="mdi:shield-account"
              class="text-purple-500 mr-2"
            ></iconify-icon>
            隐私设置
          </h3>
          <div class="space-y-3">
            <div
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div class="font-medium text-gray-800">学习进度公开</div>
                <div class="text-sm text-gray-600">允许好友查看学习进度</div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  class="sr-only peer"
                  v-model="settingsForm.privacy.showStudyData"
                />
                <div
                  class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
                ></div>
              </label>
            </div>

            <div
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div class="font-medium text-gray-800">在线状态显示</div>
                <div class="text-sm text-gray-600">
                  显示在线状态给其他用户
                </div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  class="sr-only peer"
                  v-model="settingsForm.privacy.showProfile"
                />
                <div
                  class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
                ></div>
              </label>
            </div>

            <div
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div class="font-medium text-gray-800">邮箱展示</div>
                <div class="text-sm text-gray-600">
                  在个人资料中展示邮箱地址
                </div>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  class="sr-only peer"
                  v-model="settingsForm.privacy.showEmail"
                />
                <div
                  class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
                ></div>
              </label>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex gap-3 pt-4">
          <button
            @click="showSettings = false"
            class="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg font-medium hover:bg-gray-300"
          >
            取消
          </button>
          <button
            @click="saveSettings"
            class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed"
            :disabled="settingsSaving"
          >
            {{ settingsSaving ? "保存中..." : "保存设置" }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
  import { ElMessage } from "element-plus";
  import {
    useCurrentUser,
    DEFAULT_USER_ID,
  } from "@/composables/useCurrentUser";
  import {
    checkInUser,
    getUserAchievements,
    getUserPointsLedger,
    getUserSettings,
    updateUserSettings,
  } from "@/api/modules/user";

  const CHECK_IN_STORAGE_KEY = "task_sidebar_daily_check_in";

  const STUDY_GOAL_OPTIONS = [
    { label: "30分钟", value: 30 },
    { label: "1小时", value: 60 },
    { label: "2小时", value: 120 },
    { label: "3小时", value: 180 },
    { label: "4小时", value: 240 },
  ];

  const PREFERRED_PERIOD_OPTIONS = [
    { label: "清晨（6-9点）", value: "morning" },
    { label: "白天（10-17点）", value: "day" },
    { label: "傍晚（18-21点）", value: "evening" },
    { label: "夜间（22点后）", value: "night" },
  ];

  function createDefaultSettingsForm() {
    return {
      notifications: {
        email: true,
        sms: false,
        inApp: true,
        summary: true,
      },
      privacy: {
        showEmail: false,
        showProfile: true,
        showStudyData: true,
      },
      studyHabits: {
        dailyGoalMinutes: 60,
        preferredPeriod: "evening",
        focusMode: false,
      },
    };
  }

  function mapSettingsResponse(data) {
    const base = createDefaultSettingsForm();
    if (!data) return base;
    return {
      notifications: {
        email: data.notifications?.email ?? base.notifications.email,
        sms: data.notifications?.sms ?? base.notifications.sms,
        inApp: data.notifications?.in_app ?? base.notifications.inApp,
        summary: data.notifications?.weekly_summary ?? base.notifications.summary,
      },
      privacy: {
        showEmail: data.privacy?.show_email ?? base.privacy.showEmail,
        showProfile: data.privacy?.show_profile ?? base.privacy.showProfile,
        showStudyData: data.privacy?.show_study_data ?? base.privacy.showStudyData,
      },
      studyHabits: {
        dailyGoalMinutes:
          data.study_habits?.daily_goal_minutes ?? base.studyHabits.dailyGoalMinutes,
        preferredPeriod:
          data.study_habits?.preferred_period ?? base.studyHabits.preferredPeriod,
        focusMode: data.study_habits?.focus_mode ?? base.studyHabits.focusMode,
      },
    };
  }

  function buildSettingsPayload(form) {
    const safeForm = form || createDefaultSettingsForm();
    return {
      notifications: {
        email: !!safeForm.notifications.email,
        sms: !!safeForm.notifications.sms,
        in_app: !!safeForm.notifications.inApp,
        weekly_summary: !!safeForm.notifications.summary,
      },
      privacy: {
        show_email: !!safeForm.privacy.showEmail,
        show_profile: !!safeForm.privacy.showProfile,
        show_study_data: !!safeForm.privacy.showStudyData,
      },
      study_habits: {
        daily_goal_minutes: Number(safeForm.studyHabits.dailyGoalMinutes) || 0,
        preferred_period: safeForm.studyHabits.preferredPeriod,
        focus_mode: !!safeForm.studyHabits.focusMode,
      },
    };
  }

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

      const showAchievements = ref(false);
      const achievementsLoading = ref(false);
      const achievementsError = ref("");
      const achievementsData = ref({
        unlocked: [],
        upcoming: [],
        progress: {},
      });
      const expandedHistory = ref(new Set());

      const showSettings = ref(false);
      const settingsForm = ref(createDefaultSettingsForm());
      const settingsLoading = ref(false);
      const settingsLoadError = ref("");
      const settingsSaving = ref(false);

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

      const loadAchievements = async () => {
        if (achievementsLoading.value) return;
        achievementsLoading.value = true;
        achievementsError.value = "";
        try {
          const res = await getUserAchievements(profile.value?.id ?? DEFAULT_USER_ID);
          achievementsData.value = res?.data || { unlocked: [], upcoming: [], progress: {} };
        } catch (error) {
          console.error("加载成就失败:", error);
          achievementsError.value = error?.message || "加载成就失败";
        } finally {
          achievementsLoading.value = false;
        }
      };

      watch(showAchievements, (val) => {
        if (val) {
          loadAchievements();
        }
      });

      const openAchievements = () => {
        showAchievements.value = true;
      };

      const openSettings = () => {
        showSettings.value = true;
        if (!settingsLoading.value) {
          fetchUserSettings();
        }
      };

      const fetchUserSettings = async () => {
        if (settingsLoading.value) return;
        settingsLoading.value = true;
        settingsLoadError.value = "";
        try {
          const response = await getUserSettings(profile.value?.id ?? DEFAULT_USER_ID);
          settingsForm.value = mapSettingsResponse(response?.data);
        } catch (error) {
          console.error("获取用户设置失败:", error);
          settingsLoadError.value = error?.message || "获取用户设置失败";
        } finally {
          settingsLoading.value = false;
        }
      };

      const saveSettings = async () => {
        if (settingsSaving.value) return;
        settingsSaving.value = true;
        try {
          const payload = buildSettingsPayload(settingsForm.value);
          await updateUserSettings(profile.value?.id ?? DEFAULT_USER_ID, payload);
          ElMessage.success("设置已保存");
          showSettings.value = false;
        } catch (error) {
          console.error("保存设置失败:", error);
          ElMessage.error(error?.message || "保存设置失败");
        } finally {
          settingsSaving.value = false;
        }
      };

      const achievementStats = computed(() => {
        const unlocked = achievementsData.value?.unlocked?.length || 0;
        const upcoming = achievementsData.value?.upcoming?.length || 0;
        const total = unlocked + upcoming || 1;
        const completionRate = Math.round((unlocked / total) * 100);
        return { unlocked, upcoming, completionRate };
      });

      const progressPercent = (item) => {
        if (!item || !item.target_value) return 0;
        const percent =
          (Number(item.current_value || 0) / Number(item.target_value)) * 100;
        return Number.isFinite(percent) ? percent : 0;
      };

      const toggleHistory = (code) => {
        const next = new Set(expandedHistory.value);
        if (next.has(code)) {
          next.delete(code);
        } else {
          next.add(code);
        }
        expandedHistory.value = next;
      };

      const isHistoryExpanded = (code) => expandedHistory.value.has(code);

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
        showAchievements,
        achievementsLoading,
        achievementsError,
        achievementsData,
        achievementStats,
        progressPercent,
        expandedHistory,
        toggleHistory,
        isHistoryExpanded,
        openAchievements,
        showSettings,
        openSettings,
        settingsForm,
        settingsLoading,
        settingsLoadError,
        settingsSaving,
        fetchUserSettings,
        saveSettings,
        studyGoalOptions: STUDY_GOAL_OPTIONS,
        preferredPeriodOptions: PREFERRED_PERIOD_OPTIONS,
      };
    },
  };
</script>

<style scoped>
  .sidebar-shadow {
    box-shadow: 0px 0px 20px rgba(0, 0, 0, 0.05);
  }
</style>
