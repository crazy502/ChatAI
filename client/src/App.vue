<template>
  <div id="app">
    <div v-if="isLoading" class="page-loading">
      <div class="loading-spinner">
        <div class="spinner-ring"></div>
        <div class="spinner-ring"></div>
        <div class="spinner-ring"></div>
      </div>
      <div class="loading-text">GOPHERAI INITIALIZING</div>
    </div>

    <router-view v-slot="{ Component, route }">
      <transition name="page" mode="out-in">
        <component :is="Component" :key="route.fullPath" />
      </transition>
    </router-view>

    <AppToastContainer />
    <AppConfirmDialog />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AppToastContainer from './components/ui/AppToastContainer.vue'
import AppConfirmDialog from './components/ui/AppConfirmDialog.vue'

export default {
  name: 'App',
  components: {
    AppToastContainer,
    AppConfirmDialog
  },
  setup() {
    const router = useRouter()
    const isLoading = ref(true)

    onMounted(async () => {
      const startedAt = performance.now()
      await router.isReady()
      const elapsed = performance.now() - startedAt
      const remaining = Math.max(0, 320 - elapsed)

      window.setTimeout(() => {
        isLoading.value = false
      }, remaining)
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

html,
body {
  height: 100%;
  font-family: 'Rajdhani', 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: var(--bg-1);
}

#app {
  height: 100%;
}

.page-loading {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background:
    radial-gradient(circle at top right, rgba(var(--mint-2-rgb), 0.14), transparent 30%),
    radial-gradient(circle at bottom left, rgba(var(--mint-rgb), 0.12), transparent 26%),
    radial-gradient(circle at 22% 86%, rgba(var(--info-rgb), 0.08), transparent 36%),
    linear-gradient(135deg, var(--bg-1) 0%, #f2f6f4 54%, #edf4f1 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  transition: opacity 0.35s ease, visibility 0.35s ease;
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
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spinner-rotate 1.2s linear infinite;
}

.spinner-ring:nth-child(1) {
  width: 100px;
  height: 100px;
  top: 0;
  left: 0;
  border-top-color: var(--accent);
}

.spinner-ring:nth-child(2) {
  width: 75px;
  height: 75px;
  top: 12.5px;
  left: 12.5px;
  animation-direction: reverse;
  animation-duration: 0.9s;
  border-top-color: var(--accent-strong);
}

.spinner-ring:nth-child(3) {
  width: 50px;
  height: 50px;
  top: 25px;
  left: 25px;
  animation-duration: 0.6s;
  border-top-color: rgba(var(--mint-2-rgb), 0.8);
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
  color: var(--accent);
  animation: loading-text-pulse 1.5s ease-in-out infinite;
}

@keyframes loading-text-pulse {
  0%,
  100% {
    opacity: 0.4;
    text-shadow: 0 0 10px var(--mint-shadow);
  }

  50% {
    opacity: 1;
    text-shadow: 0 0 20px rgba(var(--mint-rgb), 0.45);
  }
}

.page-enter-active {
  transition: opacity 0.38s cubic-bezier(0.16, 1, 0.3, 1), transform 0.38s cubic-bezier(0.16, 1, 0.3, 1);
}

.page-leave-active {
  transition: opacity 0.22s ease, transform 0.22s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(14px) scale(0.99);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(1.005);
}

.page-enter-active::after {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--accent), transparent);
  animation: scan-line 0.45s ease-out;
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

@media (prefers-reduced-motion: reduce) {
  .page-loading,
  .page-enter-active,
  .page-leave-active {
    transition: none !important;
  }

  .spinner-ring,
  .loading-text,
  .page-enter-active::after {
    animation: none !important;
  }
}
</style>



