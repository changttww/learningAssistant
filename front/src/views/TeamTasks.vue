<template>
	<div class="w-full h-full overflow-auto px-4">
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<!-- 左侧主要内容 -->
			<div class="lg:col-span-2 space-y-6">
				<div class="card">
					<div class="flex items-center justify-between mb-6">
						<h1 class="text-2xl font-bold">团队任务</h1>
						<button
							@click="openCreateModal"
							class="bg-[#2D5BFF] text-white font-medium py-2 px-4 rounded-lg text-sm hover:bg-opacity-90 transition-colors flex items-center gap-2"
						>
							<iconify-icon
								icon="mdi:plus"
								width="16"
								height="16"
							></iconify-icon>
							创建任务
						</button>
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
						<h2 class="font-bold text-xl mb-4">团队进度概览</h2>
						<div class="chart-container" ref="teamProgressChart"></div>
					</div>

					<!-- 任务列表 -->
					<div class="space-y-4">
						<h3 class="font-bold text-lg">当前任务</h3>

						<div v-for="task in tasks" :key="task.id" class="p-4 border border-gray-200 rounded-lg">
							<div class="flex items-start justify-between">
								<div class="flex-1">
									<h4 class="font-medium">{{ task.title }}</h4>
									<p class="text-sm text-gray-600 mt-1">{{ task.description }}</p>
									<div class="flex items-center gap-4 mt-3">
										<div class="flex items-center gap-2">
											<iconify-icon icon="mdi:account" width="16" height="16" class="text-gray-500"></iconify-icon>
											<span class="text-sm">{{ task.owner_name || '未知' }}</span>
										</div>
										<div class="flex items-center gap-2">
											<iconify-icon icon="mdi:calendar" width="16" height="16" class="text-gray-500"></iconify-icon>
											<span class="text-sm">{{ task.due_date || task.created_at || '' }}</span>
										</div>
									</div>
									<div class="mt-3">
										<div class="flex items-center justify-between mb-1">
											<span class="text-sm">进度</span>
											<span class="text-sm font-medium">{{ computeTaskProgressLabel(task) }}</span>
										</div>
										<div class="progress-bar">
											<div class="progress-fill" :style="{ width: computeTaskProgressPercent(task) + '%' }"></div>
										</div>
										<div class="flex items-center gap-2 mt-2">
											<button @click="changeProgress(task, -10)" class="px-2 py-1 text-xs bg-gray-100 rounded">-10%</button>
											<button @click="changeProgress(task, 10)" class="px-2 py-1 text-xs bg-gray-100 rounded">+10%</button>
											<button @click="setCompleted(task)" class="px-2 py-1 text-xs bg-green-100 rounded">标记完成</button>
										</div>
									</div>
								</div>
								<span :class="task.status === 'completed' ? 'bg-green-100 text-green-800 px-2 py-1 rounded text-xs' : 'bg-blue-100 text-blue-800 px-2 py-1 rounded text-xs'">{{ task.status_label || task.status }}</span>
							</div>
						</div>
					</div>
				</div>
			</div>

			<!-- 右侧栏 -->
			<div class="space-y-5">
				<div class="card">
					<div class="flex justify-between items-center">
						<h3 class="font-bold text-lg">团队动态</h3>
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

				<div class="card">
					<div class="flex justify-between items-center">
						<h3 class="font-bold text-lg">团队积分排名</h3>
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

				<div class="card">
					<h3 class="font-bold text-lg mb-4">协作工具</h3>
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

export default {
	name: "TeamTasks",
	data() {
		return {
			tasks: [
				// 保留一些本地示例数据，loadTasks 会尝试覆盖
				{
					id: 1,
					title: "登录功能开发",
					description: "实现用户登录、注册和密码重置功能",
					owner_name: "王同学",
					due_date: "2024-01-18",
					status: "in-progress",
					status_label: "进行中",
					progress: 85,
				},
				{
					id: 2,
					title: "支付模块设计",
					description: "设计支付流程和UI界面",
					owner_name: "陈同学",
					due_date: "2024-01-22",
					status: "in-progress",
					status_label: "设计中",
					progress: 30,
				},
			],
			chart: null,
			// 创建任务模态相关状态
			showCreateModal: false,
			creating: false,
			newTaskForm: {
				title: "",
				description: "",
				due_date: "",
				owner_name: "",
			},
		};
	},
	mounted() {
		this.initChart();
		this.loadTasks();
		window.addEventListener("resize", () => {
			if (this.chart) this.chart.resize();
		});
	},
	methods: {
		// 从后端拉取团队任务，若失败则保留本地示例
		async loadTasks() {
			try {
				const res = await getTeamTasks();
				// response shape depends on backend: try common patterns
				const items = res?.data?.items || res?.data || res;
				if (Array.isArray(items)) {
					this.tasks = items.map((t) => {
						let status = "pending";
						if (t.status === 2 || t.status === "completed") {
							status = "completed";
						} else if (t.status === 1 || t.status === "in-progress") {
							status = "in-progress";
						}
						let progress;
						if (t.progress !== undefined) progress = t.progress;
						else if (status === "completed") progress = 100;
						else if (status === "in-progress") progress = 50;
						else progress = 0;

						return {
							id: t.id,
							title: t.title || t.name,
							description: t.description || "",
							owner_name: t.owner_name || (t.created_by_name || "未知"),
							due_date: t.due_at ? t.due_at.split("T")[0] : (t.due_date || ""),
							created_at: t.created_at || "",
							status,
							status_label: t.status_label || "",
							progress,
						};
					});
				}
			} catch (err) {
				// 保持本地示例并记录错误
				console.warn("加载团队任务失败，使用本地示例：", err);
			} finally {
				this.updateChart();
			}
		},

		computeTaskProgressPercent(task) {
			if (typeof task.progress === "number") return Math.max(0, Math.min(100, task.progress));
			// 根据状态映射为百分比
			if (task.status === "completed") return 100;
			if (task.status === "in-progress") return 50;
			return 0;
		},

		computeTaskProgressLabel(task) {
			const p = this.computeTaskProgressPercent(task);
			return `${p}%`;
		},

		// 计算团队总体进度：使用所有任务的平均进度
		computeTeamProgress() {
			if (!this.tasks || this.tasks.length === 0) return 0;
			const sum = this.tasks.reduce((s, t) => s + this.computeTaskProgressPercent(t), 0);
			return Math.round(sum / this.tasks.length);
		},

		// 计算任务完成率：已完成任务数 / 总任务数 * 100
		computeCompletionRate() {
			if (!this.tasks || this.tasks.length === 0) return 0;
			const completed = this.tasks.filter((t) => this.computeTaskProgressPercent(t) >= 100).length;
			return Math.round((completed / this.tasks.length) * 100);
		},

		initChart() {
			const chartDom = this.$refs.teamProgressChart;
			this.chart = echarts.init(chartDom);
			this.updateChart();
		},

		updateChart() {
			if (!this.chart) return;
			const teamProgress = this.computeTeamProgress();
			const completionRate = this.computeCompletionRate();
			const option = {
				tooltip: {
					trigger: "axis",
					formatter: (params) => {
						return `${params[0].name}<br/>进度: ${params[0].value}%<br/>任务完成率: ${completionRate}%`;
					},
				},
				legend: { data: ["项目进度", "任务完成率"], bottom: 0 },
				grid: { left: "3%", right: "4%", bottom: "15%", top: "10%", containLabel: true },
				xAxis: { type: "category", data: ["第一周", "第二周", "第三周", "第四周", "当前周"] },
				yAxis: { type: "value", axisLabel: { formatter: "{value}%" }, max: 100 },
				series: [
					{
						name: "项目进度",
						type: "bar",
						data: [20, 35, 45, 60, teamProgress],
						itemStyle: { color: "#2D5BFF" },
						barWidth: "40%",
					},
					{
						name: "任务完成率",
						type: "line",
						data: [25, 40, 55, 70, completionRate],
						smooth: true,
						lineStyle: { width: 3, color: "#4CAF50" },
						symbol: "circle",
						symbolSize: 8,
					},
				],
			};
			this.chart.setOption(option);
		},

		// 打开/关闭模态并提交创建任务
		openCreateModal() {
			this.showCreateModal = true;
			// reset form
			this.newTaskForm = { title: "", description: "", due_date: "", owner_name: "" };
		},
		closeCreateModal() {
			if (this.creating) return;
			this.showCreateModal = false;
		},
		async submitCreateTask() {
			if (!this.newTaskForm.title || !this.newTaskForm.title.trim()) {
				// 简单提示，可替换为内联错误提示
				alert("任务名称不能为空");
				return;
			}
			const payload = {
				title: this.newTaskForm.title.trim(),
				description: this.newTaskForm.description || "",
				due_at: this.newTaskForm.due_date || null,
				owner_name: this.newTaskForm.owner_name || null,
				task_type: 2,
			};

			this.creating = true;
			try {
				const res = await createTask(payload);
				const created = res?.data?.data || res?.data || res;
				let status = "pending";
				if (created.status === 2) status = "completed";
				else if (created.status === 1) status = "in-progress";

				let progress;
				if (created.progress !== undefined) progress = created.progress;
				else if (status === "completed") progress = 100;
				else if (status === "in-progress") progress = 50;
				else progress = 0;
				let dueDate = "";
				if (created.due_at) {
					if (created.due_at.split) dueDate = created.due_at.split("T")[0];
					else dueDate = created.due_at;
				} else {
					dueDate = created.due_date || payload.due_at || "";
				}

				const newTask = {
					id: created.id || Date.now(),
					title: created.title || payload.title,
					description: created.description || payload.description,
					owner_name: created.owner_name || payload.owner_name || "你",
					due_date: dueDate,
					status,
					status_label: created.status_label || "",
					progress,
				};
				this.tasks.unshift(newTask);
				this.showCreateModal = false;
			} catch (err) {
				console.warn("创建任务请求失败，已在本地创建：", err);
				const newTask = {
					id: Date.now(),
					title: payload.title,
					description: payload.description,
					owner_name: payload.owner_name || "你",
					due_date: payload.due_at || "",
					status: "in-progress",
					status_label: "进行中",
					progress: 0,
				};
				this.tasks.unshift(newTask);
				this.showCreateModal = false;
			} finally {
				this.creating = false;
				this.updateChart();
			}
		},

		async changeProgress(task, delta) {
			const old = typeof task.progress === 'number' ? task.progress : 0;
			const next = Math.max(0, Math.min(100, old + delta));
			// 乐观更新
			task.progress = next;
			if (next >= 100) task.status = 'completed';
			else if (next > 0) task.status = 'in-progress';
			this.updateChart();

			try {
				await updateTaskProgress(task.id, next);
				// 成功：可根据后端返回进一步更新（此处期待后端返回最新数据）
			} catch (err) {
				// 回退并提示
				task.progress = old;
				if (old >= 100) task.status = 'completed';
				else if (old > 0) task.status = 'in-progress';
				else task.status = 'pending';
				this.updateChart();
				console.warn('更新进度失败', err);
				alert('更新进度失败，请重试');
			}
		},

		async setCompleted(task) {
			const cur = typeof task.progress === 'number' ? task.progress : 0;
			await this.changeProgress(task, 100 - cur);
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

