<script lang="ts">
  import MemoListItem from '$components/MemoListItem.svelte'
  import type { Memo } from '$lib/memo'

  export let memos: Memo[]
</script>

<div>
  <ul class="flex flex-col gap-y-4">
    {#each memos as memo (memo.id)}
      <MemoListItem {memo} />
    {/each}
  </ul>
</div>
