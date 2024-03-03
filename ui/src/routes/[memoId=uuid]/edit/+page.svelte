<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import LoadingSpinner from '$components/LoadingSpinner.svelte'
  import MemoEditor from '$components/MemoEditor.svelte'
  import { getMemo, replaceMemo } from '$lib/apis/backend/memo'
  import { syncUserData } from '$lib/auth'
  import { mapToMemo, type Memo } from '$lib/memo'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'

  $: memoId = $page.params.memoId
  let memo: Memo | undefined

  onMount(async () => {
    try {
      await syncUserData()

      memo = mapToMemo(await getMemo(memoId))
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      goto('/')
      return
    }
  })

  async function onMemoSubmit(
    event: CustomEvent<{ title: string; content: string; tags: string[]; isHold?: boolean }>
  ) {
    try {
      await replaceMemo({
        id: memoId,
        title: event.detail.title,
        content: event.detail.content,
        tags: event.detail.tags,
        isPinUpdateTime: event.detail.isHold,
      })

      if (event.detail.isHold) {
        addToast('Updated the memo without updating time.', 'info', { timeout: 2000 })
      }
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      return
    }

    history.back()
  }

  function onMemoCancel() {
    goto(`/${memoId}`)
  }
</script>

{#if memo === undefined}
  <LoadingSpinner />
{:else}
  <MemoEditor
    title={memo.title}
    content={memo.content}
    tags={memo.tags}
    on:submit={onMemoSubmit}
    on:cancel={onMemoCancel}
  />
{/if}
