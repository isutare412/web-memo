<script lang="ts">
  import Cross from '$components/icons/Cross.svelte'
  import { ToastTimeout, deleteToast, type ToastLevel } from '$lib/toast'
  import { PausableTimer } from '$lib/utils/timer'
  import { onMount } from 'svelte'
  import { tweened } from 'svelte/motion'

  export let id: number
  export let message: string
  export let level: ToastLevel
  export let timeout: number = ToastTimeout.NORMAL

  const progress = tweened(0, { duration: timeout })
  let timer: PausableTimer
  let lastStop = 0

  onMount(() => {
    progress.set(1)

    timer = new PausableTimer(() => {
      deleteToast(id)
    }, timeout)
  })

  function pauseProgress() {
    lastStop = $progress
    progress.set(lastStop, { duration: 0 })
    timer.pause()
  }

  function resumeProgress() {
    progress.set(1, { duration: (1 - lastStop) * timeout })
    timer.resume()
  }
</script>

<div
  role="alertdialog"
  on:mouseenter={pauseProgress}
  on:mouseleave={resumeProgress}
  class="alert relative block overflow-hidden text-left"
  class:alert-info={level === 'info'}
  class:alert-success={level === 'success'}
  class:alert-warning={level === 'warning'}
  class:alert-error={level === 'error'}
>
  <button
    on:click={() => deleteToast(id)}
    class="btn-ghost btn-sm btn-circle btn absolute right-2 top-3"
  >
    <div class="h-4 w-4">
      <Cross />
    </div>
  </button>
  <div class="whitespace-pre-line break-words pr-6">{message}</div>
  <div class="absolute bottom-0 left-0 w-full bg-transparent">
    <div
      class="h-1 rounded-full"
      class:bg-info-content={level === 'info'}
      class:bg-success-content={level === 'success'}
      class:bg-warning-content={level === 'warning'}
      class:bg-error-content={level === 'error'}
      style={`width: ${(1 - $progress) * 100}%`}
    />
  </div>
</div>
