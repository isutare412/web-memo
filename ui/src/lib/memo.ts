import type { RawMemo } from '$lib/apis/backend/memo'

export const defaultPageSize = 10

const pageSizeKey = 'preferredPageSize'

export const reservedTags = ['published']

export interface Memo {
  id: string
  ownerId: string
  createTime: Date
  updateTime: Date
  title: string
  content: string
  isPublished: boolean
  tags: string[]
}

export interface MemoListPageData {
  page: number
  pageSize: number
  lastPage: number
  totalMemoCount: number
  memos: Memo[]
}

export function mapToMemo(memo: RawMemo): Memo {
  return {
    id: memo.id,
    ownerId: memo.ownerId,
    createTime: new Date(memo.createTime),
    updateTime: new Date(memo.updateTime),
    title: memo.title,
    content: memo.content,
    isPublished: memo.isPublished,
    tags: memo.tags,
  } satisfies Memo
}

export function getPreferredPageSize(): number | null {
  const pageSizeStr = localStorage.getItem(pageSizeKey)
  if (pageSizeStr === null) return null

  const pageSize = Number(pageSizeStr)
  if (isNaN(pageSize)) return null

  return pageSize
}

export function setPreferredPageSize(size: number) {
  localStorage.setItem(pageSizeKey, size.toString())
}
