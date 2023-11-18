import { sortedUniq } from 'lodash-es'

export function setPageOfSearchParams(params: URLSearchParams, page: number): boolean {
  const currentPage = params.get('p')
  if (page < 1) return false
  if (currentPage === null && page === 1) return false
  if (currentPage !== null && currentPage === page.toString()) return false

  if (page === 1) {
    params.delete('p')
  } else {
    params.set('p', page.toString())
  }

  params.sort()
  return true
}

export function addTagToSearchParams(params: URLSearchParams, tag: string): boolean {
  if (params.getAll('tag').includes(tag)) return false

  params.append('tag', tag)

  let tags = params.getAll('tag')
  tags = sortedUniq(tags.toSorted())

  params.delete('tag')
  tags.forEach((tag) => params.append('tag', tag))

  params.sort()
  return true
}

export function deleteTagFromSearchParams(params: URLSearchParams, tag: string): boolean {
  if (!params.getAll('tag').includes(tag)) return false

  let tags = params.getAll('tag')
  tags = sortedUniq(tags.toSorted())
  tags = tags.filter((t) => t !== tag)

  params.delete('tag')
  tags.forEach((tag) => params.append('tag', tag))

  params.sort()
  return true
}
