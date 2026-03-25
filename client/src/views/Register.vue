<template>
  <div class="auth-page">
    <div class="auth-bg">
      <div class="orb orb-a"></div>
      <div class="orb orb-b"></div>
      <div class="grid"></div>
    </div>

    <section class="auth-card register-card">
      <p class="panel-kicker">NEW ACCESS REQUEST</p>
      <h1 class="panel-title">创建你的 AgentGo 账号</h1>
      <p class="panel-subtitle">完成注册后即可进入智能对话与会话工作台。</p>

      <form class="auth-form" @submit.prevent="handleRegister">
        <label class="field">
          <span class="field-label">邮箱地址</span>
          <input v-model="registerForm.email" type="email" class="field-input" placeholder="请输入邮箱" required />
        </label>

        <div class="field">
          <span class="field-label">验证码</span>
          <div class="captcha-row">
            <input v-model="registerForm.captcha" type="text" class="field-input" placeholder="请输入验证码" required />
            <button type="button" class="captcha-btn" :disabled="countdown > 0 || codeLoading" @click="sendCode">
              {{ countdown > 0 ? `${countdown}s` : codeLoading ? '发送中...' : '获取验证码' }}
            </button>
          </div>
        </div>

        <label class="field">
          <span class="field-label">密码</span>
          <input
            v-model="registerForm.password"
            type="password"
            class="field-input"
            placeholder="请输入密码（至少 6 位）"
            minlength="6"
            required
          />
        </label>

        <label class="field">
          <span class="field-label">确认密码</span>
          <input
            v-model="registerForm.confirmPassword"
            type="password"
            class="field-input"
            placeholder="请再次输入密码"
            required
          />
          <span v-if="passwordError" class="error-text">{{ passwordError }}</span>
        </label>

        <button type="submit" class="primary-btn" :disabled="loading">
          {{ loading ? '处理中...' : '完成注册' }}
        </button>
      </form>

      <div class="panel-footer">
        <span>已经有账号？</span>
        <button type="button" class="link-btn" @click="goToLogin">返回登录</button>
      </div>
    </section>
  </div>
</template>

<script>
import { ref, reactive, computed, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'
import { useUi } from '../composables/useUi'

export default {
  name: 'RegisterView',
  setup() {
    const router = useRouter()
    const { showToast } = useUi()
    const loading = ref(false)
    const codeLoading = ref(false)
    const countdown = ref(0)
    let countdownTimer = null

    const registerForm = reactive({
      email: '',
      captcha: '',
      password: '',
      confirmPassword: ''
    })

    const passwordError = computed(() => {
      if (registerForm.confirmPassword && registerForm.password !== registerForm.confirmPassword) {
        return '两次输入的密码不一致'
      }
      return ''
    })

    const startCountdown = () => {
      if (countdownTimer) {
        clearInterval(countdownTimer)
      }

      countdown.value = 60
      countdownTimer = setInterval(() => {
        countdown.value -= 1
        if (countdown.value <= 0) {
          clearInterval(countdownTimer)
          countdownTimer = null
        }
      }, 1000)
    }

    const sendCode = async () => {
      if (!registerForm.email) {
        showToast('请输入邮箱地址', 'error')
        return
      }

      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(registerForm.email)) {
        showToast('请输入正确的邮箱格式', 'error')
        return
      }

      try {
        codeLoading.value = true
        const response = await api.post('/user/captcha', { email: registerForm.email })
        if (response.data.status_code === 1000) {
          showToast('验证码已发送', 'success')
          startCountdown()
        } else {
          showToast(response.data.status_msg || '发送失败', 'error')
        }
      } catch (error) {
        console.error('Send code error:', error)
        showToast('连接异常，请重试', 'error')
      } finally {
        codeLoading.value = false
      }
    }

    const handleRegister = async () => {
      if (!registerForm.email || !registerForm.captcha || !registerForm.password || !registerForm.confirmPassword) {
        showToast('请填写所有字段', 'error')
        return
      }

      if (registerForm.password !== registerForm.confirmPassword) {
        showToast('两次输入的密码不一致', 'error')
        return
      }

      if (registerForm.password.length < 6) {
        showToast('密码长度不能少于 6 位', 'error')
        return
      }

      try {
        loading.value = true
        const response = await api.post('/user/register', {
          email: registerForm.email,
          captcha: registerForm.captcha,
          password: registerForm.password
        })

        if (response.data.status_code === 1000) {
          showToast('注册成功，请登录', 'success')
          setTimeout(() => {
            router.push('/login')
          }, 1200)
        } else {
          showToast(response.data.status_msg || '注册失败', 'error')
        }
      } catch (error) {
        console.error('Register error:', error)
        showToast('连接异常，请重试', 'error')
      } finally {
        loading.value = false
      }
    }

    const goToLogin = () => {
      router.push('/login')
    }

    onBeforeUnmount(() => {
      if (countdownTimer) {
        clearInterval(countdownTimer)
      }
    })

    return {
      loading,
      codeLoading,
      countdown,
      registerForm,
      passwordError,
      sendCode,
      handleRegister,
      goToLogin
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
  background:
    radial-gradient(circle at 18% 18%, rgba(var(--mint-rgb), 0.16), transparent 30%),
    radial-gradient(circle at 78% 22%, rgba(var(--mint-2-rgb), 0.14), transparent 24%),
    radial-gradient(circle at 52% 100%, rgba(var(--info-rgb), 0.08), transparent 36%),
    linear-gradient(160deg, var(--bg-1) 0%, #eff5f2 42%, #ecf3ef 100%);
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  opacity: 0.32;
}

.orb-a {
  width: 380px;
  height: 380px;
  top: -120px;
  right: -120px;
  background: rgba(var(--mint-rgb), 0.14);
}

.orb-b {
  width: 340px;
  height: 340px;
  left: -120px;
  bottom: -120px;
  background: rgba(var(--mint-2-rgb), 0.12);
}

.grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(15, 23, 42, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(15, 23, 42, 0.03) 1px, transparent 1px);
  background-size: 48px 48px;
}

.auth-card {
  position: relative;
  z-index: 1;
  width: min(100%, 520px);
  padding: 34px;
  border-radius: 24px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(248, 251, 249, 0.96)),
    var(--panel);
  border: 1px solid var(--border);
  box-shadow:
    0 24px 64px rgba(15, 23, 42, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(26px);
}

.register-card {
  max-width: 520px;
}

.panel-kicker {
  margin: 0 0 12px;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--kicker);
}

.panel-title {
  margin: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: 28px;
  line-height: 1.2;
  color: var(--text-cold-1);
}

.panel-subtitle {
  margin: 14px 0 0;
  font-size: 15px;
  line-height: 1.7;
  color: var(--text-cold-2);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 24px;
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
  color: var(--kicker);
}

.field-input {
  width: 100%;
  padding: 14px 16px;
  border-radius: 14px;
  border: 1px solid rgba(var(--kicker-rgb), 0.14);
  background: var(--panel-soft);
  color: var(--text-cold-1);
  outline: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.field-input:focus {
  border-color: rgba(var(--kicker-rgb), 0.46);
  box-shadow: 0 0 0 4px var(--focus-ring), 0 0 22px rgba(var(--mint-rgb), 0.16);
}

.captcha-row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
}

.captcha-btn {
  padding: 0 18px;
  border: none;
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(var(--mint-rgb), 0.18), rgba(var(--mint-3-rgb), 0.12));
  color: #b9ffea;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  letter-spacing: 1px;
  cursor: pointer;
}

.captcha-btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.error-text {
  font-size: 12px;
  color: var(--sci-fi-danger);
}

.primary-btn {
  margin-top: 6px;
  padding: 14px 18px;
  border: none;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--mint-1), var(--mint-2) 52%, var(--mint-3));
  color: var(--btn-text-dark);
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  letter-spacing: 2px;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.primary-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 18px 30px var(--mint-shadow);
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
  color: var(--text-cold-2);
}

.link-btn {
  border: none;
  background: transparent;
  color: var(--kicker);
  font-weight: 600;
  cursor: pointer;
}

@media (max-width: 560px) {
  .auth-card {
    padding: 28px 20px;
  }

  .captcha-row {
    grid-template-columns: 1fr;
  }
}
</style>




