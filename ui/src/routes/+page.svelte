<script lang="ts">
  import MemoList from '$components/MemoList.svelte'
  import SignInStack from '$components/SignInStack.svelte'
  import { authStore } from '$lib/auth'
  import { memoStore, syncMemos } from '$lib/memo'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'

  $: user = $authStore.user
  $: memos = $memoStore.memos

  onMount(async () => {
    if (user === undefined) {
      return
    }

    try {
      await syncMemos()
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      return
    }
  })
</script>

{#if !user}
  <SignInStack />
{:else}
  <MemoList {memos} />
{/if}
