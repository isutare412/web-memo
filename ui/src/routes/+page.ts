import { getPreferredPageSize } from '$lib/memo'
import { setPageSizeOfSearchParams } from '$lib/searchParams'
import { redirect } from '@sveltejs/kit'
import type { PageLoad } from './$types'

export const load = (async (event) => {
  const { url } = event

  const preferredPageSize = getPreferredPageSize()

  const searchParams = url.searchParams
  if (
    searchParams.get('ps') === null &&
    preferredPageSize !== null &&
    setPageSizeOfSearchParams(searchParams, preferredPageSize)
  ) {
    redirect(302, `/?${searchParams.toString()}`)
  }
}) satisfies PageLoad
