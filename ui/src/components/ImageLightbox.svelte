<script lang="ts">
  import { createEventDispatcher } from 'svelte'

  export let src: string
  export let alt: string = ''

  const dispatch = createEventDispatcher()

  function close() {
    dispatch('close')
  }

  function onKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      close()
    }
  }

  function onBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      close()
    }
  }
</script>

<svelte:window on:keydown={onKeydown} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div
  class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4"
  on:click={onBackdropClick}
  role="presentation"
>
  <button
    type="button"
    class="absolute right-4 top-4 text-white hover:text-gray-300"
    on:click={close}
    aria-label="Close"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-8 w-8"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M6 18L18 6M6 6l12 12"
      />
    </svg>
  </button>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions a11y-click-events-have-key-events -->
  <img
    {src}
    {alt}
    class="max-h-[90vh] max-w-[90vw] object-contain"
    on:click|stopPropagation={() => {}}
  />
</div>
