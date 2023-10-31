<script lang="ts">
  import MemoList from '$components/MemoList.svelte'
  import SignInStack from '$components/SignInStack.svelte'
  import TagFilter from '$components/TagFilter.svelte'
  import Plus from '$components/icons/Plus.svelte'
  import { authStore } from '$lib/auth'
  import { memoStore, syncMemos, type Memo } from '$lib/memo'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'
  import { every, filter } from 'lodash-es'
  import { onMount } from 'svelte'

  $: user = $authStore.user
  $: selectedTags = $memoStore.selectedTags
  let memos: Memo[]

  $: {
    if (selectedTags.length === 0) {
      memos = $memoStore.memos
    } else {
      memos = filter($memoStore.memos, (memo) =>
        every(selectedTags, (selected) => memo.tags.includes(selected))
      )
    }
  }

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
  <div class="space-y-2">
    <div class="flex justify-between">
      <TagFilter />
      <div>
        <a href="/new" class="btn btn-circle btn-sm btn-primary">
          <div class="w-[14px]"><Plus /></div>
        </a>
      </div>
    </div>
    <MemoList {memos} />
  </div>
{/if}
