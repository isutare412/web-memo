<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import Tag from '$components/Tag.svelte'
  import BookmarkIcon from '$components/icons/BookmarkIcon.svelte'
  import PeopleIcon from '$components/icons/PeopleIcon.svelte'
  import type { UserData } from '$lib/auth'
  import type { Memo } from '$lib/memo'
  import { addTagToSearchParams, setPageOfSearchParams } from '$lib/searchParams'
  import { formatDate } from '$lib/utils/date'
  import { map } from 'lodash-es'
  import { get } from 'svelte/store'

  export let user: UserData
  export let memo: Memo
  export let showUpdateTime = false

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
  <div class="flex items-center justify-end gap-x-1">
    {#if memo.ownerId !== user.id}
      <div class="text-primary w-[12px]">
        <BookmarkIcon />
      </div>
    {:else if memo.isPublished}
      <div class="text-primary w-[16px]">
        <PeopleIcon />
      </div>
    {/if}
    <div>
      <span class="text-xs font-light opacity-75"
        >{formatDate(showUpdateTime ? memo.updateTime : memo.createTime)}</span
      >
    </div>
  </div>
</li>
