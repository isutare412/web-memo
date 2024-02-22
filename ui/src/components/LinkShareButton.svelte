<script lang="ts">
  import PeopleIcon from '$components/icons/PeopleIcon.svelte'
  import { createEventDispatcher } from 'svelte'

  export let link: string
  export let shareCount: number | undefined = undefined
  export let isShared: boolean

  const dispatch = createEventDispatcher()

  function onClickShareButton() {
    dispatch('share', { link })
  }
</script>

<div class:opacity-70={!isShared}>
  <button
    on:click={onClickShareButton}
    class="btn btn-sm btn-outline rounded-full"
    class:btn-primary={isShared}
  >
    <div class="w-[16px]">
      <PeopleIcon />
    </div>
    {#if shareCount !== undefined && shareCount > 0}
      {shareCount}
    {/if}
  </button>
</div>
