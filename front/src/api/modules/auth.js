/**
 * 认证相关API
 */

import { request } from '@/utils/request'

/**
 * 用户登录
 */
export function login(data) {
  return request.post('/auth/login', data)
}

/**
 * 用户注册
 */
export function register(data) {
  return request.post('/auth/register', data)
}

/**
 * 用户登出
 */
export function logout() {
  return request.post('/auth/logout')
}

/**
 * 刷新token
 */
export function refreshToken(refreshToken) {
  return request.post('/auth/refresh', { refreshToken })
}

/**
 * 获取用户信息
 */
export function getUserInfo() {
  return request.get('/auth/user-info')
}

/**
 * 修改密码
 */
export function changePassword(data) {
  return request.put('/auth/change-password', data)
}

/**
 * 忘记密码
 */
export function forgotPassword(email) {
  return request.post('/auth/forgot-password', { email })
}

/**
 * 重置密码
 */
export function resetPassword(data) {
  return request.post('/auth/reset-password', data)
}

/**
 * 验证邮箱
 */
export function verifyEmail(token) {
  return request.post('/auth/verify-email', { token })
}

/**
 * 发送验证码
 */
export function sendVerificationCode(email) {
  return request.post('/auth/send-verification-code', { email })
}