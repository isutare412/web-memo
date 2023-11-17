import { listMemos, type RawMemo } from '$lib/apis/backend/memo'
import { map, remove, sortBy } from 'lodash-es'
import { writable } from 'svelte/store'

export interface Memo {
  id: string
  createTime: Date
  updateTime: Date
  title: string
  content: string
  tags: string[]
}

interface PagedMemos {
  currentPage: number
  lastPage: number
  pageSize: number
  totalMemoCount: number
  memos: Memo[]
}

interface MemoState {
  pagedMemos: PagedMemos | null
  selectedTags: string[]
}

export const memoStore = writable<MemoState>({ pagedMemos: null, selectedTags: [] })

export async function fetchPagedMemos(page: number, pageSize: number, tags: string[]) {
  const response = await listMemos(page, pageSize, tags)
  const memos = map(response.memos, mapToMemo)

  memoStore.update((state) => {
    state.pagedMemos = {
      currentPage: response.page,
      lastPage: response.lastPage,
      pageSize: response.pageSize,
      totalMemoCount: response.totalMemoCount,
      memos,
    }
    return state
  })
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

export function insertTagFilter(tag: string) {
  memoStore.update((state) => {
    if (state.selectedTags.includes(tag)) return state

    state.selectedTags.push(tag)
    state.selectedTags = sortBy(state.selectedTags, (tag) => tag.toLowerCase())
    return state
  })
}

export function removeTagFilter(tag: string) {
  memoStore.update((state) => {
    remove(state.selectedTags, (t) => t === tag)
    return state
  })
}

export function clearTagFilter() {
  memoStore.update((state) => {
    state.selectedTags = []
    return state
  })
}
