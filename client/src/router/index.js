import { createRouter, createWebHistory } from 'vue-router'
import { clearAuth, isAdminToken, isTokenUsable } from '../utils/auth'
import { translate } from '../i18n'

const Login = () => import(/* webpackChunkName: "auth" */ '../views/Login.vue')
const AIChat = () => import(/* webpackChunkName: "chat" */ '../views/AIChat.vue')
const AdminMetrics = () => import(/* webpackChunkName: "admin" */ '../views/AdminMetrics.vue')
const NotFound = () => import(/* webpackChunkName: "not-found" */ '../views/NotFound.vue')

const routes = [
  {
    path: '/',
    redirect: () => (isTokenUsable(localStorage.getItem('token')) ? '/ai-chat' : '/login')
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
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound
  }
]

const routeTitleKeys = {
  Login: 'route.login',
  AIChat: 'route.chat',
  AdminMetrics: 'route.admin',
  NotFound: 'route.notFound'
}

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
	const authenticated = isTokenUsable(token)
	if (token && !authenticated) {
	  clearAuth()
	}

	if (authenticated && (to.path === '/login' || to.path === '/register')) {
    next('/ai-chat')
    return
  }

	if (to.matched.some((record) => record.meta.requiresAuth) && !authenticated) {
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
	const titleKey = routeTitleKeys[to.name]
	document.title = titleKey ? translate(titleKey) : translate('common.brand')
})

export default router
