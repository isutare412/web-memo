import { getMemo } from '$lib/apis/backend/memo'
import { mapToMemo, type Memo } from '$lib/memo'
import type { PageServerLoad } from './$types'

export const ssr = true

export const load = (async (event) => {
  const { params, url } = event

  let memo: Memo | undefined
  try {
    const base = `${url.protocol}//${url.host}`
    memo = mapToMemo(await getMemo(params.memoId, { base }))
  } catch (_) {
    // let client to fetch memo by passing undefined
  }

  return {
    memo,
  }
}) satisfies PageServerLoad
