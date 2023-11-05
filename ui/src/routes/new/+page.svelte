<script lang="ts">
  import { goto } from '$app/navigation'
  import MemoEditor from '$components/MemoEditor.svelte'
  import { createMemo } from '$lib/apis/backend/memo'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'

  async function onMemoSubmit(
    event: CustomEvent<{ title: string; content: string; tags: string[] }>
  ) {
    try {
      await createMemo(event.detail)
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      return
    }

    goto('/')
  }

  function onMemoCancel() {
    goto('/')
  }
</script>

<MemoEditor on:submit={onMemoSubmit} on:cancel={onMemoCancel} />
