<template>
  <div class="auth-page">
    <div class="auth-bg">
      <div class="orb orb-a"></div>
      <div class="orb orb-b"></div>
      <div class="orb orb-c"></div>
      <div class="grid"></div>
      <div class="scan-line"></div>
    </div>

    <section class="auth-shell">
      <div class="preview-panel">
        <div class="preview-copy">
          <p class="shell-kicker">INTELLIGENT DIALOGUE WORKSPACE</p>
          <h1 class="shell-title">让想法、提问与回答在一个自然对话界面里连续发生</h1>
          <p class="shell-subtitle">
            这是一个偏轻量、偏陪伴感的 AI 聊天预览。消息会自动出现，回答会像真实生成一样逐字展开。
          </p>
        </div>

        <AIChatPreview />
      </div>

      <div class="auth-stage" :class="{ flipped: isRegisterMode }">
        <div class="auth-card">
          <article class="card-face card-front">
            <div class="face-header">
              <p class="panel-kicker">SECURE ACCESS</p>
              <h2 class="panel-title">欢迎回到 AgentGo</h2>
              <p class="panel-subtitle">登录后继续你的智能协作会话。</p>
            </div>

            <form class="auth-form" @submit.prevent="handleLogin">
              <label class="field">
                <span class="field-label">用户名</span>
                <input
                  v-model="loginForm.username"
                  type="text"
                  class="field-input"
                  placeholder="请输入用户名"
                  autocomplete="username"
                  required
                />
              </label>

              <label class="field">
                <span class="field-label">密码</span>
                <input
                  v-model="loginForm.password"
                  type="password"
                  class="field-input"
                  placeholder="请输入密码"
                  autocomplete="current-password"
                  required
                />
              </label>

              <button type="submit" class="primary-btn" :disabled="loginLoading">
                {{ loginLoading ? '验证中...' : '登录' }}
              </button>
            </form>

            <div class="panel-footer">
              <span>还没有账号？</span>
              <button type="button" class="link-btn" @click="switchMode(true)">立即注册</button>
            </div>
          </article>

          <article class="card-face card-back">
            <div class="face-header">
              <p class="panel-kicker">NEW ACCESS REQUEST</p>
              <h2 class="panel-title">创建你的 AgentGo 账号</h2>
              <p class="panel-subtitle">注册成功后自动翻回登录，并使用邮件中的账号登录。</p>
            </div>

            <form class="auth-form" @submit.prevent="handleRegister">
              <label class="field">
                <span class="field-label">邮箱地址</span>
                <input
                  v-model="registerForm.email"
                  type="email"
                  class="field-input"
                  placeholder="请输入邮箱"
                  autocomplete="email"
                  required
                />
              </label>

              <div class="field">
                <span class="field-label">验证码</span>
                <div class="captcha-row">
                  <input
                    v-model="registerForm.captcha"
                    type="text"
                    class="field-input"
                    placeholder="请输入验证码"
                    required
                  />
                  <button
                    type="button"
                    class="captcha-btn"
                    :disabled="countdown > 0 || codeLoading"
                    @click="sendCode"
                  >
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
                  placeholder="请输入至少 6 位密码"
                  autocomplete="new-password"
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
                  autocomplete="new-password"
                  required
                />
                <span v-if="passwordError" class="error-text">{{ passwordError }}</span>
              </label>

              <button type="submit" class="primary-btn" :disabled="registerLoading">
                {{ registerLoading ? '创建中...' : '完成注册' }}
              </button>
            </form>

            <div class="panel-footer">
              <span>已经有账号？</span>
              <button type="button" class="link-btn" @click="switchMode(false)">返回登录</button>
            </div>
          </article>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import AIChatPreview from '../components/auth/AIChatPreview.vue'
import api from '../utils/api'
import { useUi } from '../composables/useUi'

export default {
  name: 'LoginView',
  components: {
    AIChatPreview
  },
  setup() {
    const router = useRouter()
    const { showToast } = useUi()
    const isRegisterMode = ref(false)
    const loginLoading = ref(false)
    const registerLoading = ref(false)
    const codeLoading = ref(false)
    const countdown = ref(0)
    let countdownTimer = null

    const loginForm = reactive({
      username: '',
      password: ''
    })

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

    const switchMode = (nextMode) => {
      isRegisterMode.value = nextMode
    }

    const resetRegisterForm = () => {
      registerForm.email = ''
      registerForm.captcha = ''
      registerForm.password = ''
      registerForm.confirmPassword = ''
    }

    const clearCountdown = () => {
      if (countdownTimer) {
        clearInterval(countdownTimer)
        countdownTimer = null
      }
      countdown.value = 0
    }

    const startCountdown = () => {
      clearCountdown()
      countdown.value = 60
      countdownTimer = setInterval(() => {
        countdown.value -= 1
        if (countdown.value <= 0) {
          clearCountdown()
        }
      }, 1000)
    }

    const handleLogin = async () => {
      if (!loginForm.username || !loginForm.password) {
        showToast('请输入用户名和密码', 'error')
        return
      }

      try {
        loginLoading.value = true
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
        loginLoading.value = false
      }
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
        registerLoading.value = true
        const response = await api.post('/user/register', {
          email: registerForm.email,
          captcha: registerForm.captcha,
          password: registerForm.password
        })

        if (response.data.status_code === 1000) {
          showToast('注册成功，请查收邮箱中的账号后登录', 'success')
          resetRegisterForm()
          clearCountdown()
          switchMode(false)
        } else {
          showToast(response.data.status_msg || '注册失败', 'error')
        }
      } catch (error) {
        console.error('Register error:', error)
        showToast('连接异常，请重试', 'error')
      } finally {
        registerLoading.value = false
      }
    }

    watch(
      isRegisterMode,
      (value) => {
        document.title = value ? 'AgentGo | 注册' : 'AgentGo | 登录'
      },
      { immediate: true }
    )

    onBeforeUnmount(() => {
      clearCountdown()
    })

    return {
      isRegisterMode,
      loginLoading,
      registerLoading,
      codeLoading,
      countdown,
      loginForm,
      registerForm,
      passwordError,
      switchMode,
      handleLogin,
      sendCode,
      handleRegister
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
  padding: 28px;
  position: relative;
  overflow: hidden;
}

.auth-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 18% 18%, rgba(var(--mint-rgb), 0.16), transparent 28%),
    radial-gradient(circle at 78% 22%, rgba(var(--mint-2-rgb), 0.14), transparent 24%),
    radial-gradient(circle at 50% 100%, rgba(var(--info-rgb), 0.08), transparent 35%),
    linear-gradient(160deg, var(--bg-1) 0%, #eff5f2 42%, #ecf3ef 100%);
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  opacity: 0.42;
}

.orb-a {
  width: 320px;
  height: 320px;
  top: -110px;
  left: -90px;
  background: rgba(var(--info-rgb), 0.1);
}

.orb-b {
  width: 360px;
  height: 360px;
  right: -120px;
  bottom: -120px;
  background: rgba(var(--mint-2-rgb), 0.14);
}

.orb-c {
  width: 220px;
  height: 220px;
  top: 45%;
  right: 24%;
  background: rgba(var(--info-rgb), 0.1);
}

.grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(15, 23, 42, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(15, 23, 42, 0.03) 1px, transparent 1px);
  background-size: 48px 48px;
  mask-image: radial-gradient(circle at center, #000 48%, transparent 100%);
}

.scan-line {
  position: absolute;
  inset: 0;
  background: linear-gradient(to bottom, transparent 0%, rgba(var(--mint-rgb), 0.05) 48%, transparent 100%);
  animation: sweep 8s linear infinite;
}

.auth-shell {
  position: relative;
  z-index: 1;
  width: min(100%, 1100px);
  display: grid;
  grid-template-columns: minmax(320px, 1.1fr) minmax(320px, 480px);
  gap: 52px;
  align-items: center;
}

.preview-panel {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.preview-copy {
  color: var(--text-cold-2);
}

.shell-kicker {
  margin: 0 0 16px;
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  letter-spacing: 3px;
  color: var(--kicker);
}

.shell-title {
  margin: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: clamp(34px, 5vw, 54px);
  line-height: 1.08;
  color: var(--text-cold-1);
  text-shadow: 0 0 20px rgba(var(--mint-rgb), 0.14);
}

.shell-subtitle {
  margin: 18px 0 0;
  max-width: 38rem;
  font-size: 18px;
  line-height: 1.75;
  color: var(--text-cold-2);
}

.auth-stage {
  perspective: 1800px;
}

.auth-card {
  position: relative;
  width: min(100%, 480px);
  min-height: 640px;
  transform-style: preserve-3d;
  transition: transform 0.9s cubic-bezier(0.22, 1, 0.36, 1);
}

.auth-stage.flipped .auth-card {
  transform: rotateY(180deg);
}

.card-face {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  padding: 34px 32px 30px;
  border-radius: 28px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(248, 251, 249, 0.96)),
    var(--panel);
  border: 1px solid var(--border);
  box-shadow:
    0 24px 64px rgba(15, 23, 42, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(26px);
  backface-visibility: hidden;
  -webkit-backface-visibility: hidden;
  overflow: hidden;
}

.card-face::before {
  content: '';
  position: absolute;
  inset: 0;
  background:
    linear-gradient(130deg, rgba(var(--mint-2-rgb), 0.14), transparent 28%),
    linear-gradient(320deg, rgba(var(--mint-rgb), 0.12), transparent 36%);
  pointer-events: none;
}

.card-back {
  transform: rotateY(180deg);
}

.face-header,
.auth-form,
.panel-footer {
  position: relative;
  z-index: 1;
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
  font-size: 30px;
  line-height: 1.18;
  color: var(--text-cold-1);
}

.panel-subtitle {
  margin: 14px 0 0;
  font-size: 15px;
  line-height: 1.72;
  color: var(--text-cold-2);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
  margin-top: 30px;
  flex: 1;
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
  border-radius: 16px;
  border: 1px solid rgba(var(--kicker-rgb), 0.14);
  background: var(--panel-soft);
  color: var(--text-cold-1);
  outline: none;
  transition: border-color 0.22s ease, box-shadow 0.22s ease, transform 0.22s ease;
}

.field-input::placeholder {
  color: var(--text-cold-3);
}

.field-input:focus {
  border-color: rgba(var(--kicker-rgb), 0.46);
  box-shadow: 0 0 0 4px var(--focus-ring), 0 0 22px rgba(var(--mint-rgb), 0.16);
  transform: translateY(-1px);
}

.captcha-row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
}

.primary-btn,
.captcha-btn {
  border: none;
  cursor: pointer;
  transition: transform 0.22s ease, box-shadow 0.22s ease, opacity 0.22s ease;
}

.primary-btn {
  margin-top: auto;
  padding: 15px 18px;
  border-radius: 18px;
  background: linear-gradient(135deg, var(--mint-1), var(--mint-2) 52%, var(--mint-3));
  color: var(--btn-text-dark);
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  letter-spacing: 2px;
  font-weight: 700;
  box-shadow: 0 18px 34px var(--mint-shadow);
}

.primary-btn:hover:not(:disabled),
.captcha-btn:hover:not(:disabled) {
  transform: translateY(-2px);
}

.primary-btn:hover:not(:disabled) {
  box-shadow: 0 22px 40px rgba(var(--mint-rgb), 0.3);
}

.captcha-btn {
  min-width: 132px;
  padding: 0 18px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(var(--mint-rgb), 0.18), rgba(var(--mint-3-rgb), 0.12));
  color: #b9ffea;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  letter-spacing: 1px;
  box-shadow: inset 0 0 0 1px rgba(var(--kicker-rgb), 0.14);
}

.primary-btn:disabled,
.captcha-btn:disabled {
  cursor: not-allowed;
  opacity: 0.62;
}

.error-text {
  font-size: 12px;
  color: var(--danger);
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
  color: var(--accent-strong);
  font-weight: 700;
  cursor: pointer;
  padding: 0;
}

@keyframes sweep {
  0% {
    transform: translateY(-100%);
  }

  100% {
    transform: translateY(100%);
  }
}

@media (max-width: 980px) {
  .auth-shell {
    grid-template-columns: 1fr;
    gap: 28px;
    justify-items: center;
  }

  .preview-copy {
    text-align: center;
  }

  .shell-subtitle {
    max-width: 36rem;
  }
}

@media (max-width: 560px) {
  .auth-page {
    padding: 18px;
  }

  .auth-card {
    min-height: 690px;
  }

  .card-face {
    padding: 28px 22px 24px;
  }

  .panel-title {
    font-size: 24px;
  }

  .captcha-row {
    grid-template-columns: 1fr;
  }

  .captcha-btn {
    min-height: 48px;
  }
}
</style>


