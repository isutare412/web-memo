import {
  createUploadURL,
  uploadToS3,
  waitForImageReady,
  type ImageFormat,
} from '$lib/apis/backend/image'

export interface ImageUploadCallbacks {
  onPlaceholder: (imageId: string, placeholder: string) => void
  onComplete: (imageId: string, displayUrl: string, originalUrl: string) => void
  onError: (imageId: string, error: Error) => void
}

function getImageFormatFromFile(file: File): ImageFormat | null {
  const mimeToFormat: Record<string, ImageFormat> = {
    'image/jpeg': 'JPEG',
    'image/png': 'PNG',
    'image/webp': 'WEBP',
    'image/avif': 'AVIF',
    'image/heic': 'HEIC',
  }
  return mimeToFormat[file.type] || null
}

function getImageFormatFromDataURL(dataURL: string): ImageFormat | null {
  const mimeMatch = dataURL.match(/^data:(image\/[^;]+);/)
  if (!mimeMatch) return null

  const mimeToFormat: Record<string, ImageFormat> = {
    'image/jpeg': 'JPEG',
    'image/png': 'PNG',
    'image/webp': 'WEBP',
    'image/avif': 'AVIF',
    'image/heic': 'HEIC',
  }
  return mimeToFormat[mimeMatch[1]] || null
}

async function dataURLToBlob(dataURL: string): Promise<Blob> {
  const response = await fetch(dataURL)
  return response.blob()
}

export async function uploadImage(file: File, callbacks: ImageUploadCallbacks): Promise<void> {
  const format = getImageFormatFromFile(file)
  if (!format) {
    callbacks.onError('', new Error(`Unsupported image format: ${file.type}`))
    return
  }

  try {
    // Step 1: Get presigned URL
    const uploadUrlResponse = await createUploadURL({
      fileName: file.name,
      format,
    })

    const { imageId, uploadUrl, uploadHeaders } = uploadUrlResponse

    // Step 2: Insert placeholder
    const placeholder = `![Uploading image...](uploading:${imageId})`
    callbacks.onPlaceholder(imageId, placeholder)

    // Step 3: Upload to S3
    await uploadToS3(uploadUrl, uploadHeaders, file)

    // Step 4: Wait for processing to complete
    const imageStatus = await waitForImageReady(imageId)

    if (imageStatus.state === 'READY') {
      const displayUrl = imageStatus.downscaled?.url ?? imageStatus.original?.url ?? ''
      const originalUrl = imageStatus.original?.url ?? displayUrl
      callbacks.onComplete(imageId, displayUrl, originalUrl)
    } else {
      callbacks.onError(
        imageId,
        new Error(`Image processing failed with state: ${imageStatus.state}`)
      )
    }
  } catch (error) {
    callbacks.onError('', error instanceof Error ? error : new Error(String(error)))
  }
}

export async function uploadImageFromDataURL(
  dataURL: string,
  fileName: string,
  callbacks: ImageUploadCallbacks
): Promise<void> {
  const format = getImageFormatFromDataURL(dataURL)
  if (!format) {
    callbacks.onError('', new Error('Unsupported image format from clipboard'))
    return
  }

  try {
    const blob = await dataURLToBlob(dataURL)

    // Step 1: Get presigned URL
    const uploadUrlResponse = await createUploadURL({
      fileName,
      format,
    })

    const { imageId, uploadUrl, uploadHeaders } = uploadUrlResponse

    // Step 2: Insert placeholder
    const placeholder = `![Uploading image...](uploading:${imageId})`
    callbacks.onPlaceholder(imageId, placeholder)

    // Step 3: Upload to S3
    await uploadToS3(uploadUrl, uploadHeaders, blob)

    // Step 4: Wait for processing to complete
    const imageStatus = await waitForImageReady(imageId)

    if (imageStatus.state === 'READY') {
      const displayUrl = imageStatus.downscaled?.url ?? imageStatus.original?.url ?? ''
      const originalUrl = imageStatus.original?.url ?? displayUrl
      callbacks.onComplete(imageId, displayUrl, originalUrl)
    } else {
      callbacks.onError(
        imageId,
        new Error(`Image processing failed with state: ${imageStatus.state}`)
      )
    }
  } catch (error) {
    callbacks.onError('', error instanceof Error ? error : new Error(String(error)))
  }
}

export function extractImageFromClipboard(clipboardData: DataTransfer): {
  file: File | null
  dataURL: string | null
} {
  // Check for file items first
  for (const item of clipboardData.items) {
    if (item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (file) {
        return { file, dataURL: null }
      }
    }
  }

  return { file: null, dataURL: null }
}

export function extractImageFromDrop(dataTransfer: DataTransfer): File | null {
  for (const file of dataTransfer.files) {
    if (file.type.startsWith('image/')) {
      return file
    }
  }
  return null
}
