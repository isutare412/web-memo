<script lang="ts">
  import BookmarkIcon from '$components/icons/BookmarkIcon.svelte'
  import { createEventDispatcher } from 'svelte'

  export let isActivated: boolean
  export let count: number

  const dispatch = createEventDispatcher()
</script>

<div class:opacity-70={!isActivated}>
  <button
    on:click={() => dispatch('click')}
    class="btn btn-outline btn-sm rounded-full"
    class:btn-primary={isActivated}
  >
    <div class="w-[12px]">
      <BookmarkIcon />
    </div>
    {count}
  </button>
</div>
