<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import CollaborateButton from '$components/CollaborateButton.svelte'
  import CollaborationApproveTable from '$components/CollaborationApproveTable.svelte'
  import LinkCopyButton from '$components/LinkCopyButton.svelte'
  import LinkShareButton from '$components/LinkShareButton.svelte'
  import Markdown from '$components/Markdown.svelte'
  import SubscribeButton from '$components/SubscribeButton.svelte'
  import SubscriptionApproveTable from '$components/SubscriptionApproveTable.svelte'
  import Tag from '$components/Tag.svelte'
  import BookmarkIcon from '$components/icons/BookmarkIcon.svelte'
  import PenIcon from '$components/icons/PenIcon.svelte'
  import Refresh from '$components/icons/Refresh.svelte'
  import TrashBinIcon from '$components/icons/TrashBinIcon.svelte'
  import { StatusError } from '$lib/apis/backend/error'
  import {
    authorizeCollaboration,
    authorizeSubscription,
    cancelCollaboration,
    deleteMemo,
    getMemo,
    listCollaborators,
    listSubscribers,
    publishMemo,
    replaceMemo,
    requestCollaboration,
    subscribeMemo,
    unsubscribeMemo,
    type Collaborator,
    type Subscriber,
  } from '$lib/apis/backend/memo'
  import { authStore, signInGoogle, syncUserData } from '$lib/auth'
  import { loading } from '$lib/stores/loading'
  import { mapToMemo, toggleCheckboxInMarkdown } from '$lib/memo'
  import { addTagToSearchParams, setPageOfSearchParams } from '$lib/searchParams'
  import { ToastTimeout, addToast } from '$lib/toast'
  import { formatDate } from '$lib/utils/date'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'
  import type { PageData } from './$types'

  // Extract first image URL from markdown content for OpenGraph
  function extractFirstImageFromMarkdown(content: string): string | undefined {
    const match = content.match(/!\[.*?\]\((.*?)\)/)
    return match?.[1]
  }

  export let data: PageData

  $: user = $authStore.user
  $: memoId = $page.params.memoId!
  $: pageUrl = $page.url
  $: ({ memo } = data)
  $: isOwner = (user && memo && user.id === memo.ownerId) ?? false
  $: canEdit = isOwner || isMemoCollaborateApproved
  $: ogImage = memo ? extractFirstImageFromMarkdown(memo.content) : undefined

  $: if (memo === undefined) loading.start()
  $: if (memo !== undefined) loading.stop()

  let subscriberCount: number | undefined = undefined
  let subscribeConfirmModal: HTMLDialogElement
  let isSubscribing = false
  let subscriptionState: 'none' | 'pending' | 'approved' = 'none'

  let collaborateConfirmModal: HTMLDialogElement
  let isRequestingCollaborate = false
  let isMemoCollaborated = false
  let isMemoCollaborateApproved = false

  let collaborateApproveModal: HTMLDialogElement
  let collaborators: Collaborator[] = []
  let subscribers: Subscriber[] = []
  let subscriptionApproveModal: HTMLDialogElement
  $: hasApprovedCollaborators = collaborators.some((c) => c.isApproved)

  let publishConfirmModal: HTMLDialogElement
  let isPublishing = false
  let pendingPublishState: 'private' | 'shared' | 'published' = 'private'

  let deleteConfirmModal: HTMLDialogElement
  let isDeleting = false

  let isCheckboxUpdating = false

  onMount(async () => {
    await syncPageData(false)
  })

  async function syncPageData(forceRefresh: boolean) {
    try {
      await syncUserData()
      await syncMemo(forceRefresh)
      if (user && memo && user.id === memo.ownerId) {
        await Promise.all([syncOwnerSubscribers(), syncOwnerCollaborators()])
      }
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      goto('/')
      return
    }
  }

  async function syncMemo(forceRefresh: boolean) {
    if (forceRefresh || memo === undefined) {
      memo = undefined
      memo = mapToMemo(await getMemo(memoId))
    }

    if (memo.viewerContext) {
      if (memo.viewerContext.subscription === null) {
        subscriptionState = 'none'
      } else {
        subscriptionState = memo.viewerContext.subscription.isApproved ? 'approved' : 'pending'
      }
      isMemoCollaborated = memo.viewerContext.collaboration !== null
      isMemoCollaborateApproved = memo.viewerContext.collaboration?.isApproved ?? false
    }
  }

  async function syncOwnerSubscribers() {
    const response = await listSubscribers(memoId)
    subscribers = response.subscribers
    subscriberCount = subscribers.length
  }

  async function syncOwnerCollaborators() {
    const response = await listCollaborators(memoId)
    collaborators = response.collaborators
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

      if (subscriptionState !== 'none') {
        await unsubscribeMemo({ memoId, userId: user.id })
        subscriptionState = 'none'
      } else {
        const response = await subscribeMemo({ memoId, userId: user.id })
        subscriptionState = response.subscription.approved ? 'approved' : 'pending'
      }

      isSubscribing = false
      subscribeConfirmModal.close()

      if (memo.publishState === 'shared') {
        await syncMemo(true)
      }
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
      await syncOwnerCollaborators()
    } catch (error: unknown) {
      addToast(getErrorMessage(error), 'error')
      return
    }
  }

  async function onSubscriptionApprove(event: CustomEvent<{ userId: string; checked: boolean }>) {
    try {
      await authorizeSubscription({
        memoId,
        userId: event.detail.userId,
        approve: event.detail.checked,
      })
      await syncOwnerSubscribers()
    } catch (error: unknown) {
      addToast(getErrorMessage(error), 'error')
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

  function onShareEvent(event: CustomEvent<{ publishState: 'private' | 'shared' | 'published' }>) {
    const newState = event.detail.publishState
    pendingPublishState = newState
    publishConfirmModal.showModal()
  }

  async function onCheckboxToggle(event: CustomEvent<{ index: number }>) {
    if (!memo || !canEdit || isCheckboxUpdating) return

    isCheckboxUpdating = true
    const newContent = toggleCheckboxInMarkdown(memo.content, event.detail.index)

    try {
      const updatedMemo = await replaceMemo({
        id: memo.id,
        version: memo.version,
        title: memo.title,
        content: newContent,
        tags: memo.tags,
        isPinUpdateTime: true,
      })
      memo = mapToMemo(updatedMemo)
    } catch (error) {
      if (error instanceof StatusError && error.status === 409) {
        addToast('Memo was modified elsewhere. Refreshing...', 'warning')
        await syncMemo(true)
      } else {
        addToast(getErrorMessage(error), 'error')
      }
    } finally {
      isCheckboxUpdating = false
    }
  }

  function onPublishConfirm() {
    if (memo === undefined) return

    const memoUrl = pageUrl.toString()
    isPublishing = true

    const publishMemoPromise = publishMemo({
      id: memoId,
      publishState: pendingPublishState,
    })
      .then((rawMemo) => {
        memo = mapToMemo(rawMemo)

        if (memo.publishState === 'private') {
          subscriberCount = undefined
          subscribers = []
          collaborators = []
        } else {
          syncOwnerSubscribers()
        }

        isPublishing = false
        publishConfirmModal.close()
      })
      .catch((error) => {
        addToast(getErrorMessage(error), 'error')
        isPublishing = false
      })

    if (pendingPublishState !== 'private') {
      if (/^((?!chrome|android).)*safari/i.test(navigator.userAgent)) {
        navigator.clipboard
          .write([
            new ClipboardItem({
              'text/plain': publishMemoPromise.then(() => memoUrl),
            }),
          ])
          .then(() => {
            addToast('Copied memo URL!', 'info', { timeout: ToastTimeout.SHORT })
          })
      } else {
        publishMemoPromise.then(() => {
          navigator.clipboard.writeText(memoUrl).then(() => {
            addToast('Copied memo URL!', 'info', { timeout: ToastTimeout.SHORT })
          })
        })
      }
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
    {#if ogImage}
      <meta property="og:image" content={ogImage} />
    {/if}
  {/if}
</svelte:head>

{#if memo !== undefined}
  {#if isOwner}
    <div class="flex justify-end gap-x-2">
      <div>
        <button on:click={onDeleteClick} class="btn btn-circle btn-primary btn-sm">
          <div class="w-[17px]"><TrashBinIcon /></div>
        </button>
      </div>
      <div>
        <button on:click={onEditClick} class="btn btn-circle btn-primary btn-sm">
          <div class="w-[18px]"><PenIcon /></div>
        </button>
      </div>
      <div>
        <button on:click={onRefreshButtonClick} class="btn btn-circle btn-primary btn-sm">
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
          <button on:click={onEditClick} class="btn btn-circle btn-primary btn-sm">
            <div class="w-[18px]"><PenIcon /></div>
          </button>
        </div>
      {/if}
      <div>
        <button on:click={onRefreshButtonClick} class="btn btn-circle btn-primary btn-sm">
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

  {#if isOwner || memo.publishState !== 'shared' || subscriptionState === 'approved'}
    <div class="mt-3">
      <Markdown
        content={memo.content}
        editable={canEdit}
        disabled={isCheckboxUpdating}
        on:checkboxToggle={onCheckboxToggle}
      />
    </div>
  {:else if subscriptionState === 'pending'}
    <div class="alert mt-3">
      <span>Your access request is pending approval by the memo owner.</span>
    </div>
  {:else}
    <div class="alert mt-3">
      <span>This memo is shared. Request access to view the content.</span>
    </div>
  {/if}
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
      {#if memo.publishState !== 'private'}
        <CollaborateButton
          disabled={!collaborators.length}
          isActivated={hasApprovedCollaborators}
          count={collaborators.length}
          on:click={onCollaborateClick}
        />
      {/if}
      {#if memo.publishState === 'shared' && subscribers.length > 0}
        <button
          on:click={() => subscriptionApproveModal.showModal()}
          class="btn btn-outline btn-sm rounded-full"
          class:btn-primary={subscribers.some((s) => s.approved)}
        >
          <div class="w-[12px]">
            <BookmarkIcon />
          </div>
          {subscribers.length}
        </button>
      {/if}
      <LinkShareButton
        link={pageUrl.toString()}
        shareCount={subscriberCount}
        publishState={memo.publishState}
        on:share={onShareEvent}
      />
      <LinkCopyButton link={pageUrl.toString()} />
    </div>

    <dialog bind:this={collaborateApproveModal} class="modal">
      <div class="modal-box w-fit min-w-80">
        <CollaborationApproveTable {collaborators} on:change={onCollaborateApprove} />
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <dialog bind:this={subscriptionApproveModal} class="modal">
      <div class="modal-box w-fit min-w-80">
        <SubscriptionApproveTable {subscribers} on:change={onSubscriptionApprove} />
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <dialog bind:this={publishConfirmModal} class="modal">
      <div class="modal-box">
        <p class="py-4">
          {#if pendingPublishState === 'private'}
            <p>
              Will you make the memo <span class="font-bold text-primary">private</span>?<br />
              Only you and collaborators will be able to access the memo.
            </p>
          {:else if pendingPublishState === 'shared'}
            <p>
              Will you <span class="font-bold text-primary">share</span> the memo?<br />
              Users with the link can request access. You will approve who can view.
            </p>
          {:else}
            <p>
              Will you <span class="font-bold text-primary">publish</span> the memo?<br />
              Anyone with the link will be able to access the memo.
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
        disabled={memo.publishState === 'shared' && subscriptionState !== 'approved'}
        on:click={onCollaborateClick}
      />
      <SubscribeButton
        isActivated={subscriptionState !== 'none'}
        isChecked={subscriptionState === 'approved'}
        on:subscribe={onSusbscribeEvent}
      />
      <LinkCopyButton link={pageUrl.toString()} />
    </div>

    <dialog bind:this={subscribeConfirmModal} class="modal">
      <div class="modal-box">
        <p class="py-4">
          {#if subscriptionState !== 'none'}
            <p>
              Will you <span class="font-bold text-primary">unsubscribe</span> the memo?<br />
              The memo will not be exposed to your memo list.
            </p>
          {:else if memo.publishState === 'shared'}
            <p>
              Will you <span class="font-bold text-primary">request access</span> to this memo?<br
              />
              The owner will need to approve your request.
            </p>
          {:else}
            <p>
              Will you <span class="font-bold text-primary">subscribe</span> to this memo?<br />
              The memo will be exposed to your memo list.
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
              Will you <span class="font-bold text-primary">cancel</span> collaboration?
            </p>
          {:else}
            <p>
              Will you <span class="font-bold text-primary">request</span> collaboration?<br />You
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
