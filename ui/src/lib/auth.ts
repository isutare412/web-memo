import { goto } from '$app/navigation'
import { getSelfUser, refreshUserToken, signOutUser, type RawUser } from '$lib/apis/backend/user'
import { get, writable } from 'svelte/store'

export const authStore = writable<AuthState>({})

export interface User {
  id: string
  userType: string
  email: string
  name: string
  givenName?: string
  familyName?: string
  photoUrl?: string
  issuedAt: Date
  expireAt: Date
}

interface AuthState {
  user?: User
  lastSync?: Date
}

export async function syncUserData() {
  const lastSync = get(authStore).lastSync?.getTime()
  if (lastSync !== undefined) {
    const hourAgo = Date.now() - 60 * 60 * 1000
    if (lastSync > hourAgo) return
  }

  let rawUser = await getSelfUser()
  if (rawUser === undefined) {
    authStore.set({ user: undefined, lastSync: undefined })
    return
  }

  let user = convertToUser(rawUser)
  const twoWeeksLater = Date.now() + 14 * 24 * 60 * 60 * 1000
  if (user.expireAt.getTime() < twoWeeksLater) {
    rawUser = await refreshUserToken()
    user = convertToUser(rawUser)
  }

  authStore.set({
    user,
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

function convertToUser(rawUser: RawUser): User {
  return {
    id: rawUser.id,
    userType: rawUser.userType,
    email: rawUser.email,
    name: rawUser.userName,
    givenName: rawUser.givenName,
    familyName: rawUser.familyName,
    photoUrl: rawUser.photoUrl,
    issuedAt: new Date(rawUser.issuedAt),
    expireAt: new Date(rawUser.expireAt),
  } satisfies User
}
