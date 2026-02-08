<script lang="ts">
  import AngleRight from '$components/icons/AngleRight.svelte'
  import AngleRightDouble from '$components/icons/AngleRightDouble.svelte'
  import { createEventDispatcher } from 'svelte'

  export let currentPage: string
  export let lastPage: string

  $: enablePrev = Number(currentPage) > 1
  $: enableNext = Number(currentPage) < Number(lastPage)
  $: pageInput = currentPage

  const dispatch = createEventDispatcher()

  function onClickFirstButton() {
    dispatchNavigateEvent(1)
  }

  function onClickLastButton() {
    dispatchNavigateEvent(Number(lastPage))
  }

  function onClickPrevButton() {
    const page = Number(currentPage) - 1
    if (isNaN(page) || page < 1 || page > Number(lastPage)) return

    dispatchNavigateEvent(page)
  }

  function onClickNextButton() {
    const page = Number(currentPage) + 1
    if (isNaN(page) || page < 1 || page > Number(lastPage)) return

    dispatchNavigateEvent(page)
  }

  function onPageInputKeyUp(
    event: KeyboardEvent & { currentTarget: EventTarget & HTMLInputElement }
  ) {
    if (event.key !== 'Enter') return

    const page = Number(pageInput)
    if (isNaN(page) || page < 1 || page > Number(lastPage)) return

    dispatchNavigateEvent(page)
  }

  function onPageInputFocusOut() {
    const page = Number(pageInput)
    if (isNaN(page) || page < 1 || page > Number(lastPage)) return

    dispatchNavigateEvent(page)
  }

  function dispatchNavigateEvent(page: number) {
    dispatch('navigate', { page })
  }
</script>

<div class="flex gap-x-1">
  <button
    on:click={onClickFirstButton}
    disabled={!enablePrev}
    class="btn btn-square btn-ghost btn-sm disabled:bg-transparent"
  >
    <div class="w-[30px] rotate-180"><AngleRightDouble /></div>
  </button>
  <button
    on:click={onClickPrevButton}
    disabled={!enablePrev}
    class="btn btn-square btn-ghost btn-sm disabled:bg-transparent"
  >
    <div class="w-[30px] rotate-180"><AngleRight /></div>
  </button>
  <div class="flex items-center gap-x-2">
    <input
      type="text"
      inputmode="numeric"
      bind:value={pageInput}
      on:keyup={onPageInputKeyUp}
      on:focusout={onPageInputFocusOut}
      class="input input-sm input-bordered w-full max-w-[44px] px-2 text-center text-base focus:border-primary focus:outline-none"
    />
    <span class="text-sm font-light opacity-75">of {lastPage}</span>
  </div>
  <button
    on:click={onClickNextButton}
    disabled={!enableNext}
    class="btn btn-square btn-ghost btn-sm disabled:bg-transparent"
  >
    <div class="w-[30px]"><AngleRight /></div>
  </button>
  <button
    on:click={onClickLastButton}
    disabled={!enableNext}
    class="btn btn-square btn-ghost btn-sm disabled:bg-transparent"
  >
    <div class="w-[30px]"><AngleRightDouble /></div>
  </button>
</div>
