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

interface CreateMemoRequest {
  title: string
  content: string
  tags: string[]
}

interface ReplaceMemoRequest {
  id: string
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

export async function createMemo(request: CreateMemoRequest): Promise<Memo> {
  const response = await fetch('/api/v1/memos', {
    method: 'POST',
    body: JSON.stringify(request),
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw error(response.status, errorResponse.msg)
  }

  return response.json()
}

export async function replaceMemo(request: ReplaceMemoRequest): Promise<Memo> {
  const response = await fetch(`/api/v1/memos/${request.id}`, {
    method: 'PUT',
    body: JSON.stringify({
      title: request.title,
      content: request.content,
      tags: request.tags,
    }),
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw error(response.status, errorResponse.msg)
  }

  return response.json()
}

export async function deleteMemo(id: string): Promise<void> {
  const response = await fetch(`/api/v1/memos/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw error(response.status, errorResponse.msg)
  }

  return
}
