<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- 顶部导航 -->
    <header class="bg-white shadow-sm z-10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <button 
            @click="$router.push({ name: 'TeamTasks', params: { teamId } })"
            class="p-2 rounded-full hover:bg-gray-100 text-gray-600 transition-colors"
            title="返回任务看板"
          >
            <iconify-icon icon="mdi:arrow-left" width="24"></iconify-icon>
          </button>
          <h1 class="text-xl font-bold text-gray-800 flex items-center gap-2">
            <iconify-icon icon="mdi:chart-box-outline" class="text-blue-600"></iconify-icon>
            团队数据报告
          </h1>
        </div>
        <div class="text-sm text-gray-500">
          数据实时统计
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="flex-1 max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8">
      
      <div v-if="loading" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        
        <!-- 任务健康度 -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <h2 class="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
            <iconify-icon icon="mdi:heart-pulse" class="text-red-500"></iconify-icon>
            任务健康度
          </h2>
          <div ref="healthChart" class="w-full h-64"></div>
          <div class="mt-4 grid grid-cols-3 gap-4 text-center text-sm">
            <div class="p-2 bg-green-50 rounded-lg">
              <div class="text-green-600 font-bold text-lg">{{ stats.completed }}</div>
              <div class="text-gray-500">已完成</div>
            </div>
            <div class="p-2 bg-blue-50 rounded-lg">
              <div class="text-blue-600 font-bold text-lg">{{ stats.inProgress }}</div>
              <div class="text-gray-500">进行中</div>
            </div>
            <div class="p-2 bg-red-50 rounded-lg">
              <div class="text-red-600 font-bold text-lg">{{ stats.overdue }}</div>
              <div class="text-gray-500">已逾期</div>
            </div>
          </div>
        </div>

        <!-- 成员贡献分布 -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <h2 class="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
            <iconify-icon icon="mdi:account-group" class="text-purple-500"></iconify-icon>
            成员贡献分布 (已完成任务)
          </h2>
          <div ref="contributionChart" class="w-full h-64"></div>
        </div>

        <!-- 燃尽图 -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 lg:col-span-2">
          <h2 class="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
            <iconify-icon icon="mdi:fire" class="text-orange-500"></iconify-icon>
            项目燃尽图 (Burndown Chart)
          </h2>
          <div ref="burndownChart" class="w-full h-80"></div>
          <p class="text-xs text-gray-400 mt-2 text-center">
            * 基于任务完成时间（更新时间）估算的剩余任务趋势
          </p>
        </div>

      </div>
    </main>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import { getTeamTasks } from '@/api/modules/task';
import { getTeamMembers } from '@/api/modules/team';

export default {
  name: 'TeamReports',
  data() {
    return {
      teamId: null,
      loading: true,
      tasks: [],
      members: [],
      stats: {
        completed: 0,
        inProgress: 0,
        overdue: 0,
        pending: 0
      },
      charts: {
        health: null,
        contribution: null,
        burndown: null
      }
    }
  },
  created() {
    this.teamId = this.$route.params.teamId;
    if (!this.teamId) {
      this.teamId = sessionStorage.getItem("currentTeamId");
    }
    if (!this.teamId) {
      this.$router.push({ name: 'TeamTasks' });
      return;
    }
    this.fetchData();
  },
  mounted() {
    window.addEventListener('resize', this.handleResize);
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.handleResize);
    Object.values(this.charts).forEach(chart => chart && chart.dispose());
  },
  methods: {
    handleResize() {
      Object.values(this.charts).forEach(chart => chart && chart.resize());
    },
    async fetchData() {
      this.loading = true;
      try {
        const [tasksRes, membersRes] = await Promise.all([
          getTeamTasks({ team_id: this.teamId }),
          getTeamMembers(this.teamId)
        ]);

        this.tasks = tasksRes.data || [];
        this.members = membersRes.data || [];

        this.$nextTick(() => {
          this.initCharts();
        });
      } catch (error) {
        console.error("获取数据失败", error);
      } finally {
        this.loading = false;
      }
    },
    initCharts() {
      this.initHealthChart();
      this.initContributionChart();
      this.initBurndownChart();
    },
    initHealthChart() {
      const chartDom = this.$refs.healthChart;
      if (!chartDom) return;
      this.charts.health = echarts.init(chartDom);

      const now = new Date();
      let completed = 0;
      let inProgress = 0;
      let overdue = 0;
      let pending = 0;

      this.tasks.forEach(task => {
        if (task.status === 2) {
          completed++;
        } else {
          const dueDate = task.due_date ? new Date(task.due_date) : null;
          if (dueDate && dueDate < now) {
            overdue++;
          } else if (task.status === 1 || (task.progress > 0 && task.progress < 100)) {
            inProgress++;
          } else {
            pending++;
          }
        }
      });

      this.stats = { completed, inProgress, overdue, pending };

      const option = {
        tooltip: {
          trigger: 'item'
        },
        legend: {
          bottom: '0%',
          left: 'center'
        },
        series: [
          {
            name: '任务状态',
            type: 'pie',
            radius: ['40%', '70%'],
            avoidLabelOverlap: false,
            itemStyle: {
              borderRadius: 10,
              borderColor: '#fff',
              borderWidth: 2
            },
            label: {
              show: false,
              position: 'center'
            },
            emphasis: {
              label: {
                show: true,
                fontSize: 20,
                fontWeight: 'bold'
              }
            },
            labelLine: {
              show: false
            },
            data: [
              { value: completed, name: '已完成', itemStyle: { color: '#10B981' } },
              { value: inProgress, name: '进行中', itemStyle: { color: '#3B82F6' } },
              { value: overdue, name: '已逾期', itemStyle: { color: '#EF4444' } },
              { value: pending, name: '待处理', itemStyle: { color: '#9CA3AF' } }
            ]
          }
        ]
      };

      this.charts.health.setOption(option);
    },
    initContributionChart() {
      const chartDom = this.$refs.contributionChart;
      if (!chartDom) return;
      this.charts.contribution = echarts.init(chartDom);

      const memberMap = {};
      this.members.forEach(m => {
        memberMap[m.user_id] = m.nickname || m.username || `User ${m.user_id}`;
      });

      const contribution = {};
      this.tasks.forEach(task => {
        if (task.status === 2 && task.assignee_id) {
          const name = memberMap[task.assignee_id] || '未知成员';
          contribution[name] = (contribution[name] || 0) + 1;
        }
      });

      const data = Object.entries(contribution).map(([name, value]) => ({ name, value }));
      
      if (data.length === 0) {
        data.push({ name: '暂无完成任务', value: 0 });
      }

      const option = {
        tooltip: {
          trigger: 'item'
        },
        legend: {
          type: 'scroll',
          orient: 'vertical',
          right: 10,
          top: 20,
          bottom: 20,
        },
        series: [
          {
            name: '完成任务数',
            type: 'pie',
            radius: '70%',
            center: ['40%', '50%'],
            data: data,
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            }
          }
        ]
      };

      this.charts.contribution.setOption(option);
    },
    initBurndownChart() {
      const chartDom = this.$refs.burndownChart;
      if (!chartDom) return;
      this.charts.burndown = echarts.init(chartDom);

      // 简单的燃尽图逻辑：
      // 1. 找出最早的任务创建时间作为起点
      // 2. 找出最晚的截止时间作为终点（理想线）
      // 3. 实际线：根据任务完成时间（updated_at for status=2）递减

      if (this.tasks.length === 0) return;

      const allDates = [];
      let totalTasks = this.tasks.length;
      
      // 收集所有相关日期
      this.tasks.forEach(t => {
        if (t.created_at) allDates.push(new Date(t.created_at));
        if (t.due_date) allDates.push(new Date(t.due_date));
        if (t.status === 2 && t.updated_at) allDates.push(new Date(t.updated_at));
      });

      if (allDates.length === 0) return;

      const minDate = new Date(Math.min(...allDates));
      const maxDate = new Date(Math.max(...allDates));
      
      // 如果时间跨度太小，至少展示最近7天
      if (maxDate - minDate < 86400000 * 7) {
        maxDate.setDate(minDate.getDate() + 7);
      }

      // 生成日期序列
      const dateList = [];
      let curr = new Date(minDate);
      while (curr <= maxDate) {
        dateList.push(new Date(curr).toISOString().split('T')[0]);
        curr.setDate(curr.getDate() + 1);
      }
      // 确保包含今天
      const todayStr = new Date().toISOString().split('T')[0];
      if (!dateList.includes(todayStr) && new Date(todayStr) < maxDate) {
         // 如果今天在范围内但没被包含（通常不会发生），或者今天超过了maxDate需要延长？
         // 简单起见，我们只画到 maxDate
      }

      // 计算实际剩余曲线
      const actualData = [];
      let remaining = totalTasks;
      
      // 按完成时间排序已完成的任务
      const completedTasks = this.tasks
        .filter(t => t.status === 2 && t.updated_at)
        .map(t => ({ date: new Date(t.updated_at).toISOString().split('T')[0] }))
        .sort((a, b) => a.date.localeCompare(b.date));

      // 初始点
      // 实际上应该根据每一天，计算当天已创建但未完成的任务数。
      // 简化模型：总任务数 - 累计完成数
      
      let completedCount = 0;
      dateList.forEach(dateStr => {
        // 这一天完成的任务
        const todayCompleted = completedTasks.filter(t => t.date === dateStr).length;
        completedCount += todayCompleted;
        
        // 如果日期在今天之后，且没有数据，就不画了（或者是预测线？）
        // 这里只画到今天
        if (dateStr <= todayStr) {
          actualData.push(totalTasks - completedCount);
        }
      });

      // 理想线：从 totalTasks 到 0
      const idealData = [];
      const totalDays = dateList.length;
      const dropPerDay = totalTasks / (totalDays - 1);
      
      for (let i = 0; i < totalDays; i++) {
        idealData.push(Math.max(0, Math.round((totalTasks - (dropPerDay * i)) * 10) / 10));
      }

      const option = {
        tooltip: {
          trigger: 'axis'
        },
        legend: {
          data: ['理想剩余', '实际剩余']
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: dateList
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            name: '理想剩余',
            type: 'line',
            data: idealData,
            lineStyle: {
              type: 'dashed',
              color: '#9CA3AF'
            },
            itemStyle: {
              color: '#9CA3AF'
            },
            showSymbol: false
          },
          {
            name: '实际剩余',
            type: 'line',
            data: actualData,
            lineStyle: {
              color: '#3B82F6',
              width: 3
            },
            itemStyle: {
              color: '#3B82F6'
            },
            areaStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: 'rgba(59, 130, 246, 0.3)' },
                { offset: 1, color: 'rgba(59, 130, 246, 0.1)' }
              ])
            }
          }
        ]
      };

      this.charts.burndown.setOption(option);
    }
  }
}
</script>

<style scoped>
/* 确保图表容器有高度 */
</style>