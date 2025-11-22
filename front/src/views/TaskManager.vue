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
import * as echarts from 'echarts'

export default {
  name: 'TaskManager',
  mounted() {
    this.initCharts()
  },
  methods: {
    initCharts() {
      // 环形进度图
      const ringChart = echarts.init(this.$refs.ringProgress)
      ringChart.setOption({
        tooltip: { show: false },
        series: [{
          type: 'gauge',
          startAngle: 180,
          endAngle: 0,
          radius: '100%',
          min: 0,
          max: 100,
          splitNumber: 10,
          pointer: { show: false },
          axisLine: {
            lineStyle: {
              width: 15,
              color: [[0.72, '#2D5BFF'], [1, '#F5F7FA']]
            }
          },
          axisLabel: { show: false },
          axisTick: { show: false },
          splitLine: { show: false },
          detail: { show: false }
        }]
      })
      
      // 柱状图
      const barChart = echarts.init(this.$refs.taskProgressChart)
      barChart.setOption({
        tooltip: {
          trigger: 'axis',
          formatter: '{b}<br/>{c}% 完成'
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          top: '5%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
          axisLine: { lineStyle: { color: '#E5E7EB' } },
          axisTick: { show: false }
        },
        yAxis: {
          type: 'value',
          max: 100,
          axisLine: { show: false },
          axisTick: { show: false },
          splitLine: {
            lineStyle: { color: '#F0F2F5' }
          },
          axisLabel: { formatter: '{value}%' }
        },
        series: [{
          data: [45, 52, 68, 73, 64, 42, 30],
          type: 'bar',
          barWidth: 24,
          itemStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [{
                offset: 0, color: '#2D5BFF'
              }, {
                offset: 1, color: '#5D8AFE'
              }]
            },
            borderRadius: [8, 8, 0, 0]
          },
          emphasis: {
            itemStyle: {
              shadowColor: 'rgba(45,91,255,0.5)',
              shadowBlur: 8
            }
          }
        }]
      })
      
      // 窗口大小变化时重绘图表
      window.addEventListener('resize', () => {
        ringChart.resize()
        barChart.resize()
      })
    }
  }
}
</script>

<style scoped>
.sidebar-shadow {
  box-shadow: 0px 0px 20px rgba(0, 0, 0, 0.05);
}

.badge {
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0% { box-shadow: 0 0 0 0 rgba(45, 91, 255, 0.4); }
  70% { box-shadow: 0 0 0 10px rgba(45, 91, 255, 0); }
  100% { box-shadow: 0 0 0 0 rgba(45, 91, 255, 0); }
}

.bubble-left:before {
  content: '';
  width: 0;
  height: 0;
  position: absolute;
  left: -10px;
  top: 10px;
  border-top: 10px solid transparent;
  border-bottom: 10px solid transparent;
  border-right: 10px solid #E5E7EB;
}

.bubble-right:before {
  content: '';
  width: 0;
  height: 0;
  position: absolute;
  right: -10px;
  top: 10px;
  border-top: 10px solid transparent;
  border-bottom: 10px solid transparent;
  border-left: 10px solid #2D5BFF;
}

.friend-card:hover {
  background-color: rgba(45, 91, 255, 0.05);
}
</style>
