<script lang="ts">
  import MemoList from '$components/MemoList.svelte'
  import SignInStack from '$components/SignInStack.svelte'
  import TagFilter from '$components/TagFilter.svelte'
  import Plus from '$components/icons/Plus.svelte'
  import Refresh from '$components/icons/Refresh.svelte'
  import { authStore } from '$lib/auth'
  import { memoStore, syncMemos, type Memo } from '$lib/memo'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'
  import { every, filter } from 'lodash-es'
  import { onMount } from 'svelte'

  $: user = $authStore.user
  $: selectedTags = $memoStore.selectedTags
  let memos: Memo[]

  let isFetchingMemo: Promise<void>

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
    isFetchingMemo = fetchMemos()
  })

  function onRefreshButtonClick() {
    isFetchingMemo = fetchMemos()
  }

  async function fetchMemos() {
    if (user === undefined) {
      return
    }

    try {
      await syncMemos()
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      return
    }
  }
</script>

{#if !user}
  <SignInStack />
{:else}
  {#await isFetchingMemo}
    <div class="mx-auto my-6 w-fit">
      <span class="loading loading-spinner loading-lg" />
    </div>
  {:then}
    <div class="space-y-2">
      <div class="flex justify-between">
        <div class="mr-2">
          <TagFilter />
        </div>
        <div class="flex gap-2">
          <div>
            <button on:click={onRefreshButtonClick} class="btn btn-circle btn-sm btn-primary">
              <div class="w-[18px]"><Refresh /></div>
            </button>
          </div>
          <div>
            <a href="/new" class="btn btn-circle btn-sm btn-primary">
              <div class="w-[14px]"><Plus /></div>
            </a>
          </div>
        </div>
      </div>
      <MemoList {memos} />
    </div>
  {/await}
{/if}
