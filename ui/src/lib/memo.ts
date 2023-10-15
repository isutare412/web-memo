import { listMemos } from '$lib/apis/backend/memo'
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

interface MemoState {
  memos: Memo[]
  selectedTags: string[]
}

export const memoStore = writable<MemoState>({ memos: [], selectedTags: [] })

export async function syncMemos() {
  const response = await listMemos()

  const memos = map(response, (memo) => {
    return {
      id: memo.id,
      createTime: new Date(memo.createTime),
      updateTime: new Date(memo.updateTime),
      title: memo.title,
      content: memo.content,
      tags: memo.tags,
    } satisfies Memo
  })

  memoStore.update((state) => {
    state.memos = memos
    return state
  })
}

export function insertTagFilter(tag: string) {
  memoStore.update((state) => {
    if (state.selectedTags.includes(tag)) return state

    state.selectedTags.push(tag)
    state.selectedTags = sortBy(state.selectedTags)
    return state
  })
}

export function removeTagFilter(tag: string) {
  memoStore.update((state) => {
    remove(state.selectedTags, (t) => t === tag)
    return state
  })
}
