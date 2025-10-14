/**
 * 通知相关API
 */

import { request } from '@/utils/request'

/**
 * 获取通知列表
 */
export function getNotificationList(params = {}) {
  return request.get('/notifications', params)
}

/**
 * 获取未读通知数量
 */
export function getUnreadNotificationCount() {
  return request.get('/notifications/unread-count')
}

/**
 * 标记通知为已读
 */
export function markNotificationAsRead(notificationId) {
  return request.put(`/notifications/${notificationId}/read`)
}

/**
 * 批量标记通知为已读
 */
export function markNotificationsAsRead(notificationIds) {
  return request.put('/notifications/batch-read', { notificationIds })
}

/**
 * 标记所有通知为已读
 */
export function markAllNotificationsAsRead() {
  return request.put('/notifications/read-all')
}

/**
 * 删除通知
 */
export function deleteNotification(notificationId) {
  return request.delete(`/notifications/${notificationId}`)
}

/**
 * 批量删除通知
 */
export function deleteNotifications(notificationIds) {
  return request.delete('/notifications/batch-delete', { data: { notificationIds } })
}

/**
 * 清空所有通知
 */
export function clearAllNotifications() {
  return request.delete('/notifications/clear-all')
}

/**
 * 获取通知设置
 */
export function getNotificationSettings() {
  return request.get('/notifications/settings')
}

/**
 * 更新通知设置
 */
export function updateNotificationSettings(settings) {
  return request.put('/notifications/settings', settings)
}

/**
 * 发送系统通知
 */
export function sendSystemNotification(data) {
  return request.post('/notifications/system', data)
}

/**
 * 发送团队通知
 */
export function sendTeamNotification(teamId, data) {
  return request.post(`/notifications/team/${teamId}`, data)
}

/**
 * 发送个人通知
 */
export function sendPersonalNotification(userId, data) {
  return request.post(`/notifications/personal/${userId}`, data)
}

/**
 * 获取通知模板
 */
export function getNotificationTemplates() {
  return request.get('/notifications/templates')
}

/**
 * 创建通知模板
 */
export function createNotificationTemplate(data) {
  return request.post('/notifications/templates', data)
}

/**
 * 更新通知模板
 */
export function updateNotificationTemplate(templateId, data) {
  return request.put(`/notifications/templates/${templateId}`, data)
}

/**
 * 删除通知模板
 */
export function deleteNotificationTemplate(templateId) {
  return request.delete(`/notifications/templates/${templateId}`)
}

/**
 * 订阅通知
 */
export function subscribeNotification(type, targetId) {
  return request.post('/notifications/subscribe', { type, targetId })
}

/**
 * 取消订阅通知
 */
export function unsubscribeNotification(type, targetId) {
  return request.post('/notifications/unsubscribe', { type, targetId })
}

/**
 * 获取订阅列表
 */
export function getSubscriptionList() {
  return request.get('/notifications/subscriptions')
}