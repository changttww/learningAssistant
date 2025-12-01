<template>
  <section>
    <div class="mb-4">
      <div class="flex space-x-1 bg-gray-100 p-1 rounded-lg w-fit">
        <button
          v-for="tab in tabMeta"
          :key="tab.key"
          class="px-4 py-2 rounded-md font-medium transition-colors"
          :class="tabClass(tab.key, tab.colorClass)"
          @click="updateTab(tab.key)"
        >
          {{ tab.label }} ({{ tabCount(tab.key) }})
        </button>
      </div>
    </div>

    <div class="space-y-4">
      <div v-if="tasksLoading" class="card p-6 text-center text-gray-500">
        正在加载今日任务...
      </div>
      <div
        v-else-if="tasksError"
        class="card p-6 text-center text-red-500 space-y-3"
      >
        <div>{{ tasksError }}</div>
        <button
          class="px-4 py-2 rounded-lg bg-red-500 text-white hover:bg-red-600"
          @click="loadTodayTasks"
        >
          重试加载
        </button>
      </div>
      <template v-else>
        <div
          v-for="tab in tabMeta"
          :key="tab.key"
          v-show="activeTab === tab.key"
          class="space-y-4"
        >
          <template v-if="tabTasks(tab.key).length">
            <article
              v-for="task in tabTasks(tab.key)"
              :key="task.id"
              class="card p-5 flex flex-col lg:flex-row lg:items-center hover:shadow-md transition-shadow"
            >
              <div
                class="w-12 h-12 rounded-xl flex items-center justify-center mr-4 mb-4 lg:mb-0"
                :class="iconWrapperClass(tab.key)"
              >
                <iconify-icon
                  :icon="taskIcon(tab.key)"
                  :class="[iconColorClass(tab.key), 'text-2xl']"
                />
              </div>
              <div class="flex-1 mr-0 lg:mr-4">
                <h4 class="font-bold text-gray-800">{{ task.title }}</h4>
                <p class="text-xs text-gray-500 mt-1">
                  {{ task.description }}
                </p>
                <div class="flex items-center mt-1 text-xs text-gray-500">
                  <span>{{ dueDateLabel(task) }}</span>
                  <span
                    class="ml-2 px-2 py-1 rounded-full font-medium"
                    :class="priorityBadgeClass(task.priority)"
                  >
                    {{ priorityLabel(task.priority) }}
                  </span>
                </div>
                <div v-if="task.tags?.length" class="flex flex-wrap gap-2 mt-2">
                  <span
                    v-for="tag in task.tags"
                    :key="`${task.id}-${tag}`"
                    class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-full"
                  >
                    {{ tag }}
                  </span>
                </div>
              </div>
              <div class="w-full lg:w-1/3 lg:ml-auto">
                <div class="flex justify-between text-sm text-gray-500 mb-1">
                  <span>{{ tab.progressLabel }}</span>
                  <span>{{ formatProgress(task.progress) }}%</span>
                </div>
                <div class="h-2 bg-gray-200 rounded-full">
                  <div
                    class="h-full rounded-full"
                    :class="progressBarClass(tab.key)"
                    :style="{ width: `${formatProgress(task.progress)}%` }"
                  ></div>
                </div>
              </div>
              <button
                class="mt-4 lg:mt-0 lg:ml-6 font-medium py-2 px-4 rounded-lg text-white"
                :class="actionButtonClass(tab.key)"
              >
                {{ actionLabel(tab.key) }}
              </button>
            </article>
          </template>
          <div v-else class="card p-6 text-center text-gray-500">
            暂无{{ tab.label }}任务
          </div>
        </div>
      </template>
    </div>
  </section>
</template>

<script>
  import { computed } from "vue";
  import { ElMessage } from "element-plus";
  import { getTodayTasks } from "@/api/modules/task";
  import {
    useCurrentUser,
    DEFAULT_USER_ID,
  } from "@/composables/useCurrentUser";

  const clampPercentage = (value) => {
    const numeric = Number(value);
    if (!Number.isFinite(numeric)) return 0;
    if (numeric < 0) return 0;
    if (numeric > 100) return 100;
    return Math.round(numeric);
  };

  const formatDateString = (value) => {
    if (!value) return "";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return typeof value === "string" ? value : "";
    }
    return date.toISOString().slice(0, 10);
  };

  export default {
    name: "TaskTabsSection",
    props: {
      userId: {
        type: [Number, String],
        default: null,
      },
    },
    setup(props) {
      const { profile } = useCurrentUser();

      return {
        resolvedUserId: computed(() => {
          const parsed = Number(props.userId);
          if (Number.isFinite(parsed) && parsed > 0) {
            return parsed;
          }
          return profile.value?.id ?? DEFAULT_USER_ID;
        }),
      };
    },
    data() {
      return {
        activeTab: "inProgress",
        tasks: {
          inProgress: [],
          pending: [],
          completed: [],
        },
        tasksLoading: false,
        tasksError: "",
      };
    },
    computed: {
      normalizedTasks() {
        return {
          inProgress: this.tasks?.inProgress ?? [],
          pending: this.tasks?.pending ?? this.tasks?.toStart ?? [],
          completed: this.tasks?.completed ?? [],
        };
      },
      tabMeta() {
        return [
          {
            key: "inProgress",
            label: "进行中",
            colorClass: "text-blue-600",
            progressLabel: "学习进度",
          },
          {
            key: "pending",
            label: "待开始",
            colorClass: "text-orange-600",
            progressLabel: "预计进度",
          },
          {
            key: "completed",
            label: "已完成",
            colorClass: "text-green-600",
            progressLabel: "完成进度",
          },
        ];
      },
    },
    watch: {
      resolvedUserId() {
        this.loadTodayTasks();
      },
    },
    mounted() {
      this.loadTodayTasks();
    },
    methods: {
      async loadTodayTasks() {
        this.tasksLoading = true;
        this.tasksError = "";
        try {
          const res = await getTodayTasks(this.resolvedUserId);
          const data = res?.data || {};
          const mapTasks = (list) =>
            (Array.isArray(list) ? list : [])
              .map((item) => this.normalizeTaskFromApi(item))
              .filter(Boolean);

          this.tasks = {
            inProgress: mapTasks(data.in_progress),
            pending: mapTasks(data.not_started),
            completed: mapTasks(data.completed),
          };
          this.ensureActiveTabExists();
        } catch (error) {
          this.tasksError = error?.message || "加载今日任务失败";
          console.error("加载今日任务失败", error);
          ElMessage.error(this.tasksError);
        } finally {
          this.tasksLoading = false;
        }
      },
      ensureActiveTabExists() {
        const order = ["inProgress", "pending", "completed"];
        if (this.tabTasks(this.activeTab).length > 0) {
          return;
        }
        const nextTab = order.find((key) => this.tabTasks(key).length > 0);
        this.activeTab = nextTab || "inProgress";
      },
      updateTab(tab) {
        this.activeTab = tab;
      },
      tabClass(tab, activeColorClass) {
        return this.activeTab === tab
          ? `bg-white shadow-sm ${activeColorClass}`
          : "text-gray-600 hover:text-gray-800";
      },
      tabTasks(key) {
        return this.normalizedTasks[key] ?? [];
      },
      tabCount(key) {
        return this.tabTasks(key).length;
      },
      priorityLabel(priority) {
        const map = {
          high: "高优先级",
          medium: "中优先级",
          low: "低优先级",
        };
        return map[priority] || "一般";
      },
      priorityBadgeClass(priority) {
        const map = {
          high: "bg-red-100 text-red-700",
          medium: "bg-orange-100 text-orange-700",
          low: "bg-green-100 text-green-700",
        };
        return map[priority] || "bg-gray-100 text-gray-600";
      },
      dueDateLabel(task) {
        if (!task?.dueDate) {
          return "无截止日期";
        }
        return `截止：${task.dueDate}`;
      },
      formatProgress(progress) {
        return clampPercentage(progress);
      },
      progressBarClass(tabKey) {
        const map = {
          inProgress: "bg-blue-600",
          pending: "bg-orange-500",
          completed: "bg-green-500",
        };
        return map[tabKey] || "bg-gray-400";
      },
      iconWrapperClass(tabKey) {
        const map = {
          inProgress: "bg-blue-50",
          pending: "bg-orange-50",
          completed: "bg-green-50",
        };
        return map[tabKey] || "bg-gray-100";
      },
      iconColorClass(tabKey) {
        const map = {
          inProgress: "text-blue-600",
          pending: "text-orange-600",
          completed: "text-green-600",
        };
        return map[tabKey] || "text-gray-500";
      },
      taskIcon(tabKey) {
        const map = {
          inProgress: "mdi:code-tags",
          pending: "mdi:play-circle",
          completed: "mdi:check-circle",
        };
        return map[tabKey] || "mdi:book-open";
      },
      actionButtonClass(tabKey) {
        const map = {
          inProgress: "bg-blue-600 hover:bg-blue-700",
          pending: "bg-orange-500 hover:bg-orange-600",
          completed: "bg-gray-600 hover:bg-gray-700",
        };
        return map[tabKey] || "bg-gray-500";
      },
      actionLabel(tabKey) {
        const map = {
          inProgress: "继续",
          pending: "开始",
          completed: "查看",
        };
        return map[tabKey] || "查看";
      },
      normalizeTaskFromApi(task) {
        if (!task) return null;
        const priorityValue = Number(task.priority);
        const priority =
          priorityValue >= 2
            ? "high"
            : priorityValue === 1
            ? "medium"
            : "low";
        const progressValue = clampPercentage(
          Number(task.progress ?? task.Progress ?? 0)
        );
        const tags = [];
        const categoryName =
          task.category?.name || task.category?.Name || task.Category?.name;
        if (categoryName) {
          tags.push(categoryName);
        }
        const id =
          task.id ||
          task.ID ||
          task.Id ||
          `task-${Date.now()}-${Math.random().toString(16).slice(2)}`;

        return {
          id,
          title: task.title || "未命名任务",
          description: task.description || "暂无描述",
          progress: progressValue,
          priority,
          dueDate: formatDateString(task.due_at || task.dueAt),
          tags,
        };
      },
    },
  };
</script>
