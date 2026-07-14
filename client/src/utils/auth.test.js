import { describe, expect, it } from 'vitest'
import { isAdminToken, isTokenUsable, parseTokenPayload } from './auth'

const encodePayload = (payload) => {
  const encoded = btoa(JSON.stringify(payload)).replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_')
  return `header.${encoded}.signature`
}

describe('token helpers', () => {
  it('accepts an unexpired token', () => {
    const now = Date.UTC(2026, 0, 1)
    const token = encodePayload({ exp: now / 1000 + 60, username: 'user@qq.com' })

    expect(isTokenUsable(token, now)).toBe(true)
    expect(parseTokenPayload(token)?.username).toBe('user@qq.com')
  })

  it('rejects an expired token', () => {
    const now = Date.UTC(2026, 0, 1)
    const token = encodePayload({ exp: now / 1000 - 1, is_admin: true })

    expect(isTokenUsable(token, now)).toBe(false)
  })

  it('does not grant administrator access to an expired token', () => {
    const token = encodePayload({ exp: 1, is_admin: true })

    expect(isAdminToken(token)).toBe(false)
  })
})
