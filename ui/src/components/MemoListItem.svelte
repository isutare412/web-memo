<script lang="ts">
  import Tag from '$components/Tag.svelte'
  import { informUpdate, insertTagFilter, memoStore, type Memo } from '$lib/memo'
  import { formatDate } from '$lib/utils/date'
  import { map } from 'lodash-es'

  export let memo: Memo

  $: tags = map(memo.tags, (tag) => ({
    name: tag,
    filtered: $memoStore.selectedTags.includes(tag),
  }))

  function selectTag(event: CustomEvent<{ name: string }>) {
    insertTagFilter(event.detail.name)
    informUpdate()
  }
</script>

<li class="flex flex-col gap-y-1">
  <div class="flex items-center">
    <a href={`/${memo.id}`} class="link link-hover inline-block max-w-full truncate text-lg">
      {memo.title}
    </a>
  </div>
  {#if tags.length > 0}
    <div class="flex flex-wrap gap-1">
      {#each tags as tag (tag.name)}
        <Tag value={tag.name} outline={!tag.filtered} on:click={selectTag} />
      {/each}
    </div>
  {/if}
  <div class="flex justify-end">
    <span class="text-xs font-light opacity-75">{formatDate(memo.createTime)}</span>
  </div>
</li>
