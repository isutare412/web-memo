<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import LoadingSpinner from '$components/LoadingSpinner.svelte'
  import MemoList from '$components/MemoList.svelte'
  import PageNavigator from '$components/PageNavigator.svelte'
  import PageSizeSelector from '$components/PageSizeSelector.svelte'
  import SignInStack from '$components/SignInStack.svelte'
  import TagFilter from '$components/TagFilter.svelte'
  import Plus from '$components/icons/Plus.svelte'
  import Refresh from '$components/icons/Refresh.svelte'
  import { listMemos } from '$lib/apis/backend/memo'
  import { authStore, syncUserData } from '$lib/auth'
  import {
      defaultPageSize,
      mapToMemo,
      setPreferredPageSize,
      type MemoListPageData,
  } from '$lib/memo'
  import { setPageOfSearchParams, setPageSizeOfSearchParams } from '$lib/searchParams'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'

  $: user = $authStore.user
  $: tags = $page.url.searchParams.getAll('tag')
  $: searchParams = $page.url.searchParams

  let currentPage = 1
  let pageSize = defaultPageSize
  let listData: MemoListPageData | undefined

  $: {
    const rawPageStr = searchParams.get('p') ?? '1'
    const rawPage = Number(rawPageStr)
    if (isNaN(rawPage)) {
      addToast(`page '${rawPageStr}' is invalid`, 'error')
    } else {
      currentPage = rawPage
    }

    const rawPageSizeStr = searchParams.get('ps') ?? defaultPageSize.toString()
    const rawPageSize = Number(rawPageSizeStr)
    if (isNaN(rawPageSize)) {
      addToast(`page size '${rawPageSizeStr}' is invalid`, 'error')
    } else {
      pageSize = rawPageSize
    }
  }

  onMount(async () => {
    try {
      await syncUserData()
      await fetchMemos()
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      return
    }
  })

  function onRefreshButtonClick() {
    fetchMemos()
  }

  function onNavigateEvent(event: CustomEvent<{ page: number }>) {
    if (!setPageOfSearchParams(searchParams, event.detail.page)) return

    goto(`/?${searchParams.toString()}`)
  }

  function onPageSizeSelectChange(event: CustomEvent<{ pageSize: number }>) {
    if (!setPageSizeOfSearchParams(searchParams, event.detail.pageSize)) return

    setPreferredPageSize(event.detail.pageSize)
    setPageOfSearchParams(searchParams, 1)

    goto(`/?${searchParams.toString()}`)
  }

  async function fetchMemos() {
    if (user === undefined) {
      return
    }

    listData = undefined

    try {
      const response = await listMemos(currentPage, pageSize, tags)
      listData = {
        page: response.page,
        pageSize: response.pageSize,
        lastPage: response.lastPage,
        totalMemoCount: response.totalMemoCount,
        memos: response.memos.map(mapToMemo),
      }
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
    }
  }
</script>

<svelte:head>
  <meta property="og:title" content="Web Memo" />
</svelte:head>

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
    <div class="flex justify-center">
      <PageNavigator
        currentPage={listData.page.toString()}
        lastPage={listData.lastPage.toString()}
        on:navigate={onNavigateEvent}
      />
    </div>
    <div class="flex justify-end">
      <PageSizeSelector currentSize={pageSize} on:change={onPageSizeSelectChange} />
    </div>
  </div>
{/if}
