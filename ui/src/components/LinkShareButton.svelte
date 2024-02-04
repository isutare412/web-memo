<script lang="ts">
  import PeopleIcon from '$components/icons/PeopleIcon.svelte'
  import Share from '$components/icons/Share.svelte'
  import { createEventDispatcher } from 'svelte'

  export let link: string
  export let isShared: boolean

  const dispatch = createEventDispatcher()

  function onClickShareButton() {
    dispatch('share', { link })
  }
</script>

<div class="opacity-70">
  <button on:click={onClickShareButton} class="btn btn-sm btn-outline rounded-full">
    {#if isShared}
      <div class="w-[16px]">
        <PeopleIcon />
      </div>
    {:else}
      <div class="w-[16px]">
        <Share />
      </div>
    {/if}
  </button>
</div>
