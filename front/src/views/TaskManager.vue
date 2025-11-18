<template>
  <div class="bg-gray-50 min-h-screen py-8">
    <div class="w-full h-full">
      <div class="flex gap-6 h-full items-stretch">
        <TaskSidebar
          @show-achievements="showAchievements = true"
          @show-settings="openSettingsModal"
        />

        <main class="flex-1 flex flex-col h-full">
          <LearningBoardHeader
            :active-filter="activeTimeFilter"
            @update:activeFilter="activeTimeFilter = $event"
          />

          <TaskProgressOverview
            :current-time-data="currentTimeData"
            :active-time-filter="activeTimeFilter"
            @show-details="showTaskDetails"
          />

          <AnalysisEntryGrid
            @show-efficiency="showEfficiencyAnalysis"
            @show-summary="showSmartSummary"
            @show-check-in="showCheckInAnalysis"
          />

          <TaskTabsSection
            :active-tab="activeTab"
            :tasks="tasks"
            @update:activeTab="activeTab = $event"
          />
        </main>

        <InteractionPanel
          :is-chat-expanded="isChatExpanded"
          :current-motivational-quote="currentMotivationalQuote"
          @toggle-chat="toggleChatList"
          @open-chat="goToChatHistory"
        />
      </div>
    </div>
  </div>

<!-- 我的成就弹窗 -->
    <div
      v-if="showAchievements"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="showAchievements = false"
    >
      <div
        class="bg-white rounded-2xl p-6 w-[600px] max-h-[80vh] overflow-y-auto"
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

        <!-- 成就统计 -->
        <div class="grid grid-cols-3 gap-4 mb-6">
          <div
            class="bg-gradient-to-br from-yellow-50 to-orange-50 p-4 rounded-xl text-center"
          >
            <iconify-icon
              icon="mdi:medal"
              class="text-3xl text-yellow-500 mb-2"
            ></iconify-icon>
            <div class="text-2xl font-bold text-gray-800">12</div>
            <div class="text-sm text-gray-600">已获得成就</div>
          </div>
          <div
            class="bg-gradient-to-br from-blue-50 to-indigo-50 p-4 rounded-xl text-center"
          >
            <iconify-icon
              icon="mdi:star"
              class="text-3xl text-blue-500 mb-2"
            ></iconify-icon>
            <div class="text-2xl font-bold text-gray-800">3860</div>
            <div class="text-sm text-gray-600">成就积分</div>
          </div>
          <div
            class="bg-gradient-to-br from-green-50 to-emerald-50 p-4 rounded-xl text-center"
          >
            <iconify-icon
              icon="mdi:target"
              class="text-3xl text-green-500 mb-2"
            ></iconify-icon>
            <div class="text-2xl font-bold text-gray-800">85%</div>
            <div class="text-sm text-gray-600">完成度</div>
          </div>
        </div>

        <!-- 成就列表 -->
        <div class="space-y-4">
          <h3 class="text-lg font-bold text-gray-700 mb-3">最近获得</h3>

          <!-- 成就项目 -->
          <div
            class="flex items-center p-4 bg-yellow-50 border border-yellow-200 rounded-xl"
          >
            <div
              class="w-12 h-12 bg-yellow-500 rounded-full flex items-center justify-center mr-4"
            >
              <iconify-icon
                icon="mdi:school"
                class="text-white text-xl"
              ></iconify-icon>
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-gray-800">学习达人</h4>
              <p class="text-sm text-gray-600">连续学习30天</p>
            </div>
            <div class="text-right">
              <div class="text-yellow-600 font-bold">+500积分</div>
              <div class="text-xs text-gray-500">2024-01-15</div>
            </div>
          </div>

          <div
            class="flex items-center p-4 bg-blue-50 border border-blue-200 rounded-xl"
          >
            <div
              class="w-12 h-12 bg-blue-500 rounded-full flex items-center justify-center mr-4"
            >
              <iconify-icon
                icon="mdi:code-tags"
                class="text-white text-xl"
              ></iconify-icon>
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-gray-800">编程高手</h4>
              <p class="text-sm text-gray-600">完成10个编程项目</p>
            </div>
            <div class="text-right">
              <div class="text-blue-600 font-bold">+300积分</div>
              <div class="text-xs text-gray-500">2024-01-10</div>
            </div>
          </div>

          <div
            class="flex items-center p-4 bg-green-50 border border-green-200 rounded-xl"
          >
            <div
              class="w-12 h-12 bg-green-500 rounded-full flex items-center justify-center mr-4"
            >
              <iconify-icon
                icon="mdi:account-group"
                class="text-white text-xl"
              ></iconify-icon>
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-gray-800">团队协作者</h4>
              <p class="text-sm text-gray-600">参与5个团队项目</p>
            </div>
            <div class="text-right">
              <div class="text-green-600 font-bold">+200积分</div>
              <div class="text-xs text-gray-500">2024-01-05</div>
            </div>
          </div>

          <!-- 未获得成就 -->
          <h3 class="text-lg font-bold text-gray-700 mb-3 mt-6">待解锁</h3>

          <div
            class="flex items-center p-4 bg-gray-50 border border-gray-200 rounded-xl opacity-60"
          >
            <div
              class="w-12 h-12 bg-gray-400 rounded-full flex items-center justify-center mr-4"
            >
              <iconify-icon
                icon="mdi:lightning-bolt"
                class="text-white text-xl"
              ></iconify-icon>
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-gray-800">速度之王</h4>
              <p class="text-sm text-gray-600">单日完成5个任务</p>
              <div class="w-full bg-gray-200 rounded-full h-2 mt-2">
                <div
                  class="bg-blue-500 h-2 rounded-full"
                  style="width: 60%"
                ></div>
              </div>
              <div class="text-xs text-gray-500 mt-1">进度: 3/5</div>
            </div>
            <div class="text-right">
              <div class="text-gray-500 font-bold">+800积分</div>
            </div>
          </div>
        </div>
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

  <!-- 学习效率分析弹窗 -->
  <div
    v-if="showEfficiencyModal"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="closeEfficiencyModal"
  >
    <div
      class="bg-white rounded-2xl p-6 w-[800px] max-h-[80vh] overflow-y-auto"
      @click.stop
    >
      <div class="flex justify-between items-center mb-6">
        <h3 class="text-2xl font-bold text-gray-800">学习效率分析报告</h3>
        <button
          @click="closeEfficiencyModal"
          class="text-gray-500 hover:text-gray-700"
        >
          <iconify-icon icon="mdi:close" class="text-2xl"></iconify-icon>
        </button>
      </div>

      <!-- 效率概览 -->
      <div class="grid grid-cols-3 gap-4 mb-6">
        <div class="bg-gradient-to-br from-purple-50 to-pink-50 p-4 rounded-xl">
          <div class="flex items-center mb-2">
            <iconify-icon
              icon="mdi:clock-outline"
              class="text-purple-600 text-xl mr-2"
            ></iconify-icon>
            <span class="text-gray-600 text-sm">本周学习时长</span>
          </div>
          <div class="text-2xl font-bold text-purple-600">
            {{ efficiencyData.weeklyStudyTime }}小时
          </div>
        </div>
        <div class="bg-gradient-to-br from-green-50 to-teal-50 p-4 rounded-xl">
          <div class="flex items-center mb-2">
            <iconify-icon
              icon="mdi:target"
              class="text-green-600 text-xl mr-2"
            ></iconify-icon>
            <span class="text-gray-600 text-sm">专注度评分</span>
          </div>
          <div class="text-2xl font-bold text-green-600">
            {{ efficiencyData.focusScore }}分
          </div>
        </div>
        <div class="bg-gradient-to-br from-blue-50 to-cyan-50 p-4 rounded-xl">
          <div class="flex items-center mb-2">
            <iconify-icon
              icon="mdi:check-circle-outline"
              class="text-blue-600 text-xl mr-2"
            ></iconify-icon>
            <span class="text-gray-600 text-sm">任务完成率</span>
          </div>
          <div class="text-2xl font-bold text-blue-600">
            {{ efficiencyData.taskCompletionRate }}%
          </div>
        </div>
      </div>

      <!-- 学习趋势图表 -->
      <div class="mb-6">
        <h4 class="text-lg font-bold text-gray-800 mb-4">学习趋势分析</h4>
        <div class="grid grid-cols-2 gap-4">
          <div class="bg-gray-50 p-4 rounded-xl">
            <h5 class="text-sm font-medium text-gray-600 mb-3">
              每日学习时长 (小时)
            </h5>
            <div class="h-32" ref="studyTrendChart"></div>
          </div>
          <div class="bg-gray-50 p-4 rounded-xl">
            <h5 class="text-sm font-medium text-gray-600 mb-3">
              每日专注度评分
            </h5>
            <div class="h-32" ref="focusTrendChart"></div>
          </div>
        </div>
      </div>

      <!-- 智能建议 -->
      <div class="mb-6">
        <h4 class="text-lg font-bold text-gray-800 mb-4">智能建议</h4>
        <div class="space-y-3">
          <div
            v-for="(suggestion, index) in efficiencyData.suggestions"
            :key="index"
            class="flex items-start p-3 rounded-lg"
            :class="{
              'bg-green-50 border-l-4 border-green-500':
                suggestion.type === 'positive',
              'bg-yellow-50 border-l-4 border-yellow-500':
                suggestion.type === 'warning',
              'bg-blue-50 border-l-4 border-blue-500':
                suggestion.type === 'tip',
            }"
          >
            <iconify-icon
              :icon="
                suggestion.type === 'positive'
                  ? 'mdi:thumb-up'
                  : suggestion.type === 'warning'
                  ? 'mdi:alert'
                  : 'mdi:lightbulb'
              "
              :class="{
                'text-green-600': suggestion.type === 'positive',
                'text-yellow-600': suggestion.type === 'warning',
                'text-blue-600': suggestion.type === 'tip',
              }"
              class="text-xl mr-3 mt-0.5"
            >
            </iconify-icon>
            <span class="text-gray-700">{{ suggestion.message }}</span>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex gap-3">
        <button
          @click="generateReport"
          class="flex-1 bg-blue-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-blue-700 flex items-center justify-center"
        >
          <iconify-icon icon="mdi:download" class="mr-2"></iconify-icon>
          生成详细报告
        </button>
        <button
          @click="closeEfficiencyModal"
          class="flex-1 bg-gray-200 text-gray-700 py-3 px-4 rounded-lg font-medium hover:bg-gray-300"
        >
          关闭
        </button>
      </div>
    </div>
  </div>

  <!-- 智能总结与复习弹窗 -->
  <div
    v-if="showSummaryModal"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="closeSummaryModal"
  >
    <div
      class="bg-white rounded-2xl p-6 w-[700px] max-h-[80vh] overflow-y-auto"
      @click.stop
    >
      <div class="flex justify-between items-center mb-6">
        <h3 class="text-2xl font-bold text-gray-800">智能总结与复习</h3>
        <button
          @click="closeSummaryModal"
          class="text-gray-500 hover:text-gray-700"
        >
          <iconify-icon icon="mdi:close" class="text-2xl"></iconify-icon>
        </button>
      </div>

      <!-- 知识掌握概览 -->
      <div
        class="bg-gradient-to-br from-blue-50 to-indigo-50 p-4 rounded-xl mb-6"
      >
        <h4 class="text-lg font-bold text-gray-800 mb-3">知识掌握情况</h4>
        <div class="grid grid-cols-3 gap-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-green-600">
              {{ summaryData.knowledgeMap.mastered }}%
            </div>
            <div class="text-sm text-gray-600">已掌握</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-blue-600">
              {{ summaryData.knowledgeMap.learning }}%
            </div>
            <div class="text-sm text-gray-600">学习中</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-orange-600">
              {{ summaryData.knowledgeMap.toLearn }}%
            </div>
            <div class="text-sm text-gray-600">待学习</div>
          </div>
        </div>
      </div>

      <!-- 待复习内容 -->
      <div class="mb-6">
        <h4 class="text-lg font-bold text-gray-800 mb-4">待复习内容</h4>
        <div class="space-y-3">
          <div
            v-for="(item, index) in summaryData.reviewItems"
            :key="index"
            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <div class="flex items-center">
              <div
                class="w-3 h-3 rounded-full mr-3"
                :class="{
                  'bg-red-500': item.priority === 'high',
                  'bg-yellow-500': item.priority === 'medium',
                  'bg-green-500': item.priority === 'low',
                }"
              ></div>
              <div>
                <div class="font-medium text-gray-800">{{ item.subject }}</div>
                <div class="text-sm text-gray-500">
                  进度: {{ item.progress }}% | 复习时间: {{ item.dueDate }}
                </div>
              </div>
            </div>
            <button
              @click="startReview(item)"
              class="bg-blue-600 text-white px-3 py-1 rounded-lg text-sm hover:bg-blue-700"
            >
              开始复习
            </button>
          </div>
        </div>
      </div>

      <!-- 复习提醒 -->
      <div class="mb-6">
        <h4 class="text-lg font-bold text-gray-800 mb-4">复习提醒</h4>
        <div class="space-y-3">
          <div
            v-for="(reminder, index) in summaryData.reminders"
            :key="index"
            class="flex items-center justify-between p-3 bg-yellow-50 rounded-lg border-l-4 border-yellow-500"
          >
            <div class="flex items-center">
              <iconify-icon
                icon="mdi:bell-outline"
                class="text-yellow-600 text-xl mr-3"
              ></iconify-icon>
              <div>
                <div class="font-medium text-gray-800">
                  {{ reminder.content }}
                </div>
                <div class="text-sm text-gray-500">{{ reminder.time }}</div>
              </div>
            </div>
            <button
              @click="setReminder(reminder)"
              class="bg-yellow-600 text-white px-3 py-1 rounded-lg text-sm hover:bg-yellow-700"
            >
              设置提醒
            </button>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex gap-3">
        <button
          @click="closeSummaryModal"
          class="flex-1 bg-gray-200 text-gray-700 py-3 px-4 rounded-lg font-medium hover:bg-gray-300"
        >
          关闭
        </button>
      </div>
    </div>
  </div>

  <!-- 学习打卡分析弹窗 -->
  <div
    v-if="showCheckInModal"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="closeCheckInModal"
  >
    <div
      class="bg-white rounded-2xl p-6 w-[600px] max-h-[80vh] overflow-y-auto"
      @click.stop
    >
      <div class="flex justify-between items-center mb-6">
        <h3 class="text-2xl font-bold text-gray-800">学习打卡分析</h3>
        <button
          @click="closeCheckInModal"
          class="text-gray-500 hover:text-gray-700"
        >
          <iconify-icon icon="mdi:close" class="text-2xl"></iconify-icon>
        </button>
      </div>

      <!-- 打卡统计 -->
      <div class="grid grid-cols-2 gap-4 mb-6">
        <div
          class="bg-gradient-to-br from-green-50 to-teal-50 p-4 rounded-xl text-center"
        >
          <iconify-icon
            icon="mdi:calendar-check"
            class="text-3xl text-green-600 mb-2"
          ></iconify-icon>
          <div class="text-2xl font-bold text-green-600">
            {{ checkInData.consecutiveDays }}
          </div>
          <div class="text-sm text-gray-600">连续打卡天数</div>
        </div>
        <div
          class="bg-gradient-to-br from-blue-50 to-cyan-50 p-4 rounded-xl text-center"
        >
          <iconify-icon
            icon="mdi:keyboard"
            class="text-3xl text-blue-600 mb-2"
          ></iconify-icon>
          <div class="text-2xl font-bold text-blue-600">
            {{ checkInData.avgTypingSpeed }}
          </div>
          <div class="text-sm text-gray-600">平均打字速度 (WPM)</div>
        </div>
      </div>

      <!-- 学习习惯分析 -->
      <div class="mb-6">
        <h4 class="text-lg font-bold text-gray-800 mb-4">学习习惯分析</h4>
        <div class="bg-gray-50 p-4 rounded-xl space-y-3">
          <div class="flex justify-between">
            <span class="text-gray-600">最佳学习时段</span>
            <span class="font-medium text-blue-600">{{
              checkInData.studyHabits.bestTime
            }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600">平均学习时长</span>
            <span class="font-medium text-green-600">{{
              checkInData.studyHabits.avgSession
            }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600">周学习目标</span>
            <span class="font-medium text-purple-600">{{
              checkInData.studyHabits.weeklyGoal
            }}</span>
          </div>
        </div>
      </div>

      <!-- 激励建议 -->
      <div class="mb-6">
        <h4 class="text-lg font-bold text-gray-800 mb-4">激励建议</h4>
        <div class="space-y-3">
          <div
            v-for="(suggestion, index) in checkInData.suggestions"
            :key="index"
            class="flex items-start p-3 bg-green-50 rounded-lg border-l-4 border-green-500"
          >
            <iconify-icon
              icon="mdi:star"
              class="text-green-600 text-xl mr-3 mt-0.5"
            ></iconify-icon>
            <span class="text-gray-700">{{ suggestion }}</span>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex gap-3">
        <button
          @click="closeCheckInModal"
          class="flex-1 bg-gray-200 text-gray-700 py-3 px-4 rounded-lg font-medium hover:bg-gray-300"
        >
          关闭
        </button>
      </div>
    </div>
  </div>
</template>

<script>
  import { computed } from "vue";
  import * as echarts from "echarts";
  import { ElMessage } from "element-plus";
  import {
    useCurrentUser,
    DEFAULT_USER_ID,
  } from "@/composables/useCurrentUser";
  import { getUserSettings, updateUserSettings } from "@/api/modules/user";
  import TaskSidebar from "@/components/TaskManager/TaskSidebar.vue";
  import LearningBoardHeader from "@/components/TaskManager/LearningBoardHeader.vue";
  import TaskProgressOverview from "@/components/TaskManager/TaskProgressOverview.vue";
  import AnalysisEntryGrid from "@/components/TaskManager/AnalysisEntryGrid.vue";
  import TaskTabsSection from "@/components/TaskManager/TaskTabsSection.vue";
  import InteractionPanel from "@/components/TaskManager/InteractionPanel.vue";

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
    if (!data) {
      return base;
    }

    return {
      notifications: {
        email: data.notifications?.email ?? base.notifications.email,
        sms: data.notifications?.sms ?? base.notifications.sms,
        inApp: data.notifications?.in_app ?? base.notifications.inApp,
        summary:
          data.notifications?.weekly_summary ?? base.notifications.summary,
      },
      privacy: {
        showEmail: data.privacy?.show_email ?? base.privacy.showEmail,
        showProfile: data.privacy?.show_profile ?? base.privacy.showProfile,
        showStudyData:
          data.privacy?.show_study_data ?? base.privacy.showStudyData,
      },
      studyHabits: {
        dailyGoalMinutes:
          data.study_habits?.daily_goal_minutes ??
          base.studyHabits.dailyGoalMinutes,
        preferredPeriod:
          data.study_habits?.preferred_period ??
          base.studyHabits.preferredPeriod,
        focusMode:
          data.study_habits?.focus_mode ?? base.studyHabits.focusMode,
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

  export default {
    name: "TaskManager",
    components: {
      TaskSidebar,
      LearningBoardHeader,
      TaskProgressOverview,
      AnalysisEntryGrid,
      TaskTabsSection,
      InteractionPanel,
    },
    setup() {
      const { profile } = useCurrentUser();

      return {
        currentUserId: computed(() => profile.value?.id ?? DEFAULT_USER_ID),
      };
    },
    data() {
      return {
        showAchievements: false,
        showSettings: false,
        settingsForm: createDefaultSettingsForm(),
        settingsLoading: false,
        settingsLoadError: "",
        settingsSaving: false,
        studyGoalOptions: STUDY_GOAL_OPTIONS,
        preferredPeriodOptions: PREFERRED_PERIOD_OPTIONS,
        activeTab: "inProgress",
        // 聊天相关状态
        isChatExpanded: true, // 聊天列表展开状态
        currentMotivationalQuote:
          "每一次努力都是成长的阶梯，坚持下去，你会看到不一样的自己！", // 当前励志语录
        motivationalQuotes: [
          "每一次努力都是成长的阶梯，坚持下去，你会看到不一样的自己！",
          "学习不是为了证明什么，而是为了成为更好的自己。",
          "今天的努力，是为了明天的从容不迫。",
          "知识是唯一不会贬值的投资，学习是最好的成长方式。",
          "不怕慢，只怕停。每天进步一点点，就是成功的开始。",
          "困难是成长的垫脚石，挑战是能力的试金石。",
          "相信自己，你比想象中更强大，比昨天更优秀。",
          "学习的路上没有捷径，但每一步都算数。",
        ],
        // 学习效率分析相关状态
        showEfficiencyModal: false,
        showSummaryModal: false,
        showCheckInModal: false,
        // 学习效率分析数据
        efficiencyData: {
          weeklyStudyTime: 28.5,
          focusScore: 85,
          taskCompletionRate: 92,
          studyTrend: [6.2, 4.8, 5.1, 3.9, 4.5, 2.8, 1.2], // 每日学习时长
          focusTrend: [88, 82, 90, 78, 85, 92, 80], // 每日专注度
          suggestions: [
            { type: "positive", message: "本周学习时长超过目标，继续保持！" },
            { type: "warning", message: "周末学习时间较少，建议合理安排" },
            { type: "tip", message: "下午2-4点是您的高效学习时段" },
          ],
        },
        // 智能总结数据
        summaryData: {
          reviewItems: [
            {
              subject: "JavaScript ES6",
              priority: "high",
              dueDate: "今天",
              progress: 60,
            },
            {
              subject: "Vue组件通信",
              priority: "medium",
              dueDate: "明天",
              progress: 75,
            },
            {
              subject: "CSS Grid布局",
              priority: "low",
              dueDate: "后天",
              progress: 40,
            },
          ],
          reminders: [
            { content: "复习Promise和async/await语法", time: "今天 14:00" },
            { content: "完成Vue项目实战练习", time: "明天 10:00" },
            { content: "整理CSS学习笔记", time: "后天 16:00" },
          ],
          knowledgeMap: {
            mastered: 78,
            learning: 15,
            toLearn: 7,
          },
        },
        // 打卡分析数据
        checkInData: {
          consecutiveDays: 28,
          avgTypingSpeed: 65,
          studyHabits: {
            bestTime: "14:00-16:00",
            avgSession: "2.5小时",
            weeklyGoal: "30小时",
          },
          motivationLevel: "high",
          suggestions: [
            "您的学习习惯很好，建议继续保持",
            "可以尝试在最佳时段安排重要任务",
            "打字速度不错，可以提高编程效率",
          ],
        },
        activeTimeFilter: "week", // 修正数据属性名称
        // 不同时间段的数据
        timeFilterData: {
          week: {
            chartData: [45, 52, 68, 73, 64, 42, 30],
            chartLabels: [
              "周一",
              "周二",
              "周三",
              "周四",
              "周五",
              "周六",
              "周日",
            ],
            completionRate: 72,
            completedTasks: 86,
            totalTasks: 120,
          },
          month: {
            chartData: [
              65, 72, 58, 83, 76, 69, 74, 81, 67, 79, 85, 78, 72, 88, 91, 69,
              75, 82, 77, 84, 73, 86, 79, 81, 75, 83, 78, 80, 76, 89,
            ],
            chartLabels: [
              "1日",
              "2日",
              "3日",
              "4日",
              "5日",
              "6日",
              "7日",
              "8日",
              "9日",
              "10日",
              "11日",
              "12日",
              "13日",
              "14日",
              "15日",
              "16日",
              "17日",
              "18日",
              "19日",
              "20日",
              "21日",
              "22日",
              "23日",
              "24日",
              "25日",
              "26日",
              "27日",
              "28日",
              "29日",
              "30日",
            ],
            completionRate: 78,
            completedTasks: 234,
            totalTasks: 300,
          },
          quarter: {
            chartData: [68, 74, 82],
            chartLabels: ["1月", "2月", "3月"],
            completionRate: 81,
            completedTasks: 486,
            totalTasks: 600,
          },
        },
        tasks: {
          inProgress: [
            {
              id: 1,
              title: "Vue.js 组件开发",
              description: "学习Vue组件的高级用法",
              progress: 75,
              priority: "high",
              dueDate: "2024-01-20",
              tags: ["前端", "Vue"],
            },
            {
              id: 2,
              title: "JavaScript ES6+",
              description: "掌握现代JavaScript语法",
              progress: 60,
              priority: "medium",
              dueDate: "2024-01-25",
              tags: ["JavaScript", "基础"],
            },
            {
              id: 3,
              title: "CSS Grid 布局",
              description: "学习CSS Grid布局系统",
              progress: 40,
              priority: "low",
              dueDate: "2024-01-30",
              tags: ["CSS", "布局"],
            },
          ],
          pending: [
            {
              id: 4,
              title: "React Hooks",
              description: "学习React Hooks的使用",
              progress: 0,
              priority: "medium",
              dueDate: "2024-02-01",
              tags: ["React", "前端"],
            },
            {
              id: 5,
              title: "Node.js 后端开发",
              description: "构建RESTful API",
              progress: 0,
              priority: "high",
              dueDate: "2024-02-05",
              tags: ["Node.js", "后端"],
            },
          ],
          completed: [
            {
              id: 6,
              title: "HTML5 基础",
              description: "掌握HTML5新特性",
              progress: 100,
              priority: "low",
              dueDate: "2024-01-10",
              tags: ["HTML", "基础"],
            },
            {
              id: 7,
              title: "Git 版本控制",
              description: "学习Git基本操作",
              progress: 100,
              priority: "medium",
              dueDate: "2024-01-15",
              tags: ["Git", "工具"],
            },
          ],
        },
      };
    },
    computed: {
      // 当前时间筛选器对应的数据
      currentTimeData() {
        return this.timeFilterData[this.activeTimeFilter];
      },
    },
    methods: {
      openSettingsModal() {
        this.showSettings = true;
        if (!this.settingsLoading) {
          this.fetchUserSettings();
        }
      },
      async fetchUserSettings() {
        if (this.settingsLoading) {
          return;
        }
        this.settingsLoading = true;
        this.settingsLoadError = "";
        try {
          const response = await getUserSettings(this.currentUserId);
          this.settingsForm = mapSettingsResponse(response?.data);
        } catch (error) {
          console.error("获取用户设置失败:", error);
          this.settingsLoadError = error?.message || "获取用户设置失败";
        } finally {
          this.settingsLoading = false;
        }
      },
      showTaskDetails() {
        // 点击环形图显示任务详情的联动功能
        console.log("显示任务详情");
      },
      // 显示学习效率分析
      showEfficiencyAnalysis() {
        this.showEfficiencyModal = true;
        this.$nextTick(() => {
          this.initEfficiencyCharts();
        });
      },
      closeEfficiencyModal() {
        this.showEfficiencyModal = false;
      },
      // 智能总结与复习方法
      showSmartSummary() {
        this.showSummaryModal = true;
      },
      closeSummaryModal() {
        this.showSummaryModal = false;
      },
      // 学习打卡分析方法
      showCheckInAnalysis() {
        this.showCheckInModal = true;
      },
      closeCheckInModal() {
        this.showCheckInModal = false;
      },
      // 生成学习报告
      generateReport() {
        console.log("生成学习效率报告");
        // 这里可以添加生成PDF报告的逻辑
      },
      // 开始复习
      startReview(item) {
        console.log("开始复习:", item.subject);
        // 这里可以添加跳转到具体复习内容的逻辑
      },
      // 设置提醒
      setReminder(reminder) {
        console.log("设置提醒:", reminder.content);
        // 这里可以添加设置系统提醒的逻辑
      },
      // 初始化效率分析图表
      initEfficiencyCharts() {
        this.$nextTick(() => {
          // 学习时长趋势图
          if (this.$refs.studyTrendChart) {
            const studyChart = echarts.init(this.$refs.studyTrendChart);
            studyChart.setOption({
              tooltip: {
                trigger: "axis",
                formatter: "{b}<br/>学习时长: {c}小时",
              },
              grid: {
                left: "10%",
                right: "10%",
                bottom: "15%",
                top: "10%",
              },
              xAxis: {
                type: "category",
                data: ["周一", "周二", "周三", "周四", "周五", "周六", "周日"],
                axisLine: { lineStyle: { color: "#E5E7EB" } },
                axisTick: { show: false },
                axisLabel: { fontSize: 10 },
              },
              yAxis: {
                type: "value",
                axisLine: { show: false },
                axisTick: { show: false },
                splitLine: { lineStyle: { color: "#F0F2F5" } },
                axisLabel: { formatter: "{value}h", fontSize: 10 },
              },
              series: [
                {
                  data: this.efficiencyData.studyTrend,
                  type: "bar",
                  barWidth: 20,
                  itemStyle: {
                    color: {
                      type: "linear",
                      x: 0,
                      y: 0,
                      x2: 0,
                      y2: 1,
                      colorStops: [
                        { offset: 0, color: "#8B5CF6" },
                        { offset: 1, color: "#A78BFA" },
                      ],
                    },
                    borderRadius: [4, 4, 0, 0],
                  },
                },
              ],
            });
          }

          // 专注度趋势图
          if (this.$refs.focusTrendChart) {
            const focusChart = echarts.init(this.$refs.focusTrendChart);
            focusChart.setOption({
              tooltip: {
                trigger: "axis",
                formatter: "{b}<br/>专注度: {c}分",
              },
              grid: {
                left: "10%",
                right: "10%",
                bottom: "15%",
                top: "10%",
              },
              xAxis: {
                type: "category",
                data: ["周一", "周二", "周三", "周四", "周五", "周六", "周日"],
                axisLine: { lineStyle: { color: "#E5E7EB" } },
                axisTick: { show: false },
                axisLabel: { fontSize: 10 },
              },
              yAxis: {
                type: "value",
                min: 70,
                max: 100,
                axisLine: { show: false },
                axisTick: { show: false },
                splitLine: { lineStyle: { color: "#F0F2F5" } },
                axisLabel: { formatter: "{value}", fontSize: 10 },
              },
              series: [
                {
                  data: this.efficiencyData.focusTrend,
                  type: "line",
                  smooth: true,
                  symbol: "circle",
                  symbolSize: 6,
                  lineStyle: {
                    width: 3,
                    color: "#10B981",
                  },
                  itemStyle: {
                    color: "#10B981",
                    borderColor: "#fff",
                    borderWidth: 2,
                  },
                  areaStyle: {
                    color: {
                      type: "linear",
                      x: 0,
                      y: 0,
                      x2: 0,
                      y2: 1,
                      colorStops: [
                        { offset: 0, color: "rgba(16, 185, 129, 0.3)" },
                        { offset: 1, color: "rgba(16, 185, 129, 0.1)" },
                      ],
                    },
                  },
                },
              ],
            });
          }
        });
      },
      async saveSettings() {
        if (this.settingsSaving) {
          return;
        }
        this.settingsSaving = true;
        try {
          const payload = buildSettingsPayload(this.settingsForm);
          await updateUserSettings(this.currentUserId, payload);
          ElMessage.success("设置已保存");
          this.showSettings = false;
        } catch (error) {
          console.error("保存设置失败:", error);
          if (!error?.message) {
            ElMessage.error("保存设置失败");
          }
        } finally {
          this.settingsSaving = false;
        }
      },

      // 聊天相关方法
      toggleChatList() {
        this.isChatExpanded = !this.isChatExpanded;
      },

      // 跳转到历史聊天界面
      goToChatHistory(friendName) {
        // 这里可以使用Vue Router进行页面跳转
        // 假设有一个聊天页面路由，传递好友名称作为参数
        console.log(`跳转到与 ${friendName} 的聊天界面`);

        // 示例：使用Vue Router跳转到聊天页面
        // this.$router.push({
        //   name: 'ChatHistory',
        //   params: { friendName: friendName },
        //   query: { autoFocus: true } // 自动聚焦到输入框
        // });

        // 临时实现：显示提示信息
        alert(`即将跳转到与 ${friendName} 的聊天界面，并自动聚焦到消息输入框`);
      },

    },
  };
</script>

<style scoped></style>
