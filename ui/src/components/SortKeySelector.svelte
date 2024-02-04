<script lang="ts">
  import { SortOrder } from '$lib/memo'
  import { createEventDispatcher } from 'svelte'

  export let sortKey: SortOrder
  let selectedKey: SortOrder

  const dispatch = createEventDispatcher()

  $: options = [SortOrder.CREATE_TIME, SortOrder.UPDATE_TIME].map((key) => ({
    key,
    selected: key === sortKey,
  }))

  const displayValues = {
    [SortOrder.CREATE_TIME]: 'Create Time',
    [SortOrder.UPDATE_TIME]: 'Update Time',
  }

  function onSelectChange() {
    dispatch('change', { sortKey: selectedKey })
  }
</script>

<div class="flex items-center">
  <select
    bind:value={selectedKey}
    on:change={onSelectChange}
    class="select select-sm select-bordered rounded-none border-0 border-b focus:outline-none"
  >
    {#each options as opt}
      <option value={opt.key} selected={opt.selected}>{displayValues[opt.key]}</option>
    {/each}
  </select>
</div>
