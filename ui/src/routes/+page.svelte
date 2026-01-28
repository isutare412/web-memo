<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import MemoList from '$components/MemoList.svelte'
  import PageNavigator from '$components/PageNavigator.svelte'
  import PageSizeSelector from '$components/PageSizeSelector.svelte'
  import SignInStack from '$components/SignInStack.svelte'
  import SortKeySelector from '$components/SortKeySelector.svelte'
  import TagFilter from '$components/TagFilter.svelte'
  import Plus from '$components/icons/Plus.svelte'
  import Refresh from '$components/icons/Refresh.svelte'
  import { listMemos } from '$lib/apis/backend/memo'
  import { authStore, syncUserData } from '$lib/auth'
  import {
    SortOrder,
    defaultPageSize,
    defaultSortOrder,
    mapToMemo,
    setPreferredPageSize,
    setPreferredSortOrder,
    type MemoListPageData,
  } from '$lib/memo'
  import {
    setPageOfSearchParams,
    setPageSizeOfSearchParams,
    setSortOrderOfSearchParams,
  } from '$lib/searchParams'
  import { loading } from '$lib/stores/loading'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'

  $: if (user && listData === undefined) loading.start()
  $: if (!user || listData !== undefined) loading.stop()

  $: user = $authStore.user
  $: tags = $page.url.searchParams.getAll('tag')
  $: searchParams = $page.url.searchParams

  let currentPage = 1
  let pageSize = defaultPageSize
  let sortOrder = defaultSortOrder
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

    const sortOrderStr = searchParams.get('sort') ?? defaultSortOrder.valueOf()
    sortOrder =
      Object.values(SortOrder).find((v) => v.valueOf() === sortOrderStr) ?? defaultSortOrder
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

  function onSortOrderChange(event: CustomEvent<{ sortKey: SortOrder }>) {
    if (!setSortOrderOfSearchParams(searchParams, event.detail.sortKey)) return

    setPreferredSortOrder(event.detail.sortKey)

    goto(`/?${searchParams.toString()}`)
  }

  async function fetchMemos() {
    if (user === undefined) {
      return
    }

    listData = undefined

    try {
      const response = await listMemos(currentPage, pageSize, sortOrder, tags)
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
{:else if listData !== undefined}
  <div class="space-y-2">
    <TagFilter>
      <div class="flex gap-2">
        <div>
          <a href="/new" class="btn btn-circle btn-primary btn-sm">
            <div class="w-[13px]"><Plus /></div>
          </a>
        </div>
        <div>
          <button on:click={onRefreshButtonClick} class="btn btn-circle btn-primary btn-sm">
            <div class="w-[16px]"><Refresh /></div>
          </button>
        </div>
      </div>
    </TagFilter>
    <div class="flex justify-end">
      <SortKeySelector sortKey={sortOrder} on:change={onSortOrderChange} />
    </div>
    <MemoList {user} memos={listData.memos} showUpdateTime={sortOrder === SortOrder.UPDATE_TIME} />
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
