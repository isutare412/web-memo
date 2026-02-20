import type { RawMemo } from '$lib/apis/backend/memo'

const storageKeyPageSize = 'preferredPageSize'
const storageKeySortOrder = 'preferredSortOrder'

export enum SortOrder {
  CREATE_TIME = 'create',
  UPDATE_TIME = 'update',
}

export const defaultPageSize = 10
export const defaultSortOrder = SortOrder.UPDATE_TIME

export const reservedTags = ['published', 'shared']

export interface MemoViewerContext {
  subscription: { isApproved: boolean } | null
  collaboration: { isApproved: boolean } | null
}

export interface Memo {
  id: string
  ownerId: string
  version: number
  createTime: Date
  updateTime: Date
  title: string
  content: string
  publishState: 'private' | 'shared' | 'published'
  tags: string[]
  scores: { rrf: number; semantic: number; bm25: number } | null
  viewerContext: MemoViewerContext | null
}

export interface MemoListPageData {
  page: number | null
  pageSize: number | null
  lastPage: number | null
  totalMemoCount: number | null
  memos: Memo[]
}

export function mapToMemo(memo: RawMemo): Memo {
  return {
    id: memo.id,
    ownerId: memo.ownerId,
    version: memo.version,
    createTime: new Date(memo.createTime),
    updateTime: new Date(memo.updateTime),
    title: memo.title,
    content: memo.content,
    publishState: memo.publishState,
    tags: memo.tags,
    scores: memo.scores,
    viewerContext: memo.viewerContext ?? null,
  } satisfies Memo
}

export function getPreferredPageSize(): number | null {
  const pageSizeStr = localStorage.getItem(storageKeyPageSize)
  if (pageSizeStr === null) return null

  const pageSize = Number(pageSizeStr)
  if (isNaN(pageSize)) return null

  return pageSize
}

export function setPreferredPageSize(size: number) {
  localStorage.setItem(storageKeyPageSize, size.toString())
}

export function getPreferredSortOrder(): SortOrder | null {
  const order = localStorage.getItem(storageKeySortOrder)
  if (order === null) return null

  return Object.values(SortOrder).find((v) => v.valueOf() === order) ?? null
}

export function setPreferredSortOrder(order: SortOrder) {
  localStorage.setItem(storageKeySortOrder, order.toString())
}

export function toggleCheckboxInMarkdown(content: string, index: number): string {
  const checkboxRegex = /^(\s*[-*+]\s*)\[([ xX])\]/gm
  let currentIndex = 0

  return content.replace(checkboxRegex, (match, prefix, checked) => {
    if (currentIndex === index) {
      currentIndex++
      const newState = checked === ' ' ? 'x' : ' '
      return `${prefix}[${newState}]`
    }
    currentIndex++
    return match
  })
}
