import type { RawMemo } from '$lib/apis/backend/memo'
import { remove, sortBy } from 'lodash-es'
import { get, writable } from 'svelte/store'

export const defaultPageSize = 10

export const depKeys = {
  memoList: 'memo:list',
}

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

interface MemoState {
  currentPage: number
  pageSize: number
  selectedTags: string[]
  informUpdate: () => void
}

export const memoStore = writable<MemoState>({
  currentPage: 1,
  pageSize: defaultPageSize,
  selectedTags: [],
  informUpdate: () => {},
})

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

export function updateCurrentPage(page: number) {
  memoStore.update((state) => {
    state.currentPage = page
    return state
  })
}

export function updatePageSize(size: number) {
  memoStore.update((state) => {
    state.pageSize = size
    return state
  })
}

export function insertTagFilter(tag: string) {
  memoStore.update((state) => {
    if (state.selectedTags.includes(tag)) return state

    state.selectedTags.push(tag)
    state.selectedTags = sortBy(state.selectedTags, (tag) => tag.toLowerCase())
    state.currentPage = 1
    return state
  })
}

export function removeTagFilter(tag: string) {
  memoStore.update((state) => {
    remove(state.selectedTags, (t) => t === tag)
    state.currentPage = 1
    return state
  })
}

export function informUpdate() {
  get(memoStore).informUpdate()
}

export function setUpdateInformer(fn: () => void) {
  memoStore.update((state) => {
    state.informUpdate = fn
    return state
  })
}

export function clearUpdateInformer() {
  memoStore.update((state) => {
    state.informUpdate = () => {}
    return state
  })
}
