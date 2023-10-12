import { goto } from '$app/navigation'
import { getSelfUser } from '$lib/apis/backend/user'
import Cookies from 'js-cookie'
import { writable } from 'svelte/store'

export const authStore = writable<AuthState>({})

const appTokenKey = 'wmToken'

interface UserData {
  id: string
  userType: string
  email: string
  name: string
  givenName?: string
  familyName?: string
  photoUrl?: string
}

interface AuthState {
  user?: UserData
}

export async function syncUserData() {
  if (!Cookies.get(appTokenKey)) {
    return undefined
  }

  const response = await getSelfUser()
  authStore.set({
    user: {
      id: response.id,
      userType: response.userType,
      email: response.email,
      name: response.userName,
      givenName: response.givenName,
      familyName: response.familyName,
      photoUrl: response.photoUrl,
    },
  })
}

export function signInGoogle() {
  window.location.assign('/api/v1/google/sign-in')
}

export function signOut() {
  Cookies.remove('wmToken')
  authStore.set({})
  goto('/')
}
