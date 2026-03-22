<template>
  <div id="app">
    <!-- 页面加载动画 -->
    <div v-if="isLoading" class="page-loading">
      <div class="loading-spinner">
        <div class="spinner-ring"></div>
        <div class="spinner-ring"></div>
        <div class="spinner-ring"></div>
      </div>
      <div class="loading-text">SYSTEM INITIALIZING</div>
    </div>
    
    <router-view v-slot="{ Component }">
      <transition name="page" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'

export default {
  name: 'App',
  setup() {
    const isLoading = ref(true)
    
    onMounted(() => {
      // 模拟系统初始化
      setTimeout(() => {
        isLoading.value = false
      }, 1500)
    })
    
    return {
      isLoading
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: 'Rajdhani', 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: #0a0a0f;
}

#app {
  height: 100%;
}

/* 页面加载动画 */
.page-loading {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #0a0a0f 0%, #12121a 50%, #0a0a0f 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  transition: opacity 0.6s ease, visibility 0.6s ease;
}

.page-loading.hidden {
  opacity: 0;
  visibility: hidden;
}

.loading-spinner {
  position: relative;
  width: 100px;
  height: 100px;
}

.spinner-ring {
  position: absolute;
  border: 3px solid transparent;
  border-top-color: #00d4ff;
  border-radius: 50%;
  animation: spinner-rotate 1.2s linear infinite;
}

.spinner-ring:nth-child(1) {
  width: 100px;
  height: 100px;
  top: 0;
  left: 0;
  border-top-color: #00d4ff;
}

.spinner-ring:nth-child(2) {
  width: 75px;
  height: 75px;
  top: 12.5px;
  left: 12.5px;
  animation-direction: reverse;
  animation-duration: 0.9s;
  border-top-color: #7b2cbf;
}

.spinner-ring:nth-child(3) {
  width: 50px;
  height: 50px;
  top: 25px;
  left: 25px;
  animation-duration: 0.6s;
  border-top-color: #ff006e;
}

@keyframes spinner-rotate {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.loading-text {
  margin-top: 40px;
  font-family: 'Orbitron', sans-serif;
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 6px;
  color: #00d4ff;
  animation: loading-text-pulse 1.5s ease-in-out infinite;
}

@keyframes loading-text-pulse {
  0%, 100% {
    opacity: 0.4;
    text-shadow: 0 0 10px rgba(0, 212, 255, 0.3);
  }
  50% {
    opacity: 1;
    text-shadow: 0 0 20px rgba(0, 212, 255, 0.8);
  }
}

/* 科幻页面切换动画 */
.page-enter-active {
  transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}

.page-leave-active {
  transition: all 0.4s cubic-bezier(0.7, 0, 0.84, 0);
}

.page-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.98);
  filter: blur(10px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(1.02);
  filter: blur(10px);
}

/* 扫描线效果 */
.page-enter-active::after {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, #00d4ff, transparent);
  animation: scan-line 0.5s ease-out;
  pointer-events: none;
  z-index: 9998;
}

@keyframes scan-line {
  0% {
    transform: translateY(-100vh);
    opacity: 1;
  }
  100% {
    transform: translateY(100vh);
    opacity: 0;
  }
}
</style>
