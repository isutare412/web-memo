<script lang="ts">
  import LoadingSpinner from '$components/LoadingSpinner.svelte'
  import MemoList from '$components/MemoList.svelte'
  import PageNavigator from '$components/PageNavigator.svelte'
  import SignInStack from '$components/SignInStack.svelte'
  import TagFilter from '$components/TagFilter.svelte'
  import Plus from '$components/icons/Plus.svelte'
  import Refresh from '$components/icons/Refresh.svelte'
  import { listMemos } from '$lib/apis/backend/memo'
  import { authStore } from '$lib/auth'
  import {
      clearUpdateInformer,
      informUpdate,
      mapToMemo,
      memoStore,
      setUpdateInformer,
      updateCurrentPage,
      type MemoListPageData,
  } from '$lib/memo'
  import { onDestroy, onMount } from 'svelte'
  import { get } from 'svelte/store'

  $: user = $authStore.user
  let listData: MemoListPageData | undefined

  onMount(() => {
    setUpdateInformer(fetchMemos)
    fetchMemos()
  })

  onDestroy(() => {
    clearUpdateInformer()
  })

  function onRefreshButtonClick() {
    informUpdate()
  }

  function onNavigateEvent(event: CustomEvent<{ page: number }>) {
    updateCurrentPage(event.detail.page)
    informUpdate()
  }

  async function fetchMemos() {
    if (user === undefined) {
      return
    }

    listData = undefined

    const { currentPage, pageSize, selectedTags } = get(memoStore)
    const response = await listMemos(currentPage, pageSize, selectedTags)
    listData = {
      page: response.page,
      pageSize: response.pageSize,
      lastPage: response.lastPage,
      totalMemoCount: response.totalMemoCount,
      memos: response.memos.map(mapToMemo),
    }
  }
</script>

{#if !user}
  <SignInStack />
{:else if listData === undefined}
  <LoadingSpinner />
{:else}
  <div class="space-y-2">
    <TagFilter>
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
    </TagFilter>
    <MemoList memos={listData.memos} />
    <div class="mt-6 flex justify-center">
      <PageNavigator
        currentPage={listData.page.toString()}
        lastPage={listData.lastPage.toString()}
        on:navigate={onNavigateEvent}
      />
    </div>
  </div>
{/if}
