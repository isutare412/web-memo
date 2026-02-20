<script lang="ts">
  import WebPublishIcon from '$components/icons/WebPublishIcon.svelte'
  import { createEventDispatcher } from 'svelte'

  export let link: string
  export let shareCount: number | undefined = undefined
  export let publishState: 'private' | 'shared' | 'published'

  const dispatch = createEventDispatcher()

  const states: { value: 'private' | 'shared' | 'published'; label: string }[] = [
    { value: 'private', label: 'Private' },
    { value: 'shared', label: 'Shared' },
    { value: 'published', label: 'Published' },
  ]

  function select(newState: 'private' | 'shared' | 'published') {
    if (newState !== publishState) {
      dispatch('share', { publishState: newState })
    }
    if (document.activeElement instanceof HTMLElement) {
      document.activeElement.blur()
    }
  }
</script>

<div class="dropdown dropdown-end">
  <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
  <!-- svelte-ignore a11y-label-has-associated-control -->
  <label
    tabindex="0"
    class="btn btn-outline btn-sm rounded-full"
    class:btn-primary={publishState !== 'private'}
  >
    <div class="w-[16px]">
      <WebPublishIcon />
    </div>
    {#if publishState !== 'private' && shareCount !== undefined && shareCount > 0}
      {shareCount}
    {/if}
  </label>
  <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
  <ul
    tabindex="0"
    class="menu dropdown-content z-[1] rounded-box border border-base-300 bg-base-100 p-2 shadow-lg"
  >
    {#each states as state}
      <li>
        <button class:active={publishState === state.value} on:click={() => select(state.value)}>
          {state.label}
        </button>
      </li>
    {/each}
  </ul>
</div>
