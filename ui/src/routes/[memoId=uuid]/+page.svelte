<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import LoadingSpinner from '$components/LoadingSpinner.svelte'
  import Markdown from '$components/Markdown.svelte'
  import Tag from '$components/Tag.svelte'
  import { deleteMemo, getMemo } from '$lib/apis/backend/memo'
  import { mapToMemo, type Memo } from '$lib/memo'
  import { addTagToSearchParams, setPageOfSearchParams } from '$lib/searchParams'
  import { addToast } from '$lib/toast'
  import { formatDate } from '$lib/utils/date'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'

  $: memoId = $page.params.memoId
  let memo: Memo | undefined

  let deleteConfirmModal: HTMLDialogElement
  let isDeleting = false

  onMount(async () => {
    try {
      memo = mapToMemo(await getMemo(memoId))
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      goto('/')
      return
    }
  })

  function onEditClick() {
    goto(`/${memoId}/edit`)
  }

  async function onDeleteClick() {
    deleteConfirmModal.showModal()
  }

  function onTagClick(event: CustomEvent<{ name: string }>) {
    const params = new URLSearchParams()
    setPageOfSearchParams(params, 1)
    addTagToSearchParams(params, event.detail.name)
    goto(`/?${params.toString()}`)
  }

  async function onDeleteConfirm() {
    isDeleting = true
    await deleteMemo(memoId)
    isDeleting = false

    goto('/', { replaceState: true })
  }
</script>

{#if memo === undefined}
  <LoadingSpinner />
{:else}
  <h1 class="break-words border-b py-2 text-3xl">{memo.title}</h1>
  {#if memo.tags.length > 0}
    <div class="mt-4 flex flex-wrap gap-1">
      {#each memo.tags as tag (tag)}
        <Tag value={tag} outline={true} on:click={onTagClick} />
      {/each}
    </div>
  {/if}
  <div class="mt-3">
    <Markdown content={memo.content} />
  </div>
  <div class="mt-3 flex flex-col gap-y-1">
    <div class="flex justify-end">
      <span class="text-xs opacity-70">Create {formatDate(memo.createTime)}</span>
    </div>
    <div class="flex justify-end">
      <span class="text-xs opacity-70">Update {formatDate(memo.updateTime)}</span>
    </div>
  </div>
  <div class="mt-4 flex justify-end gap-x-1">
    <button on:click={onEditClick} class="btn btn-outline btn-primary outline-none">Edit</button>
    <button on:click={onDeleteClick} class="btn btn-outline btn-primary outline-none">Delete</button
    >
  </div>
  <dialog bind:this={deleteConfirmModal} class="modal">
    <div class="modal-box">
      <p>Are you sure?</p>
      <div class="modal-action flex justify-end">
        <form method="dialog">
          <button class="btn btn-outline btn-primary outline-none">Cancel</button>
        </form>
        <button
          on:click={onDeleteConfirm}
          disabled={isDeleting}
          class="btn btn-primary outline-none"
        >
          {#if isDeleting}
            <span class="loading loading-spinner" />
          {:else}
            Delete
          {/if}
        </button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
{/if}
