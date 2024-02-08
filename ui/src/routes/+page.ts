import { getPreferredPageSize, getPreferredSortOrder } from '$lib/memo'
import { setPageSizeOfSearchParams, setSortOrderOfSearchParams } from '$lib/searchParams'
import { redirect } from '@sveltejs/kit'
import type { PageLoad } from './$types'

export const load = (async (event) => {
  const { url } = event

  const searchParams = url.searchParams
  let paramUpdated = false

  const pageSize = getPreferredPageSize()
  paramUpdated ||= pageSize !== null && setPageSizeOfSearchParams(searchParams, pageSize)

  const sortOrder = getPreferredSortOrder()
  paramUpdated ||= sortOrder !== null && setSortOrderOfSearchParams(searchParams, sortOrder)

  if (paramUpdated) {
    redirect(302, `/?${searchParams.toString()}`)
  }
}) satisfies PageLoad
