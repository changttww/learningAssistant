<template>
  <div class="min-h-full bg-gray-50">
    <div class="w-full py-8">
      <!-- 顶部统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-5 mb-4">
        <div
          class="stat-card bg-blue-50 rounded-lg p-4 flex flex-col items-center justify-center"
        >
          <span class="text-2xl font-bold text-blue-600">{{
            stats.total
          }}</span>
          <span class="text-gray-600 text-sm mt-1">总任务数</span>
        </div>
        <button
          type="button"
          @click="setStatusFilter('completed')"
          class="stat-card bg-green-50 rounded-lg p-4 flex flex-col items-center justify-center cursor-pointer hover:shadow focus:outline-none focus:ring-2 focus:ring-green-400 active:scale-95 transition"
          aria-label="已完成任务"
        >
          <span class="text-2xl font-bold text-green-600">{{
            stats.completed
          }}</span>
          <span class="text-gray-700 text-sm mt-1 font-medium">已完成</span>
        </button>
        <button
          type="button"
          @click="setStatusFilter('in-progress')"
          class="stat-card bg-orange-50 rounded-lg p-4 flex flex-col items-center justify-center cursor-pointer hover:shadow focus:outline-none focus:ring-2 focus:ring-orange-400 active:scale-95 transition"
          aria-label="进行中任务"
        >
          <span class="text-2xl font-bold text-orange-600">{{
            stats.inProgress
          }}</span>
          <span class="text-gray-700 text-sm mt-1 font-medium">进行中</span>
        </button>
        <button
          type="button"
          @click="setStatusFilter('overdue')"
          class="stat-card bg-red-50 rounded-lg p-4 flex flex-col items-center justify-center cursor-pointer hover:shadow focus:outline-none focus:ring-2 focus:ring-red-400 active:scale-95 transition"
          aria-label="已逾期任务"
        >
          <span class="text-2xl font-bold text-red-600">{{
            stats.overdue
          }}</span>
          <span class="text-gray-700 text-sm mt-1 font-medium">已逾期</span>
        </button>
      </div>

      <!-- 状态任务详情列表 -->
      <div v-if="statusFilter" class="mb-6">
        <div class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-bold text-gray-800 text-lg">
              {{ getStatusLabel(statusFilter) }} 任务
            </h3>
            <button
              @click="clearStatusFilter"
              class="text-sm px-3 py-1 rounded border border-gray-300 text-gray-700 hover:bg-gray-50"
              aria-label="关闭状态面板"
            >
              关闭
            </button>
          </div>
          <div
            v-if="filteredTasksByStatus.length === 0"
            class="text-gray-500 text-sm py-4"
          >
            该状态暂无任务。
          </div>
          <div v-else class="space-y-3">
            <div
              v-for="task in filteredTasksByStatus"
              :key="task.id"
              class="p-3 border border-gray-200 rounded hover:border-blue-600 hover:shadow transition"
            >
              <div class="flex items-center justify-between">
                <div class="flex-1">
                  <div class="font-medium text-gray-800">{{ task.title }}</div>
                  <div class="text-xs text-gray-500 mt-1">
                    {{ task.date }} · {{ task.time }}
                  </div>
                  <div class="text-xs text-gray-600 mt-1">
                    {{ task.description }}
                  </div>
                </div>
                <span
                  :class="[
                    'text-xs px-2 py-0.5 rounded',
                    getCategoryStyle(task.category),
                  ]"
                >
                  {{ task.category }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 日历与任务区域 -->
      <div class="flex gap-5 mb-6" style="min-height: 500px">
        <!-- 左侧日历 -->
        <div
          class="flex-1 bg-white rounded-lg border border-gray-200 overflow-hidden modern-calendar"
        >
          <!-- 日历头部 -->
          <div
            class="flex items-center justify-between px-4 h-14 border-b border-gray-200 bg-gray-50"
          >
            <button
              @click="previousMonth"
              class="text-blue-600 text-xl font-bold hover:bg-blue-50 w-8 h-8 rounded-full flex items-center justify-center"
            >
              &lt;
            </button>
            <div class="font-bold text-base text-gray-800">
              {{ currentMonthYear }}
            </div>
            <button
              @click="nextMonth"
              class="text-blue-600 text-xl font-bold hover:bg-blue-50 w-8 h-8 rounded-full flex items-center justify-center"
            >
              &gt;
            </button>
          </div>

          <!-- 星期标题 -->
          <div class="grid grid-cols-7 bg-white py-3 border-b border-gray-100">
            <div
              class="flex items-center justify-center text-sm text-red-500 font-medium"
            >
              日
            </div>
            <div
              class="flex items-center justify-center text-sm text-gray-600 font-medium"
            >
              一
            </div>
            <div
              class="flex items-center justify-center text-sm text-gray-600 font-medium"
            >
              二
            </div>
            <div
              class="flex items-center justify-center text-sm text-gray-600 font-medium"
            >
              三
            </div>
            <div
              class="flex items-center justify-center text-sm text-gray-600 font-medium"
            >
              四
            </div>
            <div
              class="flex items-center justify-center text-sm text-gray-600 font-medium"
            >
              五
            </div>
            <div
              class="flex items-center justify-center text-sm text-red-500 font-medium"
            >
              六
            </div>
          </div>

          <!-- 日期网格 -->
          <div class="grid grid-cols-7 grid-rows-6 gap-0 bg-white">
            <div
              v-for="date in calendarDates"
              :key="date.dateString"
              @click="selectDate(date)"
              :class="[
                'modern-date-cell cursor-pointer relative h-20 transition-all flex flex-col items-center justify-center p-1',
                {
                  'bg-blue-50 border-2 border-blue-600': date.isSelected,
                  'hover:bg-blue-50': !date.isSelected,
                  today: date.isToday,
                },
              ]"
            >
              <div
                :class="[
                  'date-number w-8 h-8 flex items-center justify-center rounded-full text-xs',
                  {
                    'bg-blue-600 text-white': date.isToday,
                    'text-gray-400': !date.isCurrentMonth,
                    'text-gray-800': date.isCurrentMonth && !date.isToday,
                  },
                ]"
              >
                {{ date.day }}
              </div>
              <div v-if="date.tasks && date.tasks.length > 0" class="mt-1 flex">
                <span
                  v-for="(task, index) in date.tasks.slice(0, 3)"
                  :key="index"
                  :class="[
                    'task-dot inline-block w-1.5 h-1.5 rounded-full mx-0.5',
                    {
                      'bg-green-500': task.status === 'completed',
                      'bg-orange-500': task.status === 'in-progress',
                      'bg-red-500': task.status === 'overdue',
                    },
                  ]"
                ></span>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧任务区域 -->
        <div
          class="w-96 bg-white rounded-lg border border-gray-200 flex flex-col"
        >
          <!-- 任务头 -->
          <div
            class="h-12 flex items-center justify-between px-4 border-b border-gray-200"
          >
            <div class="font-bold text-base text-gray-800">
              {{ selectedDateFormatted }} 任务
            </div>
            <div class="flex items-center gap-2">
              <button
                @click="sortMode = 'time'"
                :class="[
                  'text-sm py-1 px-2 border rounded',
                  sortMode === 'time'
                    ? 'text-blue-600 border-blue-600 bg-blue-50'
                    : 'text-gray-600 border-gray-300 hover:border-blue-600 hover:text-blue-600',
                ]"
              >
                按开始时间排序
              </button>
              <button
                @click="sortMode = 'category'"
                :class="[
                  'text-sm py-1 px-2 border rounded',
                  sortMode === 'category'
                    ? 'text-blue-600 border-blue-600 bg-blue-50'
                    : 'text-gray-600 border-gray-300 hover:border-blue-600 hover:text-blue-600',
                ]"
              >
                按任务类别排序
              </button>
              <button
                @click="openTaskModalSelected"
                class="text-sm text-blue-600 py-1 px-2 border border-blue-600 rounded hover:bg-blue-50"
              >
                + 添加任务
              </button>
            </div>
          </div>

          <!-- 任务列表 -->
          <div class="flex-1 bg-gray-50 rounded-b-lg p-3 overflow-auto">
            <div
              v-if="selectedDateTasks.length === 0"
              class="text-center text-gray-500 py-8"
            >
              该日期暂无任务
            </div>

            <!-- 任务项 -->
            <div
              v-for="task in sortedSelectedDateTasks"
              :key="task.id"
              class="bg-white border border-gray-200 rounded p-3 mb-3"
            >
              <div class="flex items-start">
                <div
                  @click="toggleTaskComplete(task)"
                  :class="[
                    'w-4 h-4 rounded border border-gray-300 flex items-center justify-center mr-2 cursor-pointer',
                    {
                      'bg-green-500': task.status === 'completed',
                    },
                  ]"
                >
                  <svg
                    v-if="task.status === 'completed'"
                    width="10"
                    height="8"
                    viewBox="0 0 10 8"
                    fill="none"
                  >
                    <path
                      d="M9 1L3.5 6.5L1 4"
                      stroke="white"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </div>
                <div class="flex-1">
                  <div class="flex items-center">
                    <span
                      :class="[
                        'text-sm font-medium',
                        {
                          'text-gray-800': task.status !== 'completed',
                          'text-gray-500 line-through':
                            task.status === 'completed',
                        },
                      ]"
                      >{{ task.title }}</span
                    >
                    <span class="text-xs text-gray-500 ml-2"
                      >· {{ task.time }}</span
                    >
                    <span
                      :class="[
                        'text-xs ml-2 px-2 py-0.5 rounded',
                        getCategoryStyle(task.category),
                      ]"
                      >{{ task.category }}</span
                    >
                  </div>
                  <p class="text-xs text-gray-600 italic mt-1">
                    {{ task.description }}
                  </p>
                </div>
              </div>

              <div class="mt-3">
                <div class="text-sm font-bold text-blue-600 mb-1">笔记本</div>
                <textarea
                  v-model="task.notes"
                  class="w-full h-20 p-2 text-xs border border-gray-300 rounded text-gray-800 focus:border-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-200"
                  placeholder="输入笔记或心得..."
                ></textarea>
                <div class="flex items-center mt-2">
                  <button
                    @click="openNotebookModal(task)"
                    class="text-xs text-blue-600 hover:underline"
                  >
                    打开笔记本
                  </button>
                  <button class="text-xs text-gray-500 ml-4">评论</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 笔记列表 -->
      <div class="mt-4 mb-2">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-bold text-gray-800">我的笔记</h2>
          <div class="flex space-x-4">
            <button
              :class="[
                'text-sm py-1 px-3 border rounded',
                notesSortBy === 'category'
                  ? 'text-blue-600 border-blue-600 bg-blue-50'
                  : 'text-gray-600 border-gray-300 hover:border-blue-600 hover:text-blue-600',
              ]"
              @click="notesSortBy = 'category'"
            >
              按主题分类
            </button>
            <button
              :class="[
                'text-sm py-1 px-3 border rounded',
                notesSortBy === 'time'
                  ? 'text-blue-600 border-blue-600 bg-blue-50'
                  : 'text-gray-600 border-gray-300 hover:border-blue-600 hover:text-blue-600',
              ]"
              @click="notesSortBy = 'time'"
            >
              按时间排序
            </button>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
          <div
            v-for="note in sortedNotes"
            :key="note.id"
            class="bg-white border border-gray-200 rounded-lg p-4 hover:border-blue-600 hover:shadow-md transition-all cursor-pointer"
            @click="openNotebookModal(note)"
          >
            <div class="flex items-center justify-between mb-2">
              <div class="font-medium">{{ note.title }}</div>
              <span class="text-xs text-gray-500">{{ note.date }}</span>
            </div>
            <p class="text-xs text-gray-600 mb-2 line-clamp-2">
              {{ note.content }}
            </p>
            <div class="flex items-center justify-between">
              <span
                :class="[
                  'text-xs px-2 py-0.5 rounded-full',
                  getCategoryStyle(note.category),
                ]"
                >{{ note.category }}</span
              >
              <button class="text-xs text-blue-600">查看详情</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 新建任务按钮（系统日期） -->
    <button
      @click="openTaskModalSystem"
      class="fixed top-20 right-[5%] bg-blue-600 text-white text-sm px-4 py-2 rounded shadow-lg hover:bg-blue-700 z-50"
    >
      + 新建任务
    </button>

    <!-- 任务弹窗 -->
    <div
      v-if="showTaskModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="closeTaskModal"
    >
      <div
        class="bg-white rounded-lg shadow-xl w-full max-w-md max-h-screen overflow-y-auto"
        @click.stop
      >
        <div
          class="flex items-center justify-between p-4 border-b border-gray-200"
        >
          <h2 class="text-lg font-bold text-gray-800">创建新任务</h2>
          <button
            @click="closeTaskModal"
            class="text-gray-500 hover:text-gray-700"
          >
            <iconify-icon
              icon="mdi:close"
              width="20"
              height="20"
            ></iconify-icon>
          </button>
        </div>
        <div class="p-4">
          <!-- 自然语言输入框 -->
          <div class="mb-4 border border-blue-600 rounded bg-blue-50 p-3">
            <label class="block text-sm text-blue-600 mb-1 font-medium"
              >自然语言输入</label
            >
            <div class="flex">
              <input
                v-model="naturalLanguageInput"
                type="text"
                class="flex-1 border-0 bg-transparent p-1 text-sm outline-none"
                placeholder="例如：明天下午3点完成数学作业第三章"
              />
              <button
                @click="parseNaturalLanguage"
                class="text-white bg-blue-600 px-3 py-1 rounded text-sm hover:bg-blue-700"
              >
                <iconify-icon
                  icon="mdi:wand"
                  width="16"
                  height="16"
                  class="mr-1"
                ></iconify-icon>
                解析
              </button>
            </div>
          </div>

          <!-- 表单输入 -->
          <div class="space-y-4">
            <div>
              <label class="block text-sm text-gray-600 mb-1"
                >任务名称 <span class="text-red-500">*</span></label
              >
              <input
                v-model="newTask.title"
                type="text"
                class="w-full border border-gray-300 p-2 rounded text-sm focus:border-blue-600 focus:outline-none"
                placeholder="输入任务名称"
              />
            </div>

            <div>
              <label class="block text-sm text-gray-600 mb-1">任务描述</label>
              <textarea
                v-model="newTask.description"
                class="w-full border border-gray-300 p-2 rounded text-sm h-20 focus:border-blue-600 focus:outline-none"
                placeholder="输入任务详情描述"
              ></textarea>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-gray-600 mb-1"
                  >开始日期 <span class="text-red-500">*</span></label
                >
                <input
                  v-model="newTask.startDate"
                  type="date"
                  class="w-full border border-gray-300 p-2 rounded text-sm focus:border-blue-600 focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-sm text-gray-600 mb-1">开始时间</label>
                <input
                  v-model="newTask.startTime"
                  type="time"
                  class="w-full border border-gray-300 p-2 rounded text-sm focus:border-blue-600 focus:outline-none"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-gray-600 mb-1"
                  >结束日期 <span class="text-red-500">*</span></label
                >
                <input
                  v-model="newTask.endDate"
                  type="date"
                  class="w-full border border-gray-300 p-2 rounded text-sm focus:border-blue-600 focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-sm text-gray-600 mb-1">结束时间</label>
                <input
                  v-model="newTask.endTime"
                  type="time"
                  class="w-full border border-gray-300 p-2 rounded text-sm focus:border-blue-600 focus:outline-none"
                />
              </div>
            </div>

            <div>
              <label class="block text-sm text-gray-600 mb-1">任务分类</label>
              <select
                v-model="newTask.category"
                class="w-full border border-gray-300 p-2 rounded text-sm focus:border-blue-600 focus:outline-none"
              >
                <option value="">请选择分类</option>
                <option value="study">学习</option>
                <option value="exam">考试</option>
                <option value="project">项目</option>
                <option value="reading">阅读</option>
                <option value="other">其他</option>
              </select>
            </div>
          </div>
        </div>
        <div class="flex justify-end p-4 border-t border-gray-200">
          <button
            @click="closeTaskModal"
            class="text-sm text-gray-600 border border-gray-300 py-1.5 px-4 rounded mr-3 hover:bg-gray-50"
          >
            取消
          </button>
          <button
            @click="saveTask"
            class="text-sm text-white bg-blue-600 py-1.5 px-4 rounded hover:bg-blue-700"
          >
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 笔记本弹窗 -->
    <div
      v-if="showNotebookModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="closeNotebookModal"
    >
      <div
        :class="[
          'bg-white rounded-lg shadow-xl transition-all',
          isNotebookFullscreen ? 'w-11/12 h-5/6' : 'w-1/2 h-4/5',
        ]"
        @click.stop
      >
        <div
          class="flex items-center justify-between p-4 border-b border-gray-200"
        >
          <div class="flex items-center">
            <h2 class="text-lg font-bold text-gray-800">
              {{ currentNotebook.title }}
            </h2>
            <span
              :class="[
                'ml-2 text-xs px-2 py-0.5 rounded-full',
                getCategoryStyle(currentNotebook.category),
              ]"
              >{{ currentNotebook.category }}</span
            >
          </div>
          <div class="flex items-center space-x-2">
            <button
              @click="toggleNotebookFullscreen"
              class="text-gray-500 hover:text-gray-700"
            >
              <iconify-icon
                :icon="
                  isNotebookFullscreen
                    ? 'mdi:fullscreen-exit'
                    : 'mdi:fullscreen'
                "
                width="20"
                height="20"
              ></iconify-icon>
            </button>
            <button
              @click="closeNotebookModal"
              class="text-gray-500 hover:text-gray-700"
            >
              <iconify-icon
                icon="mdi:close"
                width="20"
                height="20"
              ></iconify-icon>
            </button>
          </div>
        </div>
        <div class="p-4 h-full overflow-hidden flex flex-col">
          <div class="flex items-center justify-between mb-4">
            <div class="text-sm text-gray-500">
              最后更新: {{ currentNotebook.lastUpdated }}
            </div>
            <div class="flex space-x-2">
              <button
                class="text-xs text-blue-600 py-1 px-2 border border-blue-600 rounded hover:bg-blue-50"
              >
                <iconify-icon
                  icon="mdi:format-bold"
                  width="16"
                  height="16"
                ></iconify-icon>
              </button>
              <button
                class="text-xs text-blue-600 py-1 px-2 border border-blue-600 rounded hover:bg-blue-50"
              >
                <iconify-icon
                  icon="mdi:format-italic"
                  width="16"
                  height="16"
                ></iconify-icon>
              </button>
              <button
                class="text-xs text-blue-600 py-1 px-2 border border-blue-600 rounded hover:bg-blue-50"
              >
                <iconify-icon
                  icon="mdi:format-list-bulleted"
                  width="16"
                  height="16"
                ></iconify-icon>
              </button>
            </div>
          </div>

          <textarea
            v-model="currentNotebook.content"
            class="w-full flex-1 p-3 border border-gray-300 rounded text-sm resize-none focus:border-blue-600 focus:outline-none"
          ></textarea>

          <div class="mt-4 flex items-center justify-between">
            <div v-if="currentNotebook.relatedTask">
              <span class="text-sm text-gray-600 mr-2">关联任务:</span>
              <span class="text-sm text-blue-600">{{
                currentNotebook.relatedTask
              }}</span>
            </div>
            <div>
              <button
                class="text-sm text-blue-600 py-1 px-3 border border-blue-600 rounded hover:bg-blue-50"
              >
                插入图片
              </button>
              <button
                @click="saveNotebook"
                class="ml-2 text-sm text-white bg-blue-600 py-1 px-3 rounded hover:bg-blue-700"
              >
                保存笔记
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { ref, computed, onMounted, watch } from "vue";

  export default {
    name: "PersonalTasks",
    setup() {
      // 响应式数据
      const currentDate = ref(new Date());
      const selectedDate = ref(null);
      const showTaskModal = ref(false);
      const showNotebookModal = ref(false);
      const isNotebookFullscreen = ref(false);
      const notesSortBy = ref("category");
      const naturalLanguageInput = ref("");
      const statusFilter = ref(null);
      const modalDateMode = ref('system');

      // 统计数据
      const stats = ref({
        total: 12,
        completed: 8,
        inProgress: 3,
        overdue: 1,
      });

      // 新任务表单
      const newTask = ref({
        title: "",
        description: "",
        startDate: "",
        startTime: "",
        endDate: "",
        endTime: "",
        category: "",
      });

      // 当前笔记本
      const currentNotebook = ref({
        title: "",
        category: "",
        content: "",
        lastUpdated: "",
        relatedTask: "",
      });

      // 任务数据
      const tasks = ref([
        {
          id: 1,
          title: "完成数学作业",
          description: "复习高数第三章知识点",
          date: "2024-03-05",
          time: "13:30前",
          status: "completed",
          notes: "今天掌握了导数应用部分，重点练习了极值求法题目。",
          category: "数学",
        },
        {
          id: 2,
          title: "准备英语报告",
          description: "关于气候变化的严重影响与应对",
          date: "2024-03-05",
          time: "15:00前",
          status: "in-progress",
          notes: "",
          category: "英语",
        },
        {
          id: 3,
          title: "物理实验预习",
          description: "波动光学实验操作流程",
          date: "2024-03-05",
          time: "17:00前",
          status: "pending",
          notes: "",
          category: "物理",
        },
        {
          id: 4,
          title: "阅读文献",
          description: "机器学习相关论文",
          date: "2024-03-06",
          time: "10:00前",
          status: "overdue",
          notes: "",
          category: "研究",
        },
        // 示例任务：2025年10月5日
        {
          id: 1001,
          title: "项目规划会议",
          description: "讨论季度目标与里程碑安排",
          date: "2025-10-05",
          time: "09:00 - 11:00",
          status: "in-progress",
          notes: "",
          category: "工作",
        },
        // 示例任务：2025年10月5日
        {
          id: 1002,
          title: "编程课",
          description: "学习循环与函数基础，完成课堂练习",
          date: "2025-10-05",
          time: "08:00 - 09:00",
          status: "in-progress",
          notes: "",
          category: "学习",
        },
        // 示例任务：2025年10月5日
        {
          id: 1003,
          title: "项目管理",
          description: "协调团队成员，分配任务，监控进度",
          date: "2025-10-05",
          time: "10:00 - 12:00",
          status: "in-progress",
          notes: "",
          category: "其它",
        },
      ]);

      // 笔记数据
      const notes = ref([
        {
          id: 1,
          title: "数学笔记",
          content: "今天掌握了导数应用部分，重点练习了极值求法题目。",
          date: "3月5日",
          category: "数学",
          lastUpdated: "2024年3月5日 14:35",
          relatedTask: "完成数学作业",
        },
        {
          id: 2,
          title: "英语演讲准备",
          content: "气候变化的报告的关键论点和数据整理。需要补充更多实例。",
          date: "3月3日",
          category: "英语",
          lastUpdated: "2024年3月3日 16:20",
          relatedTask: "准备英语报告",
        },
        {
          id: 3,
          title: "物理实验记录",
          content:
            "波动光学实验的预习材料和关键步骤记录。需要注意的实验误差来源。",
          date: "2月28日",
          category: "物理",
          lastUpdated: "2024年2月28日 09:15",
          relatedTask: "物理实验预习",
        },
      ]);

      // 排序模式：none | time | category（默认按开始时间）
      const sortMode = ref("time");

      // 解析开始时间为分钟数
      const getStartMinutes = (task) => {
        const t = task.time || "";
        // 格式：HH:MM - HH:MM
        const rangeMatch = t.match(/^(\d{2}):(\d{2})\s*-\s*(\d{2}):(\d{2})$/);
        if (rangeMatch) {
          const h = parseInt(rangeMatch[1], 10);
          const m = parseInt(rangeMatch[2], 10);
          return h * 60 + m;
        }
        // 格式：HH:MM 或 HH:MM前
        const singleMatch = t.match(/^(\d{2}):(\d{2})/);
        if (singleMatch) {
          const h = parseInt(singleMatch[1], 10);
          const m = parseInt(singleMatch[2], 10);
          return h * 60 + m;
        }
        // 全天或无法解析，排在最后
        return Number.POSITIVE_INFINITY;
      };

      // 本地日期格式化（避免 toISOString 的 UTC 偏移）
      const formatLocalDate = (d) => {
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, "0");
        const day = String(d.getDate()).padStart(2, "0");
        return `${year}-${month}-${day}`;
      };

      // 计算属性
      const currentMonthYear = computed(() => {
        return `${currentDate.value.getFullYear()}年${
          currentDate.value.getMonth() + 1
        }月`;
      });

      const selectedDateFormatted = computed(() => {
        if (!selectedDate.value) return "未选择日期";
        return `${selectedDate.value.getFullYear()}年${
          selectedDate.value.getMonth() + 1
        }月${selectedDate.value.getDate()}日`;
      });

      const calendarDates = computed(() => {
        const year = currentDate.value.getFullYear();
        const month = currentDate.value.getMonth();
        const firstDay = new Date(year, month, 1);
        const lastDay = new Date(year, month + 1, 0);
        const startDate = new Date(firstDay);
        startDate.setDate(startDate.getDate() - firstDay.getDay());

        const dates = [];
        const today = new Date();

        for (let i = 0; i < 42; i++) {
          const date = new Date(startDate);
          date.setDate(startDate.getDate() + i);

          const dateString = formatLocalDate(date);
          const dateTasks = tasks.value.filter(
            (task) => task.date === dateString
          );

          dates.push({
            date: new Date(date),
            day: date.getDate(),
            dateString,
            isCurrentMonth: date.getMonth() === month,
            isToday: date.toDateString() === today.toDateString(),
            isSelected: selectedDate.value
              ? date.toDateString() === selectedDate.value.toDateString()
              : false,
            tasks: dateTasks,
          });
        }

        return dates;
      });

      const selectedDateTasks = computed(() => {
        if (!selectedDate.value) return [];
        const dateString = formatLocalDate(selectedDate.value);
        return tasks.value.filter((task) => task.date === dateString);
      });

      // 根据排序模式返回排序后的任务
      const sortedSelectedDateTasks = computed(() => {
        const list = [...selectedDateTasks.value];
        if (sortMode.value === "time") {
          return list.sort((a, b) => getStartMinutes(a) - getStartMinutes(b));
        }
        if (sortMode.value === "category") {
          return list.sort((a, b) =>
            (a.category || "").localeCompare(b.category || "")
          );
        }
        return list;
      });

      const sortedNotes = computed(() => {
        if (notesSortBy.value === "time") {
          return [...notes.value].sort(
            (a, b) => new Date(b.lastUpdated) - new Date(a.lastUpdated)
          );
        }
        return [...notes.value].sort((a, b) =>
          a.category.localeCompare(b.category)
        );
      });

      const filteredTasksByStatus = computed(() => {
        if (!statusFilter.value) return [];
        return tasks.value.filter((t) => t.status === statusFilter.value);
      });

      // 方法
      const previousMonth = () => {
        currentDate.value = new Date(
          currentDate.value.getFullYear(),
          currentDate.value.getMonth() - 1,
          1
        );
      };

      const nextMonth = () => {
        currentDate.value = new Date(
          currentDate.value.getFullYear(),
          currentDate.value.getMonth() + 1,
          1
        );
      };

      const selectDate = (date) => {
        selectedDate.value = new Date(date.date);
      };

      const openTaskModalSystem = () => {
        modalDateMode.value = 'system';
        const today = formatLocalDate(new Date());
        newTask.value = {
          title: "",
          description: "",
          startDate: today,
          startTime: "",
          endDate: today,
          endTime: "",
          category: "",
        };
        showTaskModal.value = true;
      };

      const openTaskModalSelected = () => {
        if (!selectedDate.value) {
          alert("请先在日历中选择日期");
          return;
        }
        modalDateMode.value = 'selected';
        const dateString = formatLocalDate(selectedDate.value);
        newTask.value = {
          title: "",
          description: "",
          startDate: dateString,
          startTime: "",
          endDate: dateString,
          endTime: "",
          category: "",
        };
        showTaskModal.value = true;
      };

      // 弹窗打开时，与日历选中日期保持同步（若用户更换选中日期）
      watch(selectedDate, (d) => {
        if (!showTaskModal.value || !d) return;
        if (modalDateMode.value !== 'selected') return;
        const ds = formatLocalDate(d);
        newTask.value.startDate = ds;
        newTask.value.endDate = ds;
      });

        const closeTaskModal = () => {
          showTaskModal.value = false;
          naturalLanguageInput.value = "";
          modalDateMode.value = 'system';
        };

      const parseNaturalLanguage = () => {
        // 简单的自然语言解析示例
        const input = naturalLanguageInput.value.toLowerCase();
        if (input.includes("数学")) {
          newTask.value.category = "study";
          newTask.value.title = "数学作业";
        }
        if (input.includes("明天")) {
          const tomorrow = new Date();
          tomorrow.setDate(tomorrow.getDate() + 1);
          newTask.value.startDate = formatLocalDate(tomorrow);
          newTask.value.endDate = formatLocalDate(tomorrow);
        }
        if (input.includes("下午3点") || input.includes("15:00")) {
          newTask.value.endTime = "15:00";
        }
      };

      const saveTask = () => {
        if (
          !newTask.value.title ||
          !newTask.value.startDate ||
          !newTask.value.endDate
        ) {
          alert("请填写必填项");
          return;
        }

        const task = {
          id: Date.now(),
          title: newTask.value.title,
          description: newTask.value.description,
          date: newTask.value.endDate,
          time: newTask.value.endTime || "全天",
          status: "pending",
          notes: "",
          category: newTask.value.category || "其他",
        };

        tasks.value.push(task);
        stats.value.total++;
        closeTaskModal();
      };

      const toggleTaskComplete = (task) => {
        if (task.status === "completed") {
          task.status = "pending";
          stats.value.completed--;
        } else {
          task.status = "completed";
          stats.value.completed++;
          if (task.status === "in-progress") {
            stats.value.inProgress--;
          }
        }
      };

      const openNotebookModal = (item) => {
        if (item.title) {
          // 从任务打开
          currentNotebook.value = {
            title: item.title,
            category: item.category || "学习",
            content: item.notes || "",
            lastUpdated: new Date().toLocaleString("zh-CN"),
            relatedTask: item.title,
          };
        } else {
          // 从笔记列表打开
          currentNotebook.value = { ...item };
        }
        showNotebookModal.value = true;
      };

      const closeNotebookModal = () => {
        showNotebookModal.value = false;
        isNotebookFullscreen.value = false;
      };

      const toggleNotebookFullscreen = () => {
        isNotebookFullscreen.value = !isNotebookFullscreen.value;
      };

      const saveNotebook = () => {
        // 保存笔记逻辑
        alert("笔记已保存");
      };

      const getCategoryStyle = (category) => {
        const styles = {
          // 中文类别
          数学: "bg-blue-50 text-blue-600",
          英语: "bg-orange-50 text-orange-600",
          物理: "bg-red-50 text-red-600",
          研究: "bg-purple-50 text-purple-600",
          学习: "bg-blue-50 text-blue-600",
          工作: "bg-teal-50 text-teal-600",
          其他: "bg-gray-50 text-gray-600",
          // 英文代码类别（表单值）
          study: "bg-blue-50 text-blue-600",
          exam: "bg-red-50 text-red-600",
          project: "bg-purple-50 text-purple-600",
          reading: "bg-green-50 text-green-600",
          other: "bg-gray-50 text-gray-600",
        };
        return styles[category] || "bg-gray-50 text-gray-600";
      };

      const setStatusFilter = (status) => {
        statusFilter.value = status;
      };
      const clearStatusFilter = () => {
        statusFilter.value = null;
      };
      const getStatusLabel = (status) => {
        const map = {
          completed: "已完成",
          "in-progress": "进行中",
          overdue: "已逾期",
        };
        return map[status] || "任务";
      };

      // 初始化
      onMounted(() => {
        selectedDate.value = new Date();
      });

      return {
        currentDate,
        selectedDate,
        showTaskModal,
        showNotebookModal,
        isNotebookFullscreen,
        notesSortBy,
        naturalLanguageInput,
        sortMode,
        stats,
        newTask,
        currentNotebook,
        tasks,
        notes,
        currentMonthYear,
        selectedDateFormatted,
        calendarDates,
        selectedDateTasks,
        sortedSelectedDateTasks,
        sortedNotes,
        previousMonth,
        nextMonth,
        selectDate,
        closeTaskModal,
        parseNaturalLanguage,
        saveTask,
        toggleTaskComplete,
        openNotebookModal,
        closeNotebookModal,
        toggleNotebookFullscreen,
        saveNotebook,
        getCategoryStyle,
        // 状态筛选
        statusFilter,
        filteredTasksByStatus,
        setStatusFilter,
        clearStatusFilter,
        getStatusLabel,
        // 新建任务打开模式
        modalDateMode,
        openTaskModalSystem,
        openTaskModalSelected,
      };
    },
  };
</script>

<style scoped>
  .modern-calendar {
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  }

  .modern-date-cell {
    transition: all 0.2s;
  }

  .modern-date-cell:hover {
    background-color: #f0f7ff;
  }

  .modern-date-cell.today .date-number {
    background-color: #1e88e5;
    color: white;
  }

  .task-dot {
    transition: all 0.2s;
  }

  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .stat-card {
    transition: transform 0.2s;
  }
</style>
