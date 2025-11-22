import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import './style.css'
import App from './App.vue'

// 导入页面组件
import Home from './views/Home.vue'
import PersonalTasks from './views/PersonalTasks.vue'
import TeamTasks from './views/TeamTasks.vue'
import StudyRoom from './views/StudyRoom.vue'
import Profile from './views/Profile.vue'
import TaskManager from './views/TaskManager.vue'

// 配置路由
const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/personal-tasks', name: 'PersonalTasks', component: PersonalTasks },
  { path: '/team-tasks', name: 'TeamTasks', component: TeamTasks },
  { path: '/study-room', name: 'StudyRoom', component: StudyRoom },
  { path: '/profile', name: 'Profile', component: Profile },
  { path: '/task-manager', name: 'TaskManager', component: TaskManager }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const app = createApp(App)
app.use(router)
app.mount('#app')
