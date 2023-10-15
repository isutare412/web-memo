import { listMemos } from '$lib/apis/backend/memo'
import { map } from 'lodash-es'
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
}

export const memoStore = writable<MemoState>({ memos: [] })

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
