<template>
  <div class="auth-page">
    <div class="auth-bg">
      <div class="orb orb-a"></div>
      <div class="orb orb-b"></div>
      <div class="grid"></div>
    </div>

    <section class="auth-card">
      <p class="panel-kicker">SECURE ACCESS</p>
      <h1 class="panel-title">欢迎回到 GopherAI</h1>
      <p class="panel-subtitle">登录后继续你的智能协作会话。</p>

      <form class="auth-form" @submit.prevent="handleLogin">
        <label class="field">
          <span class="field-label">用户名</span>
          <input v-model="loginForm.username" type="text" class="field-input" placeholder="请输入用户名" required />
        </label>

        <label class="field">
          <span class="field-label">密码</span>
          <input v-model="loginForm.password" type="password" class="field-input" placeholder="请输入密码" required />
        </label>

        <button type="submit" class="primary-btn" :disabled="loading">
          {{ loading ? '验证中...' : '登录' }}
        </button>
      </form>

      <div class="panel-footer">
        <span>还没有账号？</span>
        <button type="button" class="link-btn" @click="goToRegister">立即注册</button>
      </div>
    </section>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'
import { useUi } from '../composables/useUi'

export default {
  name: 'LoginView',
  setup() {
    const router = useRouter()
    const { showToast } = useUi()
    const loading = ref(false)
    const loginForm = reactive({
      username: '',
      password: ''
    })

    const handleLogin = async () => {
      if (!loginForm.username || !loginForm.password) {
        showToast('请输入用户名和密码', 'error')
        return
      }

      try {
        loading.value = true
        const response = await api.post('/user/login', {
          username: loginForm.username,
          password: loginForm.password
        })

        if (response.data.status_code === 1000) {
          localStorage.setItem('token', response.data.token)
          localStorage.setItem('isAdmin', response.data.isAdmin ? 'true' : 'false')
          showToast('身份验证通过', 'success')
          setTimeout(() => {
            router.push('/menu')
          }, 800)
        } else {
          showToast(response.data.status_msg || '验证失败', 'error')
        }
      } catch (error) {
        console.error('Login error:', error)
        showToast('连接异常，请重试', 'error')
      } finally {
        loading.value = false
      }
    }

    const goToRegister = () => {
      router.push('/register')
    }

    return {
      loading,
      loginForm,
      handleLogin,
      goToRegister
    }
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  position: relative;
  overflow: hidden;
}

.auth-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.7), rgba(236, 247, 240, 0.9));
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.3;
}

.orb-a {
  width: 360px;
  height: 360px;
  top: -120px;
  right: -80px;
  background: rgba(16, 185, 129, 0.35);
}

.orb-b {
  width: 320px;
  height: 320px;
  left: -100px;
  bottom: -120px;
  background: rgba(45, 212, 191, 0.3);
}

.grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(16, 185, 129, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(16, 185, 129, 0.05) 1px, transparent 1px);
  background-size: 48px 48px;
}

.auth-card {
  position: relative;
  z-index: 1;
  width: min(100%, 440px);
  padding: 36px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.88);
  border: 1px solid rgba(16, 185, 129, 0.16);
  box-shadow: 0 28px 60px rgba(15, 23, 42, 0.12);
  backdrop-filter: blur(18px);
}

.panel-kicker {
  margin: 0 0 12px;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--sci-fi-primary);
}

.panel-title {
  margin: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: 30px;
  line-height: 1.2;
  color: var(--sci-fi-text-primary);
}

.panel-subtitle {
  margin: 14px 0 0;
  font-size: 15px;
  line-height: 1.7;
  color: var(--sci-fi-text-secondary);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
  margin-top: 28px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field-label {
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  letter-spacing: 1px;
  color: var(--sci-fi-text-secondary);
}

.field-input {
  width: 100%;
  padding: 14px 16px;
  border-radius: 14px;
  border: 1px solid rgba(16, 185, 129, 0.18);
  background: rgba(255, 255, 255, 0.9);
  color: var(--sci-fi-text-primary);
  outline: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.field-input:focus {
  border-color: rgba(16, 185, 129, 0.4);
  box-shadow: 0 0 0 4px rgba(16, 185, 129, 0.08);
}

.primary-btn {
  margin-top: 6px;
  padding: 14px 18px;
  border: none;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-secondary));
  color: #fff;
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  letter-spacing: 2px;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.primary-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 18px 30px rgba(16, 185, 129, 0.2);
}

.primary-btn:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.panel-footer {
  margin-top: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  color: var(--sci-fi-text-secondary);
}

.link-btn {
  border: none;
  background: transparent;
  color: var(--sci-fi-primary);
  font-weight: 600;
  cursor: pointer;
}

@media (max-width: 480px) {
  .auth-card {
    padding: 28px 22px;
  }

  .panel-title {
    font-size: 24px;
  }
}
</style>

