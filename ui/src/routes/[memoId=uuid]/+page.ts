import type { PageLoad } from './$types'

export const load = (async (event) => {
  const { parent } = event
  await parent()
}) satisfies PageLoad
