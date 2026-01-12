/**
 * 任务相关API
 */

import { request } from "@/utils/request";

/**
 * 触发全局任务更新事件
 */
function emitTaskUpdateEvent(eventType = "taskUpdated") {
  if (typeof window !== "undefined") {
    const event = new CustomEvent(eventType, {
      detail: { timestamp: Date.now() },
    });
    window.dispatchEvent(event);
    console.log(`[任务事件] 已触发: ${eventType}`);
  }
}

/**
 * 获取任务列表
 */
export function getTaskList(params = {}) {
  return request.get("/tasks", params);
}

/**
 * 获取个人任务列表
 */
export function getPersonalTasks(params = {}) {
  return request.get("/tasks/personal", params);
}

/**
 * 获取团队任务列表
 */
export function getTeamTasks(params = {}) {
  return request.get("/tasks/team", params);
}

/**
 * 获取任务详情
 */
export function getTaskDetail(taskId) {
  return request.get(`/tasks/${taskId}`);
}

/**
 * 创建任务
 */
export function createTask(data) {
  return request.post("/tasks", data).then((res) => {
    emitTaskUpdateEvent("taskCreated");
    return res;
  });
}

/**
 * 更新任务
 */
export function updateTask(taskId, data) {
  return request.put(`/tasks/${taskId}`, data).then((res) => {
    emitTaskUpdateEvent("taskUpdated");
    return res;
  });
}

/**
 * 添加任务评论
 */
export function addTaskComment(taskId, content) {
  return request.post(`/tasks/${taskId}/comments`, { content });
}

/**
 * 删除任务
 */
export function deleteTask(taskId) {
  return request.delete(`/tasks/${taskId}`).then((res) => {
    emitTaskUpdateEvent("taskUpdated");
    return res;
  });
}

/**
 * 完成任务
 */
export function completeTask(taskId) {
  return request.post(`/tasks/${taskId}/complete`).then((res) => {
    emitTaskUpdateEvent("taskCompleted");
    return res;
  });
}

/**
 * 获取任务分类
 */
export function getTaskCategories() {
  return request.get("/tasks/categories");
}

/**
 * 取消完成任务
 */
export function uncompleteTask(taskId) {
  return request.post(`/tasks/${taskId}/uncomplete`).then((res) => {
    emitTaskUpdateEvent("taskUpdated");
    return res;
  });
}

/**
 * 获取任务进度
 */
export function getTaskProgress(taskId) {
  return request.get(`/tasks/${taskId}/progress`);
}

/**
 * 更新任务进度
 */
export function updateTaskProgress(taskId, progress) {
  return request.put(`/tasks/${taskId}`, { progress }).then((res) => {
    emitTaskUpdateEvent("taskUpdated");
    return res;
  });
}

/**
 * 获取任务统计
 */
export function getTaskStatistics(params = {}) {
  return request.get("/tasks/statistics", params);
}

/**
 * 获取任务标签
 */
export function getTaskTags() {
  return request.get("/tasks/tags");
}

/**
 * 批量操作任务
 */
export function batchOperateTasks(operation, taskIds) {
  return request.post("/tasks/batch", {
    operation,
    taskIds,
  }).then((res) => {
    emitTaskUpdateEvent("taskUpdated");
    return res;
  });
}

/**
 * 搜索任务
 */
export function searchTasks(keyword, params = {}) {
  return request.get("/tasks/search", {
    keyword,
    ...params,
  });
}

/**
 * AI 解析自然语言任务
 */
export function parseTaskWithAI(input) {
  return request.post("/tasks/ai/parse", { input }, { timeout: 60000 });
}

/**
 * 获取任务指导和相关资源
 * @param {Object|string} titleOrParams - 任务标题或包含 title, description, category 的对象
 * @param {string} [description] - 任务描述（当第一个参数为字符串时使用）
 * @param {string} [category] - 任务分类（当第一个参数为字符串时使用）
 */
export function getTaskGuidance(titleOrParams, description, category) {
  // 支持对象参数和独立参数两种调用方式
  if (typeof titleOrParams === 'object' && titleOrParams !== null) {
    const { title, description: desc, category: cat } = titleOrParams;
    return request.post("/tasks/ai/guidance", { title, description: desc, category: cat }, { timeout: 60000 });
  }
  return request.post("/tasks/ai/guidance", { title: titleOrParams, description, category }, { timeout: 60000 });
}

/**
 * AI 生成智能测验
 */
export function generateQuiz(data) {
  return request.post("/tasks/ai/quiz", data, { timeout: 60000 });
}

/**
 * 获取今日任务
 */
export function getTodayTasks(userId) {
  return request.get(`/tasks/users/${userId}/today`);
}

/**
 * 获取任务柱状统计
 */
export function getTaskBarStats(range = "week") {
  return request.get("/tasks/stats/bar", { range });
}

/**
 * 获取近期月度完成率
 */
export function getMonthlyCompletionStats() {
  return request.get("/tasks/stats/monthly-completion");
}

/**
 * 获取任务热力图数据
 */
export function getTaskHeatmapStats() {
  return request.get("/tasks/stats/heatmap");
}
