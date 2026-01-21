import { StatusError, getErrorResponse } from '$lib/apis/backend/error'

export type ImageFormat = 'JPEG' | 'PNG' | 'WEBP' | 'AVIF' | 'HEIC'

export type ImageState = 'UPLOAD_PENDING' | 'UPLOAD_EXPIRED' | 'FAILED' | 'READY'

export interface CreateUploadURLRequest {
  fileName: string
  format: ImageFormat
}

export interface CreateUploadURLResponse {
  imageId: string
  uploadUrl: string
  uploadHeaders: Record<string, string>
  expiresAt: string
}

export interface ImageData {
  url: string
  format: ImageFormat
}

export interface ImageStatusResponse {
  id: string
  state: ImageState
  original: ImageData | null
  downscaled: ImageData | null
}

export async function createUploadURL(
  request: CreateUploadURLRequest
): Promise<CreateUploadURLResponse> {
  const response = await fetch('/api/v1/images/upload-url', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  })
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return response.json()
}

export async function uploadToS3(
  url: string,
  headers: Record<string, string>,
  file: Blob
): Promise<void> {
  const response = await fetch(url, {
    method: 'PUT',
    headers,
    body: file,
  })
  if (!response.ok) {
    throw new Error(`S3 upload failed with status ${response.status}`)
  }
}

export async function waitForImageReady(imageId: string): Promise<ImageStatusResponse> {
  const response = await fetch(`/api/v1/images/${imageId}/status?waitUntilProcessed=true`)
  if (!response.ok) {
    const errorResponse = await getErrorResponse(response)
    throw new StatusError(response.status, errorResponse.msg)
  }

  return response.json()
}
