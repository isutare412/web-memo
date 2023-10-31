<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import Tag from '$components/Tag.svelte'
  import { deleteMemo } from '$lib/apis/backend/memo'
  import { memoStore, syncMemo } from '$lib/memo'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'

  $: memoId = $page.params.memoId
  $: memo = $memoStore.memos.find((memo) => memo.id === memoId)

  onMount(async () => {
    try {
      await syncMemo(memoId)
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      goto('/')
      return
    }
  })

  function onEditClick() {
    if (memo === undefined) return

    goto(`/${memoId}/edit`)
  }

  async function onDeleteClick() {
    if (memo === undefined) return

    await deleteMemo(memoId)
    goto('/', { replaceState: true })
  }
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
  <div class="mt-4 flex justify-end gap-x-1">
    <button on:click={onEditClick} class="btn btn-primary btn-sm">Edit</button>
    <button on:click={onDeleteClick} class="btn btn-primary btn-sm">Delete</button>
  </div>
{/if}
