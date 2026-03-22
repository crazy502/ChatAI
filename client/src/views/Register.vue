<template>
  <div class="register-container">
    <!-- 背景动画元素 -->
    <div class="bg-animation">
      <div class="bg-grid"></div>
      <div class="bg-glow bg-glow-1"></div>
      <div class="bg-glow bg-glow-2"></div>
      <div class="bg-particles">
        <span v-for="n in 20" :key="n" class="particle"></span>
      </div>
    </div>

    <!-- 注册卡片 -->
    <div class="register-card sci-fi-animate-slide-up">
      <!-- 角标装饰 -->
      <div class="corner top-left"></div>
      <div class="corner top-right"></div>
      <div class="corner bottom-left"></div>
      <div class="corner bottom-right"></div>
      
      <!-- 扫描线效果 -->
      <div class="scan-line"></div>

      <div class="card-header">
        <div class="logo">
          <span class="logo-icon">◈</span>
          <span class="logo-text">Agent Go</span>
        </div>
        <h1 class="title">新用户注册</h1>
        <p class="subtitle">NEW USER REGISTRATION // ACCESS REQUEST</p>
      </div>

      <form class="register-form" @submit.prevent="handleRegister">
        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">◉</span>
            邮箱地址
          </label>
          <input
            v-model="registerForm.email"
            type="email"
            class="sci-fi-input"
            placeholder="请输入邮箱"
            required
          />
          <div class="input-glow"></div>
        </div>

        <div class="form-group captcha-group">
          <label class="form-label">
            <span class="label-icon">◈</span>
            验证码
          </label>
          <div class="captcha-wrapper">
            <input
              v-model="registerForm.captcha"
              type="text"
              class="sci-fi-input captcha-input"
              placeholder="请输入验证码"
              required
            />
            <button
              type="button"
              class="captcha-btn"
              :disabled="countdown > 0 || codeLoading"
              @click="sendCode"
            >
              <span class="btn-text">
                {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
              </span>
              <span class="btn-glow"></span>
            </button>
          </div>
          <div class="input-glow"></div>
        </div>

        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">◉</span>
            密码
          </label>
          <input
            v-model="registerForm.password"
            type="password"
            class="sci-fi-input"
            placeholder="请输入密码（至少6位）"
            required
            minlength="6"
          />
          <div class="input-glow"></div>
        </div>

        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">◈</span>
            确认密码
          </label>
          <input
            v-model="registerForm.confirmPassword"
            type="password"
            class="sci-fi-input"
            placeholder="请再次输入密码"
            required
          />
          <div class="input-glow"></div>
          <span v-if="passwordError" class="error-text">{{ passwordError }}</span>
        </div>

        <button
          type="submit"
          class="register-btn"
          :class="{ 'loading': loading }"
          :disabled="loading"
        >
          <span class="btn-text">{{ loading ? '处理中...' : '完 成 注 册' }}</span>
          <span class="btn-glow"></span>
          <div class="btn-particles" v-if="loading">
            <span v-for="n in 3" :key="n" class="particle"></span>
          </div>
        </button>

        <div class="form-footer">
          <span class="footer-text">已有访问权限？</span>
          <a href="#" class="sci-fi-link" @click.prevent="goToLogin">
            返回登录 →
          </a>
        </div>
      </form>

      <!-- 装饰性数据流 -->
      <div class="data-stream">
        <span v-for="n in 8" :key="n" class="data-bit">{{ Math.random() > 0.5 ? '1' : '0' }}</span>
      </div>
    </div>

    <!-- 底部信息 -->
    <div class="footer-info">
      <span class="version">v2.0.1</span>
      <span class="divider">|</span>
      <span class="status">
        <span class="status-dot"></span>
        系统在线
      </span>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'

export default {
  name: 'RegisterView',
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const codeLoading = ref(false)
    const countdown = ref(0)
    
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

    const sendCode = async () => {
      if (!registerForm.email) {
        showMessage('请输入邮箱地址', 'error')
        return
      }
      
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(registerForm.email)) {
        showMessage('请输入正确的邮箱格式', 'error')
        return
      }

      try {
        codeLoading.value = true
        const response = await api.post('/user/captcha', { email: registerForm.email })
        if (response.data.status_code === 1000) {
          showMessage('验证码已发送', 'success')
          countdown.value = 60
          const timer = setInterval(() => {
            countdown.value--
            if (countdown.value <= 0) {
              clearInterval(timer)
            }
          }, 1000)
        } else {
          showMessage(response.data.status_msg || '发送失败', 'error')
        }
      } catch (error) {
        console.error('Send code error:', error)
        showMessage('连接异常，请重试', 'error')
      } finally {
        codeLoading.value = false
      }
    }

    const handleRegister = async () => {
      if (!registerForm.email || !registerForm.captcha || !registerForm.password || !registerForm.confirmPassword) {
        showMessage('请填写所有字段', 'error')
        return
      }

      if (registerForm.password !== registerForm.confirmPassword) {
        showMessage('两次输入的密码不一致', 'error')
        return
      }

      if (registerForm.password.length < 6) {
        showMessage('密码长度不能少于6位', 'error')
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
          showMessage('注册成功，请登录', 'success')
          setTimeout(() => {
            router.push('/login')
          }, 1500)
        } else {
          showMessage(response.data.status_msg || '注册失败', 'error')
        }
      } catch (error) {
        console.error('Register error:', error)
        showMessage('连接异常，请重试', 'error')
      } finally {
        loading.value = false
      }
    }

    const goToLogin = () => {
      router.push('/login')
    }

    const showMessage = (message, type) => {
      const msgEl = document.createElement('div')
      msgEl.className = `sci-fi-message ${type}`
      msgEl.innerHTML = `
        <span class="msg-icon">${type === 'success' ? '✓' : '✗'}</span>
        <span class="msg-text">${message}</span>
      `
      document.body.appendChild(msgEl)
      
      setTimeout(() => msgEl.classList.add('show'), 10)
      setTimeout(() => {
        msgEl.classList.remove('show')
        setTimeout(() => document.body.removeChild(msgEl), 300)
      }, 3000)
    }

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
.register-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;
  background: var(--sci-fi-bg-dark);
}

/* 背景动画 */
.bg-animation {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
}

.bg-grid {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: 
    linear-gradient(rgba(16, 185, 129, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(16, 185, 129, 0.05) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: gridMove 20s linear infinite;
}

@keyframes gridMove {
  0% { transform: perspective(500px) rotateX(60deg) translateY(0); }
  100% { transform: perspective(500px) rotateX(60deg) translateY(50px); }
}

.bg-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.15;
  animation: glowPulse 4s ease-in-out infinite;
}

.bg-glow-1 {
  width: 400px;
  height: 400px;
  background: var(--sci-fi-secondary);
  top: -100px;
  right: -100px;
}

.bg-glow-2 {
  width: 300px;
  height: 300px;
  background: var(--sci-fi-primary);
  bottom: -50px;
  left: -50px;
  animation-delay: 2s;
}

@keyframes glowPulse {
  0%, 100% { opacity: 0.1; transform: scale(1); }
  50% { opacity: 0.2; transform: scale(1.1); }
}

.bg-particles {
  position: absolute;
  width: 100%;
  height: 100%;
}

.particle {
  position: absolute;
  width: 2px;
  height: 2px;
  background: var(--sci-fi-primary);
  border-radius: 50%;
  animation: particleFloat 15s infinite;
  opacity: 0.6;
}

.particle:nth-child(1) { left: 10%; top: 20%; animation-delay: 0s; }
.particle:nth-child(2) { left: 20%; top: 80%; animation-delay: 1s; }
.particle:nth-child(3) { left: 30%; top: 40%; animation-delay: 2s; }
.particle:nth-child(4) { left: 40%; top: 60%; animation-delay: 3s; }
.particle:nth-child(5) { left: 50%; top: 30%; animation-delay: 4s; }
.particle:nth-child(6) { left: 60%; top: 70%; animation-delay: 5s; }
.particle:nth-child(7) { left: 70%; top: 50%; animation-delay: 6s; }
.particle:nth-child(8) { left: 80%; top: 20%; animation-delay: 7s; }
.particle:nth-child(9) { left: 90%; top: 80%; animation-delay: 8s; }
.particle:nth-child(10) { left: 15%; top: 60%; animation-delay: 9s; }
.particle:nth-child(11) { left: 25%; top: 30%; animation-delay: 10s; }
.particle:nth-child(12) { left: 35%; top: 70%; animation-delay: 11s; }
.particle:nth-child(13) { left: 45%; top: 40%; animation-delay: 12s; }
.particle:nth-child(14) { left: 55%; top: 80%; animation-delay: 13s; }
.particle:nth-child(15) { left: 65%; top: 20%; animation-delay: 14s; }
.particle:nth-child(16) { left: 75%; top: 60%; animation-delay: 0.5s; }
.particle:nth-child(17) { left: 85%; top: 40%; animation-delay: 1.5s; }
.particle:nth-child(18) { left: 5%; top: 50%; animation-delay: 2.5s; }
.particle:nth-child(19) { left: 95%; top: 30%; animation-delay: 3.5s; }
.particle:nth-child(20) { left: 50%; top: 90%; animation-delay: 4.5s; }

@keyframes particleFloat {
  0%, 100% { 
    transform: translateY(0) translateX(0);
    opacity: 0;
  }
  10% { opacity: 0.6; }
  90% { opacity: 0.6; }
  50% { 
    transform: translateY(-100px) translateX(50px);
  }
}

/* 注册卡片 */
.register-card {
  position: relative;
  width: 460px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 16px;
  backdrop-filter: blur(20px);
  box-shadow: 0 4px 20px rgba(16, 185, 129, 0.15);
  z-index: 1;
  overflow: hidden;
}

.register-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--sci-fi-primary), transparent);
  animation: borderGlow 3s ease-in-out infinite;
}

@keyframes borderGlow {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}

/* 角标 */
.corner {
  position: absolute;
  width: 20px;
  height: 20px;
  border: 2px solid var(--sci-fi-primary);
  transition: all 0.3s ease;
}

.corner.top-left {
  top: -1px;
  left: -1px;
  border-right: none;
  border-bottom: none;
}

.corner.top-right {
  top: -1px;
  right: -1px;
  border-left: none;
  border-bottom: none;
}

.corner.bottom-left {
  bottom: -1px;
  left: -1px;
  border-right: none;
  border-top: none;
}

.corner.bottom-right {
  bottom: -1px;
  right: -1px;
  border-left: none;
  border-top: none;
}

.register-card:hover .corner {
  width: 30px;
  height: 30px;
  box-shadow: 0 0 10px var(--sci-fi-primary);
}

/* 扫描线 */
.scan-line {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--sci-fi-primary), transparent);
  animation: scanMove 3s linear infinite;
  opacity: 0.5;
}

@keyframes scanMove {
  0% { transform: translateY(0); }
  100% { transform: translateY(500px); }
}

/* 头部 */
.card-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 20px;
}

.logo-icon {
  font-size: 28px;
  color: var(--sci-fi-primary);
  animation: iconPulse 2s ease-in-out infinite;
  text-shadow: 0 0 20px var(--sci-fi-primary);
}

@keyframes iconPulse {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.1); opacity: 0.8; }
}

.logo-text {
  font-family: 'Orbitron', sans-serif;
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 4px;
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-success));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.title {
  font-family: 'Orbitron', sans-serif;
  font-size: 24px;
  font-weight: 700;
  letter-spacing: 6px;
  color: var(--sci-fi-text-primary);
  margin-bottom: 8px;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.subtitle {
  font-family: 'Rajdhani', sans-serif;
  font-size: 11px;
  letter-spacing: 3px;
  color: var(--sci-fi-text-muted);
  text-transform: uppercase;
}

/* 表单 */
.register-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  position: relative;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'Orbitron', sans-serif;
  font-size: 10px;
  font-weight: 500;
  letter-spacing: 2px;
  text-transform: uppercase;
  color: var(--sci-fi-text-secondary);
  margin-bottom: 8px;
}

.label-icon {
  color: var(--sci-fi-primary);
  font-size: 9px;
}

.sci-fi-input {
  width: 100%;
  padding: 12px 16px;
  font-family: 'Rajdhani', sans-serif;
  font-size: 15px;
  color: var(--sci-fi-text-primary);
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-radius: 8px;
  outline: none;
  transition: all 0.3s ease;
}

.sci-fi-input:focus {
  border-color: var(--sci-fi-primary);
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1), 0 0 20px rgba(16, 185, 129, 0.2);
  background: rgba(255, 255, 255, 0.95);
}

.sci-fi-input::placeholder {
  color: var(--sci-fi-text-muted);
  font-size: 13px;
}

.input-glow {
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--sci-fi-primary), transparent);
  transition: all 0.3s ease;
  transform: translateX(-50%);
}

.sci-fi-input:focus ~ .input-glow {
  width: 100%;
}

/* 验证码组 */
.captcha-group {
  position: relative;
}

.captcha-wrapper {
  display: flex;
  gap: 12px;
}

.captcha-input {
  flex: 1;
}

.captcha-btn {
  position: relative;
  padding: 12px 20px;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 1px;
  color: var(--sci-fi-text-primary);
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.captcha-btn:hover:not(:disabled) {
  background: rgba(16, 185, 129, 0.25);
  border-color: var(--sci-fi-primary);
  box-shadow: 0 0 20px rgba(16, 185, 129, 0.2);
}

.captcha-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-glow {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.6s ease;
}

.captcha-btn:hover:not(:disabled) .btn-glow {
  left: 100%;
}

/* 错误提示 */
.error-text {
  display: block;
  margin-top: 6px;
  font-family: 'Rajdhani', sans-serif;
  font-size: 12px;
  color: var(--sci-fi-danger);
  animation: errorPulse 2s ease-in-out infinite;
}

@keyframes errorPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

/* 注册按钮 */
.register-btn {
  position: relative;
  padding: 14px 32px;
  font-family: 'Orbitron', sans-serif;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 3px;
  text-transform: uppercase;
  color: var(--sci-fi-bg-dark);
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-success));
  border: none;
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s ease;
  margin-top: 8px;
}

.register-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(16, 185, 129, 0.4);
}

.register-btn:active:not(:disabled) {
  transform: translateY(0);
}

.register-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-text {
  position: relative;
  z-index: 1;
}

.register-btn .btn-glow {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.6s ease;
}

.register-btn:hover .btn-glow {
  left: 100%;
}

.register-btn.loading .btn-text {
  animation: textPulse 1s ease-in-out infinite;
}

@keyframes textPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

.btn-particles {
  position: absolute;
  top: 50%;
  right: 20px;
  transform: translateY(-50%);
  display: flex;
  gap: 4px;
}

.btn-particles .particle {
  width: 6px;
  height: 6px;
  background: var(--sci-fi-bg-dark);
  border-radius: 50%;
  animation: btnParticle 1s ease-in-out infinite;
}

.btn-particles .particle:nth-child(2) {
  animation-delay: 0.2s;
}

.btn-particles .particle:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes btnParticle {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(0.5); opacity: 0.5; }
}

/* 表单底部 */
.form-footer {
  text-align: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.footer-text {
  color: var(--sci-fi-text-muted);
  font-size: 13px;
  margin-right: 8px;
}

.sci-fi-link {
  color: var(--sci-fi-primary);
  text-decoration: none;
  font-weight: 500;
  position: relative;
  transition: all 0.3s ease;
}

.sci-fi-link::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 1px;
  background: var(--sci-fi-primary);
  transition: width 0.3s ease;
}

.sci-fi-link:hover {
  color: var(--sci-fi-success);
}

.sci-fi-link:hover::after {
  width: 100%;
}

/* 数据流装饰 */
.data-stream {
  position: absolute;
  bottom: 16px;
  right: 20px;
  display: flex;
  gap: 4px;
  font-family: 'Orbitron', monospace;
  font-size: 10px;
  color: var(--sci-fi-text-muted);
  opacity: 0.5;
}

.data-bit {
  animation: dataFlicker 2s ease-in-out infinite;
}

.data-bit:nth-child(2) { animation-delay: 0.2s; }
.data-bit:nth-child(3) { animation-delay: 0.4s; }
.data-bit:nth-child(4) { animation-delay: 0.6s; }
.data-bit:nth-child(5) { animation-delay: 0.8s; }
.data-bit:nth-child(6) { animation-delay: 1s; }
.data-bit:nth-child(7) { animation-delay: 1.2s; }
.data-bit:nth-child(8) { animation-delay: 1.4s; }

@keyframes dataFlicker {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 1; }
}

/* 底部信息 */
.footer-info {
  position: fixed;
  bottom: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  font-family: 'Rajdhani', sans-serif;
  font-size: 12px;
  color: var(--sci-fi-text-muted);
  z-index: 1;
}

.version {
  letter-spacing: 2px;
}

.divider {
  opacity: 0.3;
}

.status {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 6px;
  height: 6px;
  background: var(--sci-fi-success);
  border-radius: 50%;
  animation: statusPulse 2s ease-in-out infinite;
}

@keyframes statusPulse {
  0%, 100% { opacity: 1; box-shadow: 0 0 5px var(--sci-fi-success); }
  50% { opacity: 0.6; box-shadow: 0 0 10px var(--sci-fi-success); }
}

/* 消息提示样式 */
:global(.sci-fi-message) {
  position: fixed;
  top: 24px;
  right: 24px;
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-family: 'Rajdhani', sans-serif;
  font-size: 14px;
  z-index: 9999;
  transform: translateX(400px);
  opacity: 0;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

:global(.sci-fi-message.show) {
  transform: translateX(0);
  opacity: 1;
}

:global(.sci-fi-message.success) {
  border-color: var(--sci-fi-success);
  box-shadow: 0 0 20px rgba(0, 245, 212, 0.2);
}

:global(.sci-fi-message.error) {
  border-color: var(--sci-fi-danger);
  box-shadow: 0 0 20px rgba(241, 91, 181, 0.2);
}

:global(.sci-fi-message .msg-icon) {
  font-size: 18px;
  font-weight: bold;
}

:global(.sci-fi-message.success .msg-icon) {
  color: var(--sci-fi-success);
}

:global(.sci-fi-message.error .msg-icon) {
  color: var(--sci-fi-danger);
}

:global(.sci-fi-message .msg-text) {
  color: var(--sci-fi-text-primary);
}

/* 响应式 */
@media (max-width: 480px) {
  .register-card {
    width: 90%;
    padding: 28px 20px;
  }
  
  .title {
    font-size: 20px;
    letter-spacing: 3px;
  }
  
  .logo-text {
    font-size: 16px;
  }
  
  .captcha-wrapper {
    flex-direction: column;
  }
  
  .captcha-btn {
    width: 100%;
  }
}
</style>
