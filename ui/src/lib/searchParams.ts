import { SortOrder, defaultPageSize, defaultSortOrder } from '$lib/memo'
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

export function setPageSizeOfSearchParams(params: URLSearchParams, size: number): boolean {
  const currentSize = params.get('ps')
  if (size < 1) return false
  if (currentSize === null && size === defaultPageSize) return false
  if (currentSize !== null && currentSize === size.toString()) return false

  if (size === defaultPageSize) {
    params.delete('ps')
  } else {
    params.set('ps', size.toString())
  }

  params.sort()
  return true
}

export function setSortOrderOfSearchParams(params: URLSearchParams, order: SortOrder): boolean {
  const currentSortOrder = params.get('sort')
  if (currentSortOrder === null && order === defaultSortOrder) return false
  if (currentSortOrder !== null && order.toString() === currentSortOrder) return false

  if (order === defaultSortOrder) {
    params.delete('sort')
  } else {
    params.set('sort', order.toString())
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
