import { createI18n } from 'vue-i18n'
import { enUS } from './en-US'
import { zhCN } from './zh-CN'

export type AppLocale = 'zh-CN' | 'en-US'

const supportedLocales: readonly AppLocale[] = ['zh-CN', 'en-US']

const resolveInitialLocale = (): AppLocale => {
  const saved = window.localStorage.getItem('locale')
  if (saved && supportedLocales.includes(saved as AppLocale)) {
    return saved as AppLocale
  }
  return window.navigator.language.toLowerCase().startsWith('zh') ? 'zh-CN' : 'en-US'
}

export const i18n = createI18n({
  legacy: false,
  locale: resolveInitialLocale(),
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS
  }
})

export const setLocale = (locale: AppLocale): void => {
  i18n.global.locale.value = locale
  window.localStorage.setItem('locale', locale)
  document.documentElement.lang = locale
}

export const toggleLocale = (): void => {
  setLocale(i18n.global.locale.value === 'zh-CN' ? 'en-US' : 'zh-CN')
}

export const translate = (key: string, params?: Record<string, string | number>): string => {
  return params ? i18n.global.t(key, params) : i18n.global.t(key)
}

document.documentElement.lang = i18n.global.locale.value
