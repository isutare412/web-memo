<script lang="ts">
  import { goto } from '$app/navigation'
  import MemoEditor from '$components/MemoEditor.svelte'
  import { createMemo } from '$lib/apis/backend/memo'
  import { syncUserData } from '$lib/auth'
  import { mapToMemo } from '$lib/memo'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'

  async function onMemoSubmit(
    event: CustomEvent<{ title: string; content: string; tags: string[] }>
  ) {
    try {
      await syncUserData()

      const memo = mapToMemo(await createMemo(event.detail))
      goto(`/${memo.id}`, { replaceState: true })
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      return
    }
  }

  function onMemoCancel() {
    history.back()
  }
</script>

<MemoEditor on:submit={onMemoSubmit} on:cancel={onMemoCancel} />
