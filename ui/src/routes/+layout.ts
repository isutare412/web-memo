import { syncUserData } from '$lib/auth'
import { getPreferredTheme, setDocumentDataTheme } from '$lib/theme'
import type { LayoutLoad } from './$types'

export const ssr = false

export const load = (async () => {
  await syncUserData()

  const theme = getPreferredTheme()
  setDocumentDataTheme(theme)

  return {
    preferredTheme: theme,
  }
}) satisfies LayoutLoad
