import { writable } from 'svelte/store'

export type ToastLevel = 'info' | 'success' | 'warning' | 'error'

interface Toast {
  id: number
  message: string
  level: ToastLevel
}

interface ToastState {
  nextId: number
  toasts: Toast[]
}

export const toastStore = writable<ToastState>({ nextId: 0, toasts: [] })

export function addToast(message: string, level: ToastLevel) {
  toastStore.update((state) => {
    state.toasts.push({
      id: state.nextId++,
      message,
      level,
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
