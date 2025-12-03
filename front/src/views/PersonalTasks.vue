<template>
  <div class="min-h-full bg-gray-50">
    <div class="w-full py-8">
      <!-- 顶部统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-5 gap-5 mb-6">
        <!-- 总任务数卡片 -->
        <div
          class="stat-card group bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl p-5 flex flex-col items-center justify-center shadow-lg hover:shadow-xl transition-all duration-300 border-2 border-blue-600"
        >
          <div class="bg-white/20 backdrop-blur-sm w-14 h-14 rounded-full flex items-center justify-center mb-3 group-hover:scale-110 transition-transform duration-300">
            <iconify-icon icon="mdi:format-list-checks" width="28" height="28" class="text-white"></iconify-icon>
          </div>
          <span class="text-3xl font-bold text-white drop-shadow-md">{{ stats.total }}</span>
          <span class="text-blue-100 text-sm mt-1.5 font-medium">总任务数</span>
        </div>

        <!-- 已完成任务卡片 -->
        <button
          type="button"
          @click="setStatusFilter('completed')"
          class="stat-card group bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl p-5 flex flex-col items-center justify-center cursor-pointer shadow-lg hover:shadow-xl focus:outline-none focus:ring-2 focus:ring-green-400 focus:ring-offset-2 active:scale-95 transition-all duration-300 transform hover:-translate-y-1 border-2 border-green-600"
          aria-label="已完成任务"
        >
          <div class="bg-white/20 backdrop-blur-sm w-14 h-14 rounded-full flex items-center justify-center mb-3 group-hover:scale-110 group-hover:rotate-12 transition-all duration-300">
            <iconify-icon icon="mdi:check-circle" width="28" height="28" class="text-white"></iconify-icon>
          </div>
          <span class="text-3xl font-bold text-white drop-shadow-md">{{ stats.completed }}</span>
          <span class="text-green-100 text-sm mt-1.5 font-medium">已完成</span>
        </button>

        <!-- 进行中任务卡片 -->
        <button
          type="button"
          @click="setStatusFilter('in-progress')"
          class="stat-card group bg-gradient-to-br from-orange-500 to-amber-600 rounded-xl p-5 flex flex-col items-center justify-center cursor-pointer shadow-lg hover:shadow-xl focus:outline-none focus:ring-2 focus:ring-orange-400 focus:ring-offset-2 active:scale-95 transition-all duration-300 transform hover:-translate-y-1 border-2 border-orange-600"
          aria-label="进行中任务"
        >
          <div class="bg-white/20 backdrop-blur-sm w-14 h-14 rounded-full flex items-center justify-center mb-3 group-hover:scale-110 transition-transform duration-300">
            <iconify-icon icon="mdi:clock-fast" width="28" height="28" class="text-white"></iconify-icon>
          </div>
          <span class="text-3xl font-bold text-white drop-shadow-md">{{ stats.inProgress }}</span>
          <span class="text-orange-100 text-sm mt-1.5 font-medium">进行中</span>
        </button>

        <!-- 待处理任务卡片 -->
        <button
          type="button"
          @click="setStatusFilter('pending')"
          class="stat-card group bg-gradient-to-br from-gray-500 to-slate-600 rounded-xl p-5 flex flex-col items-center justify-center cursor-pointer shadow-lg hover:shadow-xl focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-2 active:scale-95 transition-all duration-300 transform hover:-translate-y-1 border-2 border-gray-600"
          aria-label="待处理任务"
        >
          <div class="bg-white/20 backdrop-blur-sm w-14 h-14 rounded-full flex items-center justify-center mb-3 group-hover:scale-110 transition-transform duration-300">
            <iconify-icon icon="mdi:clock-outline" width="28" height="28" class="text-white"></iconify-icon>
          </div>
          <span class="text-3xl font-bold text-white drop-shadow-md">{{ stats.pending }}</span>
          <span class="text-gray-100 text-sm mt-1.5 font-medium">待处理</span>
        </button>

        <!-- 已逾期任务卡片 -->

        <button
          type="button"
          @click="setStatusFilter('overdue')"
          class="stat-card group bg-gradient-to-br from-red-500 to-rose-600 rounded-xl p-5 flex flex-col items-center justify-center cursor-pointer shadow-lg hover:shadow-xl focus:outline-none focus:ring-2 focus:ring-red-400 focus:ring-offset-2 active:scale-95 transition-all duration-300 transform hover:-translate-y-1 border-2 border-red-600"
          aria-label="已逾期任务"
        >
          <div class="bg-white/20 backdrop-blur-sm w-14 h-14 rounded-full flex items-center justify-center mb-3 group-hover:scale-110 group-hover:rotate-12 transition-all duration-300">
            <iconify-icon icon="mdi:alert-circle" width="28" height="28" class="text-white"></iconify-icon>
          </div>
          <span class="text-3xl font-bold text-white drop-shadow-md">{{ stats.overdue }}</span>
          <span class="text-red-100 text-sm mt-1.5 font-medium">已逾期</span>
        </button>
      </div>

      <!-- 状态任务详情列表 -->
      <div v-if="statusFilter" class="mb-6 animate-modal-enter">
        <div class="bg-white border-2 border-gray-200 rounded-2xl shadow-lg overflow-hidden">
          <!-- 列表头部 -->
          <div class="bg-gradient-to-r from-gray-50 to-gray-100 px-5 py-4 border-b-2 border-gray-200">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <div 
                  :class="[
                    'w-10 h-10 rounded-xl flex items-center justify-center shadow-md',
                    statusFilter === 'completed' ? 'bg-gradient-to-br from-green-500 to-green-600' :
                    statusFilter === 'in-progress' ? 'bg-gradient-to-br from-orange-500 to-orange-600' :
                    statusFilter === 'pending' ? 'bg-gradient-to-br from-gray-500 to-gray-600' :
                    'bg-gradient-to-br from-red-500 to-red-600'
                  ]"
                >
                  <iconify-icon 
                    :icon="
                      statusFilter === 'completed' ? 'mdi:check-circle' :
                      statusFilter === 'in-progress' ? 'mdi:clock-fast' :
                      statusFilter === 'pending' ? 'mdi:clock-outline' :
                      'mdi:alert-circle'
                    " 
                    width="24" 
                    height="24"
                    class="text-white"
                  ></iconify-icon>
                </div>
                <div>
                  <h3 class="font-bold text-gray-800 text-lg">
                    {{ getStatusLabel(statusFilter) }} 任务
                  </h3>
                  <p class="text-xs text-gray-500 mt-0.5">
                    共 {{ filteredTasksByStatus.length }} 个任务
                  </p>
                </div>
              </div>
              <button
                @click="clearStatusFilter"
                class="flex items-center gap-2 text-sm px-4 py-2 rounded-lg border-2 border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-gray-400 transition-all font-medium"
                aria-label="关闭状态面板"
              >
                <iconify-icon icon="mdi:close" width="16" height="16"></iconify-icon>
                关闭
              </button>
            </div>
          </div>

          <!-- 任务列表内容 -->
          <div class="p-5">
            <div
              v-if="filteredTasksByStatus.length === 0"
              class="flex flex-col items-center justify-center py-12 text-center"
            >
              <iconify-icon icon="mdi:inbox" width="64" height="64" class="text-gray-300 mb-3"></iconify-icon>
              <p class="text-gray-400 text-sm">该状态暂无任务</p>
            </div>
            <div v-else class="space-y-3">
              <div
                v-for="task in filteredTasksByStatus"
                :key="task.id"
                class="group bg-gradient-to-r from-white to-gray-50 border-2 border-gray-200 rounded-xl p-4 hover:border-blue-400 hover:shadow-md transition-all duration-300"
              >
                <div class="flex items-start gap-4">
                  <!-- 完成状态复选框 -->
                  <button
                    type="button"
                    @click.stop="toggleTaskComplete(task)"
                    :class="[
                      'flex-shrink-0 w-6 h-6 rounded-lg border-2 flex items-center justify-center mt-1',
                      'transition-all duration-200 hover:scale-110',
                      task.status === 'completed'
                        ? 'bg-gradient-to-br from-green-500 to-green-600 border-green-500 text-white shadow-md'
                        : 'border-gray-300 hover:border-blue-500 bg-white hover:bg-blue-50',
                    ]"
                    :title="task.status === 'completed' ? '标记为未完成' : '标记为已完成'"
                  >
                    <svg
                      v-if="task.status === 'completed'"
                      width="14"
                      height="10"
                      viewBox="0 0 12 9"
                      fill="none"
                      class="drop-shadow-sm"
                    >
                      <path
                        d="M10.5 1.5L4.5 7.5L1.5 4.5"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </button>

                  <!-- 任务主要内容 -->
                  <div class="flex-1 min-w-0">
                    <!-- 任务标题和标签 -->
                    <div class="flex items-start justify-between mb-2">
                      <div class="flex-1 min-w-0 mr-3">
                        <h4 
                          :class="[
                            'font-bold text-base mb-1.5',
                            task.status === 'completed' ? 'text-gray-500 line-through' : 'text-gray-800'
                          ]"
                        >
                          {{ task.title }}
                        </h4>
                        <div class="flex items-center gap-2 flex-wrap">
                          <!-- 分类标签 -->
                          <span
                            :class="[
                              'text-xs px-2.5 py-1 rounded-md font-medium shadow-sm',
                              getCategoryStyle(task.category),
                            ]"
                          >
                            {{ task.category }}
                          </span>
                          <!-- 状态标签 -->
                          <span 
                            :class="[
                              'text-xs px-2 py-1 rounded-md font-medium flex items-center gap-1',
                              task.status === 'completed' ? 'bg-green-100 text-green-700' :
                              getTaskActualStatus(task) === '已逾期' ? 'bg-red-100 text-red-700' :
                              getTaskActualStatus(task) === '进行中' ? 'bg-orange-100 text-orange-700' :
                              'bg-gray-100 text-gray-700'
                            ]"
                          >
                            <div 
                              :class="[
                                'w-1.5 h-1.5 rounded-full',
                                getTaskDotColor(task)
                              ]"
                            ></div>
                            {{ getTaskActualStatus(task) }}
                          </span>
                        </div>
                      </div>
                      <!-- 操作按钮 -->
                      <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button
                          @click.stop="openGuidanceModal(task)"
                          class="p-2 text-gray-400 hover:text-purple-600 hover:bg-purple-50 rounded-lg transition-colors"
                          title="获取任务指导"
                        >
                          <iconify-icon icon="mdi:lightbulb-outline" width="16" height="16"></iconify-icon>
                        </button>
                        <button
                          @click.stop="openQuizModal(task)"
                          class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                          title="智能测验"
                        >
                          <iconify-icon icon="mdi:file-question-outline" width="16" height="16"></iconify-icon>
                        </button>
                        <button
                          @click.stop="editTask(task)"
                          class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                          title="编辑任务"
                        >
                          <iconify-icon icon="mdi:pencil-outline" width="16" height="16"></iconify-icon>
                        </button>
                        <button
                          @click.stop="handleDelete(task)"
                          class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                          title="删除任务"
                        >
                          <iconify-icon icon="mdi:trash-can-outline" width="16" height="16"></iconify-icon>
                        </button>
                      </div>
                    </div>

                    <!-- 任务描述 -->
                    <p 
                      v-if="task.description"
                      :class="[
                        'text-sm leading-relaxed mb-2',
                        task.status === 'completed' ? 'text-gray-400' : 'text-gray-600'
                      ]"
                    >
                      {{ task.description }}
                    </p>

                    <!-- 任务时间信息 -->
                    <div class="flex items-center gap-4 text-xs text-gray-500">
                      <div class="flex items-center gap-1">
                        <iconify-icon icon="mdi:calendar" width="14" height="14"></iconify-icon>
                        <span>{{ task.date }}</span>
                      </div>
                      <div class="flex items-center gap-1">
                        <iconify-icon icon="mdi:clock-outline" width="14" height="14"></iconify-icon>
                        <span>{{ task.time }}</span>
                      </div>
                      <div 
                        v-if="task.startDate !== task.endDate"
                        class="flex items-center gap-1 text-blue-600 bg-blue-50 px-2 py-0.5 rounded-md"
                      >
                        <iconify-icon icon="mdi:calendar-range" width="14" height="14"></iconify-icon>
                        <span class="font-medium">{{ task.startDate }} ~ {{ task.endDate }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 日历与任务区域 -->
      <div class="flex gap-5 mb-6" style="min-height: 500px">
        <!-- 左侧日历 -->
        <div
          class="flex-1 bg-white rounded-2xl border-2 border-gray-200 overflow-hidden shadow-lg modern-calendar"
        >
          <!-- 日历头部 -->
          <div
            class="flex items-center justify-between px-6 h-16 border-b-2 border-gray-200 bg-gradient-to-r from-indigo-500 to-purple-600"
          >
            <button
              @click="previousMonth"
              class="group w-10 h-10 rounded-xl flex items-center justify-center bg-white/20 backdrop-blur-sm hover:bg-white/30 transition-all duration-200 text-white hover:scale-110"
              aria-label="上个月"
            >
              <iconify-icon icon="mdi:chevron-left" width="24" height="24"></iconify-icon>
            </button>
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center">
                <iconify-icon icon="mdi:calendar-month" width="20" height="20" class="text-white"></iconify-icon>
              </div>
              <div class="font-bold text-xl text-white drop-shadow-md">
                {{ currentMonthYear }}
              </div>
            </div>
            <button
              @click="nextMonth"
              class="group w-10 h-10 rounded-xl flex items-center justify-center bg-white/20 backdrop-blur-sm hover:bg-white/30 transition-all duration-200 text-white hover:scale-110"
              aria-label="下个月"
            >
              <iconify-icon icon="mdi:chevron-right" width="24" height="24"></iconify-icon>
            </button>
          </div>

          <!-- 星期标题 -->
          <div class="grid grid-cols-7 bg-gradient-to-b from-gray-50 to-white py-4 border-b border-gray-200">
            <div
              class="flex items-center justify-center text-sm text-red-600 font-bold"
            >
              <div class="w-8 h-8 flex items-center justify-center rounded-lg bg-red-50">
                日
              </div>
            </div>
            <div
              class="flex items-center justify-center text-sm text-gray-700 font-bold"
            >
              <div class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-gray-100 transition-colors">
                一
              </div>
            </div>
            <div
              class="flex items-center justify-center text-sm text-gray-700 font-bold"
            >
              <div class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-gray-100 transition-colors">
                二
              </div>
            </div>
            <div
              class="flex items-center justify-center text-sm text-gray-700 font-bold"
            >
              <div class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-gray-100 transition-colors">
                三
              </div>
            </div>
            <div
              class="flex items-center justify-center text-sm text-gray-700 font-bold"
            >
              <div class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-gray-100 transition-colors">
                四
              </div>
            </div>
            <div
              class="flex items-center justify-center text-sm text-gray-700 font-bold"
            >
              <div class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-gray-100 transition-colors">
                五
              </div>
            </div>
            <div
              class="flex items-center justify-center text-sm text-red-600 font-bold"
            >
              <div class="w-8 h-8 flex items-center justify-center rounded-lg bg-red-50">
                六
              </div>
            </div>
          </div>

          <!-- 日期网格 -->
          <div class="grid grid-cols-7 grid-rows-6 gap-0 bg-gradient-to-b from-white to-gray-50 p-2">
            <div
              v-for="date in calendarDates"
              :key="date.dateString"
              @click="selectDate(date)"
              :class="[
                'modern-date-cell group cursor-pointer relative h-24 transition-all duration-300 flex flex-col items-center justify-center p-2 m-0.5 rounded-xl',
                {
                  'bg-gradient-to-br from-blue-500 to-indigo-600 shadow-lg scale-105 border-2 border-blue-600': date.isSelected,
                  'hover:bg-gradient-to-br hover:from-blue-50 hover:to-indigo-50 hover:shadow-md hover:scale-105 border border-transparent hover:border-blue-200': !date.isSelected,
                  'ring-2 ring-purple-400 ring-offset-2': date.isToday && !date.isSelected,
                },
              ]"
            >
              <div
                :class="[
                  'date-number w-9 h-9 flex items-center justify-center rounded-xl text-sm font-bold transition-all duration-200',
                  {
                    'bg-gradient-to-br from-purple-500 to-pink-500 text-white shadow-md': date.isToday && !date.isSelected,
                    'text-white': date.isSelected,
                    'text-gray-400': !date.isCurrentMonth && !date.isSelected,
                    'text-gray-800 group-hover:text-blue-600': date.isCurrentMonth && !date.isToday && !date.isSelected,
                  },
                ]"
              >
                {{ date.day }}
              </div>
              <div v-if="date.tasks && date.tasks.length > 0" class="mt-1.5 flex flex-col items-center gap-1">
                <!-- 任务状态圆点 -->
                <div class="flex items-center justify-center gap-0.5">
                  <span
                    v-for="(task, index) in date.tasks.slice(0, 3)"
                    :key="index"
                    :class="[
                      'task-dot inline-block w-2 h-2 rounded-full shadow-sm transition-transform duration-200 hover:scale-150',
                      date.isSelected ? 'bg-white' : getTaskDotColor(task)
                    ]"
                    :title="`${task.title} - ${getTaskActualStatus(task)}`"
                  ></span>
                </div>
                <!-- 任务数量标记 -->
                <span 
                  v-if="date.tasks.length > 3" 
                  :class="[
                    'text-xs font-bold px-1.5 py-0.5 rounded-md',
                    date.isSelected ? 'text-white bg-white/20' : 'text-blue-600 bg-blue-100'
                  ]"
                  :title="`共${date.tasks.length}个任务`"
                >
                  +{{ date.tasks.length - 3 }}
                </span>
                <span 
                  v-else-if="date.tasks.length > 0"
                  :class="[
                    'text-xs font-medium',
                    date.isSelected ? 'text-white' : 'text-gray-500'
                  ]"
                >
                  {{ date.tasks.length }}
                </span>
              </div>
              
              <!-- 今天的特殊标记 -->
              <div 
                v-if="date.isToday && !date.isSelected"
                class="absolute top-1 right-1"
              >
                <div class="w-2 h-2 bg-gradient-to-br from-purple-500 to-pink-500 rounded-full animate-pulse"></div>
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
            class="bg-gradient-to-r from-blue-500 to-indigo-600 border-b border-blue-600 px-5 py-4 rounded-t-lg shadow-sm"
          >
            <!-- 标题行 -->
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <div class="w-1.5 h-8 bg-white rounded-full"></div>
                <div>
                  <h2 class="font-bold text-xl text-white">
                    {{ selectedDateFormatted }} 
                  </h2>
                  <div class="flex items-center gap-3 mt-1">
                    <span class="text-xs text-blue-100 font-medium">
                      共 {{ selectedDateTasks.length }} 个任务
                    </span>
                    <div class="flex items-center gap-2 text-xs">
                      <div class="flex items-center gap-1">
                        <div class="w-1.5 h-1.5 bg-green-400 rounded-full"></div>
                        <span class="text-blue-100">
                          已完成 {{ selectedDateTasks.filter(t => t.status === 'completed').length }}
                        </span>
                      </div>
                      <div class="w-px h-3 bg-blue-300"></div>
                      <div class="flex items-center gap-1">
                        <div class="w-1.5 h-1.5 bg-orange-400 rounded-full"></div>
                        <span class="text-blue-100">
                          进行中 {{ selectedDateTasks.filter(t => {
                            const today = formatLocalDate(new Date());
                            return t.status !== 'completed' && 
                                   t.startDate <= today && 
                                   today <= t.endDate;
                          }).length }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <!-- 快速添加按钮 -->
              <button
                @click="openTaskModalSelected"
                class="flex items-center gap-2 bg-white text-blue-600 px-4 py-2.5 rounded-lg shadow-md hover:shadow-lg transition-all duration-200 transform hover:scale-105 font-medium"
              >
                <iconify-icon icon="mdi:plus-circle" width="18" height="18"></iconify-icon>
                <span class="text-sm">新建任务</span>
              </button>
            </div>
          </div>

          <!-- 任务列表 -->
          <div class="flex-1 bg-gradient-to-b from-gray-50 to-white rounded-b-lg p-4 overflow-auto">
            <div
              v-if="selectedDateTasks.length === 0"
              class="flex flex-col items-center justify-center py-12 text-center"
            >
              <iconify-icon icon="mdi:calendar-check" width="48" height="48" class="text-gray-300 mb-3"></iconify-icon>
              <p class="text-gray-400 text-sm">该日期暂无任务</p>
              <button
                @click="openTaskModalSelected"
                class="mt-4 text-sm text-blue-600 hover:text-blue-700 font-medium"
              >
                + 添加第一个任务
              </button>
            </div>

            <!-- 任务项 -->
            <div
              v-for="task in selectedDateTasks"
              :key="task.id"
              :class="[
                'group relative rounded-xl p-4 mb-3 border transition-all duration-300',
                'bg-white hover:bg-gray-50',
                'shadow-sm hover:shadow-md',
                task.status === 'completed' ? 'opacity-60 border-gray-200' : 'opacity-100 border-gray-200 hover:border-blue-300'
              ]"
            >
              <!-- 顶部状态栏 -->
              <div class="flex items-center justify-between mb-3">
                <div class="flex items-center gap-2">
                  <!-- 任务分类标签 -->
                  <span
                    :class="[
                      'text-xs px-2.5 py-1 rounded-md font-medium shadow-sm',
                      getCategoryStyle(task.category),
                    ]"
                  >
                    {{ task.category }}
                  </span>
                  <!-- 任务状态标签 -->
                  <span 
                    :class="[
                      'text-xs px-2 py-1 rounded-md font-medium flex items-center gap-1',
                      task.status === 'completed' ? 'bg-green-100 text-green-700' :
                      getTaskActualStatus(task) === '已逾期' ? 'bg-red-100 text-red-700' :
                      getTaskActualStatus(task) === '进行中' ? 'bg-orange-100 text-orange-700' :
                      'bg-gray-100 text-gray-700'
                    ]"
                  >
                    <div 
                      :class="[
                        'w-1.5 h-1.5 rounded-full',
                        getTaskDotColor(task)
                      ]"
                    ></div>
                    {{ getTaskActualStatus(task) }}
                  </span>
                </div>
                <!-- 时间显示 -->
                <span class="text-xs text-gray-500 font-medium flex items-center gap-1">
                  <iconify-icon icon="mdi:clock-outline" width="14" height="14"></iconify-icon>
                  {{ task.time }}
                </span>
              </div>

              <!-- 主要内容区域 -->
              <div class="flex items-start gap-3">
                <!-- 完成状态复选框 -->
                <button
                  type="button"
                  @click.stop="toggleTaskComplete(task)"
                  :class="[
                    'flex-shrink-0 w-6 h-6 rounded-lg border-2 flex items-center justify-center',
                    'transition-all duration-200 hover:scale-110',
                    task.status === 'completed'
                      ? 'bg-gradient-to-br from-green-500 to-green-600 border-green-500 text-white shadow-md'
                      : 'border-gray-300 hover:border-blue-500 bg-white hover:bg-blue-50',
                  ]"
                  :title="task.status === 'completed' ? '标记为未完成' : '标记为已完成'"
                >
                  <svg
                    v-if="task.status === 'completed'"
                    width="14"
                    height="10"
                    viewBox="0 0 12 9"
                    fill="none"
                    class="drop-shadow-sm"
                  >
                    <path
                      d="M10.5 1.5L4.5 7.5L1.5 4.5"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </button>

                <!-- 任务内容 -->
                <div class="flex-1 min-w-0">
                  <!-- 任务标题 -->
                  <h3
                    :class="[
                      'font-bold text-base leading-tight mb-1.5',
                      {
                        'text-gray-800': task.status !== 'completed',
                        'text-gray-500 line-through': task.status === 'completed',
                      },
                    ]"
                  >
                    {{ task.title }}
                  </h3>
                  
                  <!-- 任务描述 -->
                  <p
                    v-if="task.description"
                    :class="[
                      'text-sm leading-relaxed line-clamp-2 mb-2',
                      task.status === 'completed' ? 'text-gray-400' : 'text-gray-600'
                    ]"
                  >
                    {{ task.description }}
                  </p>

                  <!-- 任务日期范围（如果跨多天） -->
                  <div 
                    v-if="task.startDate !== task.endDate" 
                    class="flex items-center gap-1.5 mt-2 text-xs text-blue-600 bg-blue-50 px-2 py-1 rounded-md w-fit"
                  >
                    <iconify-icon icon="mdi:calendar-range" width="14" height="14"></iconify-icon>
                    <span class="font-medium">{{ task.startDate }} ~ {{ task.endDate }}</span>
                  </div>
                </div>

                <!-- 操作按钮组 -->
                <div class="flex-shrink-0 flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                  <!-- 指导按钮 -->
                  <button
                    @click.stop="openGuidanceModal(task)"
                    class="p-2 text-gray-400 hover:text-purple-600 hover:bg-purple-50 rounded-lg transition-colors duration-200"
                    title="获取任务指导"
                  >
                    <iconify-icon icon="mdi:lightbulb-outline" width="16" height="16"></iconify-icon>
                  </button>
                  <!-- 智能测验按钮 -->
                  <button
                    @click.stop="openQuizModal(task)"
                    class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors duration-200"
                    title="智能测验"
                  >
                    <iconify-icon icon="mdi:file-question-outline" width="16" height="16"></iconify-icon>
                  </button>
                  <!-- 编辑按钮 -->
                  <button
                    @click.stop="editTask(task)"
                    class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors duration-200"
                    title="编辑任务"
                  >
                    <iconify-icon icon="mdi:pencil-outline" width="16" height="16"></iconify-icon>
                  </button>
                  <!-- 删除按钮 -->
                  <button
                    @click.stop="handleDelete(task)"
                    class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors duration-200"
                    title="删除任务"
                  >
                    <iconify-icon icon="mdi:trash-can-outline" width="16" height="16"></iconify-icon>
                  </button>
                </div>
              </div>

              <!-- 优先级指示器（右上角） -->
              <div 
                v-if="task.priority && task.priority > 1"
                :class="[
                  'absolute top-2 right-2 w-2 h-2 rounded-full',
                  task.priority >= 3 ? 'bg-red-500' : 'bg-orange-500'
                ]"
                :title="`优先级: ${task.priority >= 3 ? '高' : '中'}`"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 笔记列表 -->
      <div class="mt-8 mb-2">
        <div class="bg-gradient-to-r from-purple-500 to-pink-500 rounded-t-xl p-5 shadow-md">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-1.5 h-8 bg-white rounded-full"></div>
              <div>
                <h2 class="text-xl font-bold text-white flex items-center gap-2">
                  <iconify-icon icon="mdi:notebook-outline" width="24" height="24"></iconify-icon>
                  我的笔记
                </h2>
                <p class="text-xs text-purple-100 mt-1">记录学习点滴，沉淀知识精华</p>
              </div>
            </div>
            <!-- 新建笔记按钮 -->
            <button 
              @click="openNotebookModal({ title: '新笔记', category: '默认', content: '', date: new Date().toLocaleDateString() })" 
              class="bg-white text-purple-600 px-4 py-2 rounded-lg shadow-md hover:shadow-lg transition-all duration-200 transform hover:scale-105 font-medium flex items-center gap-2"
            >
              <iconify-icon icon="mdi:plus-circle" width="18" height="18"></iconify-icon>
              <span class="text-sm">新建笔记</span>
            </button>
          </div>
        </div>

        <div class="bg-white rounded-b-xl border border-t-0 border-gray-200 p-5 shadow-sm">
          <div 
            v-if="notes.length === 0" 
            class="flex flex-col items-center justify-center py-12 text-center"
          >
            <iconify-icon icon="mdi:notebook-outline" width="64" height="64" class="text-gray-300 mb-3"></iconify-icon>
            <p class="text-gray-400 text-sm mb-2">暂无笔记</p>
            <button 
              @click="openNotebookModal({ title: '新笔记', category: '默认', content: '', date: new Date().toLocaleDateString() })"
              class="text-sm text-purple-600 hover:text-purple-700 font-medium"
            >
              + 创建第一篇笔记
            </button>
          </div>

          <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
            <div
              v-for="note in notes"
              :key="note.id"
              class="group bg-gradient-to-br from-white to-gray-50 border-2 border-gray-200 rounded-xl p-5 hover:border-purple-300 hover:shadow-lg transition-all duration-300 cursor-pointer transform hover:-translate-y-1"
              @click="openNotebookModal(note)"
            >
              <!-- 笔记头部 -->
              <div class="flex items-start justify-between mb-3">
                <div class="flex-1 min-w-0">
                  <h3 class="font-bold text-gray-800 text-base mb-1 truncate group-hover:text-purple-600 transition-colors">
                    {{ note.title }}
                  </h3>
                  <div class="flex items-center gap-2 text-xs text-gray-500">
                    <iconify-icon icon="mdi:calendar-outline" width="14" height="14"></iconify-icon>
                    <span>{{ note.date }}</span>
                  </div>
                </div>
                <span
                  :class="[
                    'text-xs px-2.5 py-1 rounded-lg font-medium shadow-sm flex-shrink-0',
                    getCategoryStyle(note.category),
                  ]"
                >
                  {{ note.category }}
                </span>
              </div>

              <!-- 笔记内容预览 -->
              <div 
                v-if="note.content"
                class="text-sm text-gray-600 mb-3 line-clamp-3 leading-relaxed"
                v-html="note.content"
              ></div>
              <p v-else class="text-sm text-gray-400 italic mb-3">暂无内容</p>

              <!-- 笔记底部 -->
              <div class="flex items-center justify-between pt-3 border-t border-gray-200">
                <div class="flex items-center gap-1 text-xs text-gray-500">
                  <iconify-icon icon="mdi:clock-outline" width="14" height="14"></iconify-icon>
                  <span>{{ note.lastUpdated }}</span>
                </div>
                <button class="text-xs text-purple-600 font-medium opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1">
                  查看详情
                  <iconify-icon icon="mdi:arrow-right" width="14" height="14"></iconify-icon>
                </button>
              </div>
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
      class="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click="closeTaskModal"
    >
      <div
        class="bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-hidden flex flex-col animate-modal-enter"
        @click.stop
      >
        <!-- 弹窗头部 -->
        <div class="bg-gradient-to-r from-blue-500 to-indigo-600 px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 rounded-lg flex items-center justify-center backdrop-blur-sm">
                <iconify-icon 
                  :icon="modalDateMode === 'edit' ? 'mdi:pencil' : 'mdi:plus-circle'" 
                  width="20" 
                  height="20"
                  class="text-white"
                ></iconify-icon>
              </div>
              <div>
                <h2 class="text-xl font-bold text-white">
                  {{ modalDateMode === 'edit' ? '编辑任务' : '创建新任务' }}
                </h2>
                <p class="text-xs text-blue-100 mt-0.5">
                  {{ modalDateMode === 'edit' ? '修改任务信息' : '填写任务详细信息' }}
                </p>
              </div>
            </div>
            <button
              @click="closeTaskModal"
              class="w-8 h-8 rounded-lg flex items-center justify-center text-white hover:bg-white/20 transition-colors"
            >
              <iconify-icon icon="mdi:close" width="22" height="22"></iconify-icon>
            </button>
          </div>
        </div>

        <!-- 弹窗内容 -->
        <div class="flex-1 overflow-y-auto p-6">
          <!-- 自然语言输入框 -->
          <div class="mb-6 border-2 border-blue-400 rounded-xl bg-gradient-to-r from-blue-50 to-indigo-50 p-4 shadow-sm">
            <div class="flex items-center gap-2 mb-2">
              <iconify-icon icon="mdi:magic-staff" width="18" height="18" class="text-blue-600"></iconify-icon>
              <label class="text-sm text-blue-700 font-bold">智能解析</label>
            </div>
            <div class="flex gap-2">
              <input
                v-model="naturalLanguageInput"
                type="text"
                class="flex-1 border-0 bg-white px-3 py-2.5 rounded-lg text-sm outline-none focus:ring-2 focus:ring-blue-400 shadow-sm"
                placeholder="例如：明天下午3点完成数学作业第三章"
                :disabled="isParsing"
              />
              <button
                @click="parseNaturalLanguage"
                :disabled="isParsing"
                class="text-white bg-gradient-to-r from-blue-600 to-indigo-600 px-4 py-2.5 rounded-lg text-sm font-medium hover:shadow-lg transition-all duration-200 transform hover:scale-105 flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <iconify-icon :icon="isParsing ? 'mdi:loading' : 'mdi:wand'" width="16" height="16" :class="{ 'animate-spin': isParsing }"></iconify-icon>
                {{ isParsing ? '解析中...' : '解析' }}
                class="text-white bg-gradient-to-r from-blue-600 to-indigo-600 px-4 py-2.5 rounded-lg text-sm font-medium hover:shadow-lg transition-all duration-200 transform hover:scale-105 flex items-center gap-2"
              >
                <iconify-icon icon="mdi:wand" width="16" height="16"></iconify-icon>
                解析
              </button>
            </div>
            <p class="text-xs text-blue-600 mt-2 flex items-center gap-1">
              <iconify-icon icon="mdi:information-outline" width="14" height="14"></iconify-icon>
              支持自然语言输入，AI 将自动填充表单
            </p>
          </div>

          <!-- 表单输入 -->
          <div class="space-y-5">
            <!-- 任务名称 -->
            <div>
              <label class="flex items-center gap-2 text-sm text-gray-700 font-semibold mb-2">
                <iconify-icon icon="mdi:text" width="16" height="16" class="text-blue-600"></iconify-icon>
                任务名称 
                <span class="text-red-500">*</span>
              </label>
              <input
                v-model="newTask.title"
                type="text"
                class="w-full border-2 border-gray-200 px-4 py-2.5 rounded-lg text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-200 focus:outline-none transition-all"
                placeholder="输入任务名称"
              />
            </div>

            <!-- 任务描述 -->
            <div>
              <label class="flex items-center gap-2 text-sm text-gray-700 font-semibold mb-2">
                <iconify-icon icon="mdi:text-box-outline" width="16" height="16" class="text-blue-600"></iconify-icon>
                任务描述
              </label>
              <textarea
                v-model="newTask.description"
                class="w-full border-2 border-gray-200 px-4 py-2.5 rounded-lg text-sm h-24 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 focus:outline-none transition-all resize-none"
                placeholder="详细描述任务内容和目标"
              ></textarea>
            </div>

            <!-- 开始时间 -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="flex items-center gap-2 text-sm text-gray-700 font-semibold mb-2">
                  <iconify-icon icon="mdi:calendar-start" width="16" height="16" class="text-green-600"></iconify-icon>
                  开始日期 
                  <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="newTask.startDate"
                  type="date"
                  class="w-full border-2 border-gray-200 px-4 py-2.5 rounded-lg text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-200 focus:outline-none transition-all"
                />
              </div>
              <div>
                <label class="flex items-center gap-2 text-sm text-gray-700 font-semibold mb-2">
                  <iconify-icon icon="mdi:clock-start" width="16" height="16" class="text-green-600"></iconify-icon>
                  开始时间
                </label>
                <input
                  v-model="newTask.startTime"
                  type="time"
                  class="w-full border-2 border-gray-200 px-4 py-2.5 rounded-lg text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-200 focus:outline-none transition-all"
                />
              </div>
            </div>

            <!-- 结束时间 -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="flex items-center gap-2 text-sm text-gray-700 font-semibold mb-2">
                  <iconify-icon icon="mdi:calendar-end" width="16" height="16" class="text-red-600"></iconify-icon>
                  结束日期 
                  <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="newTask.endDate"
                  type="date"
                  class="w-full border-2 border-gray-200 px-4 py-2.5 rounded-lg text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-200 focus:outline-none transition-all"
                />
              </div>
              <div>
                <label class="flex items-center gap-2 text-sm text-gray-700 font-semibold mb-2">
                  <iconify-icon icon="mdi:clock-end" width="16" height="16" class="text-red-600"></iconify-icon>
                  结束时间
                </label>
                <input
                  v-model="newTask.endTime"
                  type="time"
                  class="w-full border-2 border-gray-200 px-4 py-2.5 rounded-lg text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-200 focus:outline-none transition-all"
                />
              </div>
            </div>

            <!-- 任务分类 -->
            <div>
              <label class="flex items-center gap-2 text-sm text-gray-700 font-semibold mb-2">
                <iconify-icon icon="mdi:tag-outline" width="16" height="16" class="text-blue-600"></iconify-icon>
                任务分类
              </label>
              <select
                v-model="newTask.category"
                class="w-full border-2 border-gray-200 px-4 py-2.5 rounded-lg text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-200 focus:outline-none transition-all bg-white cursor-pointer"
              >
                <option value="">请选择分类</option>
                <option value="study">📚 学习</option>
                <option value="exam">📝 考试</option>
                <option value="project">💼 项目</option>
                <option value="reading">📖 阅读</option>
                <option value="other">📌 其他</option>
              </select>
            </div>
          </div>
        </div>

        <!-- 弹窗底部 -->
        <div class="flex items-center justify-between px-6 py-4 border-t border-gray-200 bg-gray-50">
          <p class="text-xs text-gray-500 flex items-center gap-1">
            <iconify-icon icon="mdi:information-outline" width="14" height="14"></iconify-icon>
            标记 <span class="text-red-500">*</span> 为必填项
          </p>
          <div class="flex gap-3">
            <button
              @click="closeTaskModal"
              class="text-sm text-gray-700 bg-white border-2 border-gray-300 py-2 px-5 rounded-lg hover:bg-gray-50 transition-all font-medium"
            >
              取消
            </button>
            <button
              @click="saveTask"
              class="text-sm text-white bg-gradient-to-r from-blue-600 to-indigo-600 py-2 px-5 rounded-lg hover:shadow-lg transition-all duration-200 transform hover:scale-105 font-medium"
            >
              {{ modalDateMode === 'edit' ? '💾 保存修改' : '✨ 创建任务' }}
            </button>
          </div>
        </div>
        <div class="p-6 text-sm text-gray-700">
          <p>您确定要删除该任务及其关联的所有笔记吗？此操作不可撤销。</p>
        </div>
        <div class="px-6 py-4 border-t border-gray-200 bg-gray-50 flex justify-end gap-3">
          <button @click="cancelDeleteTask" class="text-sm text-gray-700 bg-white border-2 border-gray-300 py-2 px-4 rounded-lg hover:bg-gray-50">取消</button>
          <button @click="confirmDeleteTask" class="text-sm text-white bg-gradient-to-r from-red-600 to-pink-600 py-2 px-4 rounded-lg hover:shadow-lg">确认</button>
        </div>
      </div>
    </div>

    <!-- 笔记本弹窗 -->
    <div
      v-if="showNotebookModal && currentNote"
      class="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
    >
      <div
        class="bg-white rounded-2xl shadow-2xl w-full max-w-4xl h-[85vh] flex flex-col animate-modal-enter"
      >
        <!-- 笔记头部 -->
        <div class="bg-gradient-to-r from-purple-500 to-pink-500 px-6 py-4 rounded-t-2xl">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 rounded-lg flex items-center justify-center backdrop-blur-sm">
                <iconify-icon 
                  :icon="currentNote && currentNote.id ? 'mdi:pencil' : 'mdi:notebook-plus'" 
                  width="20" 
                  height="20"
                  class="text-white"
                ></iconify-icon>
              </div>
              <div class="flex items-center gap-3">
                <div>
                  <h2 class="text-xl font-bold text-white">
                    {{ currentNote && currentNote.id ? "编辑笔记" : "新建笔记" }}
                  </h2>
                  <p class="text-xs text-purple-100 mt-0.5">
                    记录灵感，积累知识
                  </p>
                </div>
                <span
                  v-if="currentNote && currentNote.category"
                  :class="[
                    'px-3 py-1 rounded-lg text-xs font-medium bg-white/90 backdrop-blur-sm shadow-sm',
                    getCategoryStyle(currentNote.category),
                  ]"
                >
                  {{ currentNote.category }}
                </span>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <button
                @click="toggleNotebookFullscreen"
                class="w-8 h-8 rounded-lg flex items-center justify-center text-white hover:bg-white/20 transition-colors"
                :title="isNotebookFullscreen ? '退出全屏' : '全屏显示'"
              >
                <iconify-icon
                  :icon="isNotebookFullscreen ? 'mdi:fullscreen-exit' : 'mdi:fullscreen'"
                  width="20"
                  height="20"
                ></iconify-icon>
              </button>
              <button
                @click="closeNotebookModal"
                class="w-8 h-8 rounded-lg flex items-center justify-center text-white hover:bg-white/20 transition-colors"
                aria-label="关闭笔记"
              >
                <iconify-icon icon="mdi:close" width="22"></iconify-icon>
              </button>
            </div>
          </div>
        </div>
        <!-- 笔记内容 -->
        <div class="flex-1 overflow-hidden flex flex-col p-6 bg-gray-50">
          <div class="flex-1 overflow-y-auto pr-2 space-y-5">
            <!-- 笔记标题 -->
            <div>
              <label class="flex items-center gap-2 text-sm text-gray-700 font-semibold mb-2">
                <iconify-icon icon="mdi:format-title" width="16" height="16" class="text-purple-600"></iconify-icon>
                笔记标题
              </label>
              <input
                v-model="currentNote.title"
                type="text"
                class="w-full border-2 border-gray-200 px-4 py-3 rounded-lg text-base font-medium focus:border-purple-500 focus:ring-2 focus:ring-purple-200 focus:outline-none transition-all"
                placeholder="为你的笔记起个标题"
              />
            </div>

            <!-- 笔记分类 -->
            <div>
              <label class="flex items-center gap-2 text-sm text-gray-700 font-semibold mb-2">
                <iconify-icon icon="mdi:tag-outline" width="16" height="16" class="text-purple-600"></iconify-icon>
                笔记分类
              </label>
              <select
                v-model="currentNote.category"
                class="w-full border-2 border-gray-200 px-4 py-3 rounded-lg text-sm focus:border-purple-500 focus:ring-2 focus:ring-purple-200 focus:outline-none transition-all bg-white cursor-pointer"
              >
                <option value="学习">📚 学习</option>
                <option value="工作">💼 工作</option>
                <option value="数学">🔢 数学</option>
                <option value="英语">🗣️ 英语</option>
                <option value="物理">⚛️ 物理</option>
                <option value="研究">🔬 研究</option>
                <option value="其他">📌 其他</option>
              </select>
            </div>

            <!-- 笔记内容 -->
            <div>
              <label class="flex items-center gap-2 text-sm text-gray-700 font-semibold mb-2">
                <iconify-icon icon="mdi:text-box-outline" width="16" height="16" class="text-purple-600"></iconify-icon>
                笔记内容
              </label>
              <div
                v-if="editor"
                class="border-2 border-gray-200 rounded-xl focus-within:border-purple-500 focus-within:ring-2 focus-within:ring-purple-200 transition-all shadow-sm bg-white"
              >
                <div
                  class="flex items-center p-3 border-b border-gray-200 bg-gradient-to-r from-gray-50 to-gray-100 rounded-t-xl flex-wrap gap-2"
                >
                  <button
                    @click="editor.chain().focus().toggleBold().run()"
                    :class="{
                      'bg-purple-600 text-white shadow-md': editor.isActive('bold'),
                      'text-gray-700 hover:bg-gray-200': !editor.isActive('bold'),
                    }"
                    class="p-2 rounded-lg transition-all"
                    aria-label="加粗"
                    title="加粗 (Ctrl+B)"
                  >
                    <iconify-icon icon="mdi:format-bold" width="18"></iconify-icon>
                  </button>
                  <button
                    @click="editor.chain().focus().toggleItalic().run()"
                    :class="{
                      'bg-purple-600 text-white shadow-md': editor.isActive('italic'),
                      'text-gray-700 hover:bg-gray-200': !editor.isActive('italic'),
                    }"
                    class="p-2 rounded-lg transition-all"
                    aria-label="斜体"
                    title="斜体 (Ctrl+I)"
                  >
                    <iconify-icon icon="mdi:format-italic" width="18"></iconify-icon>
                  </button>
                  <button
                    @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
                    :class="{
                      'bg-purple-600 text-white shadow-md': editor.isActive('heading', { level: 2 }),
                      'text-gray-700 hover:bg-gray-200': !editor.isActive('heading', { level: 2 }),
                    }"
                    class="p-2 rounded-lg transition-all"
                    aria-label="二级标题"
                    title="二级标题"
                  >
                    <iconify-icon icon="mdi:format-header-2" width="18"></iconify-icon>
                  </button>
                  <button
                    @click="addImage"
                    class="p-2 rounded-lg text-gray-700 hover:bg-gray-200 transition-all"
                    aria-label="插入图片"
                    title="插入图片"
                  >
                    <iconify-icon icon="mdi:image-plus" width="18"></iconify-icon>
                  </button>
                  <div class="border-l-2 border-gray-300 h-6 mx-1"></div>
                  <button
                    @click="editor.chain().focus().undo().run()"
                    :disabled="!editor.can().undo()"
                    class="p-2 rounded-lg transition-all"
                    :class="{
                      'text-gray-400 cursor-not-allowed': !editor.can().undo(),
                      'text-gray-700 hover:bg-gray-200': editor.can().undo(),
                    }"
                    aria-label="撤销"
                    title="撤销 (Ctrl+Z)"
                  >
                    <iconify-icon icon="mdi:undo" width="18"></iconify-icon>
                  </button>
                  <button
                    @click="editor.chain().focus().redo().run()"
                    :disabled="!editor.can().redo()"
                    class="p-2 rounded-lg transition-all"
                    :class="{
                      'text-gray-400 cursor-not-allowed': !editor.can().redo(),
                      'text-gray-700 hover:bg-gray-200': editor.can().redo(),
                    }"
                    aria-label="重做"
                    title="重做 (Ctrl+Y)"
                  >
                    <iconify-icon icon="mdi:redo" width="18"></iconify-icon>
                  </button>
                </div>
                <editor-content
                  :editor="editor"
                  class="p-5 min-h-[350px] bg-white rounded-b-xl focus:outline-none prose prose-sm max-w-none"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 笔记底部 -->
        <div class="flex items-center justify-between px-6 py-4 border-t-2 border-gray-200 bg-gradient-to-r from-gray-50 to-white rounded-b-2xl">
          <div class="flex items-center gap-2 text-sm text-gray-500">
            <iconify-icon icon="mdi:clock-outline" width="16" height="16"></iconify-icon>
            <span>最后更新: {{ currentNote.lastUpdated }}</span>
          </div>
          <div class="flex gap-3">
            <button
              @click="addImage"
              class="text-sm text-purple-600 py-2 px-4 border-2 border-purple-600 rounded-lg hover:bg-purple-50 transition-all font-medium flex items-center gap-2"
            >
              <iconify-icon icon="mdi:image-plus" width="16" height="16"></iconify-icon>
              插入图片
            </button>
            <button
              @click="closeAndSaveNote"
              class="text-sm text-white bg-gradient-to-r from-purple-600 to-pink-600 py-2 px-5 rounded-lg hover:shadow-lg transition-all duration-200 transform hover:scale-105 font-medium flex items-center gap-2"
            >
              <iconify-icon icon="mdi:content-save" width="16" height="16"></iconify-icon>
              保存并关闭
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 任务指导弹窗 -->
    <div
      v-if="showGuidanceModal"
      class="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click="closeGuidanceModal"
    >
      <div
        class="bg-white rounded-2xl shadow-2xl w-full max-w-2xl max-h-[85vh] overflow-hidden flex flex-col animate-modal-enter"
        @click.stop
      >
        <!-- 弹窗头部 -->
        <div class="bg-gradient-to-r from-purple-500 to-pink-600 px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 rounded-lg flex items-center justify-center backdrop-blur-sm">
                <iconify-icon icon="mdi:lightbulb-on" width="20" height="20" class="text-white"></iconify-icon>
              </div>
              <div>
                <h3 class="text-lg font-bold text-white">任务指导</h3>
                <p class="text-purple-100 text-sm">{{ guidanceTask?.title }}</p>
              </div>
            </div>
            <button
              @click="closeGuidanceModal"
              class="w-8 h-8 rounded-lg bg-white/20 hover:bg-white/30 flex items-center justify-center transition-colors"
            >
              <iconify-icon icon="mdi:close" width="20" height="20" class="text-white"></iconify-icon>
            </button>
          </div>
        </div>

        <!-- 弹窗内容 -->
        <div class="flex-1 overflow-y-auto p-6">
          <!-- 加载状态 -->
          <div v-if="isLoadingGuidance" class="flex flex-col items-center justify-center py-12">
            <iconify-icon icon="mdi:loading" width="48" height="48" class="text-purple-500 animate-spin"></iconify-icon>
            <p class="text-gray-500 mt-4">AI 正在生成任务指导...</p>
          </div>

          <!-- 指导内容 -->
          <div v-else-if="taskGuidance" class="space-y-6">
            <!-- 执行步骤 -->
            <div class="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl p-5 border border-blue-200">
              <div class="flex items-center gap-2 mb-4">
                <iconify-icon icon="mdi:format-list-numbered" width="20" height="20" class="text-blue-600"></iconify-icon>
                <h4 class="font-bold text-blue-800">执行步骤</h4>
              </div>
              <ol class="space-y-3">
                <li v-for="(step, index) in taskGuidance.steps" :key="index" class="flex items-start gap-3">
                  <span class="flex-shrink-0 w-6 h-6 bg-blue-500 text-white rounded-full flex items-center justify-center text-sm font-bold">
                    {{ index + 1 }}
                  </span>
                  <span class="text-gray-700 leading-relaxed">{{ step }}</span>
                </li>
              </ol>
            </div>

            <!-- 学习技巧 -->
            <div class="bg-gradient-to-r from-green-50 to-emerald-50 rounded-xl p-5 border border-green-200">
              <div class="flex items-center gap-2 mb-4">
                <iconify-icon icon="mdi:lightbulb" width="20" height="20" class="text-green-600"></iconify-icon>
                <h4 class="font-bold text-green-800">学习技巧</h4>
              </div>
              <ul class="space-y-2">
                <li v-for="(tip, index) in taskGuidance.tips" :key="index" class="flex items-start gap-2">
                  <iconify-icon icon="mdi:check-circle" width="18" height="18" class="text-green-500 mt-0.5 flex-shrink-0"></iconify-icon>
                  <span class="text-gray-700">{{ tip }}</span>
                </li>
              </ul>
            </div>

            <!-- 时间建议 -->
            <div class="bg-gradient-to-r from-orange-50 to-amber-50 rounded-xl p-5 border border-orange-200">
              <div class="flex items-center gap-2 mb-3">
                <iconify-icon icon="mdi:clock-outline" width="20" height="20" class="text-orange-600"></iconify-icon>
                <h4 class="font-bold text-orange-800">时间建议</h4>
              </div>
              <p class="text-gray-700">{{ taskGuidance.timeAdvice }}</p>
            </div>

            <!-- 相关资源 -->
            <div class="bg-gradient-to-r from-purple-50 to-pink-50 rounded-xl p-5 border border-purple-200">
              <div class="flex items-center gap-2 mb-4">
                <iconify-icon icon="mdi:link-variant" width="20" height="20" class="text-purple-600"></iconify-icon>
                <h4 class="font-bold text-purple-800">相关资源</h4>
              </div>
              <div class="grid grid-cols-1 gap-3">
                <a
                  v-for="(resource, index) in taskGuidance.resources"
                  :key="index"
                  :href="resource.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="flex items-center gap-3 p-3 bg-white rounded-lg border border-purple-200 hover:border-purple-400 hover:shadow-md transition-all group"
                >
                  <div :class="[
                    'w-10 h-10 rounded-lg flex items-center justify-center',
                    resource.type === 'video' ? 'bg-red-100' :
                    resource.type === 'course' ? 'bg-blue-100' :
                    resource.type === 'article' ? 'bg-green-100' :
                    'bg-gray-100'
                  ]">
                    <iconify-icon 
                      :icon="
                        resource.type === 'video' ? 'mdi:play-circle' :
                        resource.type === 'course' ? 'mdi:school' :
                        resource.type === 'article' ? 'mdi:file-document' :
                        'mdi:tools'
                      " 
                      width="20" 
                      height="20"
                      :class="[
                        resource.type === 'video' ? 'text-red-600' :
                        resource.type === 'course' ? 'text-blue-600' :
                        resource.type === 'article' ? 'text-green-600' :
                        'text-gray-600'
                      ]"
                    ></iconify-icon>
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="font-medium text-gray-800 group-hover:text-purple-600 truncate">{{ resource.title }}</p>
                    <p class="text-xs text-gray-500 truncate">{{ resource.url }}</p>
                  </div>
                  <iconify-icon icon="mdi:open-in-new" width="16" height="16" class="text-gray-400 group-hover:text-purple-600"></iconify-icon>
                </a>
              </div>
            </div>
          </div>

          <!-- 无数据状态 -->
          <div v-else class="flex flex-col items-center justify-center py-12">
            <iconify-icon icon="mdi:alert-circle-outline" width="48" height="48" class="text-gray-400"></iconify-icon>
            <p class="text-gray-500 mt-4">暂无指导信息</p>
          </div>
        </div>

        <!-- 底部按钮 -->
        <div class="px-6 py-4 border-t border-gray-200 bg-gray-50">
          <div class="flex justify-end gap-3">
            <button
              @click="refreshGuidance"
              :disabled="isLoadingGuidance"
              class="px-4 py-2 text-purple-600 border border-purple-300 rounded-lg hover:bg-purple-50 transition-colors flex items-center gap-2 disabled:opacity-50"
            >
              <iconify-icon icon="mdi:refresh" width="16" height="16" :class="{ 'animate-spin': isLoadingGuidance }"></iconify-icon>
              重新生成
            </button>
            <button
              @click="closeGuidanceModal"
              class="px-4 py-2 bg-gradient-to-r from-purple-500 to-pink-600 text-white rounded-lg hover:shadow-lg transition-all"
            >
              关闭
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 智能测验弹窗 -->
    <div
      v-if="showQuizModal"
      class="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click="closeQuizModal"
    >
      <div
        class="bg-white rounded-2xl shadow-2xl w-full max-w-3xl max-h-[85vh] overflow-hidden flex flex-col animate-modal-enter"
        @click.stop
      >
        <!-- 弹窗头部 -->
        <div class="bg-gradient-to-r from-blue-500 to-cyan-600 px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 rounded-lg flex items-center justify-center backdrop-blur-sm">
                <iconify-icon icon="mdi:file-question" width="20" height="20" class="text-white"></iconify-icon>
              </div>
              <div>
                <h3 class="text-lg font-bold text-white">智能测验</h3>
                <p class="text-blue-100 text-sm">{{ quizTask?.title }}</p>
              </div>
            </div>
            <button
              @click="closeQuizModal"
              class="w-8 h-8 rounded-lg bg-white/20 hover:bg-white/30 flex items-center justify-center transition-colors"
            >
              <iconify-icon icon="mdi:close" width="20" height="20" class="text-white"></iconify-icon>
            </button>
          </div>
        </div>

        <!-- 弹窗内容 -->
        <div class="flex-1 overflow-y-auto p-6">
          <!-- 加载状态 -->
          <div v-if="isLoadingQuiz" class="flex flex-col items-center justify-center py-12">
            <iconify-icon icon="mdi:loading" width="48" height="48" class="text-blue-500 animate-spin"></iconify-icon>
            <p class="text-gray-500 mt-4">AI 正在生成智能测验...</p>
          </div>

          <!-- 测验内容 -->
          <div v-else-if="quiz" class="space-y-6">
            <!-- 选择题部分 -->
            <div v-for="(question, index) in quiz.questions" :key="index" class="bg-white rounded-xl p-5 border-2 border-gray-200 hover:border-blue-300 transition-colors">
              <div class="flex items-start gap-3 mb-4">
                <span class="flex-shrink-0 w-8 h-8 bg-blue-500 text-white rounded-full flex items-center justify-center font-bold">
                  {{ index + 1 }}
                </span>
                <div class="flex-1">
                  <p class="text-gray-800 font-medium leading-relaxed">{{ question.question }}</p>
                  <span :class="[
                    'inline-block mt-2 px-2 py-1 text-xs rounded-full',
                    question.difficulty === 'easy' ? 'bg-green-100 text-green-700' :
                    question.difficulty === 'medium' ? 'bg-yellow-100 text-yellow-700' :
                    'bg-red-100 text-red-700'
                  ]">
                    {{ question.difficulty === 'easy' ? '简单' : question.difficulty === 'medium' ? '中等' : '困难' }}
                  </span>
                </div>
              </div>

              <!-- 选项 -->
              <div class="space-y-2 ml-11">
                <label
                  v-for="option in ['A', 'B', 'C', 'D']"
                  :key="option"
                  :class="[
                    'flex items-start gap-3 p-3 rounded-lg border-2 cursor-pointer transition-all',
                    userAnswers[index] === option
                      ? quizSubmitted
                        ? option === question.correctAnswer
                          ? 'border-green-500 bg-green-50'
                          : 'border-red-500 bg-red-50'
                        : 'border-blue-500 bg-blue-50'
                      : quizSubmitted && option === question.correctAnswer
                        ? 'border-green-500 bg-green-50'
                        : 'border-gray-200 hover:border-blue-300 hover:bg-blue-50/50'
                  ]"
                >
                  <input
                    type="radio"
                    :name="`question-${index}`"
                    :value="option"
                    v-model="userAnswers[index]"
                    :disabled="quizSubmitted"
                    class="mt-1 w-4 h-4 text-blue-600"
                  />
                  <span class="flex-1 text-gray-700">
                    <span class="font-semibold">{{ option }}.</span> {{ question.options[option] }}
                  </span>
                  <iconify-icon
                    v-if="quizSubmitted && option === question.correctAnswer"
                    icon="mdi:check-circle"
                    width="20"
                    height="20"
                    class="text-green-600 mt-1"
                  ></iconify-icon>
                  <iconify-icon
                    v-else-if="quizSubmitted && userAnswers[index] === option && option !== question.correctAnswer"
                    icon="mdi:close-circle"
                    width="20"
                    height="20"
                    class="text-red-600 mt-1"
                  ></iconify-icon>
                </label>
              </div>

              <!-- 答案解析 -->
              <div v-if="quizSubmitted" class="ml-11 mt-4 p-4 bg-blue-50 rounded-lg border border-blue-200">
                <div class="flex items-center gap-2 mb-2">
                  <iconify-icon icon="mdi:lightbulb-on" width="18" height="18" class="text-blue-600"></iconify-icon>
                  <span class="font-semibold text-blue-800">答案解析</span>
                </div>
                <p class="text-gray-700 leading-relaxed">{{ question.explanation }}</p>
              </div>
            </div>

            <!-- 问答题部分 -->
            <div v-if="quiz.essayQuestion" class="bg-white rounded-xl p-5 border-2 border-gray-200 hover:border-purple-300 transition-colors">
              <div class="flex items-start gap-3 mb-4">
                <span class="flex-shrink-0 w-8 h-8 bg-purple-500 text-white rounded-full flex items-center justify-center font-bold">
                  {{ quiz.questions.length + 1 }}
                </span>
                <div class="flex-1">
                  <p class="text-gray-800 font-medium leading-relaxed">{{ quiz.essayQuestion.question }}</p>
                  <span class="inline-block mt-2 px-2 py-1 text-xs rounded-full bg-purple-100 text-purple-700">
                    问答题
                  </span>
                </div>
              </div>

              <!-- 文本框 -->
              <div class="pl-11">
                <textarea
                  v-model="essayAnswer"
                  :disabled="quizSubmitted"
                  rows="6"
                  placeholder="请在此输入你的答案..."
                  class="w-full p-4 border-2 border-gray-200 rounded-lg focus:border-purple-500 focus:ring-2 focus:ring-purple-200 outline-none transition-all resize-none disabled:bg-gray-50"
                ></textarea>
              </div>

              <!-- 学习建议 -->
              <div v-if="quizSubmitted" class="pl-11 mt-4">
                <div class="p-4 bg-purple-50 rounded-lg border border-purple-200">
                  <div class="flex items-center gap-2 mb-2">
                    <iconify-icon icon="mdi:school" width="18" height="18" class="text-purple-600"></iconify-icon>
                    <span class="font-semibold text-purple-800">学习建议</span>
                  </div>
                  <p class="text-gray-700 leading-relaxed">{{ quiz.essayQuestion.studySuggestion }}</p>
                </div>
              </div>
            </div>

            <!-- 测验结果 -->
            <div v-if="quizSubmitted" class="bg-gradient-to-r from-blue-50 to-cyan-50 rounded-xl p-6 border border-blue-200">
              <div class="flex items-center justify-between">
                <div>
                  <div class="flex items-center gap-2 mb-2">
                    <iconify-icon icon="mdi:trophy" width="24" height="24" class="text-yellow-500"></iconify-icon>
                    <h4 class="text-lg font-bold text-gray-800">测验完成</h4>
                  </div>
                  <p class="text-gray-600">
                    选择题得分: <span class="font-bold text-blue-600">{{ quizScore }}</span> / {{ quiz.questions.length }}
                  </p>
                  <p class="text-sm text-gray-500 mt-1">
                    正确率: {{ ((quizScore / quiz.questions.length) * 100).toFixed(1) }}%
                  </p>
                </div>
                <div class="text-right">
                  <div :class="[
                    'text-5xl font-bold',
                    (quizScore / quiz.questions.length) >= 0.8 ? 'text-green-500' :
                    (quizScore / quiz.questions.length) >= 0.6 ? 'text-yellow-500' :
                    'text-red-500'
                  ]">
                    {{ ((quizScore / quiz.questions.length) * 100).toFixed(0) }}
                  </div>
                  <div class="text-sm text-gray-500 mt-1">分</div>
                </div>
              </div>
            </div>
          </div>

          <!-- 无数据状态 -->
          <div v-else class="flex flex-col items-center justify-center py-12">
            <iconify-icon icon="mdi:alert-circle-outline" width="48" height="48" class="text-gray-400"></iconify-icon>
            <p class="text-gray-500 mt-4">暂无测验信息</p>
          </div>
        </div>

        <!-- 底部按钮 -->
        <div class="px-6 py-4 border-t border-gray-200 bg-gray-50">
          <div class="flex justify-between items-center">
            <button
              v-if="quizSubmitted"
              @click="regenerateQuiz"
              :disabled="isLoadingQuiz"
              class="px-4 py-2 text-blue-600 border border-blue-300 rounded-lg hover:bg-blue-50 transition-colors flex items-center gap-2 disabled:opacity-50"
            >
              <iconify-icon icon="mdi:refresh" width="16" height="16" :class="{ 'animate-spin': isLoadingQuiz }"></iconify-icon>
              重新生成
            </button>
            <div v-else></div>
            <div class="flex gap-3">
              <button
                v-if="!quizSubmitted"
                @click="submitQuiz"
                :disabled="!canSubmitQuiz"
                class="px-6 py-2 bg-gradient-to-r from-blue-500 to-cyan-600 text-white rounded-lg hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                提交答案
              </button>
              <button
                @click="closeQuizModal"
                class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
              >
                关闭
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useEditor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import Image from "@tiptap/extension-image";
import { createTask, getPersonalTasks, completeTask, uncompleteTask, deleteTask, parseTaskWithAI, getTaskGuidance, generateQuiz } from "@/api/modules/task";

// Name
defineOptions({
  name: "PersonalTasks",
});

// 响应式数据
const currentDate = ref(new Date());
const selectedDate = ref(new Date());
const showTaskModal = ref(false);
const showNotebookModal = ref(false);
const showCompleteConfirm = ref(false);
const confirmingTask = ref(null);
const showDeleteConfirm = ref(false);
const deletingTask = ref(null);
const isNotebookFullscreen = ref(false);
const naturalLanguageInput = ref("");
const statusFilter = ref(null);
const modalDateMode = ref('system');
const tasks = ref([]);
const notes = ref([]);

// 任务指导相关
const showGuidanceModal = ref(false);
const guidanceTask = ref(null);
const taskGuidance = ref(null);
const isLoadingGuidance = ref(false);

// 智能测验相关
const showQuizModal = ref(false);
const quizTask = ref(null);
const quiz = ref(null);
const isLoadingQuiz = ref(false);
const userAnswers = ref([]);
const essayAnswer = ref('');
const quizSubmitted = ref(false);
const quizScore = ref(0);

const newTask = ref({
  title: "",
  description: "",
  startDate: "",
  startTime: "",
  endDate: "",
  endTime: "",
  category: "",
});
const teamTasks = ref([]);
const teamTasksLoading = ref(false);
const teamTasksError = ref("");
const TEAM_TASK_PREVIEW_LIMIT = 4;
const teamTaskPreview = computed(() => teamTasks.value.slice(0, TEAM_TASK_PREVIEW_LIMIT));

// 当前笔记
const currentNote = ref(null);
const isNoteDirty = ref(false);
const isNoteSaving = ref(false);

const loadNotes = async () => {
  try {
    const res = await getStudyNotes();
    if (res && (res.code === 0 || res.code === 200)) {
      const items = res.data || res.items || [];
      notes.value = (items || []).map((n) => ({
        id: n.id,
        title: n.title,
        content: n.content || "",
        category: "学习",
        date: n.created_at ? new Date(n.created_at).toLocaleString("zh-CN") : "",
        lastUpdated: n.updated_at ? new Date(n.updated_at).toLocaleString("zh-CN") : "",
        taskId: n.task_id || null,
      }));
    }
  } catch (e) {
    console.error("加载笔记失败", e);
  }
};

const editor = useEditor({
  content: "",
  extensions: [StarterKit, Image],
  editable: true,
  onUpdate: ({ editor }) => {
    if (currentNote.value) {
      const next = editor.getHTML();
      if (currentNote.value.content !== next) {
        currentNote.value.content = next;
        isNoteDirty.value = true;
      }
    }
  },
});

const addImage = () => {
  const url = window.prompt("请输入图片URL");
  if (url && editor.value) {
    editor.value.chain().focus().setImage({ src: url }).run();
  }
};

watch(
  () => currentNote.value,
  (newNote) => {
    if (newNote && editor.value) {
      const currentContent = editor.value.getHTML();
      if (currentContent !== newNote.content) {
        editor.value.commands.setContent(newNote.content || "");
      }
    }
  },
  { deep: true }
);

watch(
  () => currentNote.value && currentNote.value.title,
  () => {
    if (currentNote.value) {
      isNoteDirty.value = true;
    }
  }
);

const normalizeTeamTaskStatus = (status) => {
  const numeric = Number(status);
  if (!Number.isNaN(numeric)) {
    if (numeric === 2) return "completed";
    if (numeric === 1) return "in-progress";
    return "pending";
  }
  const lowered = String(status ?? "").toLowerCase();
  if (lowered.includes("complete")) return "completed";
  if (lowered.includes("progress") || lowered.includes("doing")) return "in-progress";
  return "pending";
};

const getTeamTaskStatusLabel = (status) => {
  const map = {
    completed: "已完成",
    "in-progress": "进行中",
    pending: "待处理",
  };
  return map[status] || "待处理";
};

const getTeamTaskBadgeClass = (status) => {
  if (status === "completed") return "bg-green-100 text-green-700";
  if (status === "in-progress") return "bg-orange-100 text-orange-700";
  return "bg-gray-100 text-gray-700";
};

const clampProgress = (value) => {
  const n = Number(value);
  if (Number.isNaN(n)) return 0;
  return Math.max(0, Math.min(100, Math.round(n)));
};

const calcTeamTaskProgress = (task) => {
  if (typeof task?.progress === "number" && !Number.isNaN(task.progress)) {
    return clampProgress(task.progress);
  }
  const status = normalizeTeamTaskStatus(task?.status);
  if (status === "completed") return 100;
  if (status === "in-progress") return 60;
  return 0;
};

const formatISODate = (value) => {
  if (!value) return "";
  try {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "";
    return date.toISOString().split("T")[0];
  } catch (error) {
    console.warn("格式化日期失败：", error);
    return "";
  }
};

const loadTeamTasks = async () => {
  teamTasksLoading.value = true;
  teamTasksError.value = "";
  try {
    const response = await getTeamTasks();
    let rawList = [];
    if (Array.isArray(response?.data)) {
      rawList = response.data;
    } else if (Array.isArray(response?.items)) {
      rawList = response.items;
    } else if (Array.isArray(response)) {
      rawList = response;
    }
    teamTasks.value = rawList
      .filter((task) => task.owner_team_id)
      .map((task) => ({
        id: task.id,
        title: task.title,
        description: task.description,
        teamId: task.owner_team_id,
        dueDate: formatISODate(task.due_at || task.due_date),
        status: normalizeTeamTaskStatus(task.status),
        progress: calcTeamTaskProgress(task),
      }));
  } catch (error) {
    console.error("加载团队任务失败:", error);
    teamTasksError.value = error?.message || "加载团队任务失败";
  } finally {
    teamTasksLoading.value = false;
  }
};

// 计算属性
const currentMonthYear = computed(() => {
  return currentDate.value.toLocaleString("default", {
    month: "long",
    year: "numeric",
  });
});

const selectedDateFormatted = computed(() => {
  return selectedDate.value
    ? selectedDate.value.toLocaleDateString()
    : "未选择日期";
});

const formatLocalDate = (d) => {
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
};

const calendarDates = computed(() => {
  const year = currentDate.value.getFullYear();
  const month = currentDate.value.getMonth();
  const firstDay = new Date(year, month, 1);
  const startDate = new Date(firstDay);
  startDate.setDate(startDate.getDate() - firstDay.getDay());

  const dates = [];
  const today = new Date();

  for (let i = 0; i < 42; i++) {
    const date = new Date(startDate);
    date.setDate(startDate.getDate() + i);

    const dateString = formatLocalDate(date);
    // 修改：任务在其持续周期内的所有日期都显示
    const dateTasks = tasks.value.filter((t) => {
      const taskStart = t.startDate || t.date;
      const taskEnd = t.endDate || t.date;
      return dateString >= taskStart && dateString <= taskEnd;
    });

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
  const dateStr = formatLocalDate(selectedDate.value);
  // 修改：显示在任务持续周期内的所有任务
  return tasks.value.filter((task) => {
    const taskStart = task.startDate || task.date;
    const taskEnd = task.endDate || task.date;
    return dateStr >= taskStart && dateStr <= taskEnd;
  });
});

const filteredTasksByStatus = computed(() => {
  if (!statusFilter.value) return [];
  
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const todayStr = formatLocalDate(today);
  
  return tasks.value.filter((task) => {
    const taskStartDate = task.startDate || task.date;
    const taskEndDate = task.endDate || task.date;
    
    // 根据实际状态过滤
    if (statusFilter.value === 'completed') {
      return task.status === 'completed';
    } else if (statusFilter.value === 'overdue') {
      return task.status !== 'completed' && taskEndDate < todayStr;
    } else if (statusFilter.value === 'in-progress') {
      return task.status !== 'completed' && taskStartDate <= todayStr && todayStr <= taskEndDate;
    } else if (statusFilter.value === 'pending') {
      return task.status !== 'completed' && taskStartDate > todayStr;
    }
    return false;
  });
});

const stats = computed(() => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const todayStr = formatLocalDate(today);
  
  const total = tasks.value.length;
  
  // 使用动态计算的实际状态进行统计
  let completed = 0;
  let inProgress = 0;
  let pending = 0;
  let overdue = 0;
  
  tasks.value.forEach((task) => {
    const taskStartDate = task.startDate || task.date;
    const taskEndDate = task.endDate || task.date;
    
    if (task.status === 'completed') {
      completed++;
    } else if (taskEndDate < todayStr) {
      overdue++;
    } else if (taskStartDate <= todayStr && todayStr <= taskEndDate) {
      inProgress++;
    } else if (taskStartDate > todayStr) {
      pending++;
    }
  });
  
  return { total, completed, inProgress, pending, overdue };
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
    const today = new Date();
    selectedDate.value = today;
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

const closeTaskModal = () => {
  showTaskModal.value = false;
};

// 任务指导相关函数
const openGuidanceModal = async (task) => {
  guidanceTask.value = task;
  showGuidanceModal.value = true;
  await fetchTaskGuidance();
};

const closeGuidanceModal = () => {
  showGuidanceModal.value = false;
  guidanceTask.value = null;
  taskGuidance.value = null;
};

const fetchTaskGuidance = async () => {
  if (!guidanceTask.value) return;
  
  isLoadingGuidance.value = true;
  try {
    const response = await getTaskGuidance(
      guidanceTask.value.title,
      guidanceTask.value.description || '',
      guidanceTask.value.category || 'other'
    );
    console.log('任务指导响应:', response);
    
    // response 已经是拦截器处理后的 data 对象
    if (response && response.data) {
      taskGuidance.value = response.data;
    } else if (response && response.steps) {
      // 如果 response 直接就是指导数据
      taskGuidance.value = response;
    }
  } catch (error) {
    console.error('获取任务指导失败:', error);
  } finally {
    isLoadingGuidance.value = false;
  }
};

const refreshGuidance = () => {
  fetchTaskGuidance();
};

// 智能测验相关函数
const openQuizModal = (task) => {
  quizTask.value = task;
  showQuizModal.value = true;
  userAnswers.value = [];
  essayAnswer.value = '';
  quizSubmitted.value = false;
  quizScore.value = 0;
  quiz.value = null;
  fetchQuiz();
};

const closeQuizModal = () => {
  showQuizModal.value = false;
  quizTask.value = null;
  quiz.value = null;
  userAnswers.value = [];
  essayAnswer.value = '';
  quizSubmitted.value = false;
  quizScore.value = 0;
};

const fetchQuiz = async () => {
  if (!quizTask.value) return;
  
  isLoadingQuiz.value = true;
  try {
    const response = await generateQuiz({
      topic: quizTask.value.title,
      content: quizTask.value.description || '',
      difficulty: 'medium', // 可以根据任务难度动态调整
      quizCount: 3,
      includeEssay: true
    });
    console.log('智能测验响应:', response);
    
    // response 已经是拦截器处理后的 data 对象
    if (response && response.data) {
      quiz.value = response.data;
      // 初始化答案数组
      userAnswers.value = new Array(quiz.value.questions.length).fill(null);
    } else if (response && response.questions) {
      // 如果 response 直接就是测验数据
      quiz.value = response;
      userAnswers.value = new Array(quiz.value.questions.length).fill(null);
    }
  } catch (error) {
    console.error('生成智能测验失败:', error);
  } finally {
    isLoadingQuiz.value = false;
  }
};

const regenerateQuiz = () => {
  userAnswers.value = [];
  essayAnswer.value = '';
  quizSubmitted.value = false;
  quizScore.value = 0;
  fetchQuiz();
};

const canSubmitQuiz = computed(() => {
  if (!quiz.value) return false;
  // 检查所有选择题是否都已作答
  const allAnswered = userAnswers.value.every(answer => answer !== null && answer !== undefined);
  return allAnswered;
});

const submitQuiz = () => {
  if (!canSubmitQuiz.value) return;
  
  // 计算得分
  let score = 0;
  quiz.value.questions.forEach((question, index) => {
    if (userAnswers.value[index] === question.correctAnswer) {
      score++;
    }
  });
  
  quizScore.value = score;
  quizSubmitted.value = true;
  
  // 滚动到顶部查看结果
  const modalContent = document.querySelector('.overflow-y-auto');
  if (modalContent) {
    modalContent.scrollTop = 0;
  }
};

const isParsing = ref(false);

const parseNaturalLanguage = async () => {
  const input = (naturalLanguageInput.value || "").trim();
  if (!input) return;

  isParsing.value = true;
  try {
    const response = await parseTaskWithAI(input);
    console.log('AI解析响应:', response);
    
    if (response.code === 0 && response.data) {
      const parsed = response.data;
      if (parsed.title) newTask.value.title = parsed.title;
      if (parsed.description) newTask.value.description = parsed.description;
      if (parsed.startDate) newTask.value.startDate = parsed.startDate;
      if (parsed.startTime) newTask.value.startTime = parsed.startTime;
      if (parsed.endDate) newTask.value.endDate = parsed.endDate;
      if (parsed.endTime) newTask.value.endTime = parsed.endTime;
      if (parsed.category) newTask.value.category = parsed.category;
    }
  } catch (error) {
    console.error('AI解析失败:', error);
  } finally {
    isParsing.value = false;
  }
};

// 加载个人任务
const loadPersonalTasks = async () => {
  try {
    const response = await getPersonalTasks();
    console.log('API响应:', response); // 调试日志
    if (response.code === 0) {
      // 将API返回的任务转换为前端格式
      const apiTasks = response.data || [];
      console.log('API返回的任务数量:', apiTasks.length); // 调试日志
      tasks.value = apiTasks.map(task => ({
        id: task.id,
        title: task.title,
        description: task.description,
        date: task.start_at ? new Date(task.start_at).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
        startDate: task.start_at ? new Date(task.start_at).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
        endDate: task.due_at ? new Date(task.due_at).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
        time: task.start_at ? new Date(task.start_at).toLocaleTimeString('zh-CN', {hour: '2-digit', minute: '2-digit'}) : "全天",
        endTime: task.due_at ? new Date(task.due_at).toLocaleTimeString('zh-CN', {hour: '2-digit', minute: '2-digit'}) : "全天",
        status: task.status === 2 ? "completed" : "pending", // 2=已完成, 1=进行中, 0=待处理
        notes: "",
        category: task.category?.name || "其他",
      }));
      console.log('转换后的任务数据:', tasks.value); // 调试日志
    }
  } catch (error) {
    console.error('加载任务失败:', error);
  }
};

const saveTask = async () => {
  if (!newTask.value.title || !newTask.value.startDate || !newTask.value.endDate) {
    return;
  }
  
  try {
    // 1. 先拼接成一个本地时间字符串
    const localStartStr = newTask.value.startTime 
      ? `${newTask.value.startDate}T${newTask.value.startTime}:00`
      : `${newTask.value.startDate}T09:00:00`;
    
    const localEndStr = newTask.value.endTime 
      ? `${newTask.value.endDate}T${newTask.value.endTime}:00`
      : `${newTask.value.endDate}T18:00:00`;

    // 2. ✅ 关键修改：转换为标准 ISO 8601 格式 (带时区)
    // 比如：它会把 "2025-11-28T09:00:00" 变成 "2025-11-28T01:00:00.000Z"
    const isoStartTime = new Date(localStartStr).toISOString();
    const isoEndTime = new Date(localEndStr).toISOString();
    
    // 准备API数据
    const taskData = {
      title: newTask.value.title,
      description: newTask.value.description,
      task_type: 1, // 个人任务
      priority: 1, // 默认优先级
      effort_points: 5, // 默认工作量
      start_at: isoStartTime, // ✅ 发送标准格式
      due_at: isoEndTime,     // ✅ 发送标准格式
      // 如果有分类，可以设置 category_id
    };
    
    let response;
    
    // 判断是编辑还是新建
    if (modalDateMode.value === 'edit' && newTask.value.id) {
      // 编辑现有任务 - 这里需要后端提供更新API
      // response = await updateTask(newTask.value.id, taskData);
      
      // 暂时使用前端更新
      const taskIndex = tasks.value.findIndex(t => t.id === newTask.value.id);
      if (taskIndex !== -1) {
        tasks.value[taskIndex] = {
          ...tasks.value[taskIndex],
          title: newTask.value.title,
          description: newTask.value.description,
          startDate: newTask.value.startDate,
          endDate: newTask.value.endDate,
          time: newTask.value.endTime || "全天",
          category: newTask.value.category || "其他",
        };
        
        closeTaskModal();
        naturalLanguageInput.value = "";
        modalDateMode.value = 'system';
        alert("✅ 任务已更新");
        return;
      }
    } else {
      // 调用API创建任务
      response = await createTask(taskData);
    }
    
    if (response && response.code === 0) {
      // 将API返回的任务转换为前端格式
      const apiTask = response.data;
      const task = {
        id: apiTask.id,
        title: apiTask.title,
        description: apiTask.description,
        date: newTask.value.startDate, // 保留兼容性
        startDate: newTask.value.startDate, // 任务开始日期
        endDate: newTask.value.endDate, // 任务结束日期
        time: newTask.value.startTime || "09:00",
        endTime: newTask.value.endTime || "18:00",
        status: "pending", // 转换状态
        notes: "",
        category: newTask.value.category || "其他",
      };
      
      // 添加到本地任务列表
      tasks.value.push(task);
      
      closeTaskModal();
      naturalLanguageInput.value = "";
      modalDateMode.value = 'system';
    } else {
      console.error('创建任务失败:', response);
      alert('创建任务失败，请重试');
    }
  } catch (error) {
    console.error('保存任务失败:', error);
    alert('保存任务失败，请检查网络连接');
  }
};

const handleDelete = async (task) => {
  if(!confirm("确定要删除此任务吗？")) {
    return;
  }
  try {
    const res = await deleteTask(task.id);
    if(res.code === 0){
      tasks.value = tasks.value.filter(t => t.id !== task.id);
      alert("✅ 任务已删除");
    }
  } catch (error) {
    console.error('删除任务失败:', error);
    alert('删除任务失败，请检查网络连接');
  }
};

const editTask = (task) => {
  // 填充表单数据
  newTask.value = {
    id: task.id,
    title: task.title,
    description: task.description,
    startDate: task.startDate,
    startTime: task.time !== "全天" ? task.time.split('-')[0]?.trim() || "" : "",
    endDate: task.endDate,
    endTime: task.time !== "全天" ? task.time.split('-')[1]?.trim() || task.time : "",
    category: task.category,
  };
  
  modalDateMode.value = 'edit';
  showTaskModal.value = true;
};

const toggleTaskComplete = async (task) => {
  try {
    const wasCompleted = task.status === "completed";
    if (wasCompleted) {
      const response = await uncompleteTask(task.id);
      if (response.code === 0) {
        task.status = "pending";
      } else {
        throw new Error(response.msg || "取消完成失败");
      }
    } else {
      confirmingTask.value = task;
      showCompleteConfirm.value = true;
    }
  } catch (error) {
    console.error("更新任务状态失败:", error);
    alert("更新任务状态失败，请重试");
    await loadPersonalTasks();
  }
};

const confirmCompleteWithNote = async () => {
  if (!confirmingTask.value) return;
  try {
    const response = await completeTaskWithNote(confirmingTask.value.id);
    if (response.code === 0) {
      confirmingTask.value.status = "completed";
      const note = response.data?.note || response.note || response.data;
      if (note) {
        currentNote.value = {
          id: note.id,
          title: note.title,
          content: note.content || "",
          category: confirmingTask.value.category || "学习",
          date: note.created_at ? new Date(note.created_at).toLocaleString("zh-CN") : new Date().toLocaleString("zh-CN"),
          lastUpdated: new Date().toLocaleString("zh-CN"),
          taskId: confirmingTask.value.id,
        };
        notes.value.push({ ...currentNote.value });
        showNotebookModal.value = true;
        setTimeout(() => {
          if (editor.value) {
            editor.value.chain().focus().run();
          }
        }, 0);
      }
    } else {
      throw new Error(response.msg || "完成任务并创建笔记失败");
    }
  } catch (e) {
    console.error("完成并创建笔记失败:", e);
    alert("笔记创建失败，任务状态已回滚。请稍后重试");
    await loadPersonalTasks();
  } finally {
    showCompleteConfirm.value = false;
    confirmingTask.value = null;
  }
};

const cancelCompleteWithoutNote = async () => {
  if (!confirmingTask.value) return;
  try {
    const response = await completeTask(confirmingTask.value.id);
    if (response.code === 0) {
      confirmingTask.value.status = "completed";
    } else {
      throw new Error(response.msg || "完成任务失败");
    }
  } catch (e) {
    console.error("仅完成任务失败:", e);
    alert("完成任务失败，请重试");
    await loadPersonalTasks();
  } finally {
    showCompleteConfirm.value = false;
    confirmingTask.value = null;
  }
};

const openNotebookModal = (note = null) => {
  if (note) {
    currentNote.value = { ...note };
  } else {
    currentNote.value = {
      id: null,
      title: "新笔记",
      content: "",
      category: "默认",
      tags: [],
      lastUpdated: "",
    };
  }
  showNotebookModal.value = true;
  isNoteDirty.value = false;
};

const closeNotebookModal = () => {
  showNotebookModal.value = false;
  currentNote.value = null;
};

const toggleNotebookFullscreen = () => {
  isNotebookFullscreen.value = !isNotebookFullscreen.value;
};

const saveNote = async () => {
  if (!currentNote.value) return false;
  if (isNoteSaving.value) return false;
  isNoteSaving.value = true;
  try {
    const payload = {
      title: currentNote.value.title || "",
      content: currentNote.value.content || "",
    };
    let savedId = currentNote.value.id;
    if (currentNote.value.id) {
      const res = await updateStudyNote(currentNote.value.id, payload);
      const updated = res?.data || res;
      if (updated?.updated_at) {
        currentNote.value.lastUpdated = new Date(updated.updated_at).toLocaleString("zh-CN");
      }
    } else {
      const res = await createStudyNote({
        title: payload.title,
        content: payload.content,
        task_id: currentNote.value.taskId || undefined,
      });
      const created = res?.data || res?.note || res;
      savedId = created?.id || savedId;
      if (savedId) {
        currentNote.value.id = savedId;
      }
      if (created?.updated_at) {
        currentNote.value.lastUpdated = new Date(created.updated_at).toLocaleString("zh-CN");
      }
    }
    const idx = notes.value.findIndex((n) => n.id === currentNote.value.id);
    if (idx !== -1) {
      notes.value[idx] = { ...currentNote.value };
    } else {
      notes.value.push({ ...currentNote.value });
    }
    isNoteDirty.value = false;
    // ElMessage.success("笔记已保存");
    return true;
  } catch (e) {
    console.error("保存笔记失败", e);
    ElMessage.error(e?.message || "保存失败，请稍后重试");
    return false;
  } finally {
    isNoteSaving.value = false;
  }
};

const closeAndSaveNote = async () => {
  const ok = await saveNote();
  if (ok) {
    closeNotebookModal();
  }
};


const getTaskNote = (taskId) => {
  return notes.value.find(n => n.taskId === taskId);
};

const getRelatedTask = (taskId) => {
  if (!taskId) return null;
  return tasks.value.find((t) => t.id === taskId);
};

const getCategoryStyle = (category) => {
  const styles = {
    数学: "bg-blue-50 text-blue-600",
    英语: "bg-orange-50 text-orange-600",
    物理: "bg-red-50 text-red-600",
    研究: "bg-purple-50 text-purple-600",
    学习: "bg-blue-50 text-blue-600",
    工作: "bg-teal-50 text-teal-600",
    其他: "bg-gray-50 text-gray-600",
    study: "bg-blue-50 text-blue-600",
    exam: "bg-red-50 text-red-600",
    project: "bg-purple-50 text-purple-600",
    reading: "bg-green-50 text-green-600",
    other: "bg-gray-50 text-gray-600",
  };
  return styles[category] || "bg-gray-50 text-gray-600";
};

const getTaskCardBackground = (category) => {
  const style = getCategoryStyle(category);
  return style.split(' ')[0];
};

const getStatusText = (status) => {
  const statusMap = {
    'completed': '已完成',
    'in-progress': '进行中',
    'pending': '待处理',
    'overdue': '已逾期',
  };
  return statusMap[status] || '未知状态';
};

// 动态计算任务的实际状态和颜色
const getTaskDotColor = (task) => {
  const today = new Date();
  today.setHours(0, 0, 0, 0); // 设置为当天的开始时间
  
  const todayStr = formatLocalDate(today);
  const taskStartDate = task.startDate || task.date;
  const taskEndDate = task.endDate || task.date;
  
  // 如果任务已完成，永远显示绿色
  if (task.status === 'completed') {
    return 'bg-green-500';
  }
  
  // 如果任务结束时间小于当前时间且未完成，显示红色（已逾期）
  if (taskEndDate < todayStr) {
    return 'bg-red-500';
  }
  
  // 如果当前时间处于任务起始时间和结束时间之间且未完成，显示橙色（进行中）
  if (taskStartDate <= todayStr && todayStr <= taskEndDate) {
    return 'bg-orange-500';
  }
  
  // 如果任务起始时间晚于当前时间，显示灰色（待处理）
  if (taskStartDate > todayStr) {
    return 'bg-gray-500';
  }
  
  // 默认紫色
  return 'bg-gray-500';
};

// 获取任务的实际状态文本（用于tooltip）
const getTaskActualStatus = (task) => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  
  const todayStr = formatLocalDate(today);
  const taskStartDate = task.startDate || task.date;
  const taskEndDate = task.endDate || task.date;
  
  if (task.status === 'completed') {
    return '已完成';
  }
  
  if (taskEndDate < todayStr) {
    return '已逾期';
  }
  
  if (taskStartDate <= todayStr && todayStr <= taskEndDate) {
    return '进行中';
  }
  
  if (taskStartDate > todayStr) {
    return '待处理';
  }
  
  return '未知状态';
};

const formatTaskEndTime = (task) => {
  if (!task.endDate) return task.endTime;
  const [year, month, day] = task.endDate.split('-');
  return `${parseInt(month)}月${parseInt(day)}日 ${task.endTime}`;
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
    pending: "待处理",
    overdue: "已逾期",
  };
  return map[status] || "任务";
};

// 初始化
onMounted(async () => {
  // 为测试目的设置一个mock token
  if (!localStorage.getItem('token')) {
    localStorage.setItem('token', 'mock-token-3-test');
  }
  
  // 加载个人任务
  await loadPersonalTasks();

  // 加载团队任务
  await loadTeamTasks();

  // 加载学习笔记
  await loadNotes();
});

watch(selectedDate, (d) => {
  if (!showTaskModal.value || !d) return;
  if (modalDateMode.value !== 'selected') return;
  const ds = formatLocalDate(d);
  newTask.value.startDate = ds;
  newTask.value.endDate = ds;
});
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
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }

  .task-dot:hover {
    transform: scale(1.3);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .stat-card {
    transition: transform 0.2s;
  }

  /* 弹窗动画 */
  @keyframes modal-enter {
    from {
      opacity: 0;
      transform: scale(0.95) translateY(-10px);
    }
    to {
      opacity: 1;
      transform: scale(1) translateY(0);
    }
  }

  .animate-modal-enter {
    animation: modal-enter 0.3s ease-out;
  }

  /* 富文本编辑器样式优化 */
  :deep(.ProseMirror) {
    outline: none;
  }

  :deep(.ProseMirror p) {
    margin: 0.75em 0;
  }

  :deep(.ProseMirror h2) {
    font-size: 1.5em;
    font-weight: bold;
    margin-top: 1em;
    margin-bottom: 0.5em;
    color: #374151;
  }

  :deep(.ProseMirror img) {
    max-width: 100%;
    height: auto;
    border-radius: 0.5rem;
    margin: 1em 0;
  }

  :deep(.ProseMirror strong) {
    font-weight: 700;
    color: #1f2937;
  }

  :deep(.ProseMirror em) {
    font-style: italic;
    color: #4b5563;
  }
</style>
