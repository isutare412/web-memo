import { browser } from '$app/environment'
import { getMemo } from '$lib/apis/backend/memo'
import { mapToMemo, type Memo } from '$lib/memo'
import { getErrorMessage } from '$lib/utils/error'
import type { PageLoad } from './$types'

export const ssr = true

export const load = (async (event) => {
  const { params, fetch } = event

  let memo: Memo | undefined
  let errorMessage: string | undefined
  try {
    memo = mapToMemo(await getMemo(params.memoId, { fetch }))
  } catch (error) {
    errorMessage = getErrorMessage(error)
  }

  if (!browser) {
    if (memo !== undefined) {
      console.log(`rendered memo(${params.memoId}) page`)
    }

    if (errorMessage !== undefined) {
      console.log(`error during memo(${params.memoId}) page rendering: ${errorMessage}`)
    }
  }

  return {
    memo,
    errorMessage,
  }
}) satisfies PageLoad
