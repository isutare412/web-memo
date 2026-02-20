<script lang="ts">
  import BookmarkIcon from '$components/icons/BookmarkIcon.svelte'
  import CheckIcon from '$components/icons/CheckIcon.svelte'
  import { createEventDispatcher } from 'svelte'

  export let isActivated: boolean
  export let isChecked = false

  const dispatch = createEventDispatcher()

  function onClickSubscribeButton() {
    dispatch('subscribe', { isSubscribe: !isActivated })
  }
</script>

<div class:opacity-70={!isActivated}>
  <button
    on:click={onClickSubscribeButton}
    class="btn btn-outline btn-sm rounded-full"
    class:btn-primary={isActivated}
  >
    <div class="w-[12px]">
      <BookmarkIcon />
    </div>
    {#if isChecked}
      <div class="w-[12px]">
        <CheckIcon />
      </div>
    {/if}
  </button>
</div>
