<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import Tag from '$components/Tag.svelte'
  import type { Memo } from '$lib/memo'
  import { addTagToSearchParams, setPageOfSearchParams } from '$lib/searchParams'
  import { formatDate } from '$lib/utils/date'
  import { map } from 'lodash-es'
  import { get } from 'svelte/store'

  export let memo: Memo

  $: tags = map(memo.tags, (tag) => ({
    name: tag,
    filtered: $page.url.searchParams.getAll('tag').includes(tag),
  }))

  function selectTag(event: CustomEvent<{ name: string }>) {
    const searchParams = get(page).url.searchParams
    if (!addTagToSearchParams(searchParams, event.detail.name)) return

    setPageOfSearchParams(searchParams, 1)
    goto(`/?${searchParams.toString()}`)
  }
</script>

<li class="flex flex-col gap-y-1">
  <div class="flex items-center">
    <a
      href={`/${memo.id}`}
      class="link link-hover inline-block max-w-full flex-auto truncate text-lg"
    >
      {memo.title}
    </a>
  </div>
  {#if tags.length > 0}
    <div class="flex flex-wrap gap-1">
      {#each tags as tag (tag.name)}
        <Tag
          value={tag.name}
          outline={!tag.filtered}
          isButton={!tag.filtered}
          on:click={selectTag}
        />
      {/each}
    </div>
  {/if}
  <div class="flex justify-end">
    <span class="text-xs font-light opacity-75">{formatDate(memo.createTime)}</span>
  </div>
</li>
