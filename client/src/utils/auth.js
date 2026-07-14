const decodeBase64Url = (value) => {
  const normalized = value.replace(/-/g, '+').replace(/_/g, '/')
  const padding = normalized.length % 4
  const base64 = padding ? normalized.padEnd(normalized.length + (4 - padding), '=') : normalized

  return atob(base64)
}

export const parseTokenPayload = (token) => {
  if (!token) {
    return null
  }

  const parts = token.split('.')
  if (parts.length < 2) {
    return null
  }

  try {
    return JSON.parse(decodeBase64Url(parts[1]))
  } catch (error) {
    console.error('Parse token payload error:', error)
    return null
  }
}

export const isTokenUsable = (token, now = Date.now()) => {
  const payload = parseTokenPayload(token)
  if (!payload || typeof payload.exp !== 'number') {
    return false
  }

  const nowSeconds = Math.floor(now / 1000)
  if (payload.exp <= nowSeconds) {
    return false
  }
  if (typeof payload.nbf === 'number' && payload.nbf > nowSeconds) {
    return false
  }
  return true
}

export const clearAuth = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('isAdmin')
}

export const isAdminToken = (token) => {
  if (!isTokenUsable(token)) {
    return false
  }
  const payload = parseTokenPayload(token)
  return Boolean(payload?.is_admin ?? payload?.isAdmin)
}

export const getTokenUserName = (token) => {
  if (!isTokenUsable(token)) {
    return ''
  }
  const payload = parseTokenPayload(token)
  return payload?.username || ''
}
