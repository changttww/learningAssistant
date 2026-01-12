<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- 顶部导航 -->
    <header class="bg-white shadow-sm z-10 sticky top-0">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-3">
             <button @click="goBack" class="flex items-center gap-1 bg-gray-100 hover:bg-gray-200 text-gray-600 px-3 py-1.5 rounded-lg transition-colors">
               <iconify-icon icon="mdi:arrow-left" width="20"></iconify-icon>
               <span class="text-sm font-medium">返回团队任务</span>
             </button>
             <div class="h-4 w-px bg-gray-300"></div>
            <h1 class="text-xl font-bold text-gray-800 flex items-center gap-2">
              <iconify-icon icon="mdi:chart-box-outline" class="text-blue-600"></iconify-icon>
              团队数据报告
            </h1>
          </div>
        </div>
        <div class="flex items-center gap-4">
           <div class="text-sm text-gray-500 flex items-center gap-2">
             <span v-if="lastUpdated">更新于 {{ formatTime(lastUpdated) }}</span>
             <button 
               @click="fetchData" 
               class="p-1.5 rounded-full hover:bg-gray-100 text-blue-600 transition-colors"
               :class="{ 'animate-spin': loading }"
               title="刷新数据"
             >
               <iconify-icon icon="mdi:refresh" width="18"></iconify-icon>
             </button>
           </div>
           <div class="px-3 py-1 bg-blue-50 text-blue-700 rounded-full text-sm font-medium">
             {{ tasks.length }} 个任务
           </div>
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="flex-1 max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8">
      
      <div v-if="loading" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>

      <div v-else class="space-y-8">
        <!-- 概览卡片 -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
           <!-- 完成率 -->
           <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
              <div class="text-gray-500 text-sm mb-2">任务完成率</div>
              <div class="flex items-baseline gap-2">
                 <span class="text-3xl font-bold text-gray-800">{{ completionRate }}%</span>
                 <span class="text-xs text-gray-400">共 {{ stats.completed }} / {{ tasks.length }}</span>
              </div>
              <div class="w-full bg-gray-100 h-1.5 rounded-full mt-3 overflow-hidden">
                 <div class="h-full bg-green-500" :style="{ width: `${completionRate}%` }"></div>
              </div>
           </div>
           
           <!-- 总积分 -->
           <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
              <div class="text-gray-500 text-sm mb-2">累计积分产出</div>
              <div class="flex items-baseline gap-2">
                 <span class="text-3xl font-bold text-yellow-600">{{ stats.totalPoints }}</span>
                 <span class="text-xs text-gray-400">pts</span>
              </div>
              <div class="text-xs text-gray-400 mt-3">
                 平均每任务 {{ avgPoints }} 分
              </div>
           </div>

           <!-- 进行中 -->
           <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
              <div class="text-gray-500 text-sm mb-2">正在进行</div>
              <div class="flex items-baseline gap-2">
                 <span class="text-3xl font-bold text-blue-600">{{ stats.inProgress }}</span>
                 <span class="text-xs text-gray-400">任务</span>
              </div>
              <div class="text-xs text-blue-400 mt-3">
                 剩余工作量约 {{ remainingPoints }} 分
              </div>
           </div>

           <!-- 逾期 -->
           <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
              <div class="text-gray-500 text-sm mb-2">逾期任务</div>
              <div class="flex items-baseline gap-2">
                 <span class="text-3xl font-bold text-red-600">{{ stats.overdue }}</span>
                 <span class="text-xs text-gray-400">任务</span>
              </div>
              <div class="text-xs text-red-400 mt-3" v-if="stats.overdue > 0">
                 需立即关注
              </div>
              <div class="text-xs text-green-500 mt-3" v-else>
                 目前状态良好
              </div>
           </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          
          <!-- 任务状态分布 -->
          <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 flex flex-col">
            <h2 class="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
              <iconify-icon icon="mdi:chart-pie" class="text-blue-500"></iconify-icon>
              任务状态分布
            </h2>
            <div ref="healthChart" class="w-full flex-1 min-h-[300px]"></div>
          </div>

          <!-- 成员贡献排行榜 (按完成任务数) -->
          <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
            <h2 class="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
               <iconify-icon icon="mdi:trophy" class="text-yellow-500"></iconify-icon>
               成员贡献排行 (已完成任务)
            </h2>
            <div ref="contributionChart" class="w-full h-80"></div>
          </div>

          <!-- 燃尽图 -->
          <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 lg:col-span-2">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-bold text-gray-800 flex items-center gap-2">
                <iconify-icon icon="mdi:fire" class="text-red-500"></iconify-icon>
                项目燃尽图 (Burndown Chart)
              </h2>
              <div class="flex gap-2">
                 <button 
                   @click="burndownMode = 'count'; initBurndownChart()"
                   :class="['px-3 py-1 text-xs rounded-lg transition-colors', burndownMode === 'count' ? 'bg-blue-100 text-blue-700 font-medium' : 'bg-gray-100 text-gray-600']"
                 >按任务数</button>
                 <button 
                   @click="burndownMode = 'points'; initBurndownChart()"
                   :class="['px-3 py-1 text-xs rounded-lg transition-colors', burndownMode === 'points' ? 'bg-blue-100 text-blue-700 font-medium' : 'bg-gray-100 text-gray-600']"
                 >按积分</button>
              </div>
            </div>
            <div ref="burndownChart" class="w-full h-96"></div>
          </div>

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
        pending: 0,
        totalPoints: 0
      },
      pieStats: {
        completed: 0,
        inProgress: 0,
        overdue: 0,
        pending: 0
      },
      lastUpdated: null,
      refreshTimer: null,
      charts: {
        health: null,
        contribution: null,
        burndown: null
      },
      burndownMode: 'count' // 'count' or 'points'
    }
  },
  computed: {
    completionRate() {
      if (this.tasks.length === 0) return 0;
      return Math.round((this.stats.completed / this.tasks.length) * 100);
    },
    avgPoints() {
      if (this.tasks.length === 0) return 0;
      return (this.stats.totalPoints / this.tasks.length).toFixed(1);
    },
    remainingPoints() {
       return this.tasks.reduce((sum, t) => {
         if (t.status !== 2) return sum + (t.effort_points || 0);
         return sum;
       }, 0);
    }
  },
  created() {
    this.teamId = this.$route.params.teamId;
    if (!this.teamId) {
      // 尝试从query或者sessionStorage获取
      this.teamId = this.$route.query.teamId || sessionStorage.getItem("currentTeamId");
    }
    
    if (!this.teamId) {
      this.$router.replace({ name: 'TeamTasks' });
      return;
    }
    this.fetchData();
  },
  mounted() {
    window.addEventListener('resize', this.handleResize);
    // 启动自动刷新 (每15秒)
    this.refreshTimer = setInterval(() => {
      this.fetchData(true); // true 表示静默刷新，不显示 loading
    }, 15000);
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.handleResize);
    if (this.refreshTimer) clearInterval(this.refreshTimer);
    Object.values(this.charts).forEach(chart => chart && chart.dispose());
  },
  methods: {
    goBack() {
      if (this.teamId) {
        this.$router.push({ name: 'TeamTasks', query: { teamId: this.teamId } });
      } else {
        this.$router.push({ name: 'TeamTasks' });
      }
    },
    formatTime(date) {
      return new Date(date).toLocaleTimeString();
    },
    handleResize() {
      Object.values(this.charts).forEach(chart => chart && chart.resize());
    },
    async fetchData(silent = false) {
      if (!silent) this.loading = true;
      try {
        const [tasksRes, membersRes] = await Promise.all([
          getTeamTasks({ team_id: this.teamId }),
          getTeamMembers(this.teamId)
        ]);

        const rawTasks = tasksRes.data?.items || tasksRes.data || [];
        this.tasks = Array.isArray(rawTasks) ? rawTasks : [];
        this.members = membersRes.data || [];
        
        this.lastUpdated = new Date();
        this.calculateStats();

        this.$nextTick(() => {
          this.initCharts();
        });
      } catch (error) {
        console.error("获取数据失败", error);
        // alert("获取报表数据失败，请重试");
      } finally {
        this.loading = false;
      }
    },
    calculateStats() {
      const now = new Date();
      let completed = 0;
      let inProgress = 0;
      let overdue = 0;
      let pending = 0;
      let totalPoints = 0;

      // 用于图表的互斥统计
      let pieOverdue = 0;
      let pieInProgress = 0;
      let piePending = 0;

      this.tasks.forEach(task => {
        // 积分统计 (假定 status=2 是完成)
        if (task.status === 2) {
          totalPoints += (task.effort_points || 0);
          completed++;
        } else {
          // 检查逾期
          const dueDate = (task.due_at || task.due_date);
          const isOverdue = dueDate && new Date(dueDate) < now;
          
          if (isOverdue) {
            overdue++; // 卡片统计：所有逾期
            pieOverdue++; // 图表统计：优先归为逾期
          }
          
          if (task.status === 1 || (task.progress > 0 && task.progress < 100)) {
            inProgress++; // 卡片统计：所有进行中
            if (!isOverdue) pieInProgress++; // 图表统计：未逾期的进行中
          } else {
            pending++; // 卡片统计：所有待处理
            if (!isOverdue) piePending++; // 图表统计：未逾期的待处理
          }
        }
      });

      this.stats = { completed, inProgress, overdue, pending, totalPoints };
      this.pieStats = { completed, inProgress: pieInProgress, overdue: pieOverdue, pending: piePending };
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

      const option = {
        tooltip: { trigger: 'item' },
        legend: { bottom: '0%', left: 'center' },
        color: ['#10B981', '#3B82F6', '#EF4444', '#9CA3AF'],
        series: [{
          name: '任务状态',
          type: 'pie',
          radius: ['45%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
          label: { show: false, position: 'center' },
          emphasis: {
            label: { show: true, fontSize: 18, fontWeight: 'bold' }
          },
          data: [
            { value: this.pieStats.completed, name: '已完成' },
            { value: this.pieStats.inProgress, name: '进行中' },
            { value: this.pieStats.overdue, name: '已逾期' },
            { value: this.pieStats.pending, name: '待处理' }
          ]
        }]
      };
      this.charts.health.setOption(option);
    },
    getMemberName(userId) {
       const m = this.members.find(m => String(m.user_id) === String(userId));
       return m ? (m.nickname || m.username) : `用户 ${userId}`;
    },
    initContributionChart() {
      const chartDom = this.$refs.contributionChart;
      if (!chartDom) return;
      this.charts.contribution = echarts.init(chartDom);

      const counts = {};
      this.tasks.forEach(t => {
        if (t.status === 2 && t.owner_user_id) {
          const name = this.getMemberName(t.owner_user_id);
          counts[name] = (counts[name] || 0) + 1;
        }
      });

      const data = Object.entries(counts)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 10); // Top 10

      const option = {
        tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
        grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
        xAxis: { type: 'value', boundaryGap: [0, 0.01] },
        yAxis: { type: 'category', data: data.map(d => d[0]) },
        series: [{
          name: '完成任务数',
          type: 'bar',
          data: data.map(d => d[1]),
          itemStyle: { color: '#8B5CF6', borderRadius: [0, 4, 4, 0] }
        }]
      };
      this.charts.contribution.setOption(option);
    },
    initBurndownChart() {
      const chartDom = this.$refs.burndownChart;
      if (!chartDom) return;
      this.charts.burndown = echarts.init(chartDom);
      
      if (this.tasks.length === 0) {
        this.charts.burndown.clear();
        return;
      }

      // Collect all dates
      const dates = new Set();
      const taskEvents = []; // { date, change: -value } for completion, { date, change: +value } for creation

      const totalValue = this.burndownMode === 'count' 
          ? this.tasks.length 
          : this.tasks.reduce((sum, t) => sum + (t.effort_points || 0), 0);

      const todayStr = new Date().toISOString().split('T')[0];
      
      // Determine date range
      let minDateStr = todayStr;
      let maxDateStr = todayStr;

      this.tasks.forEach(t => {
         const created = t.created_at ? t.created_at.split('T')[0] : todayStr;
         const due = (t.due_at || t.due_date) ? (t.due_at || t.due_date).split('T')[0] : null;

         if (created < minDateStr) minDateStr = created;
         if (due && due > maxDateStr) maxDateStr = due;
      });

      // Generate date sequence
      const dateList = [];
      let curr = new Date(minDateStr);
      const end = new Date(maxDateStr);
      // Limit to 60 days to avoid performance issues
      if ((end - curr) > 60 * 86400000) {
         end.setTime(curr.getTime() + 60 * 86400000);
      }
      
      while (curr <= end) {
        dateList.push(curr.toISOString().split('T')[0]);
        curr.setDate(curr.getDate() + 1);
      }

      // Calculate Remaining
      const actualData = [];
      let currentVal = totalValue;

      // Simplification: Assume all tasks existed at start (Total) 
      // and we subtract when they are completed. 
      // Real burndown adds tasks when created, but usually implies sprint backlog is fixed.
      // We will stick to "Sprint" style: Total at start -> burn down.
      
      const completionMap = {};
      this.tasks.forEach(t => {
         if (t.status === 2) {
            // Prefer completed_at, fallback to update_at
            const doneTime = t.completed_at || t.updated_at;
            if (doneTime) {
               const d = doneTime.split('T')[0];
               const val = this.burndownMode === 'count' ? 1 : (t.effort_points || 0);
               completionMap[d] = (completionMap[d] || 0) + val;
            }
         }
      });

      let burned = 0;
      dateList.forEach(d => {
         if (d > todayStr) return; // Future
         if (completionMap[d]) {
            burned += completionMap[d];
         }
         actualData.push(totalValue - burned);
      });

      // Ideal Line
      const idealData = [];
      const steps = dateList.length;
      const drop = totalValue / (steps - 1);
      for(let i=0; i<steps; i++) {
        idealData.push(Math.max(0, totalValue - drop * i));
      }

      const option = {
        tooltip: { trigger: 'axis' },
        legend: { data: ['理想剩余', '实际剩余'] },
        grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
        xAxis: { type: 'category', boundaryGap: false, data: dateList },
        yAxis: { type: 'value', name: this.burndownMode === 'count' ? '任务数' : '积分' },
        series: [
          {
            name: '理想剩余',
            type: 'line',
            data: idealData,
            itemStyle: { color: '#9CA3AF' },
            lineStyle: { type: 'dashed' },
            showSymbol: false,
            smooth: true
          },
          {
            name: '实际剩余',
            type: 'line',
            data: actualData,
            itemStyle: { color: '#3B82F6' },
            areaStyle: { 
               opacity: 0.1, 
               color: new echarts.graphic.LinearGradient(0,0,0,1, [{offset:0, color:'#3B82F6'}, {offset:1, color:'#fff'}]) 
            },
            smooth: true
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