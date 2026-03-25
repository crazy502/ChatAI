<template>
  <div class="menu-page">
    <div class="menu-bg">
      <div class="halo halo-a"></div>
      <div class="halo halo-b"></div>
      <div class="grid"></div>
    </div>

    <header class="menu-header">
      <div>
        <p class="header-kicker">CONTROL CENTER</p>
        <h1 class="header-title">AgentGo 控制台</h1>
      </div>
      <button class="logout-btn" type="button" @click="handleLogout">退出登录</button>
    </header>

    <main class="menu-main">
      <section class="hero-card">
        <p class="hero-kicker">INTELLIGENT WORKSPACE</p>
        <h2 class="hero-title">把问题、思路和答案集中在一个协作界面里</h2>
        <p class="hero-desc">
          你可以进入智能对话模块进行问答和分析。
          <template v-if="hasAdminAccess">
            当前账号也拥有管理员权限，可以进入系统监控面板查看请求量、错误率与模型延迟表现。
          </template>
          <template v-else>
            当前账号暂无管理员权限，因此不会展示管理监控入口。
          </template>
        </p>
        <div class="hero-actions">
          <button class="enter-btn" type="button" @click="enterChat">进入智能对话</button>
          <button v-if="hasAdminAccess" class="secondary-btn" type="button" @click="openMetrics">查看监控面板</button>
        </div>
      </section>

      <section class="stats-grid">
        <article class="stat-card">
          <span class="stat-label">系统状态</span>
          <strong class="stat-value online">在线</strong>
          <p class="stat-desc">前端控制台与会话入口可用</p>
        </article>

        <article class="stat-card">
          <span class="stat-label">运行时间</span>
          <strong class="stat-value">{{ uptime }}</strong>
          <p class="stat-desc">用于展示当前控制台的持续运行时长</p>
        </article>

        <article class="stat-card wide">
          <span class="stat-label">模型支持</span>
          <strong class="stat-value">Qwen / DeepSeek</strong>
          <p class="stat-desc">聊天页可直接切换模型与流式响应模式。</p>
        </article>

        <article class="stat-card">
          <span class="stat-label">管理监控</span>
          <strong class="stat-value">{{ hasAdminAccess ? '实时可视化' : '管理员可用' }}</strong>
          <p class="stat-desc">{{ hasAdminAccess ? '可查看请求量、错误率、接口延迟与模型调用表现。' : '当前账号没有管理权限，监控面板入口已隐藏。' }}</p>
        </article>
      </section>
    </main>
  </div>
</template>

<script>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUi } from '../composables/useUi'
import { isAdminToken } from '../utils/auth'

export default {
  name: 'MenuView',
  setup() {
    const router = useRouter()
    const { confirmAction } = useUi()
    const uptime = ref('00:00:00')
    const hasAdminAccess = computed(() => isAdminToken(localStorage.getItem('token')))
    let uptimeInterval = null

    const updateUptime = () => {
      const startTime = new Date('2024-01-01').getTime()
      const now = new Date().getTime()
      const diff = now - startTime

      const hours = Math.floor(diff / (1000 * 60 * 60))
      const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
      const seconds = Math.floor((diff % (1000 * 60)) / 1000)

      uptime.value = `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
    }

    onMounted(() => {
      updateUptime()
      uptimeInterval = setInterval(updateUptime, 1000)
    })

    onUnmounted(() => {
      if (uptimeInterval) {
        clearInterval(uptimeInterval)
      }
    })

    const enterChat = () => {
      router.push('/ai-chat')
    }

    const openMetrics = () => {
      router.push('/admin-metrics')
    }

    const handleLogout = async () => {
      const confirmed = await confirmAction({
        title: '断开当前连接？',
        message: '退出后需要重新登录才能继续使用 AgentGo。',
        confirmText: '退出登录',
        cancelText: '继续使用',
        intent: 'danger'
      })

      if (!confirmed) {
        return
      }

      localStorage.removeItem('token')
      localStorage.removeItem('isAdmin')
      router.push('/login')
    }

    return {
      uptime,
      hasAdminAccess,
      enterChat,
      openMetrics,
      handleLogout
    }
  }
}
</script>

<style scoped>
.menu-page {
  min-height: 100vh;
  position: relative;
  padding: 28px;
  overflow: hidden;
}

.menu-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(150deg, var(--bg-1), var(--bg-2));
}

.halo {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  opacity: 0.32;
}

.halo-a {
  width: 420px;
  height: 420px;
  top: -120px;
  right: -120px;
  background: rgba(var(--mint-rgb), 0.28);
}

.halo-b {
  width: 360px;
  height: 360px;
  left: -120px;
  bottom: -140px;
  background: rgba(var(--mint-2-rgb), 0.24);
}

.grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(var(--mint-rgb), 0.06) 1px, transparent 1px),
    linear-gradient(90deg, rgba(var(--mint-rgb), 0.06) 1px, transparent 1px);
  background-size: 48px 48px;
}

.menu-header,
.menu-main {
  position: relative;
  z-index: 1;
}

.menu-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin: 0 auto 28px;
  max-width: 1160px;
  padding: 20px 24px;
  border-radius: 24px;
  background: var(--panel);
  border: 1px solid var(--border);
  box-shadow: 0 18px 42px rgba(15, 23, 42, 0.1);
  backdrop-filter: blur(18px);
}

.header-kicker,
.hero-kicker,
.stat-label,
.enter-btn,
.logout-btn {
  font-family: 'Orbitron', sans-serif;
}

.header-kicker {
  margin: 0 0 8px;
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--sci-fi-primary);
}

.header-title {
  margin: 0;
  font-size: 28px;
  color: var(--sci-fi-text-primary);
}

.logout-btn {
  padding: 12px 18px;
  border: none;
  border-radius: 14px;
  background: rgba(240, 125, 125, 0.16);
  color: var(--sci-fi-danger);
  font-size: 12px;
  letter-spacing: 1px;
  cursor: pointer;
}

.menu-main {
  max-width: 1160px;
  margin: 0 auto;
  display: grid;
  gap: 22px;
}

.hero-card,
.stat-card {
  border-radius: 28px;
  background: var(--panel);
  border: 1px solid var(--border);
  box-shadow: 0 24px 54px rgba(15, 23, 42, 0.12);
  backdrop-filter: blur(18px);
}

.hero-card {
  padding: 36px;
}

.hero-kicker {
  margin: 0 0 12px;
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--sci-fi-primary);
}

.hero-title {
  margin: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: 34px;
  line-height: 1.25;
  color: var(--sci-fi-text-primary);
}

.hero-desc {
  max-width: 680px;
  margin: 18px 0 0;
  font-size: 16px;
  line-height: 1.8;
  color: var(--sci-fi-text-secondary);
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 28px;
}

.enter-btn {
  padding: 14px 20px;
  border: none;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-secondary));
  color: var(--btn-text);
  font-size: 12px;
  letter-spacing: 2px;
  cursor: pointer;
}

.secondary-btn {
  padding: 14px 20px;
  border: 1px solid var(--border);
  border-radius: 16px;
  background: var(--panel-soft);
  color: var(--sci-fi-text-primary);
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  letter-spacing: 2px;
  cursor: pointer;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 18px;
}

.stat-card {
  padding: 26px;
}

.stat-card.wide {
  grid-column: span 1;
}

.stat-label {
  display: block;
  margin-bottom: 12px;
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--sci-fi-text-muted);
}

.stat-value {
  display: block;
  font-size: 28px;
  color: var(--sci-fi-text-primary);
}

.stat-value.online {
  color: var(--sci-fi-primary);
}

.stat-desc {
  margin: 12px 0 0;
  font-size: 14px;
  line-height: 1.7;
  color: var(--sci-fi-text-secondary);
}

@media (max-width: 900px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .menu-page {
    padding: 16px;
  }

  .menu-header,
  .hero-card,
  .stat-card {
    padding: 22px 18px;
  }

  .menu-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .hero-title {
    font-size: 28px;
  }

  .logout-btn,
  .enter-btn,
  .secondary-btn {
    width: 100%;
  }
}
</style>




