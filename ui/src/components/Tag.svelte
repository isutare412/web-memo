<script lang="ts">
  import { createEventDispatcher } from 'svelte'

  type Color = 'primary' | 'secondary'

  export let value: string
  export let color: Color = 'primary'
  export let outline = false
  export let isClose: boolean = false
  export let isButton = true

  const dispatch = createEventDispatcher()

  function dispatchTag() {
    dispatch('click', {
      name: value,
    })
  }
</script>

<button
  on:click={dispatchTag}
  disabled={!isButton}
  class="badge badge-primary max-w-[200px] text-sm transition-none"
  class:badge-primary={color === 'primary'}
  class:badge-secondary={color === 'secondary'}
  class:badge-outline={outline}
>
  <span class="truncate">
    {value}
  </span>
  {#if isClose}
    <span>&nbsp;✕</span>
  {/if}
</button>
