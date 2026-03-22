<template>
  <teleport to="body">
    <div class="toast-region" aria-live="polite" aria-atomic="false">
      <transition-group name="toast-move" tag="div" class="toast-stack">
        <div
          v-for="toast in uiState.toasts"
          :key="toast.id"
          class="app-toast"
          :class="`type-${toast.type}`"
          role="status"
        >
          <div class="toast-icon">{{ toastIcons[toast.type] || toastIcons.info }}</div>
          <div class="toast-content">
            <p class="toast-title">{{ toastTitles[toast.type] || toastTitles.info }}</p>
            <p class="toast-message">{{ toast.message }}</p>
          </div>
          <button class="toast-close" type="button" aria-label="关闭提示" @click="removeToast(toast.id)">
            ×
          </button>
        </div>
      </transition-group>
    </div>
  </teleport>
</template>

<script setup>
import { useUi } from '../../composables/useUi'

const toastIcons = {
  success: 'OK',
  error: '!',
  warning: '!',
  info: 'i'
}

const toastTitles = {
  success: '操作成功',
  error: '操作失败',
  warning: '请注意',
  info: '系统提示'
}

const { uiState, removeToast } = useUi()
</script>

<style scoped>
.toast-region {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 5000;
  pointer-events: none;
}

.toast-stack {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.app-toast {
  min-width: 280px;
  max-width: 360px;
  display: grid;
  grid-template-columns: 36px 1fr 24px;
  align-items: start;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(16, 185, 129, 0.2);
  box-shadow: 0 16px 40px rgba(15, 23, 42, 0.12);
  backdrop-filter: blur(18px);
  pointer-events: auto;
}

.app-toast.type-success {
  border-color: rgba(16, 185, 129, 0.35);
}

.app-toast.type-error {
  border-color: rgba(239, 68, 68, 0.35);
}

.app-toast.type-warning {
  border-color: rgba(245, 158, 11, 0.35);
}

.toast-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  font-family: 'Orbitron', sans-serif;
  font-size: 16px;
  font-weight: 700;
  background: rgba(16, 185, 129, 0.12);
  color: var(--sci-fi-primary);
}

.app-toast.type-error .toast-icon {
  background: rgba(239, 68, 68, 0.12);
  color: var(--sci-fi-danger);
}

.app-toast.type-warning .toast-icon {
  background: rgba(245, 158, 11, 0.12);
  color: var(--sci-fi-warning);
}

.toast-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.toast-title {
  margin: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  letter-spacing: 1px;
  text-transform: uppercase;
  color: var(--sci-fi-text-secondary);
}

.toast-message {
  margin: 0;
  font-size: 14px;
  line-height: 1.5;
  color: var(--sci-fi-text-primary);
}

.toast-close {
  appearance: none;
  border: none;
  background: transparent;
  color: var(--sci-fi-text-muted);
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
  padding: 0;
}

.toast-close:hover {
  color: var(--sci-fi-text-primary);
}

.toast-move-enter-active,
.toast-move-leave-active {
  transition: transform 0.25s ease, opacity 0.25s ease;
}

.toast-move-enter-from,
.toast-move-leave-to {
  opacity: 0;
  transform: translateX(24px) translateY(-8px);
}

@media (max-width: 640px) {
  .toast-region {
    left: 12px;
    right: 12px;
    top: 12px;
  }

  .app-toast {
    min-width: 0;
    max-width: none;
    width: 100%;
  }
}
</style>


