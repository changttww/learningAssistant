<template>
  <div class="w-full min-h-full flex flex-col px-4 py-6">
    <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center gap-4">
          <button
            @click="$router.back()"
            class="flex items-center px-3 py-2 text-sm font-medium text-gray-600 bg-gray-50 hover:bg-blue-50 hover:text-blue-600 rounded-lg transition-all duration-200"
          >
            <iconify-icon icon="mdi:arrow-left" class="mr-1"></iconify-icon>
            返回
          </button>
          <h1 class="text-2xl font-bold text-gray-800">团队日历</h1>
        </div>
        <div class="flex items-center gap-2">
          <div class="flex items-center gap-2 text-sm text-gray-500">
            <span class="w-3 h-3 rounded-full bg-blue-500"></span> 进行中
            <span class="w-3 h-3 rounded-full bg-green-500"></span> 已完成
            <span class="w-3 h-3 rounded-full bg-red-500"></span> 逾期
            <span class="w-3 h-3 rounded-full bg-yellow-500"></span> 即将到期
          </div>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>

      <el-calendar v-else v-model="currentDate">
        <template #date-cell="{ data }">
          <div class="h-full w-full flex flex-col hover:bg-gray-50 transition-colors p-1 overflow-hidden">
            <div :class="['text-sm font-medium mb-1', data.isSelected ? 'text-blue-600' : 'text-gray-700']">
              {{ data.day.split('-').slice(2).join('') }}
            </div>
            <div class="flex-1 overflow-y-auto space-y-1 custom-scrollbar">
              <div
                v-for="task in getTasksForDate(data.day)"
                :key="task.id"
                class="text-xs px-1.5 py-0.5 rounded truncate cursor-pointer transition-opacity hover:opacity-80"
                :class="getTaskClass(task)"
                :title="`${task.title} (${task.owner_name || '未分配'})`"
                @click.stop="handleTaskClick(task)"
              >
                {{ task.title }}
              </div>
            </div>
          </div>
        </template>
      </el-calendar>
    </div>

    <!-- 任务详情弹窗 -->
    <div
      v-if="selectedTask"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
      @click.self="selectedTask = null"
    >
      <div class="bg-white rounded-xl shadow-2xl w-full max-w-md p-6 transform transition-all scale-100">
        <div class="flex justify-between items-start mb-4">
          <h3 class="text-lg font-bold text-gray-800 pr-8">{{ selectedTask.title }}</h3>
          <button @click="selectedTask = null" class="text-gray-400 hover:text-gray-600">
            <iconify-icon icon="mdi:close" width="24" height="24"></iconify-icon>
          </button>
        </div>
        
        <div class="space-y-4">
          <div class="bg-gray-50 p-3 rounded-lg text-sm text-gray-600">
            {{ selectedTask.description || "暂无描述" }}
          </div>
          
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span class="text-gray-500 block mb-1">负责人</span>
              <div class="flex items-center gap-2">
                <iconify-icon icon="mdi:account" class="text-blue-500"></iconify-icon>
                <span class="font-medium">{{ selectedTask.owner_name || "未分配" }}</span>
              </div>
            </div>
            <div>
              <span class="text-gray-500 block mb-1">截止日期</span>
              <div class="flex items-center gap-2">
                <iconify-icon icon="mdi:calendar-clock" class="text-orange-500"></iconify-icon>
                <span class="font-medium">{{ formatDate(selectedTask.due_date) }}</span>
              </div>
            </div>
            <div>
              <span class="text-gray-500 block mb-1">状态</span>
              <span :class="['px-2 py-0.5 rounded text-xs font-medium', getStatusBadgeClass(selectedTask)]">
                {{ getStatusLabel(selectedTask) }}
              </span>
            </div>
            <div>
              <span class="text-gray-500 block mb-1">进度</span>
              <div class="flex items-center gap-2">
                <div class="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
                  <div class="h-full bg-blue-500" :style="{ width: `${selectedTask.progress || 0}%` }"></div>
                </div>
                <span class="text-xs">{{ selectedTask.progress || 0 }}%</span>
              </div>
            </div>
          </div>
        </div>

        <div class="mt-6 flex justify-end">
          <button
            @click="selectedTask = null"
            class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 font-medium transition-colors"
          >
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getTeamTasks } from "@/api/modules/task";
import { getTeamMembers } from "@/api/modules/team";

export default {
  name: "TeamCalendar",
  data() {
    return {
      currentDate: new Date(),
      loading: false,
      tasks: [],
      teamId: null,
      selectedTask: null,
      teamMembers: [],
    };
  },
  created() {
    this.teamId = this.$route.params.teamId || this.$route.query.teamId;
    if (this.teamId) {
      this.loadData();
    } else {
      // 尝试从 sessionStorage 获取
      const storedTeamId = sessionStorage.getItem("currentTeamId");
      if (storedTeamId) {
        this.teamId = storedTeamId;
        this.loadData();
      } else {
        alert("未找到团队信息，请返回团队任务页面重新进入");
        this.$router.push({ name: "TeamTasks" });
      }
    }
  },
  methods: {
    async loadData() {
      this.loading = true;
      try {
        await Promise.all([this.loadTasks(), this.loadMembers()]);
      } catch (error) {
        console.error("加载数据失败", error);
      } finally {
        this.loading = false;
      }
    },
    async loadMembers() {
      try {
        const res = await getTeamMembers(this.teamId);
        this.teamMembers = res.data || [];
      } catch (error) {
        console.error("获取成员失败", error);
      }
    },
    async loadTasks() {
      try {
        const res = await getTeamTasks({ team_id: this.teamId });
        const items = res?.data?.items || res?.data || res;
        if (Array.isArray(items)) {
          this.tasks = items.map(task => {
            // 补充 owner_name
            let ownerName = task.owner_name;
            if (!ownerName && task.owner_user_id) {
              const member = this.teamMembers.find(m => String(m.user_id) === String(task.owner_user_id));
              if (member) ownerName = member.nickname;
            }
            
            // 规范化 due_date
            const rawDue = task.due_at || task.due_date;
            const normalizedDueDate = rawDue ? (rawDue.includes('T') ? rawDue.split('T')[0] : rawDue) : null;

            return {
              ...task,
              owner_name: ownerName,
              due_date: normalizedDueDate // 确保使用统一的 YYYY-MM-DD 格式
            };
          });
        }
      } catch (error) {
        console.error("获取任务失败", error);
      }
    },
    getTasksForDate(dateStr) {
      // dateStr format: YYYY-MM-DD
      return this.tasks.filter(task => {
        return task.due_date === dateStr;
      });
    },
    getTaskClass(task) {
      const status = this.normalizeStatus(task.status);
      const isOverdue = this.isTaskOverdue(task);
      
      if (status === 'completed') return 'bg-green-100 text-green-700 border border-green-200';
      if (isOverdue) return 'bg-red-100 text-red-700 border border-red-200';
      if (this.isTaskDueSoon(task)) return 'bg-yellow-100 text-yellow-700 border border-yellow-200';
      return 'bg-blue-100 text-blue-700 border border-blue-200';
    },
    getStatusBadgeClass(task) {
      const status = this.normalizeStatus(task.status);
      if (status === 'completed') return 'bg-green-100 text-green-700';
      if (status === 'in-progress') return 'bg-blue-100 text-blue-700';
      return 'bg-gray-100 text-gray-600';
    },
    getStatusLabel(task) {
      const status = this.normalizeStatus(task.status);
      const map = {
        'completed': '已完成',
        'in-progress': '进行中',
        'pending': '待处理'
      };
      return map[status] || '未知';
    },
    normalizeStatus(status) {
      if (status === null || status === undefined) return "pending";
      const numeric = Number(status);
      if (!Number.isNaN(numeric)) {
        if (numeric === 2) return "completed";
        if (numeric === 1) return "in-progress";
        return "pending";
      }
      const lowered = String(status).toLowerCase();
      if (["completed", "done", "finish"].includes(lowered)) return "completed";
      if (["in-progress", "progress", "doing"].includes(lowered)) return "in-progress";
      return lowered;
    },
    isTaskOverdue(task) {
      if (this.normalizeStatus(task.status) === 'completed') return false;
      if (!task.due_date) return false;
      const due = new Date(task.due_date);
      const now = new Date();
      now.setHours(0, 0, 0, 0);
      return due < now;
    },
    isTaskDueSoon(task) {
      if (this.normalizeStatus(task.status) === 'completed') return false;
      if (!task.due_date) return false;
      const due = new Date(task.due_date);
      const now = new Date();
      now.setHours(0, 0, 0, 0);
      const diffTime = due - now;
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24)); 
      return diffDays >= 0 && diffDays <= 3;
    },
    handleTaskClick(task) {
      this.selectedTask = task;
    },
    formatDate(dateStr) {
      if (!dateStr) return "无";
      return new Date(dateStr).toLocaleDateString('zh-CN');
    }
  }
};
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #e5e7eb;
  border-radius: 2px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background-color: #d1d5db;
}
/* 覆盖 element-plus 日历样式以适应自定义内容 */
:deep(.el-calendar-table .el-calendar-day) {
  height: 120px;
  padding: 4px;
}
:deep(.el-calendar__header) {
  padding: 12px 0;
  border-bottom: 1px solid #f3f4f6;
}
</style>
