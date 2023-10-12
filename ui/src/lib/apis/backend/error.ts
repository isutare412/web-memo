export interface ErrorResponse {
  msg?: string
}

export async function getErrorResponse(response: Response): Promise<ErrorResponse> {
  return response.json()
}
