import { goto } from '$app/navigation'
import { getSelfUser, signOutUser } from '$lib/apis/backend/user'
import { get, writable } from 'svelte/store'

export const authStore = writable<AuthState>({})

export interface UserData {
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
  lastSync?: Date
}

export async function syncUserData() {
  const lastSync = get(authStore).lastSync?.getTime()
  if (lastSync !== undefined) {
    const hourAgo = Date.now() - 60 * 60 * 1000
    if (lastSync > hourAgo) return
  }

  const response = await getSelfUser()
  if (response === undefined) {
    authStore.set({ user: undefined, lastSync: undefined })
    return
  }

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
    lastSync: new Date(),
  })
}

export function signInGoogle() {
  window.location.assign('/api/v1/google/sign-in')
}

export async function signOut() {
  await signOutUser()
  authStore.set({})
  goto('/')
}
