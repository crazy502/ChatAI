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

export const isAdminToken = (token) => {
  const payload = parseTokenPayload(token)
  const userName = payload?.username || ''
  return userName === 'admin' && Boolean(payload?.is_admin ?? payload?.isAdmin)
}

export const getTokenUserName = (token) => {
  const payload = parseTokenPayload(token)
  return payload?.username || ''
}
