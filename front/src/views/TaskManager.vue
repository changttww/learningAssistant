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
            :daily-overview="dailyOverview"
            :current-time-data="currentTimeData"
            :active-time-filter="activeTimeFilter"
            @show-details="showTaskDetails"
          />

          <AnalysisEntryGrid
            :summary="efficiencySummary"
            :loading="efficiencyLoading"
            @show-efficiency="showEfficiencyAnalysis"
            @show-summary="showSmartSummary"
          />

          <TaskTabsSection
          />
        </main>
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
      <div class="relative">
        <div
          v-if="efficiencyLoading"
          class="absolute inset-0 flex flex-col items-center justify-center bg-white bg-opacity-80 rounded-2xl z-10"
        >
          <svg
            class="animate-spin h-10 w-10 text-blue-600"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
            ></path>
          </svg>
          <p class="text-gray-600 mt-3">正在获取学习效率分析...</p>
        </div>

        <div :class="{ 'opacity-40 pointer-events-none': efficiencyLoading }">
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
              <p class="text-xs text-gray-500 mt-1">
                连续打卡 {{ efficiencyData.detailStats.streakDays }} 天 · 共计
                {{ efficiencyData.detailStats.totalStudyMinutes }} 分钟
              </p>
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
              <p class="text-xs text-gray-500 mt-1">
                基于最近 14 天学习记录的综合专注度评估
              </p>
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
              <p class="text-xs text-gray-500 mt-1">
                已完成 {{ efficiencyData.detailStats.completedTasks }} 项 ·
                进行中 {{ efficiencyData.detailStats.inProgressTasks }} 项
              </p>
            </div>
          </div>

          <!-- 关键洞察 -->
          <div class="mb-6" v-if="efficiencyData.insights.length">
            <h4 class="text-lg font-bold text-gray-800 mb-4">关键洞察</h4>
            <div class="space-y-3">
              <div
                v-for="(insight, index) in efficiencyData.insights"
                :key="index"
                class="flex items-start p-4 rounded-lg bg-purple-50 border border-purple-100"
              >
                <iconify-icon
                  icon="mdi:lightbulb-on-outline"
                  class="text-purple-600 text-xl mr-3 mt-0.5"
                ></iconify-icon>
                <div class="flex-1">
                  <p class="text-gray-800 font-semibold">
                    洞察 {{ index + 1 }}
                  </p>
                  <p class="text-gray-600 text-sm leading-relaxed">
                    {{ insight }}
                  </p>
                </div>
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
                <div class="flex-1">
                  <div class="flex items-center gap-2">
                    <span class="text-gray-800 font-semibold">
                      {{ suggestion.title }}
                    </span>
                    <span
                      v-if="suggestion.impact"
                      class="text-xs px-2 py-1 rounded-full bg-white/60 text-gray-600 border border-gray-200"
                    >
                      影响：{{ suggestion.impact }}
                    </span>
                  </div>
                  <p class="text-gray-700 text-sm leading-relaxed mt-1">
                    {{ suggestion.description || suggestion.message }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex flex-wrap gap-3">
            <button
              @click="regenerateEfficiency"
              class="flex-1 border border-blue-200 text-blue-700 py-3 px-4 rounded-lg font-medium hover:border-blue-300 hover:bg-blue-50 flex items-center justify-center disabled:opacity-60 disabled:cursor-not-allowed"
              :disabled="efficiencyLoading"
            >
              <iconify-icon icon="mdi:refresh" class="mr-2"></iconify-icon>
              重新生成报告
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
        <div v-if="summaryData.reviewItems.length" class="space-y-3">
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
        <p v-else class="text-sm text-gray-500">暂无AI生成的复习清单，稍后再试</p>
      </div>

      <!-- 复习提醒 -->
      <div class="mb-6">
        <h4 class="text-lg font-bold text-gray-800 mb-4">复习提醒</h4>
        <div v-if="summaryData.reminders.length" class="space-y-3">
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
        <p v-else class="text-sm text-gray-500">暂无AI生成的复习提醒</p>
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
  import * as echarts from "echarts";
  import { ElMessage } from "element-plus";
  import { getTaskBarStats } from "@/api/modules/task";
  import { fetchEfficiencyAnalysis } from "@/api/modules/analysis";
  import { getUserInfo } from "@/utils/auth";
  import TaskSidebar from "@/components/TaskManager/TaskSidebar.vue";
  import LearningBoardHeader from "@/components/TaskManager/LearningBoardHeader.vue";
  import TaskProgressOverview from "@/components/TaskManager/TaskProgressOverview.vue";
  import AnalysisEntryGrid, {
    parseEfficiencyPayload,
    transformEfficiencyData,
  } from "@/components/TaskManager/AnalysisEntryGrid.vue";
  import TaskTabsSection from "@/components/TaskManager/TaskTabsSection.vue";

  const createEmptyTimeData = () => ({
    chartData: [],
    chartLabels: [],
    completionRate: 0,
    completedTasks: 0,
    totalTasks: 0,
  });

  const EMPTY_TIME_DATA = createEmptyTimeData();

  const clampPercentage = (value) => {
    const numeric = Number(value);
    if (!Number.isFinite(numeric)) return 0;
    if (numeric < 0) return 0;
    if (numeric > 100) return 100;
    return Math.round(numeric);
  };

  export default {
    name: "TaskManager",
    components: {
      TaskSidebar,
      LearningBoardHeader,
      TaskProgressOverview,
      AnalysisEntryGrid,
      TaskTabsSection,
    },
    data() {
      return {
        // 学习效率分析相关状态
        showEfficiencyModal: false,
        showSummaryModal: false,
        showCheckInModal: false,
        efficiencyLoading: false,
        efficiencyError: "",
        efficiencyLoaded: false,
        // 学习效率分析数据
        efficiencyData: {
          weeklyStudyTime: 0,
          focusScore: 0,
          taskCompletionRate: 0,
          studyTrend: [],
          focusTrend: [],
          trendLabels: [],
          suggestions: [],
          insights: [],
          detailStats: {
            streakDays: 0,
            totalStudyMinutes: 0,
            completedTasks: 0,
            inProgressTasks: 0,
            completionRate: 0,
          },
        },
        // 智能总结数据
        summaryData: {
          reviewItems: [],
          reminders: [],
          knowledgeMap: {
            mastered: 0,
            learning: 0,
            toLearn: 0,
          },
          summary: "",
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
        activeTimeFilter: "week",
        dailyOverview: createEmptyTimeData(),
        statsLoading: false,
        statsError: "",
        analysisPrompt: "",
        // 不同时间段的数据
        timeFilterData: {
          week: createEmptyTimeData(),
          month: createEmptyTimeData(),
        },
      };
    },
    computed: {
      // 当前时间筛选器对应的数据
      currentTimeData() {
        return this.timeFilterData[this.activeTimeFilter] || EMPTY_TIME_DATA;
      },
      // 卡片摘要展示
      efficiencySummary() {
        return {
          weeklyStudyHours: this.efficiencyData.weeklyStudyTime,
          focusScore: this.efficiencyData.focusScore,
          taskCompletionRate: this.efficiencyData.taskCompletionRate,
        };
      },
    },
    mounted() {
      this.initializeBoard();
    },
    watch: {
      activeTimeFilter(newFilter) {
        this.fetchTaskStats(newFilter);
      },
    },
    methods: {
      async initializeBoard() {
        await Promise.allSettled([
          this.fetchDailyStats(),
          this.fetchTaskStats(this.activeTimeFilter),
        ]);
      },
      async loadEfficiencyAnalysis() {
        const user = getUserInfo();
        const userId = user?.id || 1;
        this.efficiencyLoading = true;
        this.efficiencyError = "";
        try {
          const res = await fetchEfficiencyAnalysis({
            user_id: userId,
            days: 14,
            model: "qwen-plus",
          });
          console.log("[AI] raw efficiency response:", res);
          const analysis = parseEfficiencyPayload(res);
          if (!analysis) {
            throw new Error("AI 返回数据为空");
          }
          const normalized = transformEfficiencyData(analysis);
          console.log("[AI] parsed analysis:", normalized);
          if (normalized.reviewPlan) {
            console.log("[AI] review plan:", normalized.reviewPlan);
          }
          this.efficiencyData = {
            ...this.efficiencyData,
            ...normalized,
          };
          this.analysisPrompt = normalized.prompt || analysis.prompt || "";
          if (normalized.reviewPlan) {
            this.applyReviewPlan(normalized.reviewPlan);
          }
          this.efficiencyLoaded = true;
          this.$nextTick(() => {
            this.initEfficiencyCharts();
          });
        } catch (error) {
          this.efficiencyError = error?.message || "加载学习效率失败";
          console.error("加载学习效率失败", error);
          ElMessage.error(this.efficiencyError);
        } finally {
          this.efficiencyLoading = false;
        }
      },
      applyReviewPlan(plan = {}) {
        const knowledgeMap = plan.knowledgeMap || this.summaryData.knowledgeMap;
        const reviewItems =
          Array.isArray(plan.reviewItems) && plan.reviewItems.length
            ? plan.reviewItems
            : this.summaryData.reviewItems;
        const reminders =
          Array.isArray(plan.reminders) && plan.reminders.length
            ? plan.reminders
            : this.summaryData.reminders;
        this.summaryData = {
          knowledgeMap: knowledgeMap || { mastered: 0, learning: 0, toLearn: 0 },
          reviewItems,
          reminders,
          summary: plan.summary || this.summaryData.summary || "",
        };
      },
      async fetchDailyStats() {
        this.statsLoading = true;
        this.statsError = "";
        try {
          const res = await getTaskBarStats("day");
          const mapped = this.transformBarStats(res?.data);
          this.dailyOverview = {
            completionRate: mapped.completionRate,
            completedTasks: mapped.completedTasks,
            totalTasks: mapped.totalTasks,
          };
          this.efficiencyData.taskCompletionRate = mapped.completionRate;
        } catch (error) {
          this.statsError = error?.message || "加载今日完成率失败";
          console.error("加载今日完成率失败", error);
          ElMessage.error(this.statsError);
        } finally {
          this.statsLoading = false;
        }
      },
      async fetchTaskStats(rangeKey = this.activeTimeFilter) {
        if (rangeKey === "day") {
          return this.fetchDailyStats();
        }
        const normalizedRange = ["week", "month"].includes(rangeKey)
          ? rangeKey
          : "week";
        this.statsLoading = true;
        this.statsError = "";
        try {
          const res = await getTaskBarStats(normalizedRange);
          const mapped = this.transformBarStats(res?.data);
          this.timeFilterData = {
            ...this.timeFilterData,
            [normalizedRange]: mapped,
          };
        } catch (error) {
          this.statsError = error?.message || "加载任务统计失败";
          console.error("加载任务统计失败", error);
          ElMessage.error(this.statsError);
        } finally {
          this.statsLoading = false;
        }
      },
      transformBarStats(payload = {}) {
        const items = Array.isArray(payload?.data) ? payload.data : [];
        const chartLabels = items.map((item) => item.day || "");
        const chartData = items.map((item) => clampPercentage(item.rate));
        const completedTasks = items.reduce(
          (sum, item) => sum + (Number(item.completed) || 0),
          0
        );
        const totalTasks = items.reduce(
          (sum, item) => sum + (Number(item.total) || 0),
          0
        );
        const completionRate =
          totalTasks > 0
            ? clampPercentage(Math.round((completedTasks / totalTasks) * 100))
            : 0;

        return {
          chartLabels,
          chartData,
          completionRate,
          completedTasks,
          totalTasks,
        };
      },
      showTaskDetails() {
        // 点击环形图显示任务详情的联动功能
        console.log("显示任务详情");
      },
      // 显示学习效率分析
      async showEfficiencyAnalysis() {
        this.showEfficiencyModal = true;
        if (!this.efficiencyLoaded && !this.efficiencyLoading) {
          await this.loadEfficiencyAnalysis();
        } else {
          this.$nextTick(() => {
            this.initEfficiencyCharts();
          });
        }
      },
      closeEfficiencyModal() {
        this.showEfficiencyModal = false;
      },
      // 智能总结与复习方法
      async showSmartSummary() {
        if (!this.efficiencyLoaded && !this.efficiencyLoading) {
          await this.loadEfficiencyAnalysis();
        }
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
      async regenerateEfficiency() {
        await this.loadEfficiencyAnalysis();
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
            const labels =
              this.efficiencyData.trendLabels?.length > 0
                ? this.efficiencyData.trendLabels
                : ["周一", "周二", "周三", "周四", "周五", "周六", "周日"];
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
                data: labels,
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
            const labels =
              this.efficiencyData.trendLabels?.length > 0
                ? this.efficiencyData.trendLabels
                : ["周一", "周二", "周三", "周四", "周五", "周六", "周日"];
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
                data: labels,
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


    },
  };
</script>

<style scoped></style>
