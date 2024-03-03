<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import LinkCopyButton from '$components/LinkCopyButton.svelte'
  import LinkShareButton from '$components/LinkShareButton.svelte'
  import LoadingSpinner from '$components/LoadingSpinner.svelte'
  import Markdown from '$components/Markdown.svelte'
  import SubscribeButton from '$components/SubscribeButton.svelte'
  import Tag from '$components/Tag.svelte'
  import { StatusError } from '$lib/apis/backend/error'
  import {
      deleteMemo,
      getMemo,
      getSubscriber,
      listSubscribers,
      publishMemo,
      subscribeMemo,
      unsubscribeMemo,
  } from '$lib/apis/backend/memo'
  import { authStore, signInGoogle, syncUserData } from '$lib/auth'
  import { mapToMemo } from '$lib/memo'
  import { addTagToSearchParams, setPageOfSearchParams } from '$lib/searchParams'
  import { addToast } from '$lib/toast'
  import { formatDate } from '$lib/utils/date'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'
  import type { PageData } from './$types'

  export let data: PageData

  $: user = $authStore.user
  $: memoId = $page.params.memoId
  $: pageUrl = $page.url
  $: ({ memo } = data)
  $: hasWritePermission = (user && memo && user.id === memo.ownerId) ?? false

  let subscriberCount: number | undefined = undefined

  let subscribeConfirmModal: HTMLDialogElement
  let isSubscribing = false
  let isMemoSubscribed = false

  let publishConfirmModal: HTMLDialogElement
  let isPublishing = false

  let deleteConfirmModal: HTMLDialogElement
  let isDeleting = false

  onMount(async () => {
    try {
      await syncUserData()
      await syncMemo()
      await syncSubscribeStatus()
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      goto('/')
      return
    }
  })

  async function syncMemo() {
    if (memo !== undefined) return

    memo = mapToMemo(await getMemo(memoId))
  }

  async function syncSubscribeStatus() {
    if (user === undefined || memo === undefined) return

    try {
      if (memo.ownerId === user.id) {
        const response = await listSubscribers(memoId)
        subscriberCount = response.subscribers.length
      } else {
        await getSubscriber({ memoId, userId: user.id })
        isMemoSubscribed = true
      }
    } catch (error: unknown) {
      if (!(error instanceof StatusError)) {
        addToast(getErrorMessage(error), 'error')
        return
      } else if (error.status !== 404) {
        addToast(error.message, 'error')
        return
      }
    }
  }

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

    history.back()
  }

  function onSusbscribeEvent() {
    if (user === undefined) {
      signInGoogle()
      return
    }

    subscribeConfirmModal.showModal()
  }

  async function onSubscribeConfirm() {
    if (memo === undefined || user === undefined) return

    try {
      isSubscribing = true

      if (isMemoSubscribed) {
        await unsubscribeMemo({ memoId, userId: user.id })
        isMemoSubscribed = false
      } else {
        await subscribeMemo({ memoId, userId: user.id })
        isMemoSubscribed = true
      }

      isSubscribing = false
      subscribeConfirmModal.close()
    } catch (error: unknown) {
      addToast(getErrorMessage(error), 'error')
      return
    }
  }

  function onShareEvent() {
    publishConfirmModal.showModal()
  }

  function onPublishConfirm() {
    if (memo === undefined) {
      return
    }

    const memoUrl = pageUrl.toString()
    isPublishing = true

    // https://forums.developer.apple.com/forums/thread/691873
    // Apple WebKit does not allow async, so we build promise and pass it.
    const publishMemoPromise = publishMemo({
      id: memoId,
      publish: !memo.isPublished,
    })
      .then((rawMemo) => {
        memo = mapToMemo(rawMemo)

        if (!memo.isPublished) subscriberCount = undefined
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
            addToast('Copied memo URL!', 'info', { timeout: 1500 })
        })
    } else {
      publishMemoPromise.then(() => {
        navigator.clipboard.writeText(memoUrl).then(() => {
          if (memo !== undefined && memo.isPublished)
            addToast('Copied memo URL!', 'info', { timeout: 1500 })
        })
      })
    }
  }
</script>

<svelte:head>
  {#if memo !== undefined}
    <meta property="og:title" content={memo.title} />
    <meta
      property="og:description"
      content={memo.content.length > 200 ? `${memo.content.slice(0, 197)}...` : memo.content}
    />
  {/if}
</svelte:head>

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
  {#if hasWritePermission}
    <div class="mt-2 flex justify-end gap-x-1">
      <LinkCopyButton link={pageUrl.toString()} />
      <LinkShareButton
        link={pageUrl.toString()}
        shareCount={subscriberCount}
        isShared={memo.isPublished}
        on:share={onShareEvent}
      />
    </div>
    <div class="mt-4 flex justify-end gap-x-1">
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
  {:else}
    <div class="mt-2 flex justify-end gap-x-1">
      <SubscribeButton isActivated={isMemoSubscribed} on:subscribe={onSusbscribeEvent} />
    </div>

    <dialog bind:this={subscribeConfirmModal} class="modal">
      <div class="modal-box">
        <p>
          {#if isMemoSubscribed}
            <p>
              Will you <span class="text-primary font-bold">unsubscribe</span> the memo?<br />The
              memo will not be exposed to your memo list.
            </p>
          {:else}
            <p>
              Will you <span class="text-primary font-bold">subcribe</span> the memo?<br />The memo
              will be exposed to your memo list.
            </p>
          {/if}
        </p>
        <div class="modal-action flex justify-end">
          <form method="dialog">
            <button class="btn btn-outline btn-primary outline-none">Cancel</button>
          </form>
          <button
            on:click={onSubscribeConfirm}
            disabled={isSubscribing}
            class="btn btn-primary outline-none"
          >
            {#if isSubscribing}
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
  {/if}
{/if}
