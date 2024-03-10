<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import CollaborateButton from '$components/CollaborateButton.svelte'
  import CollaborationApproveTable from '$components/CollaborationApproveTable.svelte'
  import LinkCopyButton from '$components/LinkCopyButton.svelte'
  import LinkShareButton from '$components/LinkShareButton.svelte'
  import LoadingSpinner from '$components/LoadingSpinner.svelte'
  import Markdown from '$components/Markdown.svelte'
  import SubscribeButton from '$components/SubscribeButton.svelte'
  import Tag from '$components/Tag.svelte'
  import PenIcon from '$components/icons/PenIcon.svelte'
  import Refresh from '$components/icons/Refresh.svelte'
  import TrashBinIcon from '$components/icons/TrashBinIcon.svelte'
  import { StatusError } from '$lib/apis/backend/error'
  import {
      authorizeCollaboration,
      cancelCollaboration,
      deleteMemo,
      getCollaborator,
      getMemo,
      getSubscriber,
      listCollaborators,
      listSubscribers,
      publishMemo,
      requestCollaboration,
      subscribeMemo,
      unsubscribeMemo,
      type Collaborator,
  } from '$lib/apis/backend/memo'
  import { authStore, signInGoogle, syncUserData } from '$lib/auth'
  import { mapToMemo } from '$lib/memo'
  import { addTagToSearchParams, setPageOfSearchParams } from '$lib/searchParams'
  import { ToastTimeout, addToast } from '$lib/toast'
  import { formatDate } from '$lib/utils/date'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'
  import type { PageData } from './$types'

  export let data: PageData

  $: user = $authStore.user
  $: memoId = $page.params.memoId
  $: pageUrl = $page.url
  $: ({ memo } = data)
  $: isOwner = (user && memo && user.id === memo.ownerId) ?? false

  let subscriberCount: number | undefined = undefined
  let subscribeConfirmModal: HTMLDialogElement
  let isSubscribing = false
  let isMemoSubscribed = false

  let collaborateConfirmModal: HTMLDialogElement
  let isRequestingCollaborate = false
  let isMemoCollaborated = false
  let isMemoCollaborateApproved = false

  let collaborateApproveModal: HTMLDialogElement
  let collaborators: Collaborator[] = []
  $: hasApprovedCollaborators = collaborators.some((c) => c.isApproved)

  let publishConfirmModal: HTMLDialogElement
  let isPublishing = false

  let deleteConfirmModal: HTMLDialogElement
  let isDeleting = false

  onMount(async () => {
    await syncPageData(false)
  })

  async function syncPageData(forceRefresh: boolean) {
    try {
      await syncUserData()
      await syncMemo(forceRefresh)
      await Promise.all([syncSubscribeStatus(), syncCollborationStatus()])
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      goto('/')
      return
    }
  }

  async function syncMemo(forceRefresh: boolean) {
    if (!forceRefresh && memo !== undefined) return

    memo = undefined
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

  async function syncCollborationStatus() {
    if (user === undefined || memo === undefined) return

    try {
      if (memo.ownerId === user.id) {
        const response = await listCollaborators(memoId)
        collaborators = response.collaborators
      } else {
        const collaborator = await getCollaborator({ memoId, userId: user.id })
        isMemoCollaborated = true
        isMemoCollaborateApproved = collaborator.isApproved
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

  async function onRefreshButtonClick() {
    await syncPageData(true)
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

  function onCollaborateClick() {
    if (user === undefined) {
      signInGoogle()
      return
    }

    if (memo == undefined) return

    if (memo.ownerId !== user.id) {
      collaborateConfirmModal.showModal()
      return
    }

    collaborateApproveModal.showModal()
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

  async function onCollaborateApprove(event: CustomEvent<{ userId: string; checked: boolean }>) {
    try {
      await authorizeCollaboration({
        memoId,
        userId: event.detail.userId,
        approve: event.detail.checked,
      })
      await syncCollborationStatus()
    } catch (error: unknown) {
      addToast(getErrorMessage(error), 'error')
      return
    }
  }

  async function onCollaborateConfirm() {
    if (memo === undefined || user === undefined) return

    try {
      isRequestingCollaborate = true

      if (isMemoCollaborated) {
        await cancelCollaboration({ memoId, userId: user.id })
        isMemoCollaborated = false
        isMemoCollaborateApproved = false
      } else {
        await requestCollaboration({ memoId, userId: user.id })
        isMemoCollaborated = true
        isMemoCollaborateApproved = false
      }

      isRequestingCollaborate = false
      collaborateConfirmModal.close()
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

        if (!memo.isPublished) {
          subscriberCount = undefined
          collaborators = []
        }

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
            addToast('Copied memo URL!', 'info', { timeout: ToastTimeout.SHORT })
        })
    } else {
      publishMemoPromise.then(() => {
        navigator.clipboard.writeText(memoUrl).then(() => {
          if (memo !== undefined && memo.isPublished)
            addToast('Copied memo URL!', 'info', { timeout: ToastTimeout.SHORT })
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
  {#if isOwner}
    <div class="flex justify-end gap-x-2">
      <div>
        <button on:click={onDeleteClick} class="btn btn-circle btn-sm btn-primary">
          <div class="w-[17px]"><TrashBinIcon /></div>
        </button>
      </div>
      <div>
        <button on:click={onEditClick} class="btn btn-circle btn-sm btn-primary">
          <div class="w-[18px]"><PenIcon /></div>
        </button>
      </div>
      <div>
        <button on:click={onRefreshButtonClick} class="btn btn-circle btn-sm btn-primary">
          <div class="w-[16px]"><Refresh /></div>
        </button>
      </div>
    </div>

    <dialog bind:this={deleteConfirmModal} class="modal">
      <div class="modal-box">
        <p class="py-4">Are you sure?</p>
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

    <!-- not owner -->
  {:else}
    <div class="flex justify-end gap-x-2">
      {#if isMemoCollaborateApproved}
        <div>
          <button on:click={onEditClick} class="btn btn-circle btn-sm btn-primary">
            <div class="w-[18px]"><PenIcon /></div>
          </button>
        </div>
      {/if}
      <div>
        <button on:click={onRefreshButtonClick} class="btn btn-circle btn-sm btn-primary">
          <div class="w-[16px]"><Refresh /></div>
        </button>
      </div>
    </div>
  {/if}

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

  {#if isOwner}
    <div class="mt-2 flex justify-end gap-x-1">
      {#if memo.isPublished}
        <CollaborateButton
          disabled={!collaborators.length}
          isActivated={hasApprovedCollaborators}
          count={collaborators.length}
          on:click={onCollaborateClick}
        />
      {/if}
      <LinkShareButton
        link={pageUrl.toString()}
        shareCount={subscriberCount}
        isShared={memo.isPublished}
        on:share={onShareEvent}
      />
      <LinkCopyButton link={pageUrl.toString()} />
    </div>

    <dialog bind:this={collaborateApproveModal} class="modal">
      <div class="modal-box min-w-80 w-fit">
        <CollaborationApproveTable {collaborators} on:change={onCollaborateApprove} />
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <dialog bind:this={publishConfirmModal} class="modal">
      <div class="modal-box">
        <p class="py-4">
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

    <!-- not owner -->
  {:else}
    <div class="mt-2 flex justify-end gap-x-1">
      <CollaborateButton
        isActivated={isMemoCollaborated}
        isChecked={isMemoCollaborateApproved}
        on:click={onCollaborateClick}
      />
      <SubscribeButton isActivated={isMemoSubscribed} on:subscribe={onSusbscribeEvent} />
      <LinkCopyButton link={pageUrl.toString()} />
    </div>

    <dialog bind:this={subscribeConfirmModal} class="modal">
      <div class="modal-box">
        <p class="py-4">
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

    <dialog bind:this={collaborateConfirmModal} class="modal">
      <div class="modal-box">
        <p class="py-4">
          {#if isMemoCollaborated}
            <p>
              Will you <span class="text-primary font-bold">cancel</span> collaboration?
            </p>
          {:else}
            <p>
              Will you <span class="text-primary font-bold">request</span> collaboration?<br />You
              will be able to modify the memo after an approval from the owner.
            </p>
          {/if}
        </p>
        <div class="modal-action flex justify-end">
          <form method="dialog">
            <button class="btn btn-outline btn-primary outline-none">No</button>
          </form>
          <button
            on:click={onCollaborateConfirm}
            disabled={isRequestingCollaborate}
            class="btn btn-primary outline-none"
          >
            {#if isRequestingCollaborate}
              <span class="loading loading-spinner" />
            {:else}
              Yes
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
