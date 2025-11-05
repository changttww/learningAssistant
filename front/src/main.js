import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import "./style.css";
import App from "./App.vue";

// 导入页面组件
import Home from "./views/Home.vue";
import PersonalTasks from "./views/PersonalTasks.vue";
import TeamTasks from "./views/TeamTasks.vue";
import StudyRoom from "./views/StudyRoom.vue";
import VideoRoom from "./views/VideoRoom.vue";
import Profile from "./views/Profile.vue";
import TaskManager from "./views/TaskManager.vue";

// 导入路由守卫
import { beforeEach, afterEach, onError } from "./router/guards.js";

// 配置路由
const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: { title: "首页", requiresAuth: false },
  },
  {
    path: "/personal-tasks",
    name: "PersonalTasks",
    component: PersonalTasks,
    meta: { title: "个人任务", requiresAuth: true },
  },
  {
    path: "/team-tasks",
    name: "TeamTasks",
    component: TeamTasks,
    meta: { title: "团队任务", requiresAuth: true, permissions: ["team:view"] },
  },
  {
    path: "/study-room",
    name: "StudyRoom",
    component: StudyRoom,
    meta: { title: "在线自习室", requiresAuth: true },
  },
  {
    path: "/video-room/:roomId",
    name: "VideoRoom",
    component: VideoRoom,
    meta: { title: "视频会议室", requiresAuth: true },
  },
  {
    path: "/profile",
    name: "Profile",
    component: Profile,
    meta: {
      title: "个人资料",
      requiresAuth: true,
      permissions: ["profile:view"],
    },
  },
  {
    path: "/task-manager",
    name: "TaskManager",
    component: TaskManager,
    meta: {
      title: "任务管理",
      requiresAuth: true,
      permissions: ["task:manage"],
    },
  },
  // 错误页面路由
  {
    path: "/403",
    name: "Forbidden",
    component: () => import("./views/error/403.vue"),
    meta: { title: "403 - 禁止访问" },
  },
  {
    path: "/404",
    name: "NotFound",
    component: () => import("./views/error/404.vue"),
    meta: { title: "404 - 页面不存在" },
  },
  // 登录相关路由
  {
    path: "/login",
    name: "Login",
    component: () => import("./views/auth/Login.vue"),
    meta: { title: "登录" },
  },
  {
    path: "/register",
    name: "Register",
    component: () => import("./views/auth/Register.vue"),
    meta: { title: "注册" },
  },
  // 捕获所有未匹配的路由
  {
    path: "/:pathMatch(.*)*",
    redirect: "/404",
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 注册路由守卫
router.beforeEach(beforeEach);
router.afterEach(afterEach);
router.onError(onError);

const app = createApp(App);
app.use(router);
app.mount("#app");
