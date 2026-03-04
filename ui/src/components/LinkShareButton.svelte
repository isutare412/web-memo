<script lang="ts">
  import LockIcon from '$components/icons/LockIcon.svelte'
  import LinkIcon from '$components/icons/LinkIcon.svelte'
  import WebPublishIcon from '$components/icons/WebPublishIcon.svelte'
  import { createEventDispatcher } from 'svelte'

  export let shareCount: number | undefined = undefined
  export let publishState: 'private' | 'shared' | 'published'
  export let isPublishing: boolean = false

  const dispatch = createEventDispatcher()

  const states: {
    value: 'private' | 'shared' | 'published'
    label: string
    description: string
  }[] = [
    { value: 'private', label: 'Private', description: 'Only you can view' },
    { value: 'shared', label: 'Shared', description: 'Approved users can view' },
    { value: 'published', label: 'Published', description: 'Anyone with the link can view' },
  ]

  let dialog: HTMLDialogElement
  let selectedState: 'private' | 'shared' | 'published' = publishState
  $: selectedState = publishState

  function openDialog() {
    selectedState = publishState
    dialog.showModal()
  }

  function onApply() {
    if (selectedState !== publishState) {
      dispatch('share', { publishState: selectedState })
    }
    dialog.close()
  }

  function onClose() {
    selectedState = publishState
  }
</script>

<button
  class="btn btn-outline btn-sm rounded-full"
  class:btn-primary={publishState !== 'private'}
  class:opacity-70={publishState === 'private'}
  on:click={openDialog}
>
  <div class="w-[16px]">
    <WebPublishIcon />
  </div>
  {#if publishState !== 'private' && shareCount !== undefined && shareCount > 0}
    {shareCount}
  {/if}
</button>

<dialog bind:this={dialog} class="modal" on:close={onClose}>
  <div class="modal-box w-fit min-w-72">
    <h3 class="mb-3 text-sm font-semibold">Visibility</h3>
    <div class="flex flex-col gap-y-1">
      {#each states as state}
        <button
          class="flex cursor-pointer items-center gap-x-3 rounded-lg p-2 text-left"
          class:bg-base-200={selectedState === state.value}
          on:click={() => (selectedState = state.value)}
        >
          <input
            type="radio"
            name="publish-state"
            class="radio-primary radio radio-sm"
            checked={selectedState === state.value}
            on:change={() => (selectedState = state.value)}
          />
          <div class="w-[16px] shrink-0">
            {#if state.value === 'private'}
              <LockIcon />
            {:else if state.value === 'shared'}
              <LinkIcon />
            {:else}
              <WebPublishIcon />
            {/if}
          </div>
          <div>
            <div class="text-sm font-medium">{state.label}</div>
            <div class="text-xs opacity-70">{state.description}</div>
          </div>
        </button>
      {/each}
    </div>
    <div class="modal-action flex justify-end">
      <form method="dialog">
        <button class="btn btn-outline btn-primary btn-sm outline-none">Cancel</button>
      </form>
      <button
        class="btn btn-primary btn-sm outline-none"
        disabled={selectedState === publishState || isPublishing}
        on:click={onApply}
      >
        {#if isPublishing}
          <span class="loading loading-spinner loading-sm" />
        {:else}
          Apply
        {/if}
      </button>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
