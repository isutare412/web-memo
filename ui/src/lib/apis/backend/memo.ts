import { buildErrorMessage, getErrorResponse } from '$lib/apis/backend/error'
import { SortOrder } from '$lib/memo'

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

export async function getMemo(id: string, option?: { fetch?: typeof fetch }): Promise<RawMemo> {
  const customFetch = option?.fetch ?? fetch

  const response = await customFetch(`/api/v1/memos/${id}`)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new Error(buildErrorMessage(response.status, errorResponse))
  }

  return response.json()
}

export async function listMemos(
  page: number,
  pageSize: number,
  sortOrder: SortOrder,
  tags: string[]
): Promise<ListMemosResponse> {
  let sort: string
  switch (sortOrder) {
    case SortOrder.CREATE_TIME:
      sort = 'createTime'
      break
    case SortOrder.UPDATE_TIME:
      sort = 'updateTime'
      break
  }

  const searchParams = new URLSearchParams()
  searchParams.append('page', page.toString())
  searchParams.append('pageSize', pageSize.toString())
  searchParams.append('sort', sort)
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
