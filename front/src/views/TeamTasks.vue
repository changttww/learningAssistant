<template>
	<div class="w-full h-full overflow-auto px-4">
		<!-- 团队选择/加入界面 -->
		<div v-if="!selectedTeam" class="relative flex flex-col items-center justify-center min-h-[80vh] overflow-hidden py-12">
			<!-- 动态背景装饰 -->
			<div class="absolute inset-0 -z-10">
				<div class="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50"></div>
				<div class="absolute top-0 left-1/4 w-96 h-96 bg-blue-200 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob"></div>
				<div class="absolute top-0 right-1/4 w-96 h-96 bg-purple-200 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-2000"></div>
				<div class="absolute -bottom-32 left-1/3 w-96 h-96 bg-pink-200 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-4000"></div>
			</div>

			<div v-if="loadingTeams" class="text-gray-500 flex flex-col items-center gap-3">
				<iconify-icon icon="eos-icons:loading" class="text-4xl text-blue-500 animate-spin"></iconify-icon>
				<span class="text-lg font-medium text-gray-600">加载团队信息...</span>
			</div>
			
			<!-- 加入或创建团队 -->
			<div v-else-if="!hasJoinedTeam || showCreateTeamPanel" class="bg-white/90 backdrop-blur-md rounded-2xl shadow-xl max-w-md w-full overflow-hidden p-8 border border-white/50 transition-all duration-500 hover:shadow-2xl">
				<!-- 加入团队面板 -->
				<div v-if="!showCreateTeamPanel" class="animate-fade-in">
					<div class="text-center mb-8">
						<div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4 text-blue-600">
							<iconify-icon icon="mdi:account-group-outline" class="text-3xl"></iconify-icon>
						</div>
						<h2 class="text-2xl font-bold text-gray-800">加入团队</h2>
						<p class="text-gray-600 mt-2">加入现有的学习小组，与伙伴共同进步。</p>
					</div>
					<div class="space-y-5">
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1.5">团队名称</label>
							<div class="relative">
								<iconify-icon icon="mdi:magnify" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-xl"></iconify-icon>
								<input 
									v-model="joinTeamName" 
									type="text" 
									placeholder="输入团队名称搜索" 
									class="w-full pl-10 pr-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
									@keyup.enter="handleJoinTeam"
								/>
							</div>
						</div>
						<button 
							@click="handleJoinTeam" 
							:disabled="!joinTeamName || joining"
							class="w-full py-3.5 bg-gradient-to-r from-blue-600 to-indigo-600 text-white font-bold rounded-xl hover:from-blue-700 hover:to-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg hover:shadow-blue-500/30 transform hover:-translate-y-0.5"
						>
							<span v-if="joining" class="flex items-center justify-center gap-2">
								<iconify-icon icon="eos-icons:loading" class="animate-spin"></iconify-icon> 加入中...
							</span>
							<span v-else>加入团队</span>
						</button>
					</div>
					<div class="mt-8 text-center">
						<button @click="showCreateTeamPanel = true" class="text-gray-500 hover:text-blue-600 text-sm transition-colors flex items-center justify-center gap-1 mx-auto group">
							还没有团队？<span class="underline decoration-dotted group-hover:decoration-solid">创建一个新团队</span>
						</button>
					</div>
				</div>

				<!-- 创建团队面板 -->
				<div v-else class="animate-fade-in">
					<div class="text-center mb-8">
						<div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4 text-green-600">
							<iconify-icon icon="mdi:plus-circle-outline" class="text-3xl"></iconify-icon>
						</div>
						<h2 class="text-2xl font-bold text-gray-800">创建团队</h2>
						<p class="text-gray-600 mt-2">创建一个新的学习小组，邀请伙伴加入。</p>
					</div>
					<div class="space-y-5">
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1.5">团队名称</label>
							<input 
								v-model="createTeamForm.name" 
								type="text" 
								placeholder="给团队起个名字" 
								class="w-full p-3 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-green-500 outline-none transition-all"
							/>
						</div>
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1.5">描述 (可选)</label>
							<textarea 
								v-model="createTeamForm.description" 
								rows="2"
								placeholder="简单介绍一下团队目标" 
								class="w-full p-3 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-green-500 outline-none transition-all resize-none"
							></textarea>
						</div>
						<button 
							@click="handleCreateTeam" 
							:disabled="!createTeamForm.name || creatingTeam"
							class="w-full py-3.5 bg-gradient-to-r from-green-600 to-emerald-600 text-white font-bold rounded-xl hover:from-green-700 hover:to-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg hover:shadow-green-500/30 transform hover:-translate-y-0.5"
						>
							<span v-if="creatingTeam" class="flex items-center justify-center gap-2">
								<iconify-icon icon="eos-icons:loading" class="animate-spin"></iconify-icon> 创建中...
							</span>
							<span v-else>创建团队</span>
						</button>
					</div>
					<div class="mt-8 text-center">
						<button @click="showCreateTeamPanel = false" class="text-gray-500 hover:text-blue-600 text-sm transition-colors flex items-center justify-center gap-1 mx-auto group">
							<span v-if="hasJoinedTeam" class="flex items-center gap-1">
								<iconify-icon icon="mdi:arrow-left"></iconify-icon> 返回团队列表
							</span>
							<span v-else>
								已有团队？<span class="underline decoration-dotted group-hover:decoration-solid">加入现有团队</span>
							</span>
						</button>
					</div>
				</div>
			</div>

			<!-- 选择团队 -->
			<div v-else class="max-w-5xl w-full z-10">
				<div class="text-center mb-12">
					<h2 class="text-3xl font-bold text-gray-800 mb-3">选择你的团队</h2>
					<p class="text-gray-500 mb-6">选择一个团队开始今天的协作任务</p>
					<button @click="showCreateTeamPanel = true" class="px-6 py-2.5 bg-white text-blue-600 font-medium rounded-full shadow-sm hover:shadow-md hover:bg-blue-50 transition-all border border-blue-100 flex items-center gap-2 mx-auto group">
						<div class="w-6 h-6 rounded-full bg-blue-100 flex items-center justify-center group-hover:bg-blue-200 transition-colors">
							<iconify-icon icon="mdi:plus" class="text-blue-600"></iconify-icon>
						</div>
						创建新团队
					</button>
				</div>
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 px-4">
					<div 
						v-for="(team, index) in teamList" 
						:key="team.id"
						@click="selectTeam(team)"
						class="bg-white/80 backdrop-blur-sm p-6 rounded-2xl shadow-sm hover:shadow-xl cursor-pointer transition-all duration-300 border border-white/60 hover:border-blue-400 hover:-translate-y-2 group relative overflow-hidden"
						:style="{ animationDelay: `${index * 100}ms` }"
					>
						<div class="absolute top-0 right-0 w-24 h-24 bg-gradient-to-br from-blue-50 to-transparent rounded-bl-full -mr-8 -mt-8 transition-transform group-hover:scale-150 duration-500"></div>
						
						<div class="flex items-center gap-5 mb-6 relative">
							<div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-blue-100 to-indigo-100 flex items-center justify-center text-blue-600 font-bold text-2xl shadow-inner group-hover:scale-110 transition-transform duration-300">
								{{ team.name ? team.name.charAt(0).toUpperCase() : 'T' }}
							</div>
							<div class="flex-1 min-w-0">
								<h3 class="font-bold text-lg text-gray-800 truncate group-hover:text-blue-600 transition-colors">{{ team.name }}</h3>
								<div class="flex items-center gap-2 text-xs text-gray-500 mt-1">
									<span class="bg-gray-100 px-2 py-0.5 rounded-full">ID: {{ team.id }}</span>
									<span>{{ formatDate(team.created_at).split(' ')[0] }}</span>
								</div>
							</div>
						</div>
						
						<p class="text-sm text-gray-600 line-clamp-2 mb-4 h-10 leading-relaxed">{{ team.description || '暂无描述' }}</p>
						
						<div class="flex items-center justify-between pt-4 border-t border-gray-100">
							<span class="text-xs font-medium text-gray-400 group-hover:text-blue-500 transition-colors">点击进入</span>
							<iconify-icon icon="mdi:arrow-right" class="text-gray-300 group-hover:text-blue-500 group-hover:translate-x-1 transition-all"></iconify-icon>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<div class="lg:col-span-3">
				<button @click="selectedTeam = null" class="mt-4 mb-2 px-4 py-2 bg-white border border-gray-200 rounded-lg shadow-sm text-gray-600 hover:text-blue-600 hover:border-blue-300 hover:shadow transition-all duration-200 flex items-center gap-2 text-sm font-medium w-fit">
					<iconify-icon icon="mdi:arrow-left"></iconify-icon> 切换团队
				</button>
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
					<div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center">
						<div class="absolute inset-0 bg-black opacity-40" @click="closeCreateModal"></div>
						<div class="bg-white rounded-lg shadow-lg z-10 w-full max-w-md p-6">
							<h3 class="text-lg font-semibold mb-4">创建新任务</h3>
							<div class="space-y-3">
								<div>
									<label for="create-title" class="text-sm">任务名称</label>
									<input id="create-title" v-model="newTaskForm.title" type="text" class="w-full mt-1 p-2 border rounded" placeholder="请输入任务名称" />
								</div>
								<div>
									<label for="create-desc" class="text-sm">描述</label>
									<textarea id="create-desc" v-model="newTaskForm.description" rows="3" class="w-full mt-1 p-2 border rounded" placeholder="任务描述（可选）"></textarea>
								</div>
								<div class="flex gap-2">
									<div class="flex-1">
										<label for="create-due" class="text-sm">到期日</label>
										<input id="create-due" v-model="newTaskForm.due_date" type="date" class="w-full mt-1 p-2 border rounded" />
									</div>
									<div class="flex-1">
										<label for="create-owner" class="text-sm">负责人</label>
										<input id="create-owner" v-model="newTaskForm.owner_name" type="text" class="w-full mt-1 p-2 border rounded" placeholder="负责人姓名（可选）" />
									</div>
								</div>
							</div>
							<div class="mt-4 flex justify-end gap-2">
								<button @click="closeCreateModal" class="px-4 py-2 rounded border">取消</button>
								<button :disabled="creating" @click="submitCreateTask" class="px-4 py-2 rounded bg-[#2D5BFF] text-white">{{ creating ? '创建中...' : '创建' }}</button>
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

						<div v-for="task in tasks" :key="task.id" :class="taskCardClass(task)">
							<div class="flex items-start justify-between gap-4">
								<div class="flex-1">
									<h4 class="task-title">{{ task.title }}</h4>
									<p class="task-desc">{{ task.description }}</p>
									<div class="task-meta">
										<div class="meta-item">
											<iconify-icon icon="mdi:account" width="16" height="16" class="text-gray-400"></iconify-icon>
											<span>{{ task.owner_name || '未知' }}</span>
										</div>
										<div class="meta-item">
											<iconify-icon icon="mdi:calendar" width="16" height="16" class="text-gray-400"></iconify-icon>
											<span>{{ task.due_date || task.created_at || '' }}</span>
										</div>
									</div>
									<div class="mt-3">
										<div class="flex items-center justify-between mb-1">
											<span class="text-sm">进度</span>
											<span class="text-sm font-medium">{{ computeTaskProgressLabel(task) }}</span>
										</div>
										<div class="progress-bar">
											<div class="progress-fill" :style="progressFillStyle(task)"></div>
										</div>
										<div class="flex flex-wrap items-center gap-2 mt-2">
											<button @click="openProgressModal(task)" class="ghost-btn">更新进度</button>
											<button @click="setCompleted(task)" class="success-btn">标记完成</button>
											<button @click="toggleTaskDetails(task.id)" class="info-btn">
												{{ isTaskExpanded(task.id) ? '收起详情' : '查看详情' }}
											</button>
										</div>
									</div>
								</div>
								<div class="flex flex-col items-end gap-2">
									<span :class="['badge', statusBadgeClass(task)]">{{ task.status_label || task.status }}</span>
									<span
										v-if="taskHealth(task)"
										:class="['badge', getHealthBadgeClass(taskHealth(task))]"
									>{{ taskHealth(task).label }}</span>
								</div>
							</div>
							<div v-if="isTaskExpanded(task.id)" class="mt-4 rounded-lg border border-dashed border-gray-200 p-4 space-y-4">
								<div>
									<div class="flex items-center justify-between mb-2">
										<h4 class="text-sm font-semibold text-gray-700">子任务</h4>
										<span class="text-xs text-gray-400">{{ getTaskDetail(task).subtasks.length }} 个</span>
									</div>
									<ul class="space-y-1">
										<li
											v-for="sub in getTaskDetail(task).subtasks"
											:key="sub.id"
											class="flex items-center justify-between text-sm"
										>
											<span>{{ sub.title }}</span>
											<span class="text-xs text-gray-500">{{ sub.status }}</span>
										</li>
									</ul>
								</div>
								<div>
									<div class="flex items-center justify-between mb-2">
										<h4 class="text-sm font-semibold text-gray-700">附件</h4>
										<span class="text-xs text-gray-400">{{ getTaskDetail(task).attachments.length }} 个</span>
									</div>
									<ul class="space-y-1 text-sm">
										<li v-for="file in getTaskDetail(task).attachments" :key="file.id" class="flex items-center justify-between">
											<span>{{ file.name }}</span>
											<span class="text-xs text-gray-500">{{ file.size }}</span>
										</li>
									</ul>
								</div>
								<div>
									<div class="flex items-center justify-between mb-2">
										<h4 class="text-sm font-semibold text-gray-700">评论</h4>
										<span class="text-xs text-gray-400">{{ getTaskDetail(task).comments.length }} 条</span>
									</div>
									<ul class="space-y-2 text-sm">
										<li v-for="comment in getTaskDetail(task).comments" :key="comment.id">
											<div class="flex items-center justify-between">
												<span class="font-medium">{{ comment.author }}</span>
												<span class="text-xs text-gray-400">{{ comment.time }}</span>
											</div>
											<p class="text-gray-600">{{ comment.content }}</p>
										</li>
									</ul>
								</div>
							</div>
						</div>
					</div>
				</div>

				<!-- 进度调整模态框 -->
				<div v-if="showProgressModal" class="fixed inset-0 z-50 flex items-center justify-center">
					<div class="absolute inset-0 bg-black opacity-40" @click="closeProgressModal"></div>
					<div class="bg-white rounded-lg shadow-lg z-10 w-full max-w-md p-6">
						<h3 class="text-lg font-semibold mb-4">更新任务进度</h3>
						<div v-if="progressTargetTask" class="space-y-4">
							<div>
								<div class="text-sm text-gray-600">{{ progressTargetTask.title }}</div>
								<div class="text-xs text-gray-400 mt-1">当前进度：{{ getTaskProgressValue(progressTargetTask) }}%</div>
							</div>
							<div>
								<label class="text-sm font-medium" for="adjust-delta-input">调整幅度（可为正负）</label>
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
								<div class="text-xs text-gray-500 mt-1">拖动滑块快速设置（-50% ~ +50%）</div>
							</div>
							<p class="text-sm text-gray-700">
								预计更新后进度：<span class="font-semibold">{{ previewAdjustedProgress }}%</span>
							</p>
						</div>
						<div class="mt-6 flex justify-end gap-2">
							<button class="px-4 py-2 rounded border" @click="closeProgressModal">取消</button>
							<button
								class="px-4 py-2 rounded bg-[#2D5BFF] text-white disabled:opacity-60"
								:disabled="!progressTargetTask || progressAdjustForm.delta === 0"
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
	</div>
</template>

<script>
import * as echarts from "echarts";
import { getTeamTasks, createTask, updateTaskProgress } from "@/api/modules/task";
import { getTeamList, joinTeamByName, createTeam } from "@/api/modules/team";

export default {
	name: "TeamTasks",
	data() {
		return {
			hasJoinedTeam: false,
			teamList: [],
			selectedTeam: null,
			joinTeamName: "",
			showCreateTeamPanel: false,
			createTeamForm: {
				name: "",
				description: ""
			},
			loadingTeams: false,
			joining: false,
			creatingTeam: false,
			tasks: [
				{
					id: 1,
					title: "登录功能开发",
					description: "实现用户登录、注册和密码重置功能",
					owner_name: "王同学",
					due_date: "2024-01-18",
					created_at: "2024-01-01",
					status: "in-progress",
					status_label: "进行中",
					progress: 85,
				},
				{
					id: 2,
					title: "支付模块设计",
					description: "设计支付流程和界面规范",
					owner_name: "陈同学",
					due_date: "2026-01-22",
					created_at: "2024-01-05",
					status: "in-progress",
					status_label: "设计中",
					progress: 30,
				},
			],
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
			newTaskForm: {
				title: "",
				description: "",
				due_date: "",
				owner_name: "",
			},
		};
	},
	computed: {
		previewAdjustedProgress() {
			if (!this.progressTargetTask) return 0;
			const base = this.getTaskProgressValue(this.progressTargetTask);
			const delta = Number(this.progressAdjustForm.delta) || 0;
			return this.clampProgress(base + delta);
		},
	},
	mounted() {
		this.fetchTeams();
		window.addEventListener("resize", this.handleResize);
	},
	beforeUnmount() {
		window.removeEventListener("resize", this.handleResize);
		if (this.chart) {
			this.chart.dispose();
			this.chart = null;
		}
	},
	methods: {
		formatDate(dateStr) {
			if (!dateStr) return '';
			return new Date(dateStr).toLocaleDateString();
		},
		async fetchTeams() {
			this.loadingTeams = true;
			try {
				const res = await getTeamList();
				const teams = res.data || res;
				if (Array.isArray(teams) && teams.length > 0) {
					this.hasJoinedTeam = true;
					this.teamList = teams;
				} else {
					this.hasJoinedTeam = false;
					this.teamList = [];
				}
			} catch (error) {
				console.error("Failed to fetch teams:", error);
				this.hasJoinedTeam = false;
			} finally {
				this.loadingTeams = false;
			}
		},
		async handleJoinTeam() {
			if (!this.joinTeamName) return;
			this.joining = true;
			try {
				await joinTeamByName(this.joinTeamName);
				this.joinTeamName = "";
				await this.fetchTeams();
			} catch (error) {
				console.error("Failed to join team:", error);
				alert("加入团队失败，请检查团队名称是否正确或是否已加入。");
			} finally {
				this.joining = false;
			}
		},
		async handleCreateTeam() {
			if (!this.createTeamForm.name) return;
			this.creatingTeam = true;
			try {
				await createTeam(this.createTeamForm);
				this.createTeamForm.name = "";
				this.createTeamForm.description = "";
				await this.fetchTeams();
				this.showCreateTeamPanel = false;
			} catch (error) {
				console.error("Failed to create team:", error);
				alert("创建团队失败，团队名称可能已存在。");
			} finally {
				this.creatingTeam = false;
			}
		},
		selectTeam(team) {
			this.selectedTeam = team;
			this.loadTasks();
			this.$nextTick(() => {
				this.initChart();
			});
		},
		goToConstellation() {
			this.$router.push("/team-tasks/constellation");
		},
		handleResize() {
			if (this.chart) this.chart.resize();
		},
		async loadTasks() {
			if (!this.selectedTeam) return;
			try {
				const res = await getTeamTasks({ team_id: this.selectedTeam.id });
				const items = res?.data?.items || res?.data || res;
				if (Array.isArray(items) && items.length) {
					this.tasks = items.map((item) => this.normalizeFetchedTask(item));
				}
			} catch (error) {
				console.warn("加载团队任务失败，使用本地示例：", error);
			} finally {
				this.updateChart();
			}
		},
		normalizeFetchedTask(raw) {
			const status = this.normalizeStatus(raw?.status);
			const progressSource = raw?.progress ?? this.statusToProgress(status);
			const progress = this.clampProgress(Number(progressSource));
			const due = raw?.due_at || raw?.due_date || "";
			return {
				id: raw?.id ?? Date.now(),
				title: raw?.title || raw?.name || "未命名任务",
				description: raw?.description || "",
				owner_name: raw?.owner_name || raw?.created_by_name || "未知",
				due_date: typeof due === "string" && due.includes("T") ? due.split("T")[0] : due,
				created_at: raw?.created_at || "",
				status,
				status_label: raw?.status_label || "",
				progress,
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
			if (["in-progress", "progress", "doing"].includes(lowered)) return "in-progress";
			return lowered;
		},
		computeTaskProgressPercent(task) {
			return this.getTaskProgressValue(task);
		},
		computeTaskProgressLabel(task) {
			return `${this.getTaskProgressValue(task)}%`;
		},
		computeTeamProgress() {
			if (!this.tasks || !this.tasks.length) return 0;
			const total = this.tasks.reduce((sum, task) => sum + this.getTaskProgressValue(task), 0);
			return Math.round(total / this.tasks.length);
		},
		computeCompletionRate() {
			if (!this.tasks || !this.tasks.length) return 0;
			const finished = this.tasks.filter((task) => this.getTaskProgressValue(task) >= 100).length;
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
			const { labels, progressSeries, completionSeries } = this.buildWeeklyChartData();
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
						const progressPoint = params.find((p) => p.seriesName === "项目进度");
						const completionPoint = params.find((p) => p.seriesName === "任务完成率");
						const progressValue = progressPoint ? `${progressPoint.value}%` : "-";
						const completionValue = completionPoint ? `${completionPoint.value}%` : "-";
						return `${params[0].axisValue}<br/>进度: ${progressValue}<br/>任务完成率: ${completionValue}`;
					},
				},
				legend: { data: ["项目进度", "进度趋势", "任务完成率"], bottom: 0, textStyle: { color: "#475569" } },
				grid: { left: "3%", right: "4%", bottom: "15%", top: "10%", containLabel: true },
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
				completionSeries: this.fillSeriesGaps(completionValues, fallbackCompletion),
			};
		},
		fillSeriesGaps(values, defaultValue) {
			let last = typeof defaultValue === "number" && !Number.isNaN(defaultValue) ? defaultValue : 0;
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
			return this.parseDateString(task?.due_date || task?.due_at || task?.created_at);
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
				return { type: "warning", label: diffDays === 0 ? "今日截止" : `剩余 ${diffDays} 天` };
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
			this.showCreateModal = true;
			this.newTaskForm = { title: "", description: "", due_date: "", owner_name: "" };
		},
		closeCreateModal() {
			if (this.creating) return;
			this.showCreateModal = false;
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
			const payload = {
				title: this.newTaskForm.title.trim(),
				description: this.newTaskForm.description?.trim?.() ? this.newTaskForm.description.trim() : this.newTaskForm.description || "",
				due_at: this.newTaskForm.due_date || null,
				owner_name: this.newTaskForm.owner_name || null,
				task_type: 2,
			};
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
				this.updateChart();
			}
		},
		buildLocalTask(payload, extra = {}) {
			return {
				id: extra.id || Date.now(),
				title: payload.title,
				description: payload.description || "",
				owner_name: payload.owner_name || "你",
				due_date: payload.due_at || payload.due_date || "",
				created_at: extra.created_at || new Date().toISOString(),
				status: extra.status || "in-progress",
				status_label: extra.status_label || "",
				progress: this.clampProgress(extra.progress ?? 0),
			};
		},
		toggleTaskDetails(taskId) {
			const index = this.expandedTaskIds.indexOf(taskId);
			if (index >= 0) {
				this.expandedTaskIds.splice(index, 1);
			} else {
				this.expandedTaskIds.push(taskId);
			}
		},
		isTaskExpanded(taskId) {
			return this.expandedTaskIds.includes(taskId);
		},
		getTaskDetail(task) {
			if (!task) {
				return { subtasks: [], attachments: [], comments: [] };
			}
			if (!this.taskDetailCache[task.id]) {
				const detail = this.buildTaskDetail(task);
				this.taskDetailCache = { ...this.taskDetailCache, [task.id]: detail };
			}
			return this.taskDetailCache[task.id];
		},
		buildTaskDetail(task) {
			const baseTitle = task.title || "任务";
			const progress = this.getTaskProgressValue(task);
			const subtasks = (task.subtasks && task.subtasks.length
				? task.subtasks
				: [
					{ id: `${task.id}-sub1`, title: `${baseTitle} - 需求梳理`, status: progress >= 20 ? "进行中" : "待开始" },
					{ id: `${task.id}-sub2`, title: `${baseTitle} - 开发实现`, status: progress >= 60 ? "进行中" : "待开始" },
					{ id: `${task.id}-sub3`, title: `${baseTitle} - 测试验收`, status: progress >= 90 ? "待验收" : "待开始" },
				]);
			const attachments = task.attachments && task.attachments.length
				? task.attachments
				: [
					{ id: `${task.id}-att1`, name: `${baseTitle}-设计稿.fig`, size: "2.3MB" },
					{ id: `${task.id}-att2`, name: `${baseTitle}-需求文档.docx`, size: "480KB" },
				];
			const comments = task.comments && task.comments.length
				? task.comments
				: [
					{ id: `${task.id}-c1`, author: "王同学", content: "今天把接口联调完毕了。", time: "今天 10:20" },
					{ id: `${task.id}-c2`, author: "李同学", content: "测试脚本已准备，等你们提测。", time: "昨天 17:45" },
				];
			return { subtasks, attachments, comments };
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
	},
};
</script>

<style scoped>
	:global(body) {
		background: #f5f7fb;
		font-family: "PingFang SC", "Segoe UI", system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
	}

	.surface-card {
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.95), rgba(245, 247, 251, 0.95));
		border-radius: 12px;
		box-shadow: 0 10px 40px rgba(15, 23, 42, 0.08);
		padding: 24px;
	}

	.card.surface-card + .surface-card {
		margin-top: 16px;
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
		background: linear-gradient(145deg, rgba(255, 255, 255, 0.95), rgba(249, 250, 251, 0.9));
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
		background: radial-gradient(circle at top, rgba(79, 70, 229, 0.08), transparent 45%);
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

	@keyframes blob {
		0% { transform: translate(0px, 0px) scale(1); }
		33% { transform: translate(30px, -50px) scale(1.1); }
		66% { transform: translate(-20px, 20px) scale(0.9); }
		100% { transform: translate(0px, 0px) scale(1); }
	}
	.animate-blob {
		animation: blob 7s infinite;
	}
	.animation-delay-2000 {
		animation-delay: 2s;
	}
	.animation-delay-4000 {
		animation-delay: 4s;
	}
	@keyframes fadeIn {
		from { opacity: 0; transform: translateY(10px); }
		to { opacity: 1; transform: translateY(0); }
	}
	.animate-fade-in {
		animation: fadeIn 0.5s ease-out forwards;
	}
</style>

