import { StatusError, getErrorResponse } from '$lib/apis/backend/error'
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

export interface Subscriber {
  id: string
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

interface GetSubscriberRequest {
  memoId: string
  userId: string
}

interface ListSubsribersResponse {
  subscribers: Subscriber[]
}

interface SubscribeMemoRequest {
  memoId: string
  userId: string
}

interface UnsubscriberMemoRequest {
  memoId: string
  userId: string
}

export async function getMemo(id: string, option?: { fetch?: typeof fetch }): Promise<RawMemo> {
  const customFetch = option?.fetch ?? fetch

  const response = await customFetch(`/api/v1/memos/${id}`)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
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
    throw new StatusError(response.status, errorResponse.msg)
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
    throw new StatusError(response.status, errorResponse.msg)
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
    throw new StatusError(response.status, errorResponse.msg)
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
    throw new StatusError(response.status, errorResponse.msg)
  }

  return response.json()
}

export async function getSubscriber({ memoId, userId }: GetSubscriberRequest): Promise<Subscriber> {
  const response = await fetch(`/api/v1/memos/${memoId}/subscribers/${userId}`)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return response.json()
}

export async function listSubscribers(memoId: string): Promise<ListSubsribersResponse> {
  const response = await fetch(`/api/v1/memos/${memoId}/subscribers`)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return response.json()
}

export async function subscribeMemo({ memoId, userId }: SubscribeMemoRequest): Promise<void> {
  const response = await fetch(`/api/v1/memos/${memoId}/subscribers/${userId}`, {
    method: 'PUT',
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return
}

export async function unsubscribeMemo({ memoId, userId }: UnsubscriberMemoRequest): Promise<void> {
  const response = await fetch(`/api/v1/memos/${memoId}/subscribers/${userId}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return
}

export async function deleteMemo(id: string): Promise<void> {
  const response = await fetch(`/api/v1/memos/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return
}

export async function listTags(keyword?: string): Promise<string[]> {
  const apiUrl = keyword ? `/api/v1/tags?kw=${keyword}` : '/api/v1/tags'

  const response = await fetch(apiUrl)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return response.json()
}
