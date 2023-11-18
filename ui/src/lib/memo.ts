import type { RawMemo } from '$lib/apis/backend/memo'

export const defaultPageSize = 10

export interface Memo {
  id: string
  createTime: Date
  updateTime: Date
  title: string
  content: string
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
    createTime: new Date(memo.createTime),
    updateTime: new Date(memo.updateTime),
    title: memo.title,
    content: memo.content,
    tags: memo.tags,
  } satisfies Memo
}
