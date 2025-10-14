/**
 * 任务相关API
 */

import { request } from "@/utils/request";

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
  return request.post("/tasks", data);
}

/**
 * 更新任务
 */
export function updateTask(taskId, data) {
  return request.put(`/tasks/${taskId}`, data);
}

/**
 * 删除任务
 */
export function deleteTask(taskId) {
  return request.delete(`/tasks/${taskId}`);
}

/**
 * 完成任务
 */
export function completeTask(taskId) {
  return request.post(`/tasks/${taskId}/complete`);
}

/**
 * 取消完成任务
 */
export function uncompleteTask(taskId) {
  return request.post(`/tasks/${taskId}/uncomplete`);
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
  return request.put(`/tasks/${taskId}/progress`, { progress });
}

/**
 * 获取任务统计
 */
export function getTaskStatistics(params = {}) {
  return request.get("/tasks/statistics", params);
}

/**
 * 获取任务分类
 */
export function getTaskCategories() {
  return request.get("/tasks/categories");
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
