<script lang="ts">
  import CheckIcon from '$components/icons/CheckIcon.svelte'
  import PeopleIcon from '$components/icons/PeopleIcon.svelte'
  import { createEventDispatcher } from 'svelte'

  export let isActivated: boolean
  export let isChecked = false
  export let disabled = false
  export let count: number | undefined = undefined

  const dispatch = createEventDispatcher()

  function onClickCollaborateButton() {
    dispatch('click')
  }
</script>

<div class:opacity-70={!isActivated}>
  <button
    {disabled}
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
    {#if count !== undefined && count > 0}
      {count}
    {/if}
  </button>
</div>
