import { getPreferredPageSize, getPreferredSortOrder } from '$lib/memo'
import { setPageSizeOfSearchParams, setSortOrderOfSearchParams } from '$lib/searchParams'
import { redirect } from '@sveltejs/kit'
import type { PageLoad } from './$types'

export const load = (async (event) => {
  const { url } = event

  const searchParams = url.searchParams
  let paramUpdated = false

  const pageSize = getPreferredPageSize()
  if (pageSize !== null && setPageSizeOfSearchParams(searchParams, pageSize)) {
    paramUpdated = true
  }

  const sortOrder = getPreferredSortOrder()
  if (sortOrder !== null && setSortOrderOfSearchParams(searchParams, sortOrder)) {
    paramUpdated = true
  }

  if (paramUpdated) {
    redirect(302, `/?${searchParams.toString()}`)
  }
}) satisfies PageLoad
