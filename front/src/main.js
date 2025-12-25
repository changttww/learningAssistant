import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import "./style.css";
import App from "./App.vue";

// 导入页面组件
import Home from "./views/Home.vue";
import PersonalTasks from "./views/PersonalTasks.vue";
import TeamTasks from "./views/TeamTasks.vue";
import TeamConstellation from "./views/TeamConstellation.vue";
import TeamMeetingRoom from "./views/TeamMeetingRoom.vue";
import TeamCalendar from "./views/TeamCalendar.vue";
import TeamDocs from "./views/TeamDocs.vue";
import TeamReports from "./views/TeamReports.vue";
import StudyRoom from "./views/StudyRoom.vue";
import VideoRoom from "./views/VideoRoom.vue";
import Profile from "./views/Profile.vue";
import TaskManager from "./views/TaskManager.vue";
import NotificationHistory from "./views/NotificationHistory.vue";
import KnowledgeBase from "./views/KnowledgeBase.vue";
import KnowledgeGraph from "./views/KnowledgeGraph.vue";
import KnowledgeChat from "./views/KnowledgeChat.vue";
import AIReport from "./views/AIReport.vue";

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
    path: "/team-tasks/constellation",
    name: "TeamConstellation",
    component: TeamConstellation,
    meta: { title: "任务星图", requiresAuth: true, permissions: ["team:view"] },
  },
  {
    path: "/team-tasks/meeting/:teamId",
    name: "TeamMeetingRoom",
    component: TeamMeetingRoom,
    meta: { title: "快速会议室", requiresAuth: true, permissions: ["team:view"] },
  },
  {
    path: "/team-tasks/calendar/:teamId",
    name: "TeamCalendar",
    component: TeamCalendar,
    meta: { title: "团队日历", requiresAuth: true, permissions: ["team:view"] },
  },
  {
    path: "/team-tasks/docs/:teamId",
    name: "TeamDocs",
    component: TeamDocs,
    meta: { title: "协作文档", requiresAuth: true, permissions: ["team:view"] },
  },
  {
    path: "/team-tasks/reports/:teamId",
    name: "TeamReports",
    component: TeamReports,
    meta: { title: "数据报告", requiresAuth: true, permissions: ["team:view"] },
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
  {
    path: "/notifications",
    name: "NotificationHistory",
    component: NotificationHistory,
    meta: {
      title: "通知中心",
      requiresAuth: true,
    },
  },
  {
    path: "/knowledge-base",
    name: "KnowledgeBase",
    component: KnowledgeBase,
    meta: {
      title: "我的知识库",
      requiresAuth: true,
    },
  },
  {
    path: "/knowledge-graph",
    name: "KnowledgeGraph",
    component: KnowledgeGraph,
    meta: {
      title: "知识图谱",
      requiresAuth: true,
    },
  },
  {
    path: "/knowledge-chat",
    name: "KnowledgeChat",
    component: KnowledgeChat,
    meta: {
      title: "知识问答",
      requiresAuth: true,
    },
  },
  {
    path: "/ai-report",
    name: "AIReport",
    component: AIReport,
    meta: {
      title: "AI学习报告",
      requiresAuth: true,
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
app.use(ElementPlus);
app.mount("#app");
