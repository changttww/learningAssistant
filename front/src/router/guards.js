/**
 * 路由守卫
 * 实现权限验证、登录状态检查、页面访问控制
 */

import {
  getToken,
  isTokenExpired,
  getUserInfo,
  hasPermission,
  clearAuth,
} from "../utils/auth";
import { ElMessage } from "element-plus";

// 白名单路由（不需要登录即可访问）
const whiteList = [
  "/login",
  "/register",
  "/forgot-password",
  "/404",
  "/403",
  "/personal-tasks",
  "/task-manager",
  "/profile",
  "/study-room",
  "/team-tasks",
  "/video-room/*",
  "/",
];

// 需要权限的路由配置
const permissionRoutes = {
  // '/team-tasks': ['team:view'],
  // '/task-manager': ['task:manage'],
  // '/profile': ['profile:view'],
  // '/video-room/:roomId': ['room:access']
};

/**
 * 前置守卫
 */
export function beforeEach(to, from, next) {
  // 开始页面加载进度条
  startProgress();

  const token = getToken();
  const userInfo = getUserInfo();

  // 如果有token
  if (token) {
    // 检查token是否过期
    if (isTokenExpired()) {
      ElMessage.warning("登录已过期，请重新登录");
      clearAuth();
      next("/login");
      return;
    }

    // 如果要去登录页，直接跳转到首页
    if (to.path === "/login") {
      next("/");
      return;
    }

    // 检查用户信息是否存在
    if (!userInfo) {
      // 获取用户信息失败，清除认证信息
      clearAuth();
      next("/login");
      return;
    }

    // 检查页面权限
    if (!checkPagePermission(to.path)) {
      ElMessage.error("您没有权限访问该页面");
      next("/403");
      return;
    }

    next();
  } else {
    // 没有token
    if (isInWhiteList(to.path)) {
      // 在白名单中，直接通过
      next();
    } else {
      // 不在白名单中，跳转到登录页
      ElMessage.warning("请先登录");
      next("/login");
    }
  }
}

/**
 * 后置守卫
 */
export function afterEach(to, from) {
  // 结束页面加载进度条
  finishProgress();

  // 设置页面标题
  setPageTitle(to);

  // 记录页面访问日志
  logPageVisit(to, from);
}

/**
 * 检查路径是否在白名单中
 * 支持通配符匹配
 */
function isInWhiteList(path) {
  return whiteList.some(whiteListPath => {
    if (whiteListPath.endsWith('/*')) {
      // 通配符匹配
      const basePath = whiteListPath.slice(0, -2);
      return path.startsWith(basePath);
    } else {
      // 精确匹配
      return path === whiteListPath;
    }
  });
}

/**
 * 检查页面权限
 */
function checkPagePermission(path) {
  // 对于动态路由（如 /video-room/:roomId），需要特殊处理
  if (path.startsWith("/video-room/")) {
    // 可以在这里添加特定的房间访问权限检查
    // 例如：检查用户是否有权限访问特定房间
    return true; // 暂时允许所有已登录用户访问
  }

  const requiredPermissions = permissionRoutes[path];

  if (!requiredPermissions) {
    // 没有配置权限要求，默认允许访问
    return true;
  }

  // 检查是否有所需权限
  return requiredPermissions.every((permission) => hasPermission(permission));
}

/**
 * 设置页面标题
 */
function setPageTitle(to) {
  const title = to.meta?.title || "学习助手";
  document.title = title;
}

/**
 * 记录页面访问日志
 */
function logPageVisit(to, from) {
  if (import.meta.env.DEV) {
    console.log(`🔄 路由跳转: ${from.path} -> ${to.path}`);
  }

  // 这里可以添加埋点统计代码
  // analytics.track('page_view', {
  //   from: from.path,
  //   to: to.path,
  //   timestamp: Date.now()
  // })
}

/**
 * 开始进度条
 */
function startProgress() {
  // 这里可以集成进度条组件，如 nprogress
  // NProgress.start()
}

/**
 * 结束进度条
 */
function finishProgress() {
  // NProgress.done()
}

/**
 * 路由错误处理
 */
export function onError(error) {
  console.error("路由错误:", error);
  // ElMessage.error('页面加载失败，请刷新重试')
  alert("页面加载失败，请刷新重试");
  finishProgress();
}

/**
 * 动态添加权限路由
 */
export function addPermissionRoutes(router, routes) {
  routes.forEach((route) => {
    router.addRoute(route);
  });
}

/**
 * 重置路由
 */
export function resetRouter(router) {
  const newRouter = createRouter({
    history: createWebHistory(),
    routes: [],
  });
  router.matcher = newRouter.matcher;
}
