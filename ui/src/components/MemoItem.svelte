<script lang="ts">
  import Tag from '$components/Tag.svelte'
  import { insertTagFilter, memoStore, type Memo } from '$lib/memo'
  import { formatDate } from '$lib/utils/date'
  import { map } from 'lodash-es'

  export let memo: Memo

  $: tags = map(memo.tags, (tag) => ({
    name: tag,
    filtered: $memoStore.selectedTags.includes(tag),
  }))

  function selectTag(event: CustomEvent<{ name: string }>) {
    insertTagFilter(event.detail.name)
  }
</script>

<li>
  {#if memo.tags.length > 0}
    <a href={`/${memo.id}`} class="link link-hover text inline-block max-w-full truncate">
      {memo.title}
    </a>
    <div class="flex gap-x-2">
      <div class="flex flex-1 flex-wrap gap-1">
        {#each tags as tag (tag.name)}
          <Tag
            value={tag.name}
            color={tag.filtered ? 'secondary' : 'primary'}
            on:click={selectTag}
          />
        {/each}
      </div>
      <span class="mt-[2px] flex-none text-xs font-light">{formatDate(memo.createTime)}</span>
    </div>
  {:else}
    <div class="flex items-center gap-x-2">
      <a
        href={`/${memo.id}`}
        class="link link-hover text inline-block max-w-full flex-auto truncate"
      >
        {memo.title}
      </a>
      <span class="mt-[2px] flex-none text-xs font-light">{formatDate(memo.createTime)}</span>
    </div>
  {/if}
</li>
