<template>
  <teleport to="body">
    <transition name="dialog-fade">
      <div v-if="uiState.inputDialog.visible" class="dialog-overlay" @click.self="cancelInput()">
        <div class="dialog-panel" role="dialog" aria-modal="true" :aria-label="uiState.inputDialog.title">
          <div class="dialog-kicker">USER INPUT</div>
          <h2 class="dialog-title">{{ uiState.inputDialog.title }}</h2>
          <input
            ref="inputRef"
            v-model="inputValue"
            type="text"
            class="dialog-input"
            :placeholder="uiState.inputDialog.placeholder"
            @keydown.enter="handleEnter"
            @keydown.escape="cancelInput()"
          />
          <div class="dialog-actions">
            <button class="dialog-btn secondary" type="button" @click="cancelInput()">
              {{ uiState.inputDialog.cancelText }}
            </button>
            <button
              class="dialog-btn primary"
              type="button"
              @click="acceptInput(inputValue)"
            >
              {{ uiState.inputDialog.confirmText }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import { ref, watch, nextTick, onBeforeUnmount } from 'vue'
import { useUi } from '../../composables/useUi'

const { uiState, acceptInput, cancelInput } = useUi()
const inputRef = ref(null)
const inputValue = ref('')

watch(
  () => uiState.inputDialog.visible,
  async (visible) => {
    if (visible) {
      inputValue.value = uiState.inputDialog.defaultValue
      await nextTick()
      inputRef.value?.focus()
      inputRef.value?.select()
      window.addEventListener('keydown', handleKeydown)
    } else {
      window.removeEventListener('keydown', handleKeydown)
    }
  }
)

const handleKeydown = (event) => {
  if (event.key === 'Escape' && uiState.inputDialog.visible) {
    cancelInput()
  }
}

const handleEnter = (event) => {
  if (event.isComposing || event.keyCode === 229) {
    return
  }

  acceptInput(inputValue.value)
}

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
  background: rgba(15, 23, 42, 0.18);
  backdrop-filter: blur(10px);
}

.dialog-panel {
  width: min(100%, 420px);
  padding: 28px;
  border-radius: 18px;
  background: var(--panel);
  border: 1px solid var(--border);
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.12);
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

.dialog-input {
  width: 100%;
  margin-top: 16px;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid var(--border);
  background: var(--panel-soft);
  font-size: 14px;
  color: var(--sci-fi-text-primary);
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
  box-sizing: border-box;
}

.dialog-input:focus {
  outline: none;
  border-color: var(--sci-fi-primary);
  box-shadow: 0 0 0 3px var(--accent-soft);
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
  border: 1px solid var(--border);
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
  background: var(--panel-soft);
  color: var(--sci-fi-text-secondary);
}

.dialog-btn.primary {
  color: var(--btn-text);
  background: linear-gradient(135deg, var(--accent), var(--accent-strong));
  border-color: transparent;
}

.dialog-btn.primary:hover {
  box-shadow: 0 12px 24px rgba(var(--mint-rgb), 0.26);
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




