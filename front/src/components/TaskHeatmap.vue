<template>
  <div class="task-heatmap">
    <!-- 标题 -->
    <h3 class="text-sm font-semibold text-gray-900 mb-2">任务活跃度</h3>
    
    <!-- GitHub 风格热力图容器 -->
    <div class="flex gap-1.5 items-stretch w-full pb-2">
      <span class="text-xs text-gray-500 whitespace-nowrap pt-3.5 flex-shrink-0">少</span>
      
      <!-- 热力图网格容器 - 填满剩余空间 -->
      <div class="flex-1 flex gap-0.5 min-w-0">
        <!-- 左侧周一到周日标签 -->
        <div class="flex flex-col gap-0.5 flex-shrink-0">
          <div class="h-4"></div>
          <div v-for="day in 7" :key="`label-${day}`" class="text-xs text-gray-500 text-center w-3 leading-none">
            {{ ['一', '二', '三', '四', '五', '六', '日'][day - 1] }}
          </div>
        </div>
        
        <!-- 按周组织的热力图 - 自动缩放到可用宽度 -->
        <div class="flex-1 flex gap-0.5 min-w-0">
          <div v-for="(week, weekIndex) in groupedWeeks" :key="`week-${weekIndex}`" class="flex flex-col gap-0.5 flex-shrink-0">
            <!-- 月份标签（每隔几周显示一次） -->
            <div v-if="weekIndex % 4 === 0" class="text-xs text-gray-500 h-4 flex items-end">
              {{ getMonthLabel(week) }}
            </div>
            <div v-else class="h-4"></div>
            
            <!-- 该周的7个日期 -->
            <div
              v-for="day in week"
              :key="day.date"
              :class="getHeatmapCellClass(day.level)"
              @mouseenter="hoveredDay = day"
              @mouseleave="hoveredDay = null"
            >
              <div
                v-if="hoveredDay && hoveredDay.date === day.date"
                class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-1.5 px-2 py-1 bg-gray-900 text-white text-xs rounded shadow-lg z-20 pointer-events-none whitespace-nowrap"
              >
                <div class="font-medium text-xs">{{ formatDate(day.date) }}</div>
                <div class="text-gray-300 text-xs">{{ day.count }} 任务 · {{ day.completed }} 完成</div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <span class="text-xs text-gray-500 whitespace-nowrap pt-3.5 flex-shrink-0">多</span>
      
      <!-- 右侧统计信息卡片 -->
      <div class="flex-shrink-0 flex gap-3 ml-2 pl-2 border-l border-gray-200">
        <div class="flex flex-col items-center">
          <div class="text-lg font-bold text-green-600">{{ currentStreak }}</div>
          <div class="text-xs text-gray-500 whitespace-nowrap">连续天数</div>
        </div>
        <div class="flex flex-col items-center">
          <div class="text-lg font-bold text-blue-600">{{ totalTasks }}</div>
          <div class="text-xs text-gray-500 whitespace-nowrap">总任务数</div>
        </div>
        <div class="flex flex-col items-center">
          <div class="text-lg font-bold text-amber-600">{{ completionRate }}</div>
          <div class="text-xs text-gray-500 whitespace-nowrap">完成率</div>
        </div>
      </div>
    </div>
    
    <!-- 图例 -->
    <div class="flex items-center justify-center gap-2 mt-2 text-xs text-gray-600">
      <span class="text-xs">活跃度:</span>
      <div class="flex gap-0.5">
        <div class="heatmap-cell heatmap-level-0 border border-gray-300"></div>
        <div class="heatmap-cell heatmap-level-1"></div>
        <div class="heatmap-cell heatmap-level-2"></div>
        <div class="heatmap-cell heatmap-level-3"></div>
        <div class="heatmap-cell heatmap-level-4"></div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from "vue";
import { getTaskHeatmapStats } from "@/api/modules/task";

export default {
  name: "TaskHeatmap",
  setup() {
    const heatmapDays = ref([]);
    const totalTasks = ref(0);
    const completedNum = ref(0);
    const currentStreak = ref(0);
    const hoveredDay = ref(null);
    const loading = ref(false);
    let refreshInterval = null;
    let handleFocus = null;

    // 将 365 天数据按周分组 (52-53周 x 7天)
    const groupedWeeks = computed(() => {
      const weeks = [];
      for (let i = 0; i < heatmapDays.value.length; i += 7) {
        weeks.push(heatmapDays.value.slice(i, i + 7));
      }
      return weeks;
    });

    // 计算完成率百分比
    const completionRate = computed(() => {
      if (totalTasks.value === 0) return "0%";
      const rate = Math.round((completedNum.value / totalTasks.value) * 100);
      return `${rate}%`;
    });

    const getMonthLabel = (week) => {
      if (week.length === 0) return "";
      const date = new Date(week[0].date + "T00:00:00");
      return ["1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"][date.getMonth()];
    };

    const fetchHeatmapStats = async () => {
      loading.value = true;
      try {
        const res = await getTaskHeatmapStats();
        if (res?.data) {
          const data = res.data;
          heatmapDays.value = data.days || [];
          totalTasks.value = data.total_tasks || 0;
          completedNum.value = data.completed_num || 0;
          currentStreak.value = data.current_streak || 0;
          
          console.log("[热力图] 数据更新成功", {
            days: heatmapDays.value.length,
            totalTasks: totalTasks.value,
            completedNum: completedNum.value,
            currentStreak: currentStreak.value,
          });
        }
      } catch (error) {
        console.error("[热力图] 获取数据失败:", error);
      } finally {
        loading.value = false;
      }
    };

    const getHeatmapCellClass = (level) => {
      const baseClass = "heatmap-cell transition-all duration-200 cursor-pointer hover:ring-2 hover:ring-offset-0.5 hover:ring-blue-400 relative";
      switch (level) {
        case 0:
          return `${baseClass} heatmap-level-0`;
        case 1:
          return `${baseClass} heatmap-level-1`;
        case 2:
          return `${baseClass} heatmap-level-2`;
        case 3:
          return `${baseClass} heatmap-level-3`;
        case 4:
          return `${baseClass} heatmap-level-4`;
        default:
          return `${baseClass} heatmap-level-0`;
      }
    };

    const formatDate = (dateStr) => {
      const date = new Date(dateStr + "T00:00:00");
      const month = date.getMonth() + 1;
      const day = date.getDate();
      const weekDays = ["日", "一", "二", "三", "四", "五", "六"];
      const week = weekDays[date.getDay()];
      return `${month}月${day}日 (周${week})`;
    };

    // 初始加载
    onMounted(async () => {
      console.log("[热力图] 组件已挂载，开始加载数据...");
      await fetchHeatmapStats();

      // 30秒自动刷新一次
      refreshInterval = setInterval(() => {
        console.log("[热力图] 自动刷新数据");
        fetchHeatmapStats();
      }, 30000);

      // 窗口获得焦点时刷新
      handleFocus = () => {
        console.log("[热力图] 窗口获得焦点，刷新数据");
        fetchHeatmapStats();
      };

      window.addEventListener("focus", handleFocus);
    });

    // 组件卸载时清理
    onUnmounted(() => {
      console.log("[热力图] 组件已卸载，清理资源");
      if (refreshInterval) {
        clearInterval(refreshInterval);
      }
      if (handleFocus) {
        window.removeEventListener("focus", handleFocus);
      }
    });

    return {
      heatmapDays,
      totalTasks,
      completedNum,
      currentStreak,
      hoveredDay,
      loading,
      groupedWeeks,
      getMonthLabel,
      getHeatmapCellClass,
      formatDate,
      completionRate,
    };
  },
};
</script>

<style scoped>
/* 热力图容器 */
.task-heatmap {
  transition: all 0.3s ease;
}

/* 热力图颜色级别 - GitHub 风格 */
.heatmap-level-0 {
  background-color: #ebedf0;
  border: 1px solid #d1d5da;
}

.heatmap-level-1 {
  background: linear-gradient(135deg, #c6e48b 0%, #b8d97f 100%);
  box-shadow: 0 1px 3px rgba(27, 154, 35, 0.1);
}

.heatmap-level-2 {
  background: linear-gradient(135deg, #7bc96f 0%, #6ba456 100%);
  box-shadow: 0 1px 3px rgba(27, 154, 35, 0.2);
}

.heatmap-level-3 {
  background: linear-gradient(135deg, #239a3b 0%, #1b7d28 100%);
  box-shadow: 0 1px 3px rgba(27, 154, 35, 0.3);
}

.heatmap-level-4 {
  background: linear-gradient(135deg, #0e4429 0%, #051d0f 100%);
  box-shadow: 0 1px 3px rgba(27, 154, 35, 0.4);
}

/* Hover 效果 */
.heatmap-level-1:hover,
.heatmap-level-2:hover,
.heatmap-level-3:hover,
.heatmap-level-4:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
}

/* 热力图单元格 - 响应式缩放 */
.heatmap-cell {
  flex-shrink: 0;
  border-radius: 0.25rem;
  position: relative;
  aspect-ratio: 1;
  min-height: 0.625rem;
  min-width: 0.625rem;
}

/* 响应式宽度 - 根据容器自动调整单元格大小 */
@media (min-width: 1920px) {
  .heatmap-cell {
    width: 0.875rem;
    height: 0.875rem;
  }
}

@media (min-width: 1440px) and (max-width: 1919px) {
  .heatmap-cell {
    width: 0.75rem;
    height: 0.75rem;
  }
}

@media (max-width: 1439px) {
  .heatmap-cell {
    width: 0.625rem;
    height: 0.625rem;
  }
}
</style>
