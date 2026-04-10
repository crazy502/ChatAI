import { createRouter, createWebHistory } from 'vue-router'
import { isAdminToken } from '../utils/auth'

const Login = () => import(/* webpackChunkName: "auth" */ '../views/Login.vue')
const AIChat = () => import(/* webpackChunkName: "chat" */ '../views/AIChat.vue')
const AdminMetrics = () => import(/* webpackChunkName: "admin" */ '../views/AdminMetrics.vue')

const routes = [
  {
    path: '/',
    redirect: () => (localStorage.getItem('token') ? '/ai-chat' : '/login')
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: 'AgentGo | 登录' }
  },
  {
    path: '/register',
    redirect: '/login'
  },
  {
    path: '/menu',
    redirect: '/ai-chat'
  },
  {
    path: '/ai-chat',
    name: 'AIChat',
    component: AIChat,
    meta: {
      requiresAuth: true,
      title: 'AgentGo | 智能对话'
    }
  },
  {
    path: '/admin-metrics',
    name: 'AdminMetrics',
    component: AdminMetrics,
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'AgentGo | 管理监控'
    }
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (token && (to.path === '/login' || to.path === '/register')) {
    next('/ai-chat')
    return
  }

  if (to.matched.some((record) => record.meta.requiresAuth) && !token) {
    next('/login')
    return
  }

  if (to.matched.some((record) => record.meta.requiresAdmin) && !isAdminToken(token)) {
    next('/ai-chat')
    return
  }

  next()
})

router.afterEach((to) => {
  document.title = to.meta.title || 'AgentGo'
})

export default router
