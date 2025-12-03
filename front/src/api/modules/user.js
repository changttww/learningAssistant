/**
 * 用户相关API
 */

import { request } from "@/utils/request";

/**
 * 获取用户详细信息
 */
export function getUserProfile(userId) {
  return request.get(`/users/${userId}`);
}

/**
 * 更新用户信息
 */
export function updateUserProfile(userId, data) {
  return request.put(`/users/${userId}`, data);
}

/**
 * 上传用户头像
 */
export function uploadAvatar(formData) {
  return request.upload("/users/avatar", formData);
}

/**
 * 获取用户学习统计
 */
export function getUserStudyStats(userId) {
  return request.get(`/users/${userId}/study-stats`);
}

/**
 * 获取用户积分账本
 */
export function getUserPointsLedger(userId, params = {}) {
  return request.get(`/users/${userId}/points/ledger`, params);
}

/**
 * 用户签到
 */
export function checkInUser(userId) {
  return request.post(`/users/${userId}/check-in`);
}

/**
 * 获取用户成就列表
 */
export function getUserAchievements(userId) {
  return request.get(`/users/${userId}/achievements`);
}

/**
 * 获取用户技能标签
 */
export function getUserSkills(userId) {
  return request.get(`/users/${userId}/skills`);
}

/**
 * 更新用户技能标签
 */
export function updateUserSkills(userId, skills) {
  return request.put(`/users/${userId}/skills`, { skills });
}

/**
 * 获取用户学习记录
 */
export function getUserStudyRecords(userId, params = {}) {
  return request.get(`/users/${userId}/study-records`, params);
}

/**
 * 获取用户设置
 */
export function getUserSettings(userId) {
  return request.get(`/users/${userId}/settings`);
}

/**
 * 更新用户设置
 */
export function updateUserSettings(userId, settings) {
  return request.put(`/users/${userId}/settings`, settings);
}

/**
 * 获取用户通知偏好
 */
export function getUserNotificationPreferences(userId) {
  return request.get(`/users/${userId}/notification-preferences`);
}

/**
 * 更新用户通知偏好
 */
export function updateUserNotificationPreferences(userId, preferences) {
  return request.put(`/users/${userId}/notification-preferences`, preferences);
}

/**
 * 学习伙伴接口
 */
export function listStudyBuddies(userId) {
  return request.get(`/users/${userId}/buddies`);
}

export function addStudyBuddy(userId, data) {
  return request.post(`/users/${userId}/buddies`, data);
}

export function updateStudyBuddy(userId, buddyId, data) {
  return request.put(`/users/${userId}/buddies/${buddyId}`, data);
}

export function deleteStudyBuddy(userId, buddyId) {
  return request.delete(`/users/${userId}/buddies/${buddyId}`);
}
