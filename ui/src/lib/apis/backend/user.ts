import { getErrorResponse } from '$lib/apis/backend/error'
import { error } from '@sveltejs/kit'

interface GetSelfUserResponse {
  id: string
  userType: string
  email: string
  userName: string
  givenName?: string
  familyName?: string
  photoUrl?: string
}

export async function getSelfUser(): Promise<GetSelfUserResponse> {
  const response = await fetch('/api/v1/users/me')
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw error(response.status, errorResponse.msg ?? response.statusText)
  }
  return response.json()
}
