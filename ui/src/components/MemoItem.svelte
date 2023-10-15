<script lang="ts">
  import Tag from '$components/Tag.svelte'
  import type { Memo } from '$lib/memo'
  import { formatDate } from '$lib/utils/date'

  export let memo: Memo
</script>

<li>
  <!-- TODO: change href -->
  <a href="/" class="link link-hover text inline-block max-w-full truncate">
    {memo.title}
  </a>
  <div class="flex gap-x-2">
    <div class="flex flex-1 flex-wrap gap-1">
      {#each memo.tags as tag (tag)}
        <Tag value={tag} />
      {/each}
    </div>
    <span class="mt-[2px] flex-none text-xs font-light">{formatDate(memo.createTime)}</span>
  </div>
</li>
