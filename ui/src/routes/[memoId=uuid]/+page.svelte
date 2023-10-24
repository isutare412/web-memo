<script lang="ts">
  import { page } from '$app/stores'
  import Tag from '$components/Tag.svelte'
  import { memoStore, syncMemo } from '$lib/memo'
  import { onMount } from 'svelte'

  $: memoId = $page.params.memoId
  $: memo = $memoStore.memos.find((memo) => memo.id === memoId)

  onMount(async () => {
    await syncMemo(memoId)
  })
</script>

{#if memo !== undefined}
  <h1 class="break-words border-b py-2 text-2xl">{memo.title}</h1>
  {#if memo.tags.length > 0}
    <div class="mt-4 flex flex-wrap gap-1">
      {#each memo.tags as tag (tag)}
        <Tag value={tag} isButton={false} />
      {/each}
    </div>
  {/if}
  <div class="mt-3">
    <span class="whitespace-pre-wrap">{memo.content}</span>
  </div>
{/if}
