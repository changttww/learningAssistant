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
                {{ totalKnowledgePointsLabel }}
              </div>
              <div class="text-xs text-gray-600">知识点总数</div>
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
            {{ totalKnowledgePointsLabel }}
          </div>
          <div class="text-gray-600 mt-1 text-sm">总知识点</div>
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

      <div class="ai-report-card mb-4">
        <div class="ai-report-main">
          <div class="ai-report-icon">
            <iconify-icon icon="mdi:chart-box-outline" width="26" height="26"></iconify-icon>
          </div>
          <div>
            <div class="ai-report-title">AI 学习报告</div>
            <div class="ai-report-desc">
              AI 已根据你的任务、知识点和学习趋势生成本周分析
            </div>
          </div>
        </div>
        <router-link to="/ai-report" class="ai-report-action">
          查看 AI 学习报告
          <iconify-icon icon="mdi:arrow-right" width="16" height="16"></iconify-icon>
        </router-link>
      </div>

      <!-- 任务热力图 -->
      <div class="card p-6 mb-6">
        <TaskHeatmap />
      </div>

      <!-- 知识点分布、技能雷达两列布局 -->
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
          <!-- 空状态 -->
          <div v-if="!knowledgeDistribution || knowledgeDistribution.length === 0" class="h-64 flex flex-col items-center justify-center text-gray-400">
            <span class="text-4xl mb-2">📚</span>
            <p class="text-sm">知识库中暂无条目</p>
            <router-link to="/knowledge-base" class="mt-2 text-blue-500 text-xs hover:underline">去添加知识点</router-link>
          </div>
          <!-- 图表 -->
          <div v-else class="chart-container h-64" ref="knowledgeDistributionChart"></div>
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

      <!-- 学习趋势 - 单列全宽 -->
      <div class="card p-6 mb-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="font-bold text-lg text-gray-900">📈 学习趋势</h2>
          <div class="flex space-x-2">
            <select
              v-model="trendRange"
              @change="handleTrendRangeChange"
              class="bg-gray-50 border border-gray-200 text-gray-700 text-xs rounded-lg focus:ring-2 focus:ring-blue-400 focus:border-transparent py-1.5 px-3 transition-all"
            >
              <option value="30">最近30天</option>
              <option value="90">最近90天</option>
              <option value="year">本年度</option>
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
  import { getTaskBarStats } from "@/api/modules/task";
  import { analyzeUserKnowledge, getSkillRadarData, getLearningTrends } from "@/api/modules/knowledge";
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

      // 总知识点数由 options(data/methods) 侧更新，这里仅作为展示占位，避免 setup 与 this 状态割裂
      const totalKnowledgePointsLabel = computed(() => "--");

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
        totalKnowledgePointsLabel,
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
        // 知识库分析数据
        knowledgeAnalysis: null,
        knowledgeDistribution: [],
        skillRadarData: [],
        learningTrends: [],

        // 总知识点（首页展示用）
        totalKnowledgePoints: 0,

        // 学习趋势范围：30/90/year
        trendRange: "30",
      };
    },
    computed: {
      // 用 options computed 覆盖 setup 同名字段（以 data 为准，且可响应更新）
      totalKnowledgePointsLabel() {
        return String(Number(this.totalKnowledgePoints) || 0);
      },
    },
    mounted() {
      // 加载图表数据
      this.fetchKnowledgeAnalysis().then(() => {
        // 将知识点总数同步给 setup 侧的 computed（来自知识库统计的分布数据）
        const totalFromDist = Array.isArray(this.knowledgeDistribution)
          ? this.knowledgeDistribution.reduce((sum, item) => sum + (Number(item?.count) || 0), 0)
          : 0;
        globalThis.__home_total_knowledge_points_from_distribution__ = totalFromDist;

        // 使用 nextTick 确保 DOM 完全就绪后再初始化图表
        this.$nextTick(() => {
          this.initCharts();
          console.log("[首页] 数据加载完成，图表已初始化");
        });
      }).catch((error) => {
        console.error("[首页] 数据加载出错:", error);
        this.$nextTick(() => {
          this.initCharts(); // 即使出错也初始化图表（使用默认数据）
        });
      });
    },
    beforeUnmount() {
      // 清理事件监听
    },
    methods: {
      unwrapArrayResponse(res) {
        const payload = res?.data ?? res;
        const arr = payload?.data ?? payload;
        return Array.isArray(arr) ? arr : null;
      },
      unwrapReportResponse(res) {
        const payload = res?.data ?? res;
        return payload?.data ?? payload;
      },
      applyKnowledgeReport(report) {
        this.knowledgeAnalysis = report;

        const distribution =
          report.knowledge_distribution || report.KnowledgeDistribution || [];
        // 不再使用默认数据，直接使用后端返回的数据（可能为空数组）
        this.knowledgeDistribution = Array.isArray(distribution) ? distribution : [];

        // 关键：总知识点 = 分布 count 求和
        this.totalKnowledgePoints = Array.isArray(this.knowledgeDistribution)
          ? this.knowledgeDistribution.reduce(
              (sum, item) => sum + (Number(item?.count) || 0),
              0
            )
          : 0;

        if (!Array.isArray(this.skillRadarData) || this.skillRadarData.length === 0) {
          const skillRadar = report.skill_radar || report.SkillRadar || [];
          this.skillRadarData =
            Array.isArray(skillRadar) && skillRadar.length > 0
              ? skillRadar
              : this.getDefaultSkillRadar();
        }

        if (!Array.isArray(this.learningTrends) || this.learningTrends.length === 0) {
          const trendsRaw = report.learning_trends || report.LearningTrends || [];
          this.learningTrends = Array.isArray(trendsRaw) ? trendsRaw : [];
        }

        if (!Array.isArray(this.skillRadarData) || this.skillRadarData.length === 0) {
          this.skillRadarData = this.getDefaultSkillRadar();
        }
      },

      // 加载知识库分析数据
      async fetchKnowledgeAnalysis() {
        const startTime = performance.now();
        try {
          console.log("[首页] 开始加载知识库分析数据");

          const [skillRes, trendRes] = await Promise.allSettled([
            getSkillRadarData(),
            getLearningTrends(this.trendRange),
          ]);

          if (skillRes.status === "fulfilled") {
            const arr = this.unwrapArrayResponse(skillRes.value);
            if (arr) this.skillRadarData = arr;
          }

          if (trendRes.status === "fulfilled") {
            const arr = this.unwrapArrayResponse(trendRes.value);
            if (arr) this.learningTrends = arr;
          }

          const res = await analyzeUserKnowledge();
          const report = this.unwrapReportResponse(res);

          if (!report) {
            console.warn("[首页] 知识库分析返回空数据");
            this.knowledgeDistribution = [];
            this.skillRadarData = this.skillRadarData?.length
              ? this.skillRadarData
              : this.getDefaultSkillRadar();
            this.learningTrends = Array.isArray(this.learningTrends)
              ? this.learningTrends
              : [];
            return;
          }

          this.applyKnowledgeReport(report);

          const loadTime = (performance.now() - startTime).toFixed(2);
          console.log(`[首页] 知识库分析已加载 (${loadTime}ms)`, {
            distribution: this.knowledgeDistribution.length,
            skillRadar: this.skillRadarData.length,
            trends: this.learningTrends.length,
          });
        } catch (error) {
          console.error("[首页] 加载知识库分析失败:", error);
          this.knowledgeDistribution = [];
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

      // 默认技能雷达数据 - 五大能力维度（与后端分类保持一致）
      getDefaultSkillRadar() {
        return [
          { skill: "理论素养", value: 50, max_value: 100, category: "理论素养" },
          { skill: "逻辑思维", value: 55, max_value: 100, category: "逻辑思维" },
          { skill: "实操应用", value: 60, max_value: 100, category: "实操应用" },
          { skill: "创新思维", value: 40, max_value: 100, category: "创新思维" },
          { skill: "沟通表达", value: 65, max_value: 100, category: "沟通表达" },
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
      handleTrendRangeChange() {
        // 范围切换：重新拉取后端聚合后的趋势数据（30=日，90=周，year=月）
        this.fetchLearningTrendsByRange(this.trendRange);
      },

      async fetchLearningTrendsByRange(range) {
        try {
          const res = await getLearningTrends(range);
          const arr = this.unwrapArrayResponse(res);
          if (arr) {
            this.learningTrends = arr;
          }
        } catch (e) {
          console.error("[首页] 加载学习趋势失败:", e);
        } finally {
          this.$nextTick(() => this.initCharts());
        }
      },

      buildTrendSeries() {
        // 返回给 ECharts 使用的 { dates, doneTasks, newNotes, newKnowledge }
        if (!this.learningTrends || this.learningTrends.length === 0) {
          return {
            dates: ["5/1", "5/3", "5/5", "5/7", "5/9", "5/11", "5/13"],
            doneTasks: [1, 0, 2, 1, 3, 1, 2],
            newNotes: [0, 1, 0, 1, 0, 2, 1],
            newKnowledge: [2, 1, 3, 2, 1, 2, 4],
          };
        }

        // 后端已按 range 做了聚合与补零，这里只需要做 label 格式化
        const recentTrends = [...this.learningTrends];

        const dates = recentTrends.map((t) => {
          const raw = t.date || t.day || t.created_at || t.createdAt;
          if (!raw) return "--";
          const s = typeof raw === "string" ? raw : String(raw);

          if (this.trendRange === "year") {
            // YYYY-MM -> M月
            if (/^\d{4}-\d{2}$/.test(s)) {
              const m = Number(s.slice(5, 7));
              return `${m}月`;
            }
          }

          if (this.trendRange === "90") {
            // YYYY-Www -> Wxx
            if (/^\d{4}-W\d{2}$/.test(s)) {
              return `W${s.slice(6, 8)}`;
            }
          }

          // 默认按天 YYYY-MM-DD -> M/D
          if (/^\d{4}-\d{2}-\d{2}/.test(s)) {
            return `${Number(s.slice(5, 7))}/${Number(s.slice(8, 10))}`;
          }

          // 兜底：尽量转 Date
          const normalized = s.includes("T") ? s : s + "T00:00:00";
          const d = new Date(normalized);
          if (Number.isNaN(d.getTime())) return s;
          return `${d.getMonth() + 1}/${d.getDate()}`;
        });

        return {
          dates,
          doneTasks: recentTrends.map((t) => Number(t.done_tasks ?? t.doneTasks ?? 0) || 0),
          newNotes: recentTrends.map((t) => Number(t.new_notes ?? t.newNotes ?? 0) || 0),
          newKnowledge: recentTrends.map((t) => Number(t.new_knowledge ?? t.newKnowledge ?? 0) || 0),
        };
      },
      initCharts() {
        console.log("[首页] 开始初始化图表...");
        console.log("[首页] 技能雷达数据:", this.skillRadarData);
        console.log("[首页] 学习趋势数据:", this.learningTrends);

        // 检查图表容器是否存在
        if (!this.$refs.studyTimeChart) {
          console.error("[首页] 学习趋势图表容器不存在");
          return;
        }
        if (!this.$refs.skillRadarChart) {
          console.error("[首页] 技能雷达图表容器不存在");
          return;
        }
        if (!this.$refs.knowledgeDistributionChart) {
          console.error("[首页] 知识分布图表容器不存在");
          return;
        }

        console.log("[首页] 图表容器检查通过，开始渲染图表...");
        console.log("[首页] 学习趋势数据长度:", this.learningTrends?.length || 0);
        console.log("[首页] 技能雷达数据长度:", this.skillRadarData?.length || 0);

        // 学习趋势图 - 使用知识库趋势数据（完成任务/创建笔记/新增知识点）
        const studyTimeChart = echarts.init(this.$refs.studyTimeChart);

        const trendSeries = this.buildTrendSeries();
        const trendDates = trendSeries.dates;
        const seriesDoneTasks = trendSeries.doneTasks;
        const seriesNewNotes = trendSeries.newNotes;
        const seriesNewKnowledge = trendSeries.newKnowledge;

        const studyTimeOption = {
          tooltip: {
            trigger: "axis",
          },
          legend: {
            top: 0,
            textStyle: { color: "#4b5563", fontSize: 12 },
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
            minInterval: 1,
            axisLine: {
              show: false,
            },
            axisLabel: {
              formatter: "{value}",
            },
            splitLine: {
              lineStyle: {
                color: "#f0f0f0",
              },
            },
          },
          series: [
            {
              name: "完成任务",
              type: "line",
              data: seriesDoneTasks,
              smooth: true,
              symbol: "circle",
              symbolSize: 7,
              itemStyle: { color: "#2D5BFF" },
              lineStyle: { width: 3, color: "#2D5BFF" },
              areaStyle: {
                color: {
                  type: "linear",
                  x: 0,
                  y: 0,
                  x2: 0,
                  y2: 1,
                  colorStops: [
                    { offset: 0, color: "rgba(45,91,255,0.18)" },
                    { offset: 1, color: "rgba(45,91,255,0.01)" },
                  ],
                },
              },
            },
            {
              name: "创建笔记",
              type: "line",
              data: seriesNewNotes,
              smooth: true,
              symbol: "circle",
              symbolSize: 6,
              itemStyle: { color: "#10B981" },
              lineStyle: { width: 2, color: "#10B981" },
            },
            {
              name: "新增知识点",
              type: "line",
              data: seriesNewKnowledge,
              smooth: true,
              symbol: "circle",
              symbolSize: 6,
              itemStyle: { color: "#F59E0B" },
              lineStyle: { width: 2, color: "#F59E0B" },
            },
          ],
        };
        studyTimeChart.setOption(studyTimeOption);

        // 技能雷达图 - 使用知识库分析的技能数据
        const skillRadarChart = echarts.init(this.$refs.skillRadarChart);

        let radarIndicators = [];
        let radarValues = [];

        const skillData =
          this.skillRadarData && this.skillRadarData.length > 0
            ? this.skillRadarData.slice(0, 6)
            : this.getDefaultSkillRadar();

        radarIndicators = skillData.map((s) => ({
          name: s.skill,
          max: Number(s.max_value ?? s.max ?? 100) || 100,
        }));
        radarValues = skillData.map((s) => Number(s.value ?? s.score ?? 0) || 0);

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

        // 知识点分布图 - 仅当有数据时才初始化
        if (this.knowledgeDistribution && this.knowledgeDistribution.length > 0 && this.$refs.knowledgeDistributionChart) {
          const knowledgeDistributionChart = echarts.init(
            this.$refs.knowledgeDistributionChart
          );
          
          // 从knowledgeDistribution中提取数据
          const distData = this.knowledgeDistribution;
          
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
        } else {
          // 没有知识分布数据时，只监听其他图表的resize
          window.addEventListener("resize", () => {
            studyTimeChart.resize();
            skillRadarChart.resize();
          });
        }
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

  .ai-report-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    padding: 18px 20px;
    border: 1px solid rgba(45, 91, 255, 0.14);
    border-radius: 0.75rem;
    background:
      linear-gradient(135deg, rgba(45, 91, 255, 0.08) 0%, rgba(16, 185, 129, 0.08) 100%),
      #ffffff;
    box-shadow: 0 2px 10px rgba(45, 91, 255, 0.06);
  }

  .ai-report-main {
    display: flex;
    align-items: center;
    gap: 14px;
    min-width: 0;
  }

  .ai-report-icon {
    width: 46px;
    height: 46px;
    flex: 0 0 auto;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 12px;
    color: #2d5bff;
    background: rgba(45, 91, 255, 0.1);
  }

  .ai-report-title {
    font-size: 16px;
    font-weight: 700;
    color: #111827;
  }

  .ai-report-desc {
    margin-top: 4px;
    font-size: 13px;
    color: #4b5563;
    line-height: 1.5;
  }

  .ai-report-action {
    flex: 0 0 auto;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    min-height: 38px;
    padding: 0 14px;
    border-radius: 10px;
    color: #ffffff;
    background: #2d5bff;
    font-size: 13px;
    font-weight: 600;
    transition: all 0.2s ease;
  }

  .ai-report-action:hover {
    background: #2448cc;
    box-shadow: 0 6px 14px rgba(45, 91, 255, 0.22);
  }

  @media (max-width: 640px) {
    .ai-report-card {
      align-items: stretch;
      flex-direction: column;
    }

    .ai-report-action {
      width: 100%;
    }
  }
</style>
