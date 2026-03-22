import axios from 'axios'

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
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api

