<script lang="ts">
  import MemoListItem from '$components/MemoListItem.svelte'
  import type { User } from '$lib/auth'
  import type { Memo } from '$lib/memo'

  export let user: User
  export let memos: Memo[]
  export let showUpdateTime = false
</script>

<div>
  <ul class="flex flex-col gap-y-4">
    {#each memos as memo (memo.id)}
      <MemoListItem {user} {memo} {showUpdateTime} />
    {/each}
  </ul>
</div>
