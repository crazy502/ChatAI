<template>
  <div class="menu-container">
    <!-- 背景动画 -->
    <div class="bg-effects">
      <div class="grid-floor"></div>
      <div class="floating-particles">
        <span v-for="n in 30" :key="n" class="particle"></span>
      </div>
      <div class="glow-orb orb-1"></div>
      <div class="glow-orb orb-2"></div>
    </div>

    <!-- 头部导航 -->
    <header class="sci-fi-header">
      <div class="header-left">
        <div class="logo">
          <span class="logo-icon">◈</span>
          <span class="logo-text">Agent Go</span>
        </div>
        <div class="header-line"></div>
      </div>
      <div class="header-right">
        <div class="user-info">
          <span class="user-status">● 在线</span>
        </div>
        <button class="logout-btn" @click="handleLogout">
          <span class="btn-icon">◉</span>
          <span class="btn-text">断开连接</span>
        </button>
      </div>
    </header>

    <!-- 主内容区 -->
    <main class="main-content">
      <div class="welcome-section sci-fi-animate-fade-in">
        <h1 class="main-title">
          <span class="title-line">智能对话系统</span>
          <span class="title-sub">INTELLIGENT DIALOGUE SYSTEM</span>
        </h1>
        <p class="description">
          基于先进的大语言模型技术，为您提供智能、流畅的对话体验
        </p>
      </div>

      <!-- 功能卡片 -->
      <div class="feature-card sci-fi-animate-slide-up" @click="enterChat">
        <!-- 卡片装饰 -->
        <div class="card-frame">
          <div class="frame-corner tl"></div>
          <div class="frame-corner tr"></div>
          <div class="frame-corner bl"></div>
          <div class="frame-corner br"></div>
        </div>
        
        <!-- 扫描线 -->
        <div class="card-scan"></div>

        <!-- 卡片内容 -->
        <div class="card-content">
          <div class="card-icon-wrapper">
            <div class="icon-ring ring-1"></div>
            <div class="icon-ring ring-2"></div>
            <div class="icon-ring ring-3"></div>
            <span class="card-icon">◈</span>
          </div>
          
          <div class="card-text">
            <h2 class="card-title">AI 智能对话</h2>
            <p class="card-subtitle">AI CHAT MODULE</p>
            <p class="card-desc">
              与通义千问、DeepSeek等先进AI模型进行实时对话，
              获得智能化的回答与建议
            </p>
          </div>

          <div class="card-action">
            <span class="action-text">进入系统</span>
            <span class="action-arrow">→</span>
          </div>
        </div>

        <!-- 数据装饰 -->
        <div class="data-decoration">
          <div class="data-row">
            <span class="data-label">STATUS</span>
            <span class="data-value active">ONLINE</span>
          </div>
          <div class="data-row">
            <span class="data-label">MODELS</span>
            <span class="data-value">QWEN / DEEPSEEK</span>
          </div>
          <div class="data-row">
            <span class="data-label">LATENCY</span>
            <span class="data-value latency">
              <span class="latency-bar"></span>
              &lt; 100ms
            </span>
          </div>
        </div>
      </div>

      <!-- 底部信息 -->
      <div class="footer-info">
        <div class="info-item">
          <span class="info-label">系统版本</span>
          <span class="info-value">v2.0.1</span>
        </div>
        <div class="info-divider"></div>
        <div class="info-item">
          <span class="info-label">安全等级</span>
          <span class="info-value secure">A+</span>
        </div>
        <div class="info-divider"></div>
        <div class="info-item">
          <span class="info-label">运行时间</span>
          <span class="info-value">{{ uptime }}</span>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'MenuView',
  setup() {
    const router = useRouter()
    const uptime = ref('00:00:00')
    let uptimeInterval = null

    // 计算运行时间
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

    const handleLogout = () => {
      if (confirm('确定要断开连接吗？')) {
        localStorage.removeItem('token')
        router.push('/login')
      }
    }

    return {
      uptime,
      enterChat,
      handleLogout
    }
  }
}
</script>

<style scoped>
.menu-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
  background: var(--sci-fi-bg-dark);
  overflow: hidden;
}

/* 背景效果 */
.bg-effects {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
}

.grid-floor {
  position: absolute;
  bottom: -50%;
  left: -50%;
  width: 200%;
  height: 100%;
  background-image: 
    linear-gradient(rgba(16, 185, 129, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(16, 185, 129, 0.05) 1px, transparent 1px);
  background-size: 60px 60px;
  transform: perspective(500px) rotateX(60deg);
  animation: gridMove 20s linear infinite;
}

@keyframes gridMove {
  0% { transform: perspective(500px) rotateX(60deg) translateY(0); }
  100% { transform: perspective(500px) rotateX(60deg) translateY(60px); }
}

.floating-particles {
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
  animation: float 20s infinite;
  opacity: 0.4;
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
.particle:nth-child(21) { left: 12%; top: 15%; animation-delay: 5.5s; }
.particle:nth-child(22) { left: 22%; top: 85%; animation-delay: 6.5s; }
.particle:nth-child(23) { left: 32%; top: 45%; animation-delay: 7.5s; }
.particle:nth-child(24) { left: 42%; top: 65%; animation-delay: 8.5s; }
.particle:nth-child(25) { left: 52%; top: 25%; animation-delay: 9.5s; }
.particle:nth-child(26) { left: 62%; top: 75%; animation-delay: 10.5s; }
.particle:nth-child(27) { left: 72%; top: 55%; animation-delay: 11.5s; }
.particle:nth-child(28) { left: 82%; top: 15%; animation-delay: 12.5s; }
.particle:nth-child(29) { left: 92%; top: 85%; animation-delay: 13.5s; }
.particle:nth-child(30) { left: 8%; top: 45%; animation-delay: 14.5s; }

@keyframes float {
  0%, 100% { 
    transform: translateY(0) translateX(0);
    opacity: 0;
  }
  10% { opacity: 0.4; }
  90% { opacity: 0.4; }
  50% { 
    transform: translateY(-150px) translateX(30px);
  }
}

.glow-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.1;
  animation: orbPulse 6s ease-in-out infinite;
}

.orb-1 {
  width: 500px;
  height: 500px;
  background: var(--sci-fi-primary);
  top: -200px;
  right: -100px;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: var(--sci-fi-secondary);
  bottom: -100px;
  left: -100px;
  animation-delay: 3s;
}

@keyframes orbPulse {
  0%, 100% { opacity: 0.08; transform: scale(1); }
  50% { opacity: 0.15; transform: scale(1.1); }
}

/* 头部导航 */
.sci-fi-header {
  position: relative;
  z-index: 10;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 40px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(16, 185, 129, 0.2);
  box-shadow: 0 2px 20px rgba(16, 185, 129, 0.1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  font-size: 28px;
  color: var(--sci-fi-primary);
  text-shadow: 0 0 20px var(--sci-fi-primary);
  animation: iconGlow 2s ease-in-out infinite;
}

@keyframes iconGlow {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
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

.header-line {
  width: 60px;
  height: 1px;
  background: linear-gradient(90deg, var(--sci-fi-primary), transparent);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.user-info {
  font-family: 'Rajdhani', sans-serif;
  font-size: 14px;
}

.user-status {
  color: var(--sci-fi-success);
  animation: statusBlink 2s ease-in-out infinite;
}

@keyframes statusBlink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  letter-spacing: 2px;
  color: var(--sci-fi-text-secondary);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.logout-btn:hover {
  border-color: var(--sci-fi-danger);
  color: var(--sci-fi-danger);
  background: rgba(241, 91, 181, 0.1);
}

.btn-icon {
  font-size: 10px;
}

/* 主内容区 */
.main-content {
  position: relative;
  z-index: 1;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 40px;
  gap: 48px;
}

/* 欢迎区域 */
.welcome-section {
  text-align: center;
}

.main-title {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.title-line {
  font-family: 'Orbitron', sans-serif;
  font-size: 42px;
  font-weight: 700;
  letter-spacing: 12px;
  color: var(--sci-fi-text-primary);
  text-shadow: 0 0 40px rgba(16, 185, 129, 0.3);
}

.title-sub {
  font-family: 'Rajdhani', sans-serif;
  font-size: 14px;
  letter-spacing: 8px;
  color: var(--sci-fi-text-muted);
}

.description {
  font-family: 'Rajdhani', sans-serif;
  font-size: 18px;
  color: var(--sci-fi-text-secondary);
  max-width: 500px;
  line-height: 1.6;
}

/* 功能卡片 */
.feature-card {
  position: relative;
  width: 600px;
  padding: 48px;
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 16px;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.4s ease;
  box-shadow: 0 4px 20px rgba(16, 185, 129, 0.15);
}

.feature-card:hover {
  border-color: rgba(16, 185, 129, 0.5);
  box-shadow: 0 0 60px rgba(16, 185, 129, 0.2);
  transform: translateY(-5px);
}

/* 卡片框架 */
.card-frame {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}

.frame-corner {
  position: absolute;
  width: 20px;
  height: 20px;
  border: 2px solid var(--sci-fi-primary);
  transition: all 0.3s ease;
}

.frame-corner.tl { top: -1px; left: -1px; border-right: none; border-bottom: none; }
.frame-corner.tr { top: -1px; right: -1px; border-left: none; border-bottom: none; }
.frame-corner.bl { bottom: -1px; left: -1px; border-right: none; border-top: none; }
.frame-corner.br { bottom: -1px; right: -1px; border-left: none; border-top: none; }

.feature-card:hover .frame-corner {
  width: 30px;
  height: 30px;
  box-shadow: 0 0 10px var(--sci-fi-primary);
}

/* 扫描线 */
.card-scan {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--sci-fi-primary), transparent);
  animation: scanMove 4s linear infinite;
  opacity: 0.5;
}

@keyframes scanMove {
  0% { transform: translateY(0); }
  100% { transform: translateY(400px); }
}

/* 卡片内容 */
.card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 24px;
}

.card-icon-wrapper {
  position: relative;
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-ring {
  position: absolute;
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 50%;
  animation: ringPulse 2s ease-in-out infinite;
}

.ring-1 { width: 60px; height: 60px; }
.ring-2 { width: 80px; height: 80px; animation-delay: 0.3s; }
.ring-3 { width: 100px; height: 100px; animation-delay: 0.6s; }

@keyframes ringPulse {
  0%, 100% { transform: scale(1); opacity: 0.5; }
  50% { transform: scale(1.1); opacity: 1; }
}

.card-icon {
  font-size: 48px;
  color: var(--sci-fi-primary);
  text-shadow: 0 0 30px var(--sci-fi-primary);
  animation: iconFloat 3s ease-in-out infinite;
}

@keyframes iconFloat {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

.card-text {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 28px;
  font-weight: 600;
  letter-spacing: 4px;
  color: var(--sci-fi-text-primary);
  margin: 0;
}

.card-subtitle {
  font-family: 'Rajdhani', sans-serif;
  font-size: 12px;
  letter-spacing: 6px;
  color: var(--sci-fi-text-muted);
  margin: 0;
}

.card-desc {
  font-family: 'Rajdhani', sans-serif;
  font-size: 16px;
  color: var(--sci-fi-text-secondary);
  line-height: 1.6;
  max-width: 400px;
  margin: 16px 0 0 0;
}

/* 卡片操作 */
.card-action {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 32px;
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-success));
  border-radius: 8px;
  transition: all 0.3s ease;
}

.feature-card:hover .card-action {
  box-shadow: 0 10px 30px rgba(16, 185, 129, 0.3);
}

.action-text {
  font-family: 'Orbitron', sans-serif;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 2px;
  color: var(--sci-fi-bg-dark);
}

.action-arrow {
  font-size: 18px;
  color: var(--sci-fi-bg-dark);
  transition: transform 0.3s ease;
}

.feature-card:hover .action-arrow {
  transform: translateX(5px);
}

/* 数据装饰 */
.data-decoration {
  position: absolute;
  bottom: 24px;
  right: 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-family: 'Rajdhani', sans-serif;
  font-size: 11px;
}

.data-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.data-label {
  color: var(--sci-fi-text-muted);
  letter-spacing: 2px;
}

.data-value {
  color: var(--sci-fi-text-secondary);
  letter-spacing: 1px;
}

.data-value.active {
  color: var(--sci-fi-success);
  animation: textPulse 2s ease-in-out infinite;
}

.latency {
  display: flex;
  align-items: center;
  gap: 6px;
}

.latency-bar {
  width: 20px;
  height: 3px;
  background: var(--sci-fi-success);
  border-radius: 2px;
  animation: barPulse 1s ease-in-out infinite;
}

@keyframes barPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* 底部信息 */
.footer-info {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 20px 40px;
  font-family: 'Rajdhani', sans-serif;
  font-size: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-label {
  color: var(--sci-fi-text-muted);
  letter-spacing: 2px;
}

.info-value {
  color: var(--sci-fi-text-secondary);
  font-weight: 500;
}

.info-value.secure {
  color: var(--sci-fi-success);
}

.info-divider {
  width: 1px;
  height: 12px;
  background: rgba(255, 255, 255, 0.1);
}

/* 响应式 */
@media (max-width: 768px) {
  .sci-fi-header {
    padding: 16px 20px;
  }
  
  .logo-text {
    font-size: 16px;
  }
  
  .title-line {
    font-size: 28px;
    letter-spacing: 6px;
  }
  
  .feature-card {
    width: 90%;
    padding: 32px 24px;
  }
  
  .card-title {
    font-size: 22px;
  }
  
  .data-decoration {
    position: relative;
    bottom: auto;
    right: auto;
    margin-top: 24px;
    align-items: center;
  }
  
  .footer-info {
    flex-wrap: wrap;
    justify-content: center;
    gap: 16px;
  }
}
</style>
