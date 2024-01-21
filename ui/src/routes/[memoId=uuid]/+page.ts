import { authStore } from '$lib/auth'
import { redirect } from '@sveltejs/kit'
import { get } from 'svelte/store'
import type { PageLoad } from './$types'

export const load = (async (event) => {
  const { parent } = event
  await parent()

  if (get(authStore).user === undefined) {
    redirect(302, '/')
  }
}) satisfies PageLoad
