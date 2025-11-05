/**
 * 认证工具
 * 管理用户认证相关的token操作
 */

import { authConfig } from "@/config";

/**
 * 获取token
 */
export function getToken() {
  return localStorage.getItem(authConfig.tokenKey);
}

/**
 * 设置token
 */
export function setToken(token) {
  localStorage.setItem(authConfig.tokenKey, token);
  // 设置过期时间
  const expireTime = Date.now() + authConfig.tokenExpireTime;
  localStorage.setItem(`${authConfig.tokenKey}_expire`, expireTime.toString());
}

/**
 * 移除token
 */
export function removeToken() {
  localStorage.removeItem(authConfig.tokenKey);
  localStorage.removeItem(`${authConfig.tokenKey}_expire`);
  localStorage.removeItem(authConfig.refreshTokenKey);
}

/**
 * 检查token是否过期
 */
export function isTokenExpired() {
  const token = getToken();
  if (!token) return true;

  const expireTime = localStorage.getItem(`${authConfig.tokenKey}_expire`);
  if (!expireTime) return true;

  return Date.now() > parseInt(expireTime);
}

/**
 * 获取刷新token
 */
export function getRefreshToken() {
  return localStorage.getItem(authConfig.refreshTokenKey);
}

/**
 * 设置刷新token
 */
export function setRefreshToken(refreshToken) {
  localStorage.setItem(authConfig.refreshTokenKey, refreshToken);
}

/**
 * 清除所有认证信息
 */
export function clearAuth() {
  removeToken();
  // 清除用户信息
  localStorage.removeItem("userInfo");
  localStorage.removeItem("permissions");
  localStorage.removeItem("roles");
}

/**
 * 获取用户信息
 */
export function getUserInfo() {
  const userInfo = localStorage.getItem("userInfo");
  return userInfo ? JSON.parse(userInfo) : null;
}

/**
 * 设置用户信息
 */
export function setUserInfo(userInfo) {
  if (!userInfo) {
    localStorage.removeItem("userInfo");
    return;
  }
  localStorage.setItem("userInfo", JSON.stringify(userInfo));
}

/**
 * 获取用户权限
 */
export function getPermissions() {
  const permissions = localStorage.getItem("permissions");
  return permissions ? JSON.parse(permissions) : [];
}

/**
 * 设置用户权限
 */
export function setPermissions(permissions) {
  localStorage.setItem("permissions", JSON.stringify(permissions));
}

/**
 * 检查是否有指定权限
 */
export function hasPermission(permission) {
  const permissions = getPermissions();
  return permissions.includes(permission);
}

/**
 * 获取用户角色
 */
export function getRoles() {
  const roles = localStorage.getItem("roles");
  return roles ? JSON.parse(roles) : [];
}

/**
 * 设置用户角色
 */
export function setRoles(roles) {
  localStorage.setItem("roles", JSON.stringify(roles));
}

/**
 * 检查是否有指定角色
 */
export function hasRole(role) {
  const roles = getRoles();
  return roles.includes(role);
}
