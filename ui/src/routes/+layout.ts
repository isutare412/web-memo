import { getPreferredTheme } from '$lib/theme'
import type { LayoutLoad } from './$types'

export const ssr = false

export const load = (() => {
  return {
    preferredTheme: getPreferredTheme(),
  }
}) satisfies LayoutLoad
