import { createRouter, createWebHistory } from 'vue-router'

const Login = () => import(/* webpackChunkName: "auth" */ '../views/Login.vue')
const Register = () => import(/* webpackChunkName: "auth" */ '../views/Register.vue')
const Menu = () => import(/* webpackChunkName: "menu" */ '../views/Menu.vue')
const AIChat = () => import(/* webpackChunkName: "chat" */ '../views/AIChat.vue')

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: 'GopherAI | 登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { title: 'GopherAI | 注册' }
  },
  {
    path: '/menu',
    name: 'Menu',
    component: Menu,
    meta: {
      requiresAuth: true,
      title: 'GopherAI | 控制台'
    }
  },
  {
    path: '/ai-chat',
    name: 'AIChat',
    component: AIChat,
    meta: {
      requiresAuth: true,
      title: 'GopherAI | 智能对话'
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
  if (to.matched.some((record) => record.meta.requiresAuth) && !token) {
    next('/login')
  } else {
    next()
  }
})

router.afterEach((to) => {
  document.title = to.meta.title || 'GopherAI'
})

export default router

