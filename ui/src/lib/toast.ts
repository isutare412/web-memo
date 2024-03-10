import { writable } from 'svelte/store'

export type ToastLevel = 'info' | 'success' | 'warning' | 'error'

export enum ToastTimeout {
  XSHORT = 3_000,
  SHORT = 5_000,
  NORMAL = 8_000,
  LONG = 10_000,
  XLONG = 15_000,
}

interface Toast {
  id: number
  message: string
  level: ToastLevel
  timeout?: number
}

interface ToastState {
  nextId: number
  toasts: Toast[]
}

export const toastStore = writable<ToastState>({ nextId: 0, toasts: [] })

export function addToast(message: string, level: ToastLevel, option?: { timeout?: ToastTimeout }) {
  toastStore.update((state) => {
    state.toasts.push({
      id: state.nextId++,
      message,
      level,
      timeout: option?.timeout,
    })
    return state
  })
}

export function deleteToast(id: number) {
  toastStore.update((state) => {
    state.toasts = state.toasts.filter((toast) => toast.id !== id)
    return state
  })
}
