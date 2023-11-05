<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import MemoEditor from '$components/MemoEditor.svelte'
  import { replaceMemo } from '$lib/apis/backend/memo'
  import { memoStore, syncMemo } from '$lib/memo'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'

  $: memoId = $page.params.memoId
  $: memo = $memoStore.memos.find((memo) => memo.id === memoId)

  $: title = memo?.title ?? ''
  $: content = memo?.content ?? ''
  $: tags = memo?.tags ?? []

  onMount(async () => {
    try {
      await syncMemo(memoId)
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      goto('/')
      return
    }
  })

  async function onMemoSubmit(
    event: CustomEvent<{ title: string; content: string; tags: string[] }>
  ) {
    try {
      await replaceMemo({
        id: memoId,
        title: event.detail.title,
        content: event.detail.content,
        tags: event.detail.tags,
      })
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      return
    }

    goto(`/${memoId}`)
  }

  function onMemoCancel() {
    goto(`/${memoId}`)
  }
</script>

<MemoEditor {title} {content} {tags} on:submit={onMemoSubmit} on:cancel={onMemoCancel} />
