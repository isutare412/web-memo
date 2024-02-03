import { buildErrorMessage, getErrorResponse } from '$lib/apis/backend/error'

export interface RawMemo {
  id: string
  ownerId: string
  createTime: string
  updateTime: string
  title: string
  content: string
  isPublished: boolean
  tags: string[]
}

interface ListMemosResponse {
  page: number
  pageSize: number
  lastPage: number
  totalMemoCount: number
  memos: RawMemo[]
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

interface PublishMemoRequest {
  id: string
  publish: boolean
}

export async function getMemo(id: string): Promise<RawMemo> {
  const response = await fetch(`/api/v1/memos/${id}`)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new Error(buildErrorMessage(response.status, errorResponse))
  }

  return response.json()
}

export async function listMemos(
  page: number,
  pageSize: number,
  tags: string[]
): Promise<ListMemosResponse> {
  const searchParams = new URLSearchParams()
  searchParams.append('page', page.toString())
  searchParams.append('pageSize', pageSize.toString())
  tags.forEach((tag) => {
    searchParams.append('tag', tag)
  })

  const response = await fetch(`/api/v1/memos?${searchParams.toString()}`)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new Error(buildErrorMessage(response.status, errorResponse))
  }

  return response.json()
}

export async function createMemo(request: CreateMemoRequest): Promise<RawMemo> {
  const response = await fetch('/api/v1/memos', {
    method: 'POST',
    body: JSON.stringify(request),
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new Error(buildErrorMessage(response.status, errorResponse))
  }

  return response.json()
}

export async function replaceMemo(request: ReplaceMemoRequest): Promise<RawMemo> {
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
    throw new Error(buildErrorMessage(response.status, errorResponse))
  }

  return response.json()
}

export async function publishMemo(request: PublishMemoRequest): Promise<RawMemo> {
  const response = await fetch(`/api/v1/memos/${request.id}/publish`, {
    method: 'POST',
    body: JSON.stringify({
      publish: request.publish,
    }),
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new Error(buildErrorMessage(response.status, errorResponse))
  }

  return response.json()
}

export async function deleteMemo(id: string): Promise<void> {
  const response = await fetch(`/api/v1/memos/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new Error(buildErrorMessage(response.status, errorResponse))
  }

  return
}

export async function listTags(keyword?: string): Promise<string[]> {
  const apiUrl = keyword ? `/api/v1/tags?kw=${keyword}` : '/api/v1/tags'

  const response = await fetch(apiUrl)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new Error(buildErrorMessage(response.status, errorResponse))
  }

  return response.json()
}
