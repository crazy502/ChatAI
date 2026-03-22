import { reactive, readonly } from 'vue'

const createDefaultConfirmState = () => ({
  visible: false,
  title: '',
  message: '',
  confirmText: '确认',
  cancelText: '取消',
  intent: 'primary'
})

const state = reactive({
  toasts: [],
  confirmDialog: createDefaultConfirmState()
})

let toastSeed = 0
let confirmResolver = null

const removeToast = (toastId) => {
  const toastIndex = state.toasts.findIndex((toast) => toast.id === toastId)
  if (toastIndex !== -1) {
    state.toasts.splice(toastIndex, 1)
  }
}

const showToast = (message, type = 'info', options = {}) => {
  const toast = {
    id: `toast-${Date.now()}-${++toastSeed}`,
    message,
    type,
    duration: options.duration ?? 3200
  }

  state.toasts.push(toast)

  if (toast.duration > 0) {
    window.setTimeout(() => {
      removeToast(toast.id)
    }, toast.duration)
  }

  return toast.id
}

const resetConfirmDialog = () => {
  Object.assign(state.confirmDialog, createDefaultConfirmState())
}

const resolveConfirm = (accepted) => {
  const resolver = confirmResolver
  confirmResolver = null
  resetConfirmDialog()
  resolver?.(accepted)
}

const confirmAction = (options = {}) => {
  if (confirmResolver) {
    confirmResolver(false)
  }

  Object.assign(state.confirmDialog, createDefaultConfirmState(), {
    visible: true,
    ...options
  })

  return new Promise((resolve) => {
    confirmResolver = resolve
  })
}

export function useUi() {
  return {
    uiState: readonly(state),
    showToast,
    removeToast,
    confirmAction,
    acceptConfirm: () => resolveConfirm(true),
    cancelConfirm: () => resolveConfirm(false)
  }
}

