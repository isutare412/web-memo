import { writable } from 'svelte/store'

function createLoadingStore() {
  const { subscribe, set } = writable(false)
  return {
    subscribe,
    start: () => set(true),
    stop: () => set(false),
  }
}

export const loading = createLoadingStore()
