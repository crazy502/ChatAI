import axios from 'axios'

import { translate } from '../i18n'
import { clearAuth } from './auth'

export const API_BASE_URL = '/api'

const normalizePath = (path) => {
  if (!path) {
    return ''
  }

  return path.startsWith('/') ? path : `/${path}`
}

export const buildApiUrl = (path) => `${API_BASE_URL}${normalizePath(path)}`

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 0
})

const AUTH_ERROR_CODES = new Set([2006, 2007])

const redirectToLogin = () => {
  clearAuth()
  if (window.location.pathname !== '/login') {
    window.location.assign('/login')
  }
}

export const getApiErrorMessage = (error, fallbackKey = 'common.unknownError') => {
  const resultCode = Number(error?.response?.data?.status_code || error?.response?.data?.error?.code || 0)
  if (resultCode) {
    const key = `errors.${resultCode}`
    const translated = translate(key)
    if (translated !== key) {
      return translated
    }
  }
  return translate(fallbackKey)
}

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => {
    if (AUTH_ERROR_CODES.has(Number(response?.data?.status_code))) {
      redirectToLogin()
    }
    return response
  },
  (error) => {
	const resultCode = Number(error?.response?.data?.status_code || error?.response?.data?.error?.code || 0)
    if (error?.response?.status === 401 || AUTH_ERROR_CODES.has(resultCode)) {
      redirectToLogin()
    }
	if (error instanceof Error) {
	  error.message = getApiErrorMessage(error)
    }
    return Promise.reject(error)
  }
)

export default api

