/**
 * 学习相关API
 */

import { request } from "@/utils/request";

/**
 * 获取学习房间列表
 */
export function getStudyRooms(params = {}) {
  return request.get("/study/rooms", params);
}

/**
 * 获取学习房间详情
 */
export function getStudyRoomDetail(roomId) {
  return request.get(`/study/rooms/${roomId}`);
}

/**
 * 创建学习房间
 */
export function createStudyRoom(data) {
  return request.post("/study/rooms", data);
}

// 加入学习房间（支持私密房间密码）
export function joinStudyRoom(roomId, data = {}) {
  return request.post(`/study/rooms/${roomId}/join`, data);
}

// 学习房间汇总数据
export function getStudySummary(params = {}) {
  return request.get("/study/summary", params);
}

// 用户学习记录
export function getStudyRecords(params = {}) {
  return request.get("/study/records", params);
}

/**
 * 离开学习房间
 */
export function leaveStudyRoom(roomId) {
  return request.post(`/study/rooms/${roomId}/leave`);
}

/**
 * 获取房间在线用户
 */
export function getRoomOnlineUsers(roomId) {
  return request.get(`/study/rooms/${roomId}/online-users`);
}

/**
 * 开始学习记录
 */
export function startStudySession(data) {
  return request.post("/study/sessions/start", data);
}

/**
 * 结束学习记录
 */
export function endStudySession(sessionId) {
  return request.post(`/study/sessions/${sessionId}/end`);
}

/**
 * 获取学习记录列表
 */
export function getStudySessions(params = {}) {
  return request.get("/study/sessions", params);
}

/**
 * 获取学习统计数据
 */
export function getStudyStatistics(params = {}) {
  return request.get("/study/statistics", params);
}

/**
 * 获取学习时长趋势
 */
export function getStudyTimeTrend(params = {}) {
  return request.get("/study/time-trend", params);
}

/**
 * 获取学习房间聊天记录
 */
export function getRoomChatHistory(roomId, params = {}) {
  return request.get(`/study/rooms/${roomId}/chat/history`, params);
}

/**
 * 发送学习房间消息
 */
export function sendRoomChatMessage(roomId, data) {
  return request.post(`/study/rooms/${roomId}/chat`, data);
}

/**
 * 获取知识点分布
 */
export function getKnowledgeDistribution(params = {}) {
  return request.get("/study/knowledge-distribution", params);
}

/**
 * 获取学习计划
 */
export function getStudyPlans(params = {}) {
  return request.get("/study/plans", params);
}

/**
 * 创建学习计划
 */
export function createStudyPlan(data) {
  return request.post("/study/plans", data);
}

/**
 * 更新学习计划
 */
export function updateStudyPlan(planId, data) {
  return request.put(`/study/plans/${planId}`, data);
}

/**
 * 删除学习计划
 */
export function deleteStudyPlan(planId) {
  return request.delete(`/study/plans/${planId}`);
}

/**
 * 获取学习笔记
 */
export function getStudyNotes(params = {}) {
  return request.get("/study/notes", params);
}

/**
 * 创建学习笔记
 */
export function createStudyNote(data) {
  return request.post("/study/notes", data);
}

/**
 * 更新学习笔记
 */
export function updateStudyNote(noteId, data) {
  return request.put(`/study/notes/${noteId}`, data);
}

/**
 * 删除学习笔记
 */
export function deleteStudyNote(noteId) {
  return request.delete(`/study/notes/${noteId}`);
}
