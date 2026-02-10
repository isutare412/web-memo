import { StatusError, getErrorResponse } from '$lib/apis/backend/error'
import { SortOrder } from '$lib/memo'

export interface RawMemo {
  id: string
  ownerId: string
  version: number
  createTime: string
  updateTime: string
  title: string
  content: string
  isPublished: boolean
  tags: string[]
  scores: { rrf: number; semantic: number; bm25: number } | null
}

export interface Subscriber {
  id: string
}

export interface Collaborator {
  id: string
  userName: string
  photoUrl: string
  isApproved: boolean
}

interface ListMemosResponse {
  page: number | null
  pageSize: number | null
  lastPage: number | null
  totalMemoCount: number | null
  memos: RawMemo[]
}

interface CreateMemoRequest {
  title: string
  content: string
  tags: string[]
}

interface ReplaceMemoRequest {
  id: string
  version: number
  title: string
  content: string
  tags: string[]
  isPinUpdateTime?: boolean
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

interface GetCollaboratorRequest {
  memoId: string
  userId: string
}

interface ListCollaboratorsResponse {
  collaborators: Collaborator[]
}

interface RequestCollaborationRequest {
  memoId: string
  userId: string
}

interface AuthorizeCollaborationRequest {
  memoId: string
  userId: string
  approve: boolean
}

interface CancelCollaborationRequest {
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

export async function searchMemos(query: string): Promise<ListMemosResponse> {
  const searchParams = new URLSearchParams()
  searchParams.append('q', query)

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
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return response.json()
}

export async function replaceMemo(request: ReplaceMemoRequest): Promise<RawMemo> {
  const searchParams = new URLSearchParams()
  searchParams.append('pinUpdateTime', request.isPinUpdateTime ? 'true' : 'false')

  const response = await fetch(`/api/v1/memos/${request.id}?${searchParams.toString()}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      title: request.title,
      content: request.content,
      tags: request.tags,
      version: request.version,
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
    headers: { 'Content-Type': 'application/json' },
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

export async function getCollaborator({
  memoId,
  userId,
}: GetCollaboratorRequest): Promise<Collaborator> {
  const response = await fetch(`/api/v1/memos/${memoId}/collaborators/${userId}`)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return response.json()
}

export async function listCollaborators(memoId: string): Promise<ListCollaboratorsResponse> {
  const response = await fetch(`/api/v1/memos/${memoId}/collaborators`)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return response.json()
}

export async function requestCollaboration({
  memoId,
  userId,
}: RequestCollaborationRequest): Promise<void> {
  const response = await fetch(`/api/v1/memos/${memoId}/collaborators/${userId}`, { method: 'PUT' })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return
}

export async function authorizeCollaboration({
  memoId,
  userId,
  approve,
}: AuthorizeCollaborationRequest): Promise<void> {
  const response = await fetch(`/api/v1/memos/${memoId}/collaborators/${userId}/authorize`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ approve }),
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return
}

export async function cancelCollaboration({
  memoId,
  userId,
}: CancelCollaborationRequest): Promise<void> {
  const response = await fetch(`/api/v1/memos/${memoId}/collaborators/${userId}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return
}
