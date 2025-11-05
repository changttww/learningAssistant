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
          <div class="text-xs text-green-600 mt-1 flex items-center">
            <iconify-icon
              icon="mdi:trending-up"
              width="12"
              height="12"
              class="mr-1"
            ></iconify-icon>
            较上月增长 15%
          </div>
        </div>
        <div class="stat-card bg-green-50 p-4">
          <div class="text-2xl font-bold text-green-600">
            {{ taskCompletionRate }}
          </div>
          <div class="text-gray-600 mt-1 text-sm">任务完成率</div>
          <div class="text-xs text-green-600 mt-1 flex items-center">
            <iconify-icon
              icon="mdi:trending-up"
              width="12"
              height="12"
              class="mr-1"
            ></iconify-icon>
            较上月增长 8%
          </div>
        </div>
        <div class="stat-card bg-orange-50 p-4">
          <div class="text-2xl font-bold text-orange-600">
            {{ tasksInProgress }}
          </div>
          <div class="text-gray-600 mt-1 text-sm">进行中任务</div>
          <div class="text-xs text-green-600 mt-1 flex items-center">
            <iconify-icon
              icon="mdi:trending-up"
              width="12"
              height="12"
              class="mr-1"
            ></iconify-icon>
            正在完成 2 个任务
          </div>
        </div>
        <div class="stat-card bg-purple-50 p-4">
          <div class="text-2xl font-bold text-purple-600">
            {{ certificatesCount }}
          </div>
          <div class="text-gray-600 mt-1 text-sm">已获得成就</div>
          <div class="text-xs text-green-600 mt-1 flex items-center">
            <iconify-icon
              icon="mdi:trending-up"
              width="12"
              height="12"
              class="mr-1"
            ></iconify-icon>
            本月新增 3 个
          </div>
        </div>
      </div>

      <!-- 本日任务和知识点分布 - 移到上方并放大 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <!-- 本日任务 -->
        <div class="card p-6">
          <div class="flex justify-between items-center mb-4">
            <h2 class="font-bold text-xl">本日任务</h2>
            <router-link
              to="/personal-tasks"
              class="text-blue-600 hover:underline text-sm"
              >查看全部</router-link
            >
          </div>
          <div class="space-y-3">
            <div
              v-if="todayTasks.length === 0"
              class="text-gray-500 text-center py-8"
            >
              今日暂无任务
            </div>
            <div
              v-for="task in todayTasks"
              :key="task.id"
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <div class="flex items-center space-x-3">
                <input
                  type="checkbox"
                  :checked="task.status === 'completed'"
                  @change="toggleTaskStatus(task)"
                  class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
                />
                <div class="flex-1">
                  <div
                    :class="[
                      'font-medium text-sm',
                      task.status === 'completed'
                        ? 'line-through text-gray-500'
                        : 'text-gray-800',
                    ]"
                  >
                    {{ task.title }}
                  </div>
                  <div class="text-xs text-gray-500 mt-1">
                    {{ task.time }} · {{ task.category }}
                  </div>
                </div>
              </div>
              <div
                :class="[
                  'px-2 py-1 rounded-full text-xs',
                  {
                    'bg-green-100 text-green-800': task.status === 'completed',
                    'bg-orange-100 text-orange-800':
                      task.status === 'in-progress',
                    'bg-red-100 text-red-800': task.status === 'overdue',
                    'bg-gray-100 text-gray-800': task.status === 'pending',
                  },
                ]"
              >
                {{ getStatusText(task.status) }}
              </div>
            </div>
          </div>
        </div>

        <!-- 知识点分布 -->
        <div class="card p-6">
          <div class="flex justify-between items-center mb-4">
            <h2 class="font-bold text-xl">知识点分布</h2>
            <select
              class="bg-gray-50 border border-gray-300 text-gray-700 text-sm rounded-lg focus:ring-blue-500 py-1.5 px-3"
            >
              <option selected>全部领域</option>
              <option>技术类</option>
              <option>管理类</option>
              <option>设计类</option>
            </select>
          </div>
          <div class="chart-container" ref="knowledgeDistributionChart"></div>
        </div>
      </div>

      <!-- 主要内容区域：左侧内容 + 右侧图表 -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 左侧内容区域 -->
        <div class="lg:col-span-2 space-y-6">
          <!-- 学习时长统计图表 -->
          <div class="card">
            <div class="flex justify-between items-center mb-4">
              <h2 class="font-bold text-xl">学习时长趋势</h2>
              <div class="flex space-x-2">
                <select
                  class="bg-gray-50 border border-gray-300 text-gray-700 text-sm rounded-lg focus:ring-blue-500 py-1.5 px-3"
                >
                  <option selected>最近30天</option>
                  <option>最近90天</option>
                  <option>本年度</option>
                </select>
              </div>
            </div>
            <div class="chart-container" ref="studyTimeChart"></div>
          </div>
        </div>

        <!-- 右侧图表区域 -->
        <div class="space-y-6">
          <!-- 技能雷达图 -->
          <div class="card">
            <div class="flex justify-between items-center mb-4">
              <h2 class="font-bold text-xl">技能雷达</h2>
              <button class="text-blue-600 hover:underline text-sm">
                自定义
              </button>
            </div>
            <div class="chart-container" ref="skillRadarChart"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { computed, onMounted } from "vue";
  import * as echarts from "echarts";
  import {
    useCurrentUser,
    DEFAULT_USER_ID,
  } from "@/composables/useCurrentUser";

  export default {
    name: "Home",
    setup() {
      const {
        profile,
        loadCurrentUser,
        loadStudyStats,
        studyStats,
        studyStatsLoaded,
      } = useCurrentUser();

      onMounted(async () => {
        try {
          const loadedProfile = await loadCurrentUser();
          await loadStudyStats(loadedProfile?.id ?? DEFAULT_USER_ID);
        } catch (error) {
          console.error("加载用户详情失败:", error);
        }
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
        const rate = studyStats.value?.task_completion_rate;
        if (rate === null || rate === undefined) return "92%";
        return `${rate}%`;
      });
      const tasksInProgress = computed(
        () => studyStats.value?.tasks_in_progress ?? 8
      );
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
      };
    },
    data() {
      return {
        // 今日任务数据
        todayTasks: [
          {
            id: 1,
            title: "完成数学作业",
            description: "微积分练习题第3章",
            date: "2024-03-05",
            time: "14:00前",
            status: "completed",
            category: "数学",
          },
          {
            id: 2,
            title: "准备英语报告",
            description: "关于气候变化的演讲",
            date: "2024-03-05",
            time: "15:00前",
            status: "in-progress",
            category: "英语",
          },
          {
            id: 3,
            title: "物理实验预习",
            description: "波动光学实验操作流程",
            date: "2024-03-05",
            time: "17:00前",
            status: "pending",
            category: "物理",
          },
        ],
      };
    },
    mounted() {
      this.initCharts();
    },
    methods: {
      // 切换任务状态
      toggleTaskStatus(task) {
        if (task.status === "completed") {
          task.status = "pending";
        } else {
          task.status = "completed";
        }
      },

      // 获取状态文本
      getStatusText(status) {
        const statusMap = {
          completed: "已完成",
          "in-progress": "进行中",
          pending: "待完成",
          overdue: "已逾期",
        };
        return statusMap[status] || "未知";
      },
      initCharts() {
        // 学习时长趋势图
        const studyTimeChart = echarts.init(this.$refs.studyTimeChart);
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
            data: ["5/1", "5/3", "5/5", "5/7", "5/9", "5/11", "5/13"],
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
              data: [3.5, 2.8, 4.2, 3.0, 3.8, 4.5, 3.2],
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

        // 技能雷达图
        const skillRadarChart = echarts.init(this.$refs.skillRadarChart);
        const skillRadarOption = {
          tooltip: {
            trigger: "item",
          },
          radar: {
            indicator: [
              { name: "前端开发", max: 100 },
              { name: "后端开发", max: 100 },
              { name: "数据分析", max: 100 },
              { name: "项目管理", max: 100 },
              { name: "UI设计", max: 100 },
              { name: "软技能", max: 100 },
            ],
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
                  value: [85, 65, 70, 90, 60, 80],
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

        // 知识点分布图
        const knowledgeDistributionChart = echarts.init(
          this.$refs.knowledgeDistributionChart
        );
        const knowledgeDistributionOption = {
          tooltip: {
            trigger: "item",
            formatter: "{b}: {c}小时 ({d}%)",
          },
          legend: {
            bottom: "0%",
            left: "center",
            itemWidth: 10,
            itemHeight: 10,
            textStyle: {
              fontSize: 11,
            },
          },
          series: [
            {
              type: "pie",
              radius: ["40%", "70%"],
              center: ["50%", "45%"],
              avoidLabelOverlap: false,
              itemStyle: {
                borderRadius: 6,
                borderColor: "#fff",
                borderWidth: 2,
              },
              label: {
                show: false,
              },
              emphasis: {
                label: {
                  show: false,
                },
              },
              labelLine: {
                show: false,
              },
              data: [
                {
                  value: 35,
                  name: "前端技术",
                  itemStyle: { color: "#2D5BFF" },
                },
                {
                  value: 20,
                  name: "后端开发",
                  itemStyle: { color: "#34C759" },
                },
                {
                  value: 15,
                  name: "数据分析",
                  itemStyle: { color: "#FF9500" },
                },
                {
                  value: 25,
                  name: "项目管理",
                  itemStyle: { color: "#AF52DE" },
                },
                { value: 5, name: "其他", itemStyle: { color: "#FF3B30" } },
              ],
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
</style>
