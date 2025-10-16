/**
 * 团队相关API
 */

import { request } from "@/utils/request";

/**
 * 获取团队列表
 */
export function getTeamList(params = {}) {
  return request.get("/teams", params);
}

/**
 * 获取团队详情
 */
export function getTeamDetail(teamId) {
  return request.get(`/teams/${teamId}`);
}

/**
 * 创建团队
 */
export function createTeam(data) {
  return request.post("/teams", data);
}

/**
 * 更新团队信息
 */
export function updateTeam(teamId, data) {
  return request.put(`/teams/${teamId}`, data);
}

/**
 * 删除团队
 */
export function deleteTeam(teamId) {
  return request.delete(`/teams/${teamId}`);
}

/**
 * 获取团队成员列表
 */
export function getTeamMembers(teamId, params = {}) {
  return request.get(`/teams/${teamId}/members`, params);
}

/**
 * 邀请成员加入团队
 */
export function inviteTeamMember(teamId, data) {
  return request.post(`/teams/${teamId}/invite`, data);
}

/**
 * 移除团队成员
 */
export function removeTeamMember(teamId, memberId) {
  return request.delete(`/teams/${teamId}/members/${memberId}`);
}

/**
 * 更新成员角色
 */
export function updateMemberRole(teamId, memberId, role) {
  return request.put(`/teams/${teamId}/members/${memberId}/role`, { role });
}

/**
 * 加入团队
 */
export function joinTeam(teamId, inviteCode) {
  return request.post(`/teams/${teamId}/join`, { inviteCode });
}

/**
 * 退出团队
 */
export function leaveTeam(teamId) {
  return request.post(`/teams/${teamId}/leave`);
}

/**
 * 获取团队进度统计
 */
export function getTeamProgress(teamId) {
  return request.get(`/teams/${teamId}/progress`);
}

/**
 * 获取团队活动记录
 */
export function getTeamActivities(teamId, params = {}) {
  return request.get(`/teams/${teamId}/activities`, params);
}

/**
 * 获取团队积分排行
 */
export function getTeamRanking(teamId) {
  return request.get(`/teams/${teamId}/ranking`);
}

/**
 * 获取团队任务统计
 */
export function getTeamTaskStatistics(teamId) {
  return request.get(`/teams/${teamId}/task-statistics`);
}

/**
 * 创建团队公告
 */
export function createTeamAnnouncement(teamId, data) {
  return request.post(`/teams/${teamId}/announcements`, data);
}

/**
 * 获取团队公告列表
 */
export function getTeamAnnouncements(teamId, params = {}) {
  return request.get(`/teams/${teamId}/announcements`, params);
}

/**
 * 获取团队设置
 */
export function getTeamSettings(teamId) {
  return request.get(`/teams/${teamId}/settings`);
}

/**
 * 更新团队设置
 */
export function updateTeamSettings(teamId, settings) {
  return request.put(`/teams/${teamId}/settings`, settings);
}
