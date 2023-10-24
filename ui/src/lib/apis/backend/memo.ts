import { getErrorResponse } from '$lib/apis/backend/error'
import { error } from '@sveltejs/kit'

interface Memo {
  id: string
  createTime: string
  updateTime: string
  title: string
  content: string
  tags: string[]
}

export async function getMemo(id: string): Promise<Memo> {
  const response = await fetch(`/api/v1/memos/${id}`)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw error(response.status, errorResponse.msg)
  }

  return response.json()
}

export async function listMemos(): Promise<Memo[]> {
  const response = await fetch('/api/v1/memos')
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw error(response.status, errorResponse.msg)
  }

  return response.json()
}
