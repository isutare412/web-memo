import { goto } from '$app/navigation'
import { getSelfUser, signOutUser } from '$lib/apis/backend/user'
import { writable } from 'svelte/store'

export const authStore = writable<AuthState>({})

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
  const response = await getSelfUser()
  if (response === undefined) {
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
