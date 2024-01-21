import { getPreferredPageSize } from '$lib/memo'
import { setPageSizeOfSearchParams } from '$lib/searchParams'
import { redirect } from '@sveltejs/kit'
import type { PageLoad } from './$types'

export const load = (async (event) => {
  const preferredPageSize = getPreferredPageSize()

  const searchParams = event.url.searchParams
  if (
    searchParams.get('ps') === null &&
    preferredPageSize !== null &&
    setPageSizeOfSearchParams(searchParams, preferredPageSize)
  ) {
    redirect(302, `/?${searchParams.toString()}`)
  }
}) satisfies PageLoad
