<script lang="ts">
  import BookmarkIcon from '$components/icons/BookmarkIcon.svelte'
  import { createEventDispatcher } from 'svelte'

  export let isActivated: boolean

  const dispatch = createEventDispatcher()

  function onClickSubscribeButton() {
    dispatch('subscribe', { isSubscribe: !isActivated })
  }
</script>

<div class:opacity-70={!isActivated}>
  <button
    on:click={onClickSubscribeButton}
    class="btn btn-sm btn-outline rounded-full"
    class:btn-primary={isActivated}
  >
    <div class="w-[12px]">
      <BookmarkIcon />
    </div>
  </button>
</div>
