<script lang="ts">
  import CheckIcon from '$components/icons/CheckIcon.svelte'
  import PeopleIcon from '$components/icons/PeopleIcon.svelte'
  import { createEventDispatcher } from 'svelte'

  export let isActivated: boolean
  export let isChecked: boolean

  const dispatch = createEventDispatcher()

  function onClickCollaborateButton() {
    dispatch('collaborate', { isSubscribe: !isActivated })
  }
</script>

<div class:opacity-70={!isActivated}>
  <button
    on:click={onClickCollaborateButton}
    class="btn btn-sm btn-outline rounded-full"
    class:btn-primary={isActivated}
  >
    <div class="w-[16px]">
      <PeopleIcon />
    </div>
    {#if isChecked}
      <div class="w-[12px]">
        <CheckIcon />
      </div>
    {/if}
  </button>
</div>
