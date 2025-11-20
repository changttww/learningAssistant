<template>
  <div class="bg-gray-50 min-h-screen py-8">
    <div class="w-full h-full">
      <div class="flex gap-6 h-full items-stretch">
        <TaskSidebar />

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
  import {
    useCurrentUser,
    DEFAULT_USER_ID,
  } from "@/composables/useCurrentUser";
  import TaskSidebar from "@/components/TaskManager/TaskSidebar.vue";
  import LearningBoardHeader from "@/components/TaskManager/LearningBoardHeader.vue";
  import TaskProgressOverview from "@/components/TaskManager/TaskProgressOverview.vue";
  import AnalysisEntryGrid from "@/components/TaskManager/AnalysisEntryGrid.vue";
  import TaskTabsSection from "@/components/TaskManager/TaskTabsSection.vue";
  import InteractionPanel from "@/components/TaskManager/InteractionPanel.vue";

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
