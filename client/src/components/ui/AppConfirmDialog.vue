<template>
  <teleport to="body">
    <transition name="dialog-fade">
      <div v-if="uiState.confirmDialog.visible" class="dialog-overlay" @click.self="cancelConfirm()">
        <div class="dialog-panel" role="dialog" aria-modal="true" :aria-label="uiState.confirmDialog.title">
          <div class="dialog-kicker">SYSTEM CONFIRM</div>
          <h2 class="dialog-title">{{ uiState.confirmDialog.title }}</h2>
          <p class="dialog-message">{{ uiState.confirmDialog.message }}</p>
          <div class="dialog-actions">
            <button class="dialog-btn secondary" type="button" @click="cancelConfirm()">
              {{ uiState.confirmDialog.cancelText }}
            </button>
            <button
              class="dialog-btn primary"
              :class="`intent-${uiState.confirmDialog.intent}`"
              type="button"
              @click="acceptConfirm()"
            >
              {{ uiState.confirmDialog.confirmText }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import { onBeforeUnmount, watch } from 'vue'
import { useUi } from '../../composables/useUi'

const { uiState, acceptConfirm, cancelConfirm } = useUi()

const handleKeydown = (event) => {
  if (event.key === 'Escape' && uiState.confirmDialog.visible) {
    cancelConfirm()
  }
}

watch(
  () => uiState.confirmDialog.visible,
  (visible) => {
    if (visible) {
      window.addEventListener('keydown', handleKeydown)
    } else {
      window.removeEventListener('keydown', handleKeydown)
    }
  }
)

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  z-index: 4900;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: rgba(15, 23, 42, 0.38);
  backdrop-filter: blur(10px);
}

.dialog-panel {
  width: min(100%, 420px);
  padding: 28px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(16, 185, 129, 0.2);
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.18);
}

.dialog-kicker {
  margin-bottom: 12px;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--sci-fi-primary);
}

.dialog-title {
  margin: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: 22px;
  letter-spacing: 2px;
  color: var(--sci-fi-text-primary);
}

.dialog-message {
  margin: 16px 0 0;
  font-size: 15px;
  line-height: 1.7;
  color: var(--sci-fi-text-secondary);
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.dialog-btn {
  min-width: 96px;
  padding: 10px 16px;
  border-radius: 10px;
  border: 1px solid rgba(16, 185, 129, 0.2);
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  letter-spacing: 1px;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.dialog-btn:hover {
  transform: translateY(-1px);
}

.dialog-btn.secondary {
  background: rgba(255, 255, 255, 0.92);
  color: var(--sci-fi-text-secondary);
}

.dialog-btn.primary {
  color: #fff;
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-secondary));
  border-color: transparent;
}

.dialog-btn.primary.intent-danger {
  background: linear-gradient(135deg, #ef4444, #f97316);
}

.dialog-btn.primary:hover {
  box-shadow: 0 12px 24px rgba(16, 185, 129, 0.2);
}

.dialog-btn.primary.intent-danger:hover {
  box-shadow: 0 12px 24px rgba(239, 68, 68, 0.2);
}

.dialog-fade-enter-active,
.dialog-fade-leave-active {
  transition: opacity 0.2s ease;
}

.dialog-fade-enter-from,
.dialog-fade-leave-to {
  opacity: 0;
}
</style>

