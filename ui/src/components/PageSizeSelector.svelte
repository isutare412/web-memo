<script lang="ts">
  import { sortedUniq } from 'lodash-es'
  import { createEventDispatcher } from 'svelte'

  export let currentSize: number
  let selectValue: number

  const dispatch = createEventDispatcher()

  const defaultOptions = [5, 6, 8, 10, 20, 50]

  $: options = sortedUniq([...defaultOptions, currentSize].toSorted((a, b) => a - b)).map(
    (size) => ({ size, selected: size === currentSize })
  )

  function onChangeSelect() {
    dispatch('change', { pageSize: selectValue })
  }
</script>

<div class="flex items-center gap-2">
  <label for="paze-sizer" class="flex-none text-xs font-light opacity-75">Items per page</label>
  <select
    bind:value={selectValue}
    on:change={onChangeSelect}
    id="page-sizer"
    class="select select-bordered select-sm w-full max-w-xs focus:outline-none"
  >
    {#each options as opt}
      <option value={opt.size} selected={opt.selected}>{opt.size}</option>
    {/each}
  </select>
</div>
