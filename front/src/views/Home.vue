<template>
  <div class="w-full h-full overflow-auto">
    <!-- 个人信息顶部卡片 -->
    <div class="card mb-4">
      <div class="flex flex-col md:flex-row items-center md:items-start gap-4">
        <div
          class="w-24 h-24 rounded-full bg-gray-200 overflow-hidden flex items-center justify-center text-3xl font-semibold text-[#2D5BFF]"
        >
          <img
            v-if="userAvatar"
            :src="userAvatar"
            :alt="displayName"
            class="w-full h-full object-cover"
          />
          <span v-else>{{ displayName.slice(0, 1) }}</span>
        </div>
        <div class="flex-1">
          <div class="flex flex-col md:flex-row md:items-center gap-3">
            <h1 class="text-xl font-bold">{{ displayName }}</h1>
            <div class="flex items-center gap-2">
              <div
                class="bg-blue-100 text-blue-800 px-2 py-1 rounded-full text-xs font-medium"
              >
                {{ userRole }}
              </div>
              <div
                class="bg-green-100 text-green-800 px-2 py-1 rounded-full text-xs font-medium flex items-center gap-1"
              >
                <span class="w-1.5 h-1.5 rounded-full bg-green-500"></span>
                在线
              </div>
            </div>
          </div>
          <p class="text-gray-600 mt-2 text-sm">
            {{ userBio }}
          </p>
          <div class="mt-3 flex flex-wrap gap-3 text-sm">
            <div class="flex items-center gap-1">
              <iconify-icon
                icon="mdi:school"
                width="16"
                height="16"
                class="text-gray-600"
              ></iconify-icon>
              <span>{{ userSchoolMajor }}</span>
            </div>
            <div class="flex items-center gap-1">
              <iconify-icon
                icon="mdi:map-marker"
                width="16"
                height="16"
                class="text-gray-600"
              ></iconify-icon>
              <span>{{ userLocation }}</span>
            </div>
            <div class="flex items-center gap-1">
              <iconify-icon
                icon="mdi:account-group"
                width="16"
                height="16"
                class="text-gray-600"
              ></iconify-icon>
              <span>已加入{{ studyGroupCount }}个学习小组</span>
            </div>
          </div>
          <div class="mt-3 flex gap-2">
            <button
              class="bg-[#2D5BFF] text-white font-medium py-1.5 px-3 rounded-lg text-xs hover:bg-opacity-90 transition-colors flex items-center gap-1"
            >
              <iconify-icon
                icon="mdi:pencil"
                width="14"
                height="14"
              ></iconify-icon>
              编辑个人资料
            </button>
            <button
              class="bg-gray-100 text-gray-700 font-medium py-1.5 px-3 rounded-lg text-xs hover:bg-gray-200 transition-colors flex items-center gap-1"
            >
              <iconify-icon
                icon="mdi:share"
                width="14"
                height="14"
              ></iconify-icon>
              分享主页
            </button>
          </div>
        </div>
        <div class="md:text-right">
          <div class="flex flex-col items-center md:items-end">
            <div
              class="bg-blue-50 text-blue-800 px-2 py-1 rounded-full text-sm font-medium"
            >
              {{ levelLabel }}
            </div>
            <div class="mt-2 text-center md:text-right">
              <div class="text-lg font-bold text-[#10B981]">
                {{ totalStudyHoursLabel }}
              </div>
              <div class="text-xs text-gray-600">学习时长</div>
              <div class="text-xs text-gray-500 mt-1">
                距离下一级还需 {{ pointsToNextLevel }} 积分
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 学习统计内容区 -->
    <div class="w-full">
      <!-- 学习数据卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-3 mb-4">
        <div class="stat-card bg-blue-50 p-4">
          <div class="text-2xl font-bold text-blue-600">
            {{ totalStudyHoursLabel }}
          </div>
          <div class="text-gray-600 mt-1 text-sm">总学习时长</div>
        </div>
        <div class="stat-card bg-green-50 p-4">
          <div class="text-2xl font-bold text-green-600">
            {{ taskCompletionRate }}
          </div>
          <div class="text-gray-600 mt-1 text-sm">周任务完成率</div>
        </div>
        <div class="stat-card bg-orange-50 p-4">
          <div class="text-2xl font-bold text-orange-600">
            {{ tasksInProgress }}
          </div>
          <div class="text-gray-600 mt-1 text-sm">进行中任务</div>
        </div>
        <div class="stat-card bg-purple-50 p-4">
          <div class="text-2xl font-bold text-purple-600">
            {{ certificatesCount }}
          </div>
          <div class="text-gray-600 mt-1 text-sm">已获得成就</div>
        </div>
      </div>

      <!-- 任务热力图 + 本日任务双列布局 -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
        <!-- 任务热力图 - 占据2列 -->
        <div class="card p-6 lg:col-span-2">
          <TaskHeatmap />
        </div>

        <!-- 本日任务 - 右列 -->
        <div class="card p-6 flex flex-col">
          <div class="flex justify-between items-center mb-4">
            <h2 class="font-bold text-lg text-gray-900">📋 本日任务</h2>
            <router-link
              to="/personal-tasks"
              class="text-blue-600 hover:text-blue-700 hover:underline text-xs font-medium"
              >全部→</router-link
            >
          </div>
          <div class="space-y-2.5 flex-1 overflow-y-auto max-h-96 pr-2">
            <div
              v-if="todayTasks.length === 0"
              class="text-gray-400 text-center py-8 text-sm"
            >
              ✨ 今日暂无任务
            </div>
            <div
              v-for="task in todayTasks"
              :key="task.id"
              class="flex items-center justify-between p-2.5 bg-gradient-to-r from-gray-50 to-transparent rounded-lg hover:from-blue-50 hover:to-transparent transition-all duration-200 border border-transparent hover:border-blue-100"
            >
              <div class="flex items-center space-x-2 flex-1 min-w-0">
                <span
                  :class="[
                    'w-1.5 h-1.5 rounded-full flex-shrink-0',
                    task.status === 'completed'
                      ? 'bg-green-500'
                      : task.status === 'in-progress'
                        ? 'bg-orange-500'
                        : 'bg-gray-300',
                  ]"
                ></span>
                <div class="flex-1 min-w-0">
                  <div
                    :class="[
                      'font-medium text-xs truncate',
                      task.status === 'completed'
                        ? 'line-through text-gray-400'
                        : 'text-gray-700',
                    ]"
                    :title="task.title"
                  >
                    {{ task.title }}
                  </div>
                </div>
              </div>
              <div
                :class="[
                  'px-1.5 py-0.5 rounded-full text-xs font-medium whitespace-nowrap ml-2 flex-shrink-0',
                  getStatusClass(task.status),
                ]"
              >
                {{ getStatusText(task.status) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 知识点分布、技能雷达二列布局 -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <!-- 知识点分布 - 左列 -->
        <div class="card p-6">
          <div class="flex justify-between items-center mb-4">
            <h2 class="font-bold text-lg text-gray-900">🎯 知识分布</h2>
            <router-link
              to="/knowledge-base"
              class="text-blue-600 hover:text-blue-700 text-xs font-medium hover:underline"
              >查看知识库→</router-link
            >
          </div>
          <div class="chart-container h-64" ref="knowledgeDistributionChart"></div>
        </div>

        <!-- 技能雷达 - 右列 -->
        <div class="card p-6">
          <div class="flex justify-between items-center mb-4">
            <h2 class="font-bold text-lg text-gray-900">⚡ 技能雷达</h2>
            <router-link
              to="/knowledge-base"
              class="text-blue-600 hover:text-blue-700 text-xs font-medium hover:underline"
              >管理技能→</router-link>
          </div>
          <div class="chart-container h-64" ref="skillRadarChart"></div>
        </div>
      </div>

      <!-- 学习时长趋势 - 单列全宽 -->
      <div class="card p-6 mb-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="font-bold text-lg text-gray-900">📈 学习趋势</h2>
          <div class="flex space-x-2">
            <select
              class="bg-gray-50 border border-gray-200 text-gray-700 text-xs rounded-lg focus:ring-2 focus:ring-blue-400 focus:border-transparent py-1.5 px-3 transition-all"
            >
              <option selected>最近30天</option>
              <option>最近90天</option>
              <option>本年度</option>
            </select>
          </div>
        </div>
        <div class="chart-container h-80" ref="studyTimeChart"></div>
      </div>
    </div>
  </div>
</template>

<script>
  import { computed, onMounted, ref } from "vue";
  import * as echarts from "echarts";
  import {
    useCurrentUser,
    DEFAULT_USER_ID,
  } from "@/composables/useCurrentUser";
  import { getTaskBarStats, getTodayTasks } from "@/api/modules/task";
  import { analyzeUserKnowledge } from "@/api/modules/knowledge";
  import { generatePieChartData, getSubjectConfig } from "@/utils/subjectConfig";
  import TaskHeatmap from "@/components/TaskHeatmap.vue";

  export default {
    name: "Home",
    components: {
      TaskHeatmap,
    },
    setup() {
      const {
        profile,
        loadCurrentUser,
        loadStudyStats,
        studyStats,
        studyStatsLoaded,
      } = useCurrentUser();

      const currentUserId = ref(DEFAULT_USER_ID);
      const taskBarStats = ref(null);

      const clampPercentage = (value) => {
        const numeric = Number(value);
        if (!Number.isFinite(numeric)) return 0;
        if (numeric < 0) return 0;
        if (numeric > 100) return 100;
        return Math.round(numeric);
      };

      const mapBarStats = (payload = {}) => {
        const items = Array.isArray(payload?.data) ? payload.data : [];
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
          completionRate,
          completedTasks,
          totalTasks,
        };
      };

      const fetchTaskBarStats = async (rangeKey = "week") => {
        try {
          const res = await getTaskBarStats(rangeKey);
          taskBarStats.value = mapBarStats(res?.data);
        } catch (error) {
          console.error("加载任务统计失败:", error);
        }
      };

      onMounted(async () => {
        let loadedProfile = null;
        try {
          loadedProfile = await loadCurrentUser();
          if (loadedProfile?.id) {
            currentUserId.value = loadedProfile.id;
          }
        } catch (error) {
          console.error("加载用户详情失败:", error);
        }

        await Promise.allSettled([
          loadStudyStats(loadedProfile?.id ?? DEFAULT_USER_ID),
          fetchTaskBarStats("week"),
        ]);
      });

      const displayName = computed(() => profile.value?.display_name || "学习者");
      const userAvatar = computed(() => profile.value?.avatar_url || "");
      const userRole = computed(() => profile.value?.role || "学习者");
      const userBio = computed(
        () =>
          profile.value?.bio ||
          "专注于自我提升，期待开启新的学习旅程。"
      );
      const userSchoolMajor = computed(() => {
        const school = profile.value?.basic_info?.school;
        const major = profile.value?.basic_info?.major;
        if (school && major) return `${school} ${major}`;
        if (school) return school;
        if (major) return major;
        return "学校与专业未填写";
      });
      const userLocation = computed(
        () => profile.value?.basic_info?.location || "所在地未填写"
      );
      const levelLabel = computed(
        () => studyStats.value?.level_label || "成长中学员"
      );
      const totalStudyHoursLabel = computed(() => {
        const hours = studyStats.value?.total_study_hours;
        if (hours === null || hours === undefined) return "--";
        return `${hours}h`;
      });
      const pointsToNextLevel = computed(
        () => studyStats.value?.distance_to_next ?? 0
      );
      const studyGroupCount = computed(
        () => studyStats.value?.study_groups ?? 0
      );
      const taskCompletionRate = computed(() => {
        const rate = taskBarStats.value?.completionRate;
        if (rate === null || rate === undefined) {
          const profileRate = studyStats.value?.task_completion_rate;
          if (profileRate === null || profileRate === undefined) return "92%";
          return `${profileRate}%`;
        }
        return `${rate}%`;
      });
      const tasksInProgress = computed(() => {
        if (taskBarStats.value) {
          const { totalTasks = 0, completedTasks = 0 } = taskBarStats.value;
          const inProgress = Math.max(totalTasks - completedTasks, 0);
          if (totalTasks || completedTasks) return inProgress;
        }
        return studyStats.value?.tasks_in_progress ?? 8;
      });
      const certificatesCount = computed(
        () => studyStats.value?.certificates_count ?? 24
      );

      return {
        displayName,
        userAvatar,
        userRole,
        userBio,
        userSchoolMajor,
        userLocation,
        levelLabel,
        totalStudyHoursLabel,
        studyGroupCount,
        taskCompletionRate,
        tasksInProgress,
        certificatesCount,
        studyStatsLoaded,
        pointsToNextLevel,
        currentUserId,
      };
    },
    data() {
      return {
        // 今日任务数据
        todayTasks: [],
        taskRefreshInterval: null,
        // 知识库分析数据
        knowledgeAnalysis: null,
        knowledgeDistribution: [],
        skillRadarData: [],
        learningTrends: [],
      };
    },
    mounted() {
      // 并行加载图表和任务数据，提高加载速度
      Promise.all([
        this.fetchKnowledgeAnalysis(),
        this.fetchTodayTasks(),
      ]).then(() => {
        // 先获取知识库数据再初始化图表
        this.initCharts();
        console.log("[首页] 数据加载完成");
      }).catch((error) => {
        console.error("[首页] 数据加载出错:", error);
        this.initCharts(); // 即使出错也初始化图表（使用默认数据）
      });
      
      // 15秒自动刷新一次今日任务
      this.taskRefreshInterval = setInterval(() => {
        console.log("[首页] 自动刷新今日任务");
        this.fetchTodayTasks();
      }, 15000);
      
      // 监听任务创建、完成等事件
      globalThis.addEventListener("taskUpdated", this.handleTaskUpdate);
      globalThis.addEventListener("taskCreated", this.handleTaskUpdate);
      globalThis.addEventListener("taskCompleted", this.handleTaskUpdate);
      window.addEventListener("focus", this.handleWindowFocus);
    },
    beforeUnmount() {
      // 清理定时器和事件监听
      if (this.taskRefreshInterval) {
        clearInterval(this.taskRefreshInterval);
      }
      globalThis.removeEventListener("taskUpdated", this.handleTaskUpdate);
      globalThis.removeEventListener("taskCreated", this.handleTaskUpdate);
      globalThis.removeEventListener("taskCompleted", this.handleTaskUpdate);
      window.removeEventListener("focus", this.handleWindowFocus);
    },
    methods: {
      // 加载知识库分析数据
      async fetchKnowledgeAnalysis() {
        try {
          console.log("[首页] 开始加载知识库分析数据");
          const startTime = performance.now();
          
          const res = await analyzeUserKnowledge();
          console.log("[首页] 知识库分析API响应:", res);
          
          // 兼容多种响应格式
          const analysisData = res?.data || res;
          
          if (analysisData) {
            this.knowledgeAnalysis = analysisData;
            
            // 解析知识分布
            const distribution = analysisData.knowledge_distribution || [];
            this.knowledgeDistribution = distribution.length > 0 ? distribution : this.getDefaultDistribution();
            
            // 解析技能雷达
            const skillRadar = analysisData.skill_radar || [];
            this.skillRadarData = skillRadar.length > 0 ? skillRadar : this.getDefaultSkillRadar();
            
            // 解析学习趋势
            this.learningTrends = analysisData.learning_trends || [];
            
            const loadTime = (performance.now() - startTime).toFixed(2);
            console.log(`[首页] 知识库分析已加载 (${loadTime}ms)`, {
              distribution: this.knowledgeDistribution.length,
              skillRadar: this.skillRadarData.length,
              trends: this.learningTrends.length,
              isDefault: distribution.length === 0
            });
          } else {
            console.warn("[首页] 知识库分析返回空数据，使用默认数据");
            this.knowledgeDistribution = this.getDefaultDistribution();
            this.skillRadarData = this.getDefaultSkillRadar();
            this.learningTrends = [];
          }
        } catch (error) {
          console.error("[首页] 加载知识库分析失败:", error);
          // 使用默认数据
          this.knowledgeDistribution = this.getDefaultDistribution();
          this.skillRadarData = this.getDefaultSkillRadar();
          this.learningTrends = [];
        }
      },

      // 默认知识分布数据 - 面向学习场景
      getDefaultDistribution() {
        return [
          { category: "数学", count: 15, percentage: 25, color: "#3b82f6", icon: "mdi:calculator-variant", gradient: "linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)" },
          { category: "语文", count: 12, percentage: 20, color: "#f59e0b", icon: "mdi:book-open-page-variant", gradient: "linear-gradient(135deg, #f59e0b 0%, #d97706 100%)" },
          { category: "英语", count: 10, percentage: 17, color: "#ec4899", icon: "mdi:alphabetical", gradient: "linear-gradient(135deg, #ec4899 0%, #db2777 100%)" },
          { category: "物理", count: 8, percentage: 13, color: "#8b5cf6", icon: "mdi:atom", gradient: "linear-gradient(135deg, #8b5cf6 0%, #6d28d9 100%)" },
          { category: "历史", count: 6, percentage: 10, color: "#92400e", icon: "mdi:castle", gradient: "linear-gradient(135deg, #92400e 0%, #78350f 100%)" },
          { category: "其他", count: 9, percentage: 15, color: "#64748b", icon: "mdi:bookshelf", gradient: "linear-gradient(135deg, #94a3b8 0%, #64748b 100%)" },
        ];
      },

      // 默认技能雷达数据 - 面向学习场景
      getDefaultSkillRadar() {
        return [
          { skill: "数学能力", value: 75, max_value: 100 },
          { skill: "语言表达", value: 80, max_value: 100 },
          { skill: "逻辑思维", value: 70, max_value: 100 },
          { skill: "阅读理解", value: 85, max_value: 100 },
          { skill: "创新思维", value: 65, max_value: 100 },
          { skill: "动手实践", value: 70, max_value: 100 },
        ];
      },

      normalizeStatus(status) {
        const normalized =
          typeof status === "string" ? status.trim().toLowerCase() : status;
        switch (normalized) {
          case 2:
          case "2":
          case "completed":
            return "completed";
          case 1:
          case "1":
          case "in-progress":
          case "in_progress":
            return "in-progress";
          case "overdue":
            return "overdue";
          case 0:
          case "0":
          case "pending":
          default:
            return "pending";
        }
      },
      async fetchTodayTasks() {
        const userId = this.currentUserId || DEFAULT_USER_ID;
        try {
          console.log("[首页] 开始加载今日任务");
          const startTime = performance.now();
          
          const res = await getTodayTasks(userId);
          const payload = res?.data || {};
          const merged = [
            ...(payload.completed || []),
            ...(payload.in_progress || []),
            ...(payload.not_started || []),
          ];

          this.todayTasks = merged.map((task) => ({
            id: task.id,
            title: task.title || "未命名任务",
            status: this.normalizeStatus(task.status),
          }));
          
          const loadTime = (performance.now() - startTime).toFixed(2);
          console.log(`[首页] 今日任务已加载: ${this.todayTasks.length} 个 (${loadTime}ms)`);
        } catch (error) {
          console.error("加载今日任务失败:", error);
          this.todayTasks = [];
        }
      },

      // 处理任务更新事件
      handleTaskUpdate() {
        console.log("[首页] 检测到任务变化，立即刷新");
        this.fetchTodayTasks();
      },

      // 处理窗口获焦事件
      handleWindowFocus() {
        console.log("[首页] 窗口获得焦点，刷新今日任务");
        this.fetchTodayTasks();
      },

      // 获取状态文本
      getStatusText(status) {
        const normalized = this.normalizeStatus(status);
        const statusMap = {
          completed: "已完成",
          "in-progress": "进行中",
          pending: "待完成",
          overdue: "已逾期",
        };
        return statusMap[normalized] || "未知";
      },
      getStatusClass(status) {
        const normalized = this.normalizeStatus(status);
        return {
          "bg-green-100 text-green-800": normalized === "completed",
          "bg-orange-100 text-orange-800": normalized === "in-progress",
          "bg-red-100 text-red-800": normalized === "overdue",
          "bg-gray-100 text-gray-800": normalized === "pending",
        };
      },
      initCharts() {
        // 学习时长趋势图 - 使用知识库趋势数据
        const studyTimeChart = echarts.init(this.$refs.studyTimeChart);
        
        // 从learningTrends中提取数据，如果没有就使用默认值
        let trendDates = [];
        let trendHours = [];
        
        if (this.learningTrends && this.learningTrends.length > 0) {
          // 取最近7-14天的数据
          const recentTrends = this.learningTrends.slice(-14);
          trendDates = recentTrends.map(t => {
            const date = new Date(t.date);
            return `${date.getMonth()+1}/${date.getDate()}`;
          });
          trendHours = recentTrends.map(t => t.study_hours || 0);
        } else {
          // 默认数据
          trendDates = ["5/1", "5/3", "5/5", "5/7", "5/9", "5/11", "5/13"];
          trendHours = [3.5, 2.8, 4.2, 3.0, 3.8, 4.5, 3.2];
        }
        
        const studyTimeOption = {
          tooltip: {
            trigger: "axis",
            formatter: function (params) {
              return `${params[0].name}<br/>学习时长: ${params[0].value}小时`;
            },
          },
          grid: {
            left: "3%",
            right: "4%",
            bottom: "3%",
            containLabel: true,
          },
          xAxis: {
            type: "category",
            boundaryGap: false,
            data: trendDates,
            axisLine: {
              lineStyle: {
                color: "#ddd",
              },
            },
          },
          yAxis: {
            type: "value",
            axisLine: {
              show: false,
            },
            axisLabel: {
              formatter: "{value}h",
            },
            splitLine: {
              lineStyle: {
                color: "#f0f0f0",
              },
            },
          },
          series: [
            {
              name: "学习时长",
              type: "line",
              data: trendHours,
              smooth: true,
              symbol: "circle",
              symbolSize: 8,
              itemStyle: {
                color: "#2D5BFF",
              },
              lineStyle: {
                width: 3,
                color: "#2D5BFF",
              },
              areaStyle: {
                color: {
                  type: "linear",
                  x: 0,
                  y: 0,
                  x2: 0,
                  y2: 1,
                  colorStops: [
                    {
                      offset: 0,
                      color: "rgba(45,91,255,0.2)",
                    },
                    {
                      offset: 1,
                      color: "rgba(45,91,255,0.01)",
                    },
                  ],
                },
              },
            },
          ],
        };
        studyTimeChart.setOption(studyTimeOption);

        // 技能雷达图 - 使用知识库分析的技能数据
        const skillRadarChart = echarts.init(this.$refs.skillRadarChart);
        
        // 从skillRadarData中提取数据，如果没有就使用默认值
        let radarIndicators = [];
        let radarValues = [];
        
        const skillData = this.skillRadarData && this.skillRadarData.length > 0 
          ? this.skillRadarData.slice(0, 6) 
          : this.getDefaultSkillRadar();
          
        radarIndicators = skillData.map(s => ({ name: s.skill, max: s.max_value || 100 }));
        radarValues = skillData.map(s => s.value || 0);
        
        const skillRadarOption = {
          tooltip: {
            trigger: "item",
          },
          radar: {
            indicator: radarIndicators,
            radius: "65%",
            splitNumber: 4,
            axisName: {
              color: "#333",
              fontSize: 12,
            },
            splitArea: {
              areaStyle: {
                color: ["#f5f5f5", "#e9e9e9", "#f0f0f0", "#fff"],
              },
            },
          },
          series: [
            {
              type: "radar",
              data: [
                {
                  value: radarValues,
                  name: "技能掌握度",
                  symbol: "circle",
                  symbolSize: 6,
                  lineStyle: {
                    width: 2,
                  },
                  areaStyle: {
                    color: "rgba(45,91,255,0.3)",
                  },
                  itemStyle: {
                    color: "#2D5BFF",
                  },
                },
              ],
            },
          ],
        };
        skillRadarChart.setOption(skillRadarOption);

        // 知识点分布图 - 使用知识库分析数据
        const knowledgeDistributionChart = echarts.init(
          this.$refs.knowledgeDistributionChart
        );
        
        // 从knowledgeDistribution中提取数据
        const distData = this.knowledgeDistribution && this.knowledgeDistribution.length > 0 
          ? this.knowledgeDistribution 
          : this.getDefaultDistribution();
        
        // 使用 subjectConfig 生成带渐变色的饼图数据
        const pieData = generatePieChartData(distData);
        
        const knowledgeDistributionOption = {
          tooltip: {
            trigger: "item",
            backgroundColor: 'rgba(255, 255, 255, 0.95)',
            borderColor: '#e5e7eb',
            borderWidth: 1,
            textStyle: {
              color: '#374151'
            },
            formatter: function(params) {
              const config = getSubjectConfig(params.name);
              return `<div style="display:flex;align-items:center;gap:8px;">
                <span style="font-size:16px;">${config.emoji}</span>
                <span style="font-weight:bold;color:${config.color}">${params.name}</span>
              </div>
              <div style="margin-top:4px;">
                📚 知识点: <b>${params.value}</b> 个<br/>
                📊 占比: <b>${params.percent}%</b>
              </div>`;
            }
          },
          legend: {
            bottom: "0%",
            left: "center",
            itemWidth: 12,
            itemHeight: 12,
            itemGap: 15,
            textStyle: {
              fontSize: 12,
              color: '#4b5563',
              fontWeight: 500
            },
            icon: 'circle'
          },
          series: [
            {
              type: "pie",
              radius: ["45%", "75%"],
              center: ["50%", "45%"],
              avoidLabelOverlap: false,
              itemStyle: {
                borderRadius: 8,
                borderColor: "#fff",
                borderWidth: 3,
                shadowBlur: 10,
                shadowColor: 'rgba(0, 0, 0, 0.1)'
              },
              label: {
                show: false,
              },
              emphasis: {
                scale: true,
                scaleSize: 8,
                itemStyle: {
                  shadowBlur: 20,
                  shadowColor: 'rgba(0, 0, 0, 0.2)'
                },
                label: {
                  show: true,
                  fontSize: 14,
                  fontWeight: 'bold',
                  formatter: '{b}\n{c}个'
                },
              },
              labelLine: {
                show: false,
              },
              data: pieData,
            },
          ],
        };
        knowledgeDistributionChart.setOption(knowledgeDistributionOption);

        // 窗口调整大小时调整图表大小
        window.addEventListener("resize", () => {
          studyTimeChart.resize();
          skillRadarChart.resize();
          knowledgeDistributionChart.resize();
        });
      },
    },
  };
</script>

<style scoped>
  .container {
    max-width: 1440px;
    margin: 0 auto;
    padding: 20px;
  }

  .chart-container {
    width: 100%;
    min-height: 300px;
  }

  .card {
    background: linear-gradient(135deg, #ffffff 0%, #fafbfc 100%);
    border: 1px solid rgba(0, 0, 0, 0.05);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    border-radius: 0.75rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .card:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    border-color: rgba(0, 0, 0, 0.08);
  }

  .stat-card {
    border: 1px solid rgba(0, 0, 0, 0.04);
    border-radius: 0.75rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .stat-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  }
</style>
