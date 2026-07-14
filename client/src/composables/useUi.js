import { reactive, readonly } from 'vue'

import { translate } from '../i18n'

const createDefaultConfirmState = () => ({
  visible: false,
  title: '',
  message: '',
  confirmText: translate('common.confirm'),
  cancelText: translate('common.cancel'),
  intent: 'primary'
})

const createDefaultInputState = () => ({
  visible: false,
  title: '',
  placeholder: '',
  defaultValue: '',
  confirmText: translate('common.confirm'),
  cancelText: translate('common.cancel')
})

const state = reactive({
  toasts: [],
  confirmDialog: createDefaultConfirmState(),
  inputDialog: createDefaultInputState()
})

let toastSeed = 0
let confirmResolver = null
let inputResolver = null

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

const resetInputDialog = () => {
  Object.assign(state.inputDialog, createDefaultInputState())
}

const resolveInput = (value) => {
  const resolver = inputResolver
  inputResolver = null
  resetInputDialog()
  resolver?.(value)
}

const promptAction = (options = {}) => {
  if (inputResolver) {
    inputResolver(null)
  }

  Object.assign(state.inputDialog, createDefaultInputState(), {
    visible: true,
    ...options
  })

  return new Promise((resolve) => {
    inputResolver = resolve
  })
}

export function useUi() {
  return {
    uiState: readonly(state),
    showToast,
    removeToast,
    confirmAction,
	promptAction,
    acceptConfirm: () => resolveConfirm(true),
	cancelConfirm: () => resolveConfirm(false),
	acceptInput: (value) => resolveInput(value),
	cancelInput: () => resolveInput(null)
  }
}

