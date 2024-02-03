<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import LoadingSpinner from '$components/LoadingSpinner.svelte'
  import Markdown from '$components/Markdown.svelte'
  import Tag from '$components/Tag.svelte'
  import Share from '$components/icons/Share.svelte'
  import { deleteMemo, getMemo, publishMemo } from '$lib/apis/backend/memo'
  import { authStore, syncUserData } from '$lib/auth'
  import { mapToMemo } from '$lib/memo'
  import { addTagToSearchParams, setPageOfSearchParams } from '$lib/searchParams'
  import { addToast } from '$lib/toast'
  import { formatDate } from '$lib/utils/date'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'
  import type { PageServerData } from './$types'

  export let data: PageServerData

  $: user = $authStore.user
  $: memoId = $page.params.memoId
  $: ({ memo } = data)

  let publishConfirmModal: HTMLDialogElement
  let isPublishing = false

  let deleteConfirmModal: HTMLDialogElement
  let isDeleting = false

  onMount(async () => {
    try {
      await syncUserData()

      if (memo === undefined) memo = mapToMemo(await getMemo(memoId))
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

  async function onPublishClick() {
    publishConfirmModal.showModal()
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

  function onPublishConfirm() {
    if (memo === undefined) {
      return
    }

    const memoUrl = $page.url.toString()
    isPublishing = true

    // https://forums.developer.apple.com/forums/thread/691873
    // Apple WebKit does not allow async, so we build promise and pass it.
    const publishMemoPromise = publishMemo({
      id: memoId,
      publish: !memo.isPublished,
    })
      .then((rawMemo) => {
        memo = mapToMemo(rawMemo)

        isPublishing = false
        publishConfirmModal.close()
      })
      .catch((error) => {
        addToast(getErrorMessage(error), 'error')
        goto('/')
      })

    if (/^((?!chrome|android).)*safari/i.test(navigator.userAgent)) {
      navigator.clipboard
        .write([
          new ClipboardItem({
            'text/plain': publishMemoPromise.then(() => memoUrl),
          }),
        ])
        .then(() => {
          if (memo !== undefined && memo.isPublished)
            addToast('Copied your memo URL to clipboard!', 'info')
        })
    } else {
      publishMemoPromise.then(() => {
        navigator.clipboard.writeText(memoUrl).then(() => {
          if (memo !== undefined && memo.isPublished)
            addToast('Copied your memo URL to clipboard!', 'info')
        })
      })
    }
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
  {#if user && (!memo.isPublished || user.id === memo.ownerId)}
    <div class="mt-4 flex justify-end gap-x-1">
      <button
        on:click={onPublishClick}
        class="btn btn-primary btn-square outline-none"
        class:btn-outline={!memo.isPublished}
      >
        <div class="w-[24px]">
          <Share />
        </div>
      </button>
      <button on:click={onEditClick} class="btn btn-outline btn-primary outline-none">Edit</button>
      <button on:click={onDeleteClick} class="btn btn-outline btn-primary outline-none"
        >Delete</button
      >
    </div>

    <dialog bind:this={publishConfirmModal} class="modal">
      <div class="modal-box">
        <p>
          {#if memo.isPublished}
            <p>
              Will you <span class="text-primary font-bold">un-publish</span> the memo?<br />Only
              you will be able to access the memo through a link.
            </p>
          {:else}
            <p>
              Will you <span class="text-primary font-bold">publish</span> the memo?<br />Anyone
              will be able to access the memo through a link.
            </p>
          {/if}
        </p>
        <div class="modal-action flex justify-end">
          <form method="dialog">
            <button class="btn btn-outline btn-primary outline-none">Cancel</button>
          </form>
          <button
            on:click={onPublishConfirm}
            disabled={isPublishing}
            class="btn btn-primary outline-none"
          >
            {#if isPublishing}
              <span class="loading loading-spinner" />
            {:else}
              OK
            {/if}
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

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
{/if}
