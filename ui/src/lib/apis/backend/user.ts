import { StatusError, getErrorResponse } from '$lib/apis/backend/error'

export interface RawUser {
  id: string
  userType: string
  email: string
  userName: string
  givenName?: string
  familyName?: string
  photoUrl?: string
  issuedAt: string
  expireAt: string
}

export async function getSelfUser(): Promise<RawUser | void> {
  const response = await fetch('/api/v1/users/me')
  if (!response.ok) {
    if (response.status == 401) {
      return
    }

    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }
  return response.json()
}

export async function refreshUserToken(): Promise<RawUser> {
  const response = await fetch('/api/v1/users/refresh-token', { method: 'POST' })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }
  return response.json()
}

export async function signOutUser(): Promise<void> {
  const response = await fetch('/api/v1/users/sign-out')
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }
}
