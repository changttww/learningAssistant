<template>
  <div class="w-full min-h-full flex flex-col px-4">
    <div
      v-if="!selectedTeam"
      class="relative w-full flex-1 min-h-[80vh] flex flex-col items-center justify-center bg-gradient-to-br from-slate-50 to-blue-50 rounded-xl"
    >
      <!-- 粒子背景 Canvas -->
      <canvas
        ref="particleCanvas"
        class="absolute inset-0 w-full h-full pointer-events-none"
      ></canvas>

      <div class="relative z-10 w-full flex flex-col items-center py-10">
        <div v-if="loadingTeams" class="text-lg text-gray-600 animate-pulse">
          正在加载团队数据...
        </div>
        <div
          v-else-if="allTeams.length === 0"
          class="bg-white/80 backdrop-blur-sm p-8 rounded-xl shadow-xl text-center max-w-md w-full border border-white/50 transform transition-all hover:scale-105 duration-300"
        >
          <div
            class="mb-6 bg-blue-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto text-blue-600"
          >
            <iconify-icon
              icon="mdi:account-group-outline"
              width="32"
              height="32"
            ></iconify-icon>
          </div>
          <h2 class="text-2xl font-bold mb-4 text-gray-800">
            欢迎加入团队协作
          </h2>
          <p class="text-gray-600 mb-8">
            您还没有加入任何团队。开启您的团队协作之旅吧！
          </p>
          <div class="flex gap-3 mb-6">
            <input
              v-model="teamNameInput"
              type="text"
              class="flex-1 p-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all bg-gray-50 focus:bg-white"
              placeholder="输入团队名称"
              @keyup.enter="handleJoinTeam"
            />
            <button
              @click="handleJoinTeam"
              class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 active:bg-blue-800 transition-colors shadow-lg shadow-blue-200 font-medium"
            >
              加入
            </button>
          </div>
          <div class="text-sm text-gray-500">
            或者
            <button
              @click="showCreateTeamModal = true"
              class="text-blue-600 hover:text-blue-700 font-medium hover:underline"
            >
              创建一个新团队
            </button>
          </div>
        </div>
        <div v-else class="w-full max-w-6xl px-4">
          <div class="text-center mb-10">
            <h2 class="text-3xl font-bold text-gray-800 mb-2">选择您的团队</h2>
            <p class="text-gray-500">点击卡片进入团队工作区</p>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <div
              v-for="team in allTeams"
              :key="team.id"
              @click="selectTeam(team)"
              class="bg-white/90 backdrop-blur-sm p-6 rounded-xl shadow-sm hover:shadow-xl transition-all cursor-pointer border border-gray-100 hover:border-blue-400 group transform hover:-translate-y-1 duration-300"
            >
              <div class="flex items-center justify-between mb-4">
                <div class="flex items-center gap-3">
                  <div
                    class="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-bold text-lg shadow-md"
                  >
                    {{ team.name.charAt(0).toUpperCase() }}
                  </div>
                  <h3
                    class="text-xl font-bold text-gray-800 group-hover:text-blue-600 transition-colors"
                  >
                    {{ team.name }}
                  </h3>
                </div>
                <span
                  :class="[
                    'text-xs px-2 py-1 rounded-full font-medium',
                    team.owner_user_id === currentUserId
                      ? 'bg-blue-100 text-blue-700'
                      : 'bg-gray-100 text-gray-600',
                  ]"
                >
                  {{ team.owner_user_id === currentUserId ? "创建者" : "成员" }}
                </span>
              </div>
              <p class="text-gray-500 text-sm line-clamp-2 mb-4 h-10">
                {{ team.description || "暂无描述" }}
              </p>

              <div
                class="grid grid-cols-2 gap-y-2 gap-x-4 text-xs text-gray-500 mb-4"
              >
                <div class="flex items-center gap-1.5" title="创建者">
                  <iconify-icon
                    icon="mdi:account-tie"
                    class="text-blue-400 text-sm"
                  ></iconify-icon>
                  <span class="truncate max-w-[80px]">{{
                    team.owner_name || "未知"
                  }}</span>
                </div>
                <div class="flex items-center gap-1.5" title="团队人数">
                  <iconify-icon
                    icon="mdi:account-group"
                    class="text-purple-400 text-sm"
                  ></iconify-icon>
                  <span>{{ team.member_count || 1 }} 人</span>
                </div>
                <div
                  class="flex items-center gap-1.5 col-span-2"
                  title="创建时间"
                >
                  <iconify-icon
                    icon="mdi:clock-outline"
                    class="text-gray-400 text-sm"
                  ></iconify-icon>
                  <span>{{ formatDate(team.created_at) }}</span>
                </div>
              </div>

              <div
                class="flex items-center justify-between text-xs text-gray-400 pt-4 border-t border-gray-100"
              >
                <span>ID: {{ team.id }}</span>
                <span
                  class="group-hover:translate-x-1 transition-transform text-blue-500 flex items-center gap-1"
                >
                  进入团队 <iconify-icon icon="mdi:arrow-right"></iconify-icon>
                </span>
              </div>
            </div>

            <!-- 操作卡片 -->
            <div
              class="bg-white/60 backdrop-blur-sm p-6 rounded-xl shadow-sm hover:shadow-md transition-all cursor-pointer border-2 border-dashed border-gray-300 hover:border-blue-400 hover:bg-blue-50/50 flex flex-col items-center justify-center min-h-[180px] group"
              @click="showJoinModal = true"
            >
              <div
                class="w-12 h-12 rounded-full bg-gray-100 group-hover:bg-blue-100 flex items-center justify-center mb-3 transition-colors text-gray-400 group-hover:text-blue-600"
              >
                <iconify-icon
                  icon="mdi:account-plus"
                  width="24"
                  height="24"
                ></iconify-icon>
              </div>
              <span class="text-gray-600 font-medium group-hover:text-blue-600"
                >加入其他团队</span
              >
              <span class="text-xs text-gray-400 mt-1"
                >输入名称加入现有团队</span
              >
            </div>

            <div
              class="bg-white/60 backdrop-blur-sm p-6 rounded-xl shadow-sm hover:shadow-md transition-all cursor-pointer border-2 border-dashed border-gray-300 hover:border-purple-400 hover:bg-purple-50/50 flex flex-col items-center justify-center min-h-[180px] group"
              @click="showCreateTeamModal = true"
            >
              <div
                class="w-12 h-12 rounded-full bg-gray-100 group-hover:bg-purple-100 flex items-center justify-center mb-3 transition-colors text-gray-400 group-hover:text-purple-600"
              >
                <iconify-icon
                  icon="mdi:plus-box"
                  width="24"
                  height="24"
                ></iconify-icon>
              </div>
              <span
                class="text-gray-600 font-medium group-hover:text-purple-600"
                >创建新团队</span
              >
              <span class="text-xs text-gray-400 mt-1">建立属于你的团队</span>
            </div>
          </div>
        </div>
      </div>

      <div
        v-if="showJoinModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm transition-opacity"
      >
        <div
          class="bg-white p-8 rounded-2xl shadow-2xl w-full max-w-md transform transition-all scale-100"
        >
          <h3
            class="text-xl font-bold mb-6 text-gray-800 flex items-center gap-2"
          >
            <iconify-icon
              icon="mdi:account-plus"
              class="text-blue-600"
            ></iconify-icon>
            加入团队
          </h3>
          <input
            v-model="teamNameInput"
            type="text"
            class="w-full p-3 border border-gray-200 rounded-lg mb-6 focus:ring-2 focus:ring-blue-500 outline-none bg-gray-50"
            placeholder="请输入准确的团队名称"
          />
          <div class="flex justify-end gap-3">
            <button
              @click="showJoinModal = false"
              class="px-5 py-2.5 text-gray-600 hover:bg-gray-100 rounded-lg transition-colors font-medium"
            >
              取消
            </button>
            <button
              @click="handleJoinTeam"
              class="px-5 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 shadow-lg shadow-blue-200 font-medium transition-all active:scale-95"
            >
              确认加入
            </button>
          </div>
        </div>
      </div>

      <div
        v-if="showCreateTeamModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm transition-opacity"
      >
        <div
          class="bg-white p-8 rounded-2xl shadow-2xl w-full max-w-md transform transition-all scale-100"
        >
          <h3
            class="text-xl font-bold mb-6 text-gray-800 flex items-center gap-2"
          >
            <iconify-icon
              icon="mdi:plus-box"
              class="text-purple-600"
            ></iconify-icon>
            创建团队
          </h3>
          <div class="space-y-5">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5"
                >团队名称</label
              >
              <input
                v-model="newTeamForm.name"
                type="text"
                class="w-full p-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-purple-500 outline-none bg-gray-50"
                placeholder="给团队起个响亮的名字"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5"
                >团队描述</label
              >
              <textarea
                v-model="newTeamForm.description"
                rows="3"
                class="w-full p-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-purple-500 outline-none bg-gray-50"
                placeholder="简单介绍一下团队的目标（可选）"
              ></textarea>
            </div>
          </div>
          <div class="flex justify-end gap-3 mt-8">
            <button
              @click="showCreateTeamModal = false"
              class="px-5 py-2.5 text-gray-600 hover:bg-gray-100 rounded-lg transition-colors font-medium"
            >
              取消
            </button>
            <button
              @click="handleCreateTeam"
              class="px-5 py-2.5 bg-purple-600 text-white rounded-lg hover:bg-purple-700 shadow-lg shadow-purple-200 font-medium transition-all active:scale-95"
            >
              立即创建
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div
        class="lg:col-span-3 flex items-center mb-6 bg-white p-4 rounded-xl shadow-sm border border-gray-100"
      >
        <button
          @click="selectedTeam = null"
          class="flex items-center px-4 py-2 text-sm font-medium text-gray-600 bg-gray-50 hover:bg-blue-50 hover:text-blue-600 rounded-lg transition-all duration-200 group"
        >
          <iconify-icon
            icon="mdi:arrow-left"
            class="mr-2 transition-transform group-hover:-translate-x-1"
          ></iconify-icon>
          切换团队
        </button>
        <div class="h-6 w-px bg-gray-200 mx-4"></div>
        <div class="flex items-center gap-3">
          <div
            class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-bold text-sm shadow-sm"
          >
            {{ selectedTeam.name.charAt(0).toUpperCase() }}
          </div>
          <span class="font-bold text-gray-800 text-lg">{{
            selectedTeam.name
          }}</span>
        </div>
        <div class="ml-auto flex items-center gap-2">
          <button
            @click="openMembersModal"
            class="flex items-center px-3 py-1.5 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors"
          >
            <iconify-icon icon="mdi:account-group" class="mr-1"></iconify-icon>
            团队成员
          </button>
          <button
            @click="openInviteModal"
            class="flex items-center px-3 py-1.5 text-sm font-medium text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors"
          >
            <iconify-icon icon="mdi:account-plus" class="mr-1"></iconify-icon>
            邀请成员
          </button>
        </div>
      </div>
      <!-- 左侧主要内容 -->
      <div class="lg:col-span-2 space-y-6">
        <div class="card surface-card">
          <div class="flex items-center justify-between mb-6">
            <h1 class="page-title">团队任务</h1>
            <div class="flex items-center gap-3">
              <button
                @click="goToConstellation"
                class="ghost-btn whitespace-nowrap"
              >
                任务概览
              </button>
              <button
                @click="openCreateModal"
                class="primary-btn flex items-center gap-2"
              >
                <iconify-icon
                  icon="mdi:plus"
                  width="16"
                  height="16"
                ></iconify-icon>
                创建任务
              </button>
            </div>
          </div>

          <!-- 创建任务模态框 -->
          <div
            v-if="showCreateModal"
            class="fixed inset-0 z-50 flex items-center justify-center"
          >
            <div
              class="absolute inset-0 bg-black opacity-40"
              @click="closeCreateModal"
            ></div>
            <div class="bg-white rounded-lg shadow-lg z-10 w-full max-w-md p-6">
              <h3 class="text-lg font-semibold mb-4">创建新任务</h3>
              <div class="space-y-3">
                <div>
                  <label for="create-title" class="text-sm">任务名称</label>
                  <input
                    id="create-title"
                    v-model="newTaskForm.title"
                    type="text"
                    class="w-full mt-1 p-2 border rounded"
                    placeholder="请输入任务名称"
                  />
                </div>
                <div>
                  <label for="create-desc" class="text-sm">描述</label>
                  <textarea
                    id="create-desc"
                    v-model="newTaskForm.description"
                    rows="3"
                    class="w-full mt-1 p-2 border rounded"
                    placeholder="任务描述（可选）"
                  ></textarea>
                </div>
                <div class="flex gap-2">
                  <div class="flex-1">
                    <label for="create-due" class="text-sm">到期日</label>
                    <input
                      id="create-due"
                      v-model="newTaskForm.due_date"
                      type="date"
                      class="w-full mt-1 p-2 border rounded"
                    />
                  </div>
                  <div class="flex-1">
                    <label for="create-team" class="text-sm"
                      >所属团队（可选）</label
                    >
                    <select
                      id="create-team"
                      v-model="newTaskForm.owner_team_id"
                      class="w-full mt-1 p-2 border rounded"
                      :disabled="ownedTeamsLoading || !ownedTeams.length"
                    >
                      <option value="">未关联团队</option>
                      <option
                        v-for="team in ownedTeams"
                        :key="team.id"
                        :value="team.id"
                      >
                        {{ team.name }}
                      </option>
                    </select>
                    <p
                      v-if="ownedTeamsLoading"
                      class="text-xs text-gray-400 mt-1"
                    >
                      加载团队列表中...
                    </p>
                    <p
                      v-else-if="!ownedTeams.length"
                      class="text-xs text-gray-400 mt-1"
                    >
                      暂无你拥有的团队
                    </p>
                  </div>
                </div>
                <div>
                  <label for="create-points" class="text-sm">任务积分</label>
                  <input
                    id="create-points"
                    v-model.number="newTaskForm.effort_points"
                    type="number"
                    min="0"
                    class="w-full mt-1 p-2 border rounded"
                    placeholder="完成任务可获得的积分"
                  />
                </div>
                <div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm">子任务</span>
                    <button
                      type="button"
                      class="text-xs text-blue-600"
                      @click="addSubtaskField"
                    >
                      + 添加子任务
                    </button>
                  </div>
                  <div class="space-y-2 mt-2">
                    <div
                      v-for="(subtask, index) in newTaskForm.subtasks"
                      :key="`new-subtask-${index}`"
                      class="flex items-center gap-2"
                    >
                      <input
                        v-model="newTaskForm.subtasks[index]"
                        type="text"
                        class="flex-1 p-2 border rounded"
                        placeholder="子任务标题"
                      />
                      <button
                        type="button"
                        class="text-xs text-gray-500 hover:text-red-500"
                        @click="removeSubtaskField(index)"
                      >
                        删除
                      </button>
                    </div>
                  </div>
                </div>
              </div>
              <div class="mt-4 flex justify-end gap-2">
                <button
                  @click="closeCreateModal"
                  class="px-4 py-2 rounded border"
                >
                  取消
                </button>
                <button
                  :disabled="creating"
                  @click="submitCreateTask"
                  class="px-4 py-2 rounded bg-[#2D5BFF] text-white"
                >
                  {{ creating ? "创建中..." : "创建" }}
                </button>
              </div>
            </div>
          </div>

          <!-- 团队任务进度图表 -->
          <div class="mb-6">
            <h2 class="section-title">团队进度概览</h2>
            <div class="chart-container" ref="teamProgressChart"></div>
          </div>

          <!-- 任务列表 -->
          <div class="space-y-4">
            <h3 class="section-title">当前任务</h3>

            <div
              v-for="task in tasks"
              :key="task.id"
              :data-task-id="task.id"
              :class="taskCardClass(task)"
            >
              <div class="flex items-start justify-between gap-4">
                <div class="flex-1">
                  <h4 class="task-title">{{ task.title }}</h4>
                  <p class="task-desc">{{ task.description }}</p>
                  <div class="task-meta">
                    <div class="meta-item">
                      <iconify-icon
                        icon="mdi:account"
                        width="16"
                        height="16"
                        class="text-gray-400"
                      ></iconify-icon>
                      <span>{{ task.owner_name || "未知" }}</span>
                    </div>
                    <div class="meta-item">
                      <iconify-icon
                        icon="mdi:calendar"
                        width="16"
                        height="16"
                        class="text-gray-400"
                      ></iconify-icon>
                      <span>{{ task.due_date || task.created_at || "" }}</span>
                    </div>
                  </div>
                  <div class="mt-3">
                    <div class="flex items-center justify-between mb-1">
                      <span class="text-sm">进度</span>
                      <span class="text-sm font-medium">{{
                        computeTaskProgressLabel(task)
                      }}</span>
                    </div>
                    <div class="progress-bar">
                      <div
                        class="progress-fill"
                        :style="progressFillStyle(task)"
                      ></div>
                    </div>
                    <div class="flex flex-wrap items-center gap-2 mt-2">
                      <button
                        @click="openProgressModal(task)"
                        class="ghost-btn"
                      >
                        更新进度
                      </button>
                      <button @click="setCompleted(task)" class="success-btn">
                        标记完成
                      </button>
                      <button
                        @click="toggleTaskDetails(task.id)"
                        class="info-btn"
                      >
                        {{ isTaskExpanded(task.id) ? "收起详情" : "查看详情" }}
                      </button>
                    </div>
                  </div>
                </div>
                <div class="flex flex-col items-end gap-2">
                  <span
                    v-if="isCreatedByCurrentUser(task)"
                    class="badge badge-info"
                    >我创建</span
                  >
                  <span :class="['badge', statusBadgeClass(task)]">{{
                    task.status_label || task.status
                  }}</span>
                  <span
                    v-if="taskHealth(task)"
                    :class="['badge', getHealthBadgeClass(taskHealth(task))]"
                    >{{ taskHealth(task).label }}</span
                  >
                </div>
              </div>
              <div
                v-if="isTaskExpanded(task.id)"
                class="mt-4 rounded-lg border border-dashed border-gray-200 p-4 space-y-4"
              >
                <div
                  v-if="isDetailLoading(task.id)"
                  class="text-xs text-gray-400"
                >
                  详情加载中...
                </div>
                <div>
                  <div class="flex items-center justify-between mb-2">
                    <h4 class="text-sm font-semibold text-gray-700">子任务</h4>
                    <span class="text-xs text-gray-400"
                      >{{ getTaskDetail(task).subtasks.length }} 个</span
                    >
                  </div>
                  <ul
                    v-if="getTaskDetail(task).subtasks.length"
                    class="space-y-1"
                  >
                    <li
                      v-for="sub in getTaskDetail(task).subtasks"
                      :key="sub.id"
                      class="flex items-center justify-between text-sm"
                    >
                      <span>{{ sub.title }}</span>
                      <span class="text-xs text-gray-500">{{
                        sub.status
                      }}</span>
                    </li>
                  </ul>
                  <p v-else class="text-xs text-gray-400">暂无子任务</p>
                </div>
                <div>
                  <div class="flex items-center justify-between mb-2">
                    <h4 class="text-sm font-semibold text-gray-700">评论</h4>
                    <span class="text-xs text-gray-400"
                      >{{ getTaskDetail(task).comments.length }} 条</span
                    >
                  </div>
                  <ul
                    v-if="getTaskDetail(task).comments.length"
                    class="space-y-2 text-sm mb-3"
                  >
                    <li
                      v-for="(comment, idx) in getTaskDetail(task).comments"
                      :key="idx"
                      class="bg-gray-50 p-2 rounded"
                    >
                      <div class="flex items-center justify-between mb-1">
                        <span class="font-medium text-xs text-blue-600">{{ comment.user_id === currentUserId ? '我' : `用户 ${comment.user_id}` }}</span>
                        <span class="text-xs text-gray-400">{{
                          formatDate(comment.created_at)
                        }}</span>
                      </div>
                      <p class="text-gray-600">{{ comment.content }}</p>
                    </li>
                  </ul>
                  <p v-else class="text-xs text-gray-400 mb-3">暂无评论</p>
                  
                  <div class="flex gap-2">
                    <input 
                      v-model="newCommentMap[task.id]" 
                      type="text" 
                      placeholder="写下你的评论..." 
                      class="flex-1 text-sm border rounded px-2 py-1 focus:outline-none focus:border-blue-500"
                      @keyup.enter="submitComment(task)"
                    >
                    <button 
                      @click="submitComment(task)" 
                      class="text-xs bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700 disabled:opacity-50"
                      :disabled="!newCommentMap[task.id]"
                    >
                      发送
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 进度调整模态框 -->
        <div
          v-if="showProgressModal"
          class="fixed inset-0 z-50 flex items-center justify-center"
        >
          <div
            class="absolute inset-0 bg-black opacity-40"
            @click="closeProgressModal"
          ></div>
          <div class="bg-white rounded-lg shadow-lg z-10 w-full max-w-md p-6">
            <h3 class="text-lg font-semibold mb-4">更新任务进度</h3>
            <div v-if="progressTargetTask" class="space-y-4">
              <div>
                <div class="text-sm text-gray-600">
                  {{ progressTargetTask.title }}
                </div>
                <div class="text-xs text-gray-400 mt-1">
                  当前进度：{{ getTaskProgressValue(progressTargetTask) }}%
                </div>
              </div>
              <div>
                <label class="text-sm font-medium" for="adjust-delta-input"
                  >调整幅度（可为正负）</label
                >
                <input
                  type="number"
                  min="-100"
                  max="100"
                  step="1"
                  v-model.number="progressAdjustForm.delta"
                  id="adjust-delta-input"
                  class="w-full mt-1 p-2 border rounded"
                />
              </div>
              <div>
                <input
                  type="range"
                  min="-50"
                  max="50"
                  step="5"
                  v-model.number="progressAdjustForm.delta"
                  class="w-full"
                />
                <div class="text-xs text-gray-500 mt-1">
                  拖动滑块快速设置（-50% ~ +50%）
                </div>
              </div>
              <p class="text-sm text-gray-700">
                预计更新后进度：<span class="font-semibold"
                  >{{ previewAdjustedProgress }}%</span
                >
              </p>
            </div>
            <div class="mt-6 flex justify-end gap-2">
              <button
                class="px-4 py-2 rounded border"
                @click="closeProgressModal"
              >
                取消
              </button>
              <button
                class="px-4 py-2 rounded bg-[#2D5BFF] text-white disabled:opacity-60"
                :disabled="
                  !progressTargetTask || progressAdjustForm.delta === 0
                "
                @click="submitProgressAdjustment"
              >
                确认更新
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧栏 -->
      <div class="space-y-5">
        <div class="card surface-card">
          <div class="flex justify-between items-center">
            <h3 class="section-title">团队动态</h3>
            <span class="text-xs text-blue-600">更多</span>
          </div>

          <div class="mt-4 space-y-3">
            <div class="flex items-start">
              <div class="relative">
                <div class="w-10 h-10 rounded-full bg-gray-300"></div>
                <div class="online-indicator"></div>
              </div>
              <div class="ml-3 flex-1">
                <div>
                  <span class="font-medium">王同学</span>
                  <span class="text-sm text-gray-500">完成了 登录功能开发</span>
                </div>
                <div class="text-xs text-gray-500 mt-1">20分钟前</div>
                <div class="flex items-center mt-1">
                  <iconify-icon
                    icon="mdi:heart-outline"
                    width="16"
                    height="16"
                    class="interaction-btn text-gray-500"
                  ></iconify-icon>
                  <span class="text-xs text-gray-500 ml-1 mr-3">8</span>
                  <iconify-icon
                    icon="mdi:comment-outline"
                    width="16"
                    height="16"
                    class="interaction-btn text-gray-500"
                  ></iconify-icon>
                  <span class="text-xs text-gray-500 ml-1">2</span>
                </div>
              </div>
            </div>

            <div class="flex items-start">
              <div class="relative">
                <div class="w-10 h-10 rounded-full bg-gray-300"></div>
              </div>
              <div class="ml-3 flex-1">
                <div>
                  <span class="font-medium">钱同学</span>
                  <span class="text-sm text-gray-500"
                    >创建了新任务: 支付模块设计</span
                  >
                </div>
                <div class="text-xs text-gray-500 mt-1">1小时前</div>
              </div>
            </div>
          </div>
        </div>

        <div class="card surface-card">
          <div class="flex justify-between items-center">
            <h3 class="section-title">团队积分排名</h3>
            <span class="text-xs text-blue-600">详情</span>
          </div>
          <div class="mt-4 space-y-3">
            <div
              class="flex items-center p-2 hover:bg-gray-50 rounded-lg transition-colors"
            >
              <div class="w-6 text-center font-bold text-[#FF6B35]">1</div>
              <div class="w-8 h-8 rounded-full bg-gray-300 ml-2"></div>
              <div class="ml-3 flex-1">
                <div class="font-medium">王同学</div>
                <div class="text-xs text-gray-500">前端开发</div>
              </div>
              <div class="text-[#FF6B35] font-semibold">1,580</div>
            </div>
            <div
              class="flex items-center p-2 hover:bg-gray-50 rounded-lg transition-colors"
            >
              <div class="w-6 text-center font-bold text-[#FF9500]">2</div>
              <div class="w-8 h-8 rounded-full bg-gray-300 ml-2"></div>
              <div class="ml-3 flex-1">
                <div class="font-medium">李同学</div>
                <div class="text-xs text-gray-500">后端开发</div>
              </div>
              <div class="text-[#FF9500] font-semibold">1,420</div>
            </div>
            <div
              class="flex items-center p-2 hover:bg-gray-50 rounded-lg transition-colors"
            >
              <div class="w-6 text-center font-bold text-[#FFC107]">3</div>
              <div class="w-8 h-8 rounded-full bg-gray-300 ml-2"></div>
              <div class="ml-3 flex-1">
                <div class="font-medium">陈同学</div>
                <div class="text-xs text-gray-500">UI设计</div>
              </div>
              <div class="text-[#FFC107] font-semibold">1,350</div>
            </div>
          </div>
        </div>

        <div class="card surface-card">
          <h3 class="section-title mb-4">协作工具</h3>
          <div class="grid grid-cols-2 gap-3">
            <button
              class="py-3 bg-blue-50 hover:bg-blue-100 rounded-lg flex flex-col items-center justify-center transition-colors"
            >
              <iconify-icon
                icon="mdi:video"
                width="24"
                height="24"
                style="color: #2d5bff"
              ></iconify-icon>
              <span class="mt-2 text-sm text-gray-700">视频会议</span>
            </button>
            <button
              class="py-3 bg-orange-50 hover:bg-orange-100 rounded-lg flex flex-col items-center justify-center transition-colors"
            >
              <iconify-icon
                icon="mdi:file-document"
                width="24"
                height="24"
                style="color: #ff9500"
              ></iconify-icon>
              <span class="mt-2 text-sm text-gray-700">协作文档</span>
            </button>
            <button
              class="py-3 bg-green-50 hover:bg-green-100 rounded-lg flex flex-col items-center justify-center transition-colors"
            >
              <iconify-icon
                icon="mdi:calendar"
                width="24"
                height="24"
                style="color: #4caf50"
              ></iconify-icon>
              <span class="mt-2 text-sm text-gray-700">团队日历</span>
            </button>
            <button
              class="py-3 bg-purple-50 hover:bg-purple-100 rounded-lg flex flex-col items-center justify-center transition-colors"
            >
              <iconify-icon
                icon="mdi:chart-pie"
                width="24"
                height="24"
                style="color: #9c27b0"
              ></iconify-icon>
              <span class="mt-2 text-sm text-gray-700">数据报告</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 邀请成员模态框 -->
    <div
      v-if="showInviteModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    >
      <div class="bg-white p-6 rounded-xl shadow-xl w-full max-w-md">
        <h3 class="text-lg font-bold mb-4">邀请成员</h3>
        <input
          v-model="inviteAccount"
          type="text"
          class="w-full p-2 border rounded mb-4"
          placeholder="请输入用户账号"
        />
        <div class="flex justify-end gap-2">
          <button
            @click="showInviteModal = false"
            class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded"
          >
            取消
          </button>
          <button
            @click="submitInvite"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            发送邀请
          </button>
        </div>
      </div>
    </div>

    <!-- 团队成员列表模态框 -->
    <div
      v-if="showMembersModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    >
      <div class="bg-white p-6 rounded-xl shadow-xl w-full max-w-lg">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-bold">团队成员</h3>
          <button
            @click="showMembersModal = false"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <iconify-icon icon="mdi:close" width="24" height="24"></iconify-icon>
          </button>
        </div>
        
        <div v-if="loadingMembers" class="text-center py-8 text-gray-500">
          <iconify-icon icon="mdi:loading" class="animate-spin" width="32" height="32"></iconify-icon>
          <p class="mt-2">加载中...</p>
        </div>
        
        <div v-else-if="teamMembers.length === 0" class="text-center py-8 text-gray-500">
          <iconify-icon icon="mdi:account-off" width="48" height="48"></iconify-icon>
          <p class="mt-2">暂无成员</p>
        </div>
        
        <div v-else class="max-h-96 overflow-y-auto">
          <div
            v-for="member in teamMembers"
            :key="member.user_id"
            class="flex items-center justify-between p-3 hover:bg-gray-50 rounded-lg transition-colors"
          >
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-bold">
                {{ (member.nickname || member.account || '?').charAt(0).toUpperCase() }}
              </div>
              <div>
                <p class="font-medium text-gray-800">{{ member.nickname || member.account }}</p>
                <p class="text-sm text-gray-500">{{ member.account }}</p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <span
                v-if="member.user_id === selectedTeam.owner_user_id"
                class="text-xs px-2 py-1 bg-yellow-100 text-yellow-700 rounded-full font-medium"
              >
                队长
              </span>
              <span
                v-else
                class="text-xs px-2 py-1 bg-gray-100 text-gray-600 rounded-full font-medium"
              >
                成员
              </span>
            </div>
          </div>
        </div>
        
        <div class="mt-4 pt-4 border-t border-gray-200 text-center text-sm text-gray-500">
          共 {{ teamMembers.length }} 位成员
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import * as echarts from "echarts";
import {
  getTeamTasks,
  createTask,
  updateTaskProgress,
  getTaskDetail,
  addTaskComment,
} from "@/api/modules/task";
import {
  getTeamList,
  joinTeamByName,
  createTeam,
  inviteMember,
  getTeamMembers,
} from "@/api/modules/team";
import { useCurrentUser } from "@/composables/useCurrentUser";

export default {
  name: "TeamTasks",
  components: {
  },
  setup() {
    const { profile, loadCurrentUser } = useCurrentUser();
    if (!profile.value) {
      loadCurrentUser().catch((error) => {
        console.warn("加载当前用户信息失败：", error);
      });
    }
    return {
      currentUserProfile: profile,
      loadCurrentUserProfile: loadCurrentUser,
    };
  },
  data() {
    return {
      tasks: [],
      chart: null,
      showCreateModal: false,
      creating: false,
      showProgressModal: false,
      progressTargetTask: null,
      progressAdjustForm: {
        delta: 0,
      },
      expandedTaskIds: [],
      taskDetailCache: {},
      newCommentMap: {},

      // Team Management
      showInviteModal: false,
      inviteAccount: "",
      showMembersModal: false,
      teamMembers: [],
      loadingMembers: false,
      detailLoadingMap: {},
      ownedTeams: [],
      ownedTeamsLoading: false,
      newTaskForm: {
        title: "",
        description: "",
        due_date: "",
        effort_points: 0,
        owner_team_id: "",
        subtasks: [""],
      },
      allTeams: [],
      selectedTeam: null,
      teamNameInput: "",
      loadingTeams: false,
      showJoinModal: false,
      showCreateTeamModal: false,
      newTeamForm: {
        name: "",
        description: "",
      },
      particles: [],
      animationFrameId: null,
    };
  },
  watch: {
    selectedTeam(val) {
      if (!val) {
        this.$nextTick(() => {
          this.initParticles();
        });
      } else {
        if (this.animationFrameId) {
          cancelAnimationFrame(this.animationFrameId);
        }
      }
    },
  },
  computed: {
    previewAdjustedProgress() {
      if (!this.progressTargetTask) return 0;
      const base = this.getTaskProgressValue(this.progressTargetTask);
      const delta = Number(this.progressAdjustForm.delta) || 0;
      return this.clampProgress(base + delta);
    },
    currentUserId() {
      const profile = this.currentUserProfile;
      return (
        profile?.id ||
        profile?.basic_info?.id ||
        profile?.basic_info?.user_id ||
        null
      );
    },
  },
  mounted() {
    this.loadAllTeams();
    window.addEventListener("resize", this.handleResize);
    if (!this.selectedTeam) {
      this.$nextTick(() => {
        this.initParticles();
      });
    }
  },
  beforeUnmount() {
    window.removeEventListener("resize", this.handleResize);
    if (this.chart) {
      this.chart.dispose();
      this.chart = null;
    }
    if (this.animationFrameId) {
      cancelAnimationFrame(this.animationFrameId);
    }
  },
  methods: {
    formatDate(dateStr) {
      if (!dateStr) return "未知时间";
      const date = new Date(dateStr);
      if (isNaN(date.getTime())) return dateStr;
      return date.toLocaleDateString("zh-CN", {
        year: "numeric",
        month: "long",
        day: "numeric",
      });
    },
    async loadAllTeams() {
      this.loadingTeams = true;
      try {
        const res = await getTeamList({ owned_only: false });
        const list = res?.data?.data || res?.data || res;
        this.allTeams = Array.isArray(list) ? list : [];
      } catch (error) {
        console.error("加载团队列表失败", error);
      } finally {
        this.loadingTeams = false;
      }
    },
    async handleJoinTeam() {
      if (!this.teamNameInput.trim()) return;
      try {
        const res = await joinTeamByName(this.teamNameInput);
        this.teamNameInput = "";
        this.showJoinModal = false;
        await this.loadAllTeams();
        const msg = res.data?.msg || "加入成功";
        alert(msg);
      } catch (error) {
        console.error("加入团队失败", error);
        alert(error.response?.data?.error || "加入失败");
      }
    },
    openInviteModal() {
      this.showInviteModal = true;
      this.inviteAccount = "";
    },
    async submitInvite() {
      if (!this.inviteAccount.trim()) return;
      try {
        await inviteMember(this.selectedTeam.id, {
          account: this.inviteAccount,
        });
        alert("邀请已发送");
        this.showInviteModal = false;
      } catch (error) {
        alert(error.response?.data?.error || "邀请失败");
      }
    },
    async openMembersModal() {
      this.showMembersModal = true;
      this.teamMembers = [];
      this.loadingMembers = true;
      try {
        const res = await getTeamMembers(this.selectedTeam.id);
        this.teamMembers = res.data || [];
      } catch (error) {
        console.error("获取团队成员失败", error);
        alert(error.response?.data?.error || "获取成员列表失败");
      } finally {
        this.loadingMembers = false;
      }
    },
    async handleCreateTeam() {
      if (!this.newTeamForm.name.trim()) {
        alert("请输入团队名称");
        return;
      }
      try {
        await createTeam(this.newTeamForm);
        this.newTeamForm.name = "";
        this.newTeamForm.description = "";
        this.showCreateTeamModal = false;
        await this.loadAllTeams();
        alert("创建成功");
      } catch (error) {
        console.error("创建团队失败", error);
        alert(error.response?.data?.error || "创建失败");
      }
    },
    selectTeam(team) {
      this.selectedTeam = team;
      this.$nextTick(() => {
        this.initChart();
        this.loadTasks();
        this.loadOwnedTeams();
      });
    },
    isCreatedByCurrentUser(task) {
      if (!task || !this.currentUserId) return false;
      return task.created_by === this.currentUserId;
    },
    focusTaskCard(taskId) {
      if (!taskId) return;
      const target = this.$el?.querySelector(`[data-task-id="${taskId}"]`);
      if (target) {
        target.classList.add("task-card--highlight");
        target.scrollIntoView({ behavior: "smooth", block: "center" });
        setTimeout(() => target.classList.remove("task-card--highlight"), 1500);
      }
      if (!this.isTaskExpanded(taskId)) {
        this.toggleTaskDetails(taskId);
      }
    },
    async loadOwnedTeams() {
      this.ownedTeamsLoading = true;
      try {
        const res = await getTeamList({ owned_only: true });
        const list = res?.data?.data || res?.data || res;
        this.ownedTeams = Array.isArray(list) ? list : [];
      } catch (error) {
        console.warn("加载团队列表失败：", error);
        this.ownedTeams = [];
      } finally {
        this.ownedTeamsLoading = false;
      }
    },
    resetNewTaskForm() {
      this.newTaskForm.title = "";
      this.newTaskForm.description = "";
      this.newTaskForm.due_date = "";
      this.newTaskForm.effort_points = 0;
      this.newTaskForm.owner_team_id = "";
      this.newTaskForm.subtasks = [""];
    },
    addSubtaskField() {
      this.newTaskForm.subtasks.push("");
    },
    removeSubtaskField(index) {
      if (this.newTaskForm.subtasks.length === 1) {
        this.newTaskForm.subtasks.splice(0, 1, "");
        return;
      }
      this.newTaskForm.subtasks.splice(index, 1);
    },
    normalizeNewTaskSubtasks() {
      return this.newTaskForm.subtasks
        .map((item) => (item || "").trim())
        .filter((item) => item.length > 0);
    },
    goToConstellation() {
      if (this.selectedTeam) {
        this.$router.push({
          path: "/team-tasks/constellation",
          query: { teamId: this.selectedTeam.id },
        });
      } else {
        this.$router.push("/team-tasks/constellation");
      }
    },
    handleResize() {
      if (this.chart) this.chart.resize();
      if (!this.selectedTeam) {
        this.initParticles();
      }
    },
    initParticles() {
      const canvas = this.$refs.particleCanvas;
      if (!canvas) return;

      const ctx = canvas.getContext("2d");
      canvas.width = canvas.parentElement.clientWidth;
      canvas.height = canvas.parentElement.clientHeight;

      this.particles = [];
      const particleCount = Math.floor((canvas.width * canvas.height) / 10000); // 增加粒子密度

      for (let i = 0; i < particleCount; i++) {
        this.particles.push({
          x: Math.random() * canvas.width,
          y: Math.random() * canvas.height,
          vx: (Math.random() - 0.5) * 0.8, // 稍微加快速度
          vy: (Math.random() - 0.5) * 0.8,
          size: Math.random() * 4 + 2, // 增大粒子尺寸
          color: `rgba(${Math.floor(Math.random() * 50 + 80)}, ${Math.floor(
            Math.random() * 50 + 120
          )}, ${Math.floor(Math.random() * 100 + 180)}, ${
            Math.random() * 0.4 + 0.2
          })`, // 增加不透明度
        });
      }

      if (this.animationFrameId) cancelAnimationFrame(this.animationFrameId);
      this.animateParticles();
    },
    animateParticles() {
      const canvas = this.$refs.particleCanvas;
      if (!canvas) return;
      const ctx = canvas.getContext("2d");

      ctx.clearRect(0, 0, canvas.width, canvas.height);

      // 更新和绘制粒子
      this.particles.forEach((p, index) => {
        p.x += p.vx;
        p.y += p.vy;

        // 边界反弹
        if (p.x < 0 || p.x > canvas.width) p.vx *= -1;
        if (p.y < 0 || p.y > canvas.height) p.vy *= -1;

        ctx.beginPath();
        ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2);
        ctx.fillStyle = p.color;
        ctx.fill();

        // 绘制连线
        for (let j = index + 1; j < this.particles.length; j++) {
          const p2 = this.particles[j];
          const dx = p.x - p2.x;
          const dy = p.y - p2.y;
          const distance = Math.sqrt(dx * dx + dy * dy);

          if (distance < 120) {
            // 增加连线距离
            ctx.beginPath();
            ctx.strokeStyle = `rgba(100, 149, 237, ${
              0.3 * (1 - distance / 120)
            })`; // 增加连线不透明度
            ctx.lineWidth = 1; // 增加连线宽度
            ctx.moveTo(p.x, p.y);
            ctx.lineTo(p2.x, p2.y);
            ctx.stroke();
          }
        }
      });

      this.animationFrameId = requestAnimationFrame(this.animateParticles);
    },
    async loadTasks() {
      try {
        const params = {};
        if (this.selectedTeam && this.selectedTeam.id) {
          params.team_id = this.selectedTeam.id;
        }
        const res = await getTeamTasks(params);
        const items = res?.data?.items || res?.data || res;
        if (Array.isArray(items) && items.length) {
          this.tasks = items.map((item) => this.normalizeFetchedTask(item));
        } else {
          this.tasks = [];
        }
      } catch (error) {
        console.warn("加载团队任务失败：", error);
        this.tasks = [];
      } finally {
        this.updateChart();
      }
    },
    normalizeFetchedTask(raw) {
      const status = this.normalizeStatus(raw?.status);
      const progressSource = raw?.progress ?? this.statusToProgress(status);
      const progress = this.clampProgress(Number(progressSource));
      const due = raw?.due_at || raw?.due_date || "";

      let ownerName = raw?.owner_name || raw?.created_by_name;
      
      if (!ownerName) {
        // 优先使用 owner_user_id，其次使用 created_by
        const userId = raw?.owner_user_id || raw?.created_by;
        if (userId) {
          if (String(userId) === String(this.currentUserId)) {
            const p = this.currentUserProfile;
            ownerName = p?.name || p?.username || p?.basic_info?.name || "我";
          } else {
            // 如果没有名字，显示用户ID
            ownerName = `用户 ${userId}`;
          }
        }
      }

      return {
        id: raw?.id ?? Date.now(),
        title: raw?.title || raw?.name || "未命名任务",
        description: raw?.description || "",
        owner_name: ownerName || "未知",
        owner_team_id: raw?.owner_team_id ?? null,
        created_by: raw?.created_by ?? null,
        due_date:
          typeof due === "string" && due.includes("T")
            ? due.split("T")[0]
            : due,
        created_at: raw?.created_at || "",
        status,
        status_label: raw?.status_label || "",
        progress,
        subtasks: this.parseSubtaskPayload(
          raw?.subtasks,
          raw?.id ?? Date.now()
        ),
      };
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
      if (["in-progress", "progress", "doing"].includes(lowered))
        return "in-progress";
      return lowered;
    },
    parseSubtaskPayload(source, taskId) {
      if (!source) return [];
      let parsed = source;
      if (typeof parsed === "string") {
        try {
          parsed = JSON.parse(parsed);
        } catch (error) {
          console.warn("解析子任务数据失败，将尝试按分隔符拆分：", error);
          parsed = parsed
            .split(/\r?\n|,/)
            .map((item) => item.trim())
            .filter(Boolean);
        }
      }
      if (
        parsed &&
        typeof parsed === "object" &&
        !Array.isArray(parsed) &&
        "data" in parsed
      ) {
        parsed = parsed.data;
      }
      if (!Array.isArray(parsed)) return [];
      return parsed
        .map((entry, index) => {
          const title =
            typeof entry === "string"
              ? entry.trim()
              : entry?.title || entry?.name || "";
          const status =
            typeof entry === "object" && entry
              ? entry.status || entry.state || ""
              : "";
          if (!title) return null;
          return { id: `${taskId}-sub-${index}`, title, status };
        })
        .filter(Boolean);
    },
    normalizeTaskDetail(raw, fallbackTask) {
      const base = raw || {};
      const taskId = fallbackTask?.id || base?.id || Date.now();
      const subtasksSource = base.subtasks ?? fallbackTask?.subtasks;
      let attachments = [];
      if (Array.isArray(base.attachments)) attachments = base.attachments;
      else if (Array.isArray(fallbackTask?.attachments))
        attachments = fallbackTask.attachments;
      let comments = [];
      if (Array.isArray(base.comments)) comments = base.comments;
      else if (Array.isArray(fallbackTask?.comments))
        comments = fallbackTask.comments;
      return {
        subtasks: this.parseSubtaskPayload(subtasksSource, taskId),
        attachments,
        comments,
      };
    },
    async fetchTaskDetail(task) {
      if (!task?.id || this.taskDetailCache[task.id]) return;
      this.detailLoadingMap = { ...this.detailLoadingMap, [task.id]: true };
      try {
        const res = await getTaskDetail(task.id);
        const payload = res?.data?.data || res?.data || res;
        const normalized = this.normalizeTaskDetail(payload, task);
        this.taskDetailCache = {
          ...this.taskDetailCache,
          [task.id]: normalized,
        };
      } catch (error) {
        console.warn(`加载任务 ${task.id} 详情失败：`, error);
        const fallback = this.normalizeTaskDetail(task, task);
        this.taskDetailCache = { ...this.taskDetailCache, [task.id]: fallback };
      } finally {
        const nextLoading = { ...this.detailLoadingMap };
        delete nextLoading[task.id];
        this.detailLoadingMap = nextLoading;
      }
    },
    isDetailLoading(taskId) {
      return Boolean(this.detailLoadingMap[taskId]);
    },
    computeTaskProgressPercent(task) {
      return this.getTaskProgressValue(task);
    },
    computeTaskProgressLabel(task) {
      return `${this.getTaskProgressValue(task)}%`;
    },
    computeTeamProgress() {
      if (!this.tasks || !this.tasks.length) return 0;
      const total = this.tasks.reduce(
        (sum, task) => sum + this.getTaskProgressValue(task),
        0
      );
      return Math.round(total / this.tasks.length);
    },
    computeCompletionRate() {
      if (!this.tasks || !this.tasks.length) return 0;
      const finished = this.tasks.filter(
        (task) => this.getTaskProgressValue(task) >= 100
      ).length;
      return Math.round((finished / this.tasks.length) * 100);
    },
    initChart() {
      const el = this.$refs.teamProgressChart;
      if (!el) return;
      this.chart = echarts.init(el);
      this.updateChart();
    },
    updateChart() {
      if (!this.chart) return;
      const { labels, progressSeries, completionSeries } =
        this.buildWeeklyChartData();
      const barGradient = new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: "#5F7BFF" },
        { offset: 1, color: "#A5B4FC" },
      ]);
      const lineGradient = new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: "rgba(139, 92, 246, 0.35)" },
        { offset: 1, color: "rgba(139, 92, 246, 0)" },
      ]);
      const option = {
        tooltip: {
          trigger: "axis",
          axisPointer: {
            type: "cross",
            label: { backgroundColor: "#1f2937" },
          },
          formatter: (params) => {
            const progressPoint = params.find(
              (p) => p.seriesName === "项目进度"
            );
            const completionPoint = params.find(
              (p) => p.seriesName === "任务完成率"
            );
            const progressValue = progressPoint
              ? `${progressPoint.value}%`
              : "-";
            const completionValue = completionPoint
              ? `${completionPoint.value}%`
              : "-";
            return `${params[0].axisValue}<br/>进度: ${progressValue}<br/>任务完成率: ${completionValue}`;
          },
        },
        legend: {
          data: ["项目进度", "进度趋势", "任务完成率"],
          bottom: 0,
          textStyle: { color: "#475569" },
        },
        grid: {
          left: "3%",
          right: "4%",
          bottom: "15%",
          top: "10%",
          containLabel: true,
        },
        xAxis: {
          type: "category",
          data: labels,
          axisLine: { lineStyle: { color: "#E2E8F0" } },
          axisLabel: { color: "#64748B" },
        },
        yAxis: {
          type: "value",
          max: 100,
          axisLabel: { formatter: "{value}%", color: "#64748B" },
          axisLine: { lineStyle: { color: "#E2E8F0" } },
          splitLine: { lineStyle: { color: "#F1F5F9" } },
        },
        series: [
          {
            name: "项目进度",
            type: "bar",
            data: progressSeries,
            barWidth: "40%",
            itemStyle: { color: barGradient, borderRadius: 0 },
            emphasis: { itemStyle: { color: "#3755F0" } },
          },
          {
            name: "进度趋势",
            type: "line",
            data: progressSeries,
            smooth: true,
            lineStyle: { width: 2, color: "#8B5CF6" },
            symbol: "circle",
            symbolSize: 8,
            areaStyle: { color: lineGradient },
            zlevel: 2,
            z: 3,
          },
          {
            name: "任务完成率",
            type: "line",
            data: completionSeries,
            smooth: true,
            lineStyle: { width: 2, color: "#4CAF50", type: "dashed" },
            symbol: "triangle",
            symbolSize: 7,
            zlevel: 1,
            z: 2,
          },
        ],
      };
      this.chart.setOption(option);
    },
    getTaskProgressValue(task) {
      if (!task) return 0;
      if (typeof task.progress === "number" && !Number.isNaN(task.progress)) {
        return this.clampProgress(task.progress);
      }
      const numeric = Number(task.progress);
      if (!Number.isNaN(numeric)) {
        return this.clampProgress(numeric);
      }
      return this.statusToProgress(task.status);
    },
    clampProgress(value) {
      if (Number.isNaN(value)) return 0;
      return Math.max(0, Math.min(100, Math.round(value)));
    },
    statusToProgress(status) {
      const normalized = this.normalizeStatus(status);
      if (normalized === "completed") return 100;
      if (normalized === "in-progress") return 50;
      return 0;
    },
    buildWeeklyChartData() {
      const buckets = this.buildWeekBuckets();
      const fallbackProgress = this.computeTeamProgress();
      const fallbackCompletion = this.computeCompletionRate();
      const referenceToday = this.startOfDay(new Date());
      const taskList = this.tasks || [];
      for (const task of taskList) {
        const referenceDate = this.getTaskReferenceDate(task) || referenceToday;
        const progress = this.getTaskProgressValue(task);
        for (const bucket of buckets) {
          if (referenceDate >= bucket.start && referenceDate <= bucket.end) {
            bucket.progresses.push(progress);
            bucket.total += 1;
            if (progress >= 100) bucket.completed += 1;
          }
        }
      }
      const progressValues = buckets.map((bucket) => {
        if (!bucket.progresses.length) return null;
        const sum = bucket.progresses.reduce((acc, cur) => acc + cur, 0);
        return Math.round(sum / bucket.progresses.length);
      });
      const completionValues = buckets.map((bucket) => {
        if (!bucket.total) return null;
        return Math.round((bucket.completed / bucket.total) * 100);
      });
      return {
        labels: buckets.map((bucket) => bucket.label),
        progressSeries: this.fillSeriesGaps(progressValues, fallbackProgress),
        completionSeries: this.fillSeriesGaps(
          completionValues,
          fallbackCompletion
        ),
      };
    },
    fillSeriesGaps(values, defaultValue) {
      let last =
        typeof defaultValue === "number" && !Number.isNaN(defaultValue)
          ? defaultValue
          : 0;
      return values.map((value) => {
        if (typeof value === "number" && !Number.isNaN(value)) {
          last = value;
          return value;
        }
        return last;
      });
    },
    buildWeekBuckets(count = 5) {
      const buckets = [];
      const today = this.startOfDay(new Date());
      for (let i = count - 1; i >= 0; i -= 1) {
        const start = new Date(today);
        start.setDate(start.getDate() - i * 7);
        const end = new Date(start);
        end.setDate(end.getDate() + 6);
        buckets.push({
          label: this.formatWeekLabel(start, end, i === 0),
          start: this.startOfDay(start),
          end: this.endOfDay(end),
          progresses: [],
          completed: 0,
          total: 0,
        });
      }
      return buckets;
    },
    formatWeekLabel(start, end, isCurrentWeek) {
      if (isCurrentWeek) return "本周";
      return `${this.formatMonthDay(start)}-${this.formatMonthDay(end)}`;
    },
    formatMonthDay(date) {
      const month = date.getMonth() + 1;
      const day = date.getDate();
      const paddedDay = day < 10 ? `0${day}` : `${day}`;
      return `${month}.${paddedDay}`;
    },
    startOfDay(date) {
      const d = new Date(date);
      d.setHours(0, 0, 0, 0);
      return d;
    },
    endOfDay(date) {
      const d = new Date(date);
      d.setHours(23, 59, 59, 999);
      return d;
    },
    getTaskReferenceDate(task) {
      return this.parseDateString(
        task?.due_date || task?.due_at || task?.created_at
      );
    },
    parseDateString(value) {
      if (!value) return null;
      if (value instanceof Date) return value;
      const direct = new Date(value);
      if (!Number.isNaN(direct.getTime())) return direct;
      if (typeof value === "string") {
        const normalized = value.replaceAll("-", "/");
        const fallback = new Date(normalized);
        if (!Number.isNaN(fallback.getTime())) return fallback;
      }
      return null;
    },
    taskHealth(task) {
      const dueDate = this.parseDateString(task?.due_date || task?.due_at);
      if (!dueDate) return null;
      const today = this.startOfDay(new Date());
      const diffMs = dueDate.getTime() - today.getTime();
      const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
      const progress = this.getTaskProgressValue(task);
      if (diffDays < 0 && progress < 100) {
        return { type: "overdue", label: `逾期 ${Math.abs(diffDays)} 天` };
      }
      if (diffDays >= 0 && diffDays <= 3 && progress < 100) {
        return {
          type: "warning",
          label: diffDays === 0 ? "今日截止" : `剩余 ${diffDays} 天`,
        };
      }
      return null;
    },
    getHealthBadgeClass(health) {
      if (!health) return "badge-info";
      if (health.type === "overdue") return "badge-danger";
      if (health.type === "warning") return "badge-warning";
      return "badge-info";
    },
    taskTone(task) {
      const health = this.taskHealth(task);
      if (health?.type === "overdue") return "danger";
      if (health?.type === "warning") return "warning";
      const normalized = this.normalizeStatus(task?.status);
      if (normalized === "completed") return "success";
      return "info";
    },
    toneColors() {
      return {
        success: { solid: "#34D399", light: "#DCFCE7" },
        info: { solid: "#3B82F6", light: "#DBEAFE" },
        warning: { solid: "#F59E0B", light: "#FEF3C7" },
        danger: { solid: "#F87171", light: "#FEE2E2" },
      };
    },
    statusBadgeClass(task) {
      return `badge-${this.taskTone(task)}`;
    },
    taskCardClass(task) {
      return ["task-card", `task-card--${this.taskTone(task)}`];
    },
    progressFillStyle(task) {
      const percent = this.computeTaskProgressPercent(task);
      const tone = this.taskTone(task);
      const toneMap = this.toneColors();
      const color = toneMap[tone]?.solid ?? toneMap.info.solid;
      return { width: `${percent}%`, background: color };
    },
    openCreateModal() {
      this.resetNewTaskForm();
      this.showCreateModal = true;
    },
    closeCreateModal() {
      if (this.creating) return;
      this.showCreateModal = false;
      this.resetNewTaskForm();
    },
    openProgressModal(task) {
      this.progressTargetTask = task;
      this.progressAdjustForm = { delta: 0 };
      this.showProgressModal = true;
    },
    closeProgressModal() {
      this.showProgressModal = false;
      this.progressTargetTask = null;
    },
    async submitProgressAdjustment() {
      if (!this.progressTargetTask) {
        this.closeProgressModal();
        return;
      }
      const delta = Number(this.progressAdjustForm.delta) || 0;
      if (delta === 0) {
        this.closeProgressModal();
        return;
      }
      await this.changeProgress(this.progressTargetTask, delta);
      this.closeProgressModal();
    },
    async submitCreateTask() {
      if (!this.newTaskForm.title || !this.newTaskForm.title.trim()) {
        alert("任务名称不能为空");
        return;
      }
      const ownerTeamId = this.newTaskForm.owner_team_id
        ? Number(this.newTaskForm.owner_team_id)
        : null;
      const subtasks = this.normalizeNewTaskSubtasks();
      
      let formattedDueDate = null;
      if (this.newTaskForm.due_date) {
        // 将 YYYY-MM-DD 转换为 ISO 8601 格式，以符合后端 time.Time 的解析要求
        // 这里简单地将其转换为当天的 00:00:00 UTC
        formattedDueDate = new Date(this.newTaskForm.due_date).toISOString();
      }

      const payload = {
        title: this.newTaskForm.title.trim(),
        description: this.newTaskForm.description?.trim?.()
          ? this.newTaskForm.description.trim()
          : this.newTaskForm.description || "",
        due_at: formattedDueDate,
        effort_points: this.newTaskForm.effort_points || 0,
        task_type: 2,
      };
      if (ownerTeamId && !Number.isNaN(ownerTeamId)) {
        payload.owner_team_id = ownerTeamId;
      }
      if (subtasks.length) {
        payload.subtasks = subtasks;
      }
      this.creating = true;
      let newTask;
      try {
        const res = await createTask(payload);
        const created = res?.data?.data || res?.data || res;
        newTask = this.normalizeFetchedTask(created);
      } catch (error) {
        console.warn("创建任务失败，使用本地数据：", error);
        newTask = this.buildLocalTask(payload);
      } finally {
        if (newTask) {
          this.tasks.unshift(newTask);
        }
        this.creating = false;
        this.showCreateModal = false;
        this.resetNewTaskForm();
        this.updateChart();
      }
    },
    buildLocalTask(payload, extra = {}) {
      return {
        id: extra.id || Date.now(),
        title: payload.title,
        description: payload.description || "",
        owner_name: "你",
        due_date: payload.due_at || payload.due_date || "",
        created_at: extra.created_at || new Date().toISOString(),
        status: extra.status || "in-progress",
        status_label: extra.status_label || "",
        progress: this.clampProgress(extra.progress ?? 0),
        owner_team_id: payload.owner_team_id || null,
        created_by: this.currentUserId,
        subtasks: this.parseSubtaskPayload(
          extra.subtasks || payload.subtasks,
          extra.id || Date.now()
        ),
      };
    },
    toggleTaskDetails(taskId) {
      const index = this.expandedTaskIds.indexOf(taskId);
      if (index >= 0) {
        this.expandedTaskIds.splice(index, 1);
      } else {
        this.expandedTaskIds.push(taskId);
        const target = this.tasks.find((item) => item.id === taskId);
        if (target) this.fetchTaskDetail(target);
      }
    },
    isTaskExpanded(taskId) {
      return this.expandedTaskIds.includes(taskId);
    },
    getTaskDetail(task) {
      if (!task) return { subtasks: [], attachments: [], comments: [] };
      return (
        this.taskDetailCache[task.id] || this.normalizeTaskDetail({}, task)
      );
    },
    async changeProgress(task, delta) {
      const previousProgress = this.getTaskProgressValue(task);
      const nextProgress = this.clampProgress(previousProgress + delta);
      const previousStatus = task.status;
      task.progress = nextProgress;
      let nextStatus = "pending";
      if (nextProgress >= 100) nextStatus = "completed";
      else if (nextProgress > 0) nextStatus = "in-progress";
      task.status = nextStatus;
      this.updateChart();
      try {
        await updateTaskProgress(task.id, nextProgress);
      } catch (error) {
        console.warn("更新进度失败：", error);
        task.progress = previousProgress;
        task.status = previousStatus;
        this.updateChart();
        alert("更新进度失败，请重试");
      }
    },
    async setCompleted(task) {
      const current = this.getTaskProgressValue(task);
      await this.changeProgress(task, 100 - current);
    },
    async submitComment(task) {
      const content = this.newCommentMap[task.id];
      if (!content || !content.trim()) return;
      
      try {
        await addTaskComment(task.id, content);
        this.newCommentMap[task.id] = "";
        // Refresh task details
        const res = await getTaskDetail(task.id);
        const payload = res?.data?.data || res?.data || res;
        const normalized = this.normalizeTaskDetail(payload, task);
        this.taskDetailCache = {
          ...this.taskDetailCache,
          [task.id]: normalized,
        };
      } catch (error) {
        console.error("发表评论失败", error);
        alert("发表评论失败");
      }
    },
  },
};
</script>

<style scoped>
:global(body) {
  background: #f5f7fb;
  font-family: "PingFang SC", "Segoe UI", system-ui, -apple-system,
    BlinkMacSystemFont, sans-serif;
}

.surface-card {
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.95),
    rgba(245, 247, 251, 0.95)
  );
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(15, 23, 42, 0.08);
  padding: 24px;
}

.card.surface-card + .surface-card {
  margin-top: 16px;
  .task-card--highlight {
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.4),
      0 20px 45px rgba(15, 23, 42, 0.08);
  }
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #0f172a;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 12px;
}

.primary-btn {
  background: linear-gradient(120deg, #4f46e5, #3b82f6);
  color: #fff;
  font-weight: 600;
  padding: 10px 18px;
  border-radius: 6px;
  font-size: 14px;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.primary-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 24px rgba(79, 70, 229, 0.25);
}

.task-card {
  background: linear-gradient(
    145deg,
    rgba(255, 255, 255, 0.95),
    rgba(249, 250, 251, 0.9)
  );
  border-radius: 12px;
  padding: 18px;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.08);
  border: 1px solid rgba(148, 163, 184, 0.25);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.task-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.12);
}

.task-card--success {
  border-color: rgba(52, 211, 153, 0.4);
}

.task-card--info {
  border-color: rgba(59, 130, 246, 0.35);
}

.task-card--warning {
  border-color: rgba(245, 158, 11, 0.45);
}

.task-card--danger {
  border-color: rgba(248, 113, 113, 0.5);
}

.task-title {
  font-size: 16px;
  font-weight: 600;
  color: #0f172a;
}

.task-desc {
  font-size: 14px;
  color: #475569;
  margin-top: 4px;
}

.task-meta {
  display: flex;
  gap: 16px;
  margin-top: 12px;
  color: #94a3b8;
  font-size: 12px;
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.progress-bar {
  width: 100%;
  height: 8px;
  border-radius: 9999px;
  background: #e2e8f0;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: inherit;
  transition: width 0.3s ease, background 0.3s ease;
}

.chart-container {
  width: 100%;
  height: 320px;
  background: radial-gradient(
    circle at top,
    rgba(79, 70, 229, 0.08),
    transparent 45%
  );
  border-radius: 12px;
}

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.25rem 0.65rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
}

.badge-success {
  background: #dcfce7;
  color: #15803d;
}

.badge-info {
  background: #dbeafe;
  color: #1d4ed8;
}

.badge-warning {
  background: #fef3c7;
  color: #b45309;
}

.badge-danger {
  background: #fee2e2;
  color: #b91c1c;
}

.ghost-btn,
.success-btn,
.info-btn {
  padding: 6px 14px;
  font-size: 12px;
  font-weight: 600;
  border-radius: 10px;
  transition: background 0.2s ease, color 0.2s ease;
  border: 1px solid transparent;
}

.ghost-btn {
  background: #f1f5f9;
  color: #0f172a;
}

.ghost-btn:hover {
  background: #e2e8f0;
}

.success-btn {
  background: #dcfce7;
  color: #15803d;
}

.success-btn:hover {
  background: #bbf7d0;
}

.info-btn {
  background: #e0e7ff;
  color: #4338ca;
}

.info-btn:hover {
  background: #c7d2fe;
}

.task-card :is(h4, p, span) {
  transition: color 0.2s ease;
}
</style>
