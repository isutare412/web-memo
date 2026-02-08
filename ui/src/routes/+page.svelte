<script lang="ts">
  import { afterNavigate, goto } from '$app/navigation'
  import { page } from '$app/stores'
  import MemoList from '$components/MemoList.svelte'
  import PageNavigator from '$components/PageNavigator.svelte'
  import PageSizeSelector from '$components/PageSizeSelector.svelte'
  import SignInStack from '$components/SignInStack.svelte'
  import SortKeySelector from '$components/SortKeySelector.svelte'
  import TagFilter from '$components/TagFilter.svelte'
  import Plus from '$components/icons/Plus.svelte'
  import Refresh from '$components/icons/Refresh.svelte'
  import SearchIcon from '$components/icons/SearchIcon.svelte'
  import { listMemos, searchMemos } from '$lib/apis/backend/memo'
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

  $: user = $authStore.user
  $: tags = $page.url.searchParams.getAll('tag')
  $: searchParams = $page.url.searchParams

  let currentPage = 1
  let pageSize = defaultPageSize
  let sortOrder = defaultSortOrder
  let listData: MemoListPageData | undefined

  $: isSearchMode = searchParams.has('q')
  $: searchQuery = searchParams.get('q') ?? ''

  let searchInputValue = ''
  let searchInput: HTMLInputElement

  afterNavigate(() => {
    searchInputValue = searchParams.get('q') ?? ''
    if (isSearchMode) {
      setTimeout(() => searchInput?.focus(), 0)
    }
  })

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
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      return
    }
  })

  // Re-fetch memos when user logs in or search params change (browse mode only)
  $: if (!isSearchMode) {
    // eslint-disable-next-line @typescript-eslint/no-unused-expressions
    user, searchParams, fetchMemos()
  }

  // Fetch search results when query changes (search mode only)
  $: if (isSearchMode && searchQuery.trim() !== '') {
    // eslint-disable-next-line @typescript-eslint/no-unused-expressions
    user, searchQuery, fetchSearchResults()
  }

  // Clear results when entering search mode with empty query
  $: if (isSearchMode && searchQuery.trim() === '') {
    listData = undefined
  }

  function enterSearchMode() {
    goto('/?q=')
  }

  function enterBrowseMode() {
    goto('/')
  }

  function onRefreshButtonClick() {
    if (isSearchMode) {
      if (searchQuery.trim() !== '') {
        fetchSearchResults()
      }
    } else {
      fetchMemos()
    }
  }

  function onSearchKeyUp(event: KeyboardEvent) {
    if (event.key === 'Enter' && searchInputValue.trim() !== '') {
      goto(`/?q=${encodeURIComponent(searchInputValue.trim())}`)
    }
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

    loading.start()
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
    } finally {
      loading.stop()
    }
  }

  async function fetchSearchResults() {
    if (user === undefined) {
      return
    }

    loading.start()
    try {
      const response = await searchMemos(searchQuery)
      listData = {
        page: null,
        pageSize: null,
        lastPage: null,
        totalMemoCount: null,
        memos: response.memos.map(mapToMemo),
      }
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
    } finally {
      loading.stop()
    }
  }
</script>

<svelte:head>
  <meta property="og:title" content="Web Memo" />
</svelte:head>

{#if !user}
  <SignInStack />
{:else}
  <div class="space-y-2">
    <div class="flex justify-between gap-3">
      <div role="tablist" class="tabs-boxed tabs tabs-sm bg-base-200">
        <button
          type="button"
          role="tab"
          class="tab"
          class:tab-active={!isSearchMode}
          on:click={enterBrowseMode}
        >
          Browse
        </button>
        <button
          type="button"
          role="tab"
          class="tab"
          class:tab-active={isSearchMode}
          on:click={enterSearchMode}
        >
          Search
        </button>
      </div>
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
    </div>
    {#if isSearchMode}
      <div class="flex items-center gap-2">
        <div class="w-4 text-primary">
          <SearchIcon />
        </div>
        <div class="w-full">
          <input
            type="text"
            placeholder="Search memos..."
            bind:this={searchInput}
            bind:value={searchInputValue}
            on:keyup={onSearchKeyUp}
            class="input input-sm input-bordered w-full text-base focus:border-primary focus:outline-none"
          />
        </div>
      </div>
      {#if listData !== undefined}
        <MemoList {user} memos={listData.memos} showUpdateTime={false} />
      {/if}
    {:else}
      <TagFilter />
      <div class="flex justify-end">
        <SortKeySelector sortKey={sortOrder} on:change={onSortOrderChange} />
      </div>
      {#if listData !== undefined}
        <MemoList
          {user}
          memos={listData.memos}
          showUpdateTime={sortOrder === SortOrder.UPDATE_TIME}
        />
        <div class="flex justify-center">
          <PageNavigator
            currentPage={(listData.page ?? 1).toString()}
            lastPage={(listData.lastPage ?? 1).toString()}
            on:navigate={onNavigateEvent}
          />
        </div>
        <div class="flex justify-end">
          <PageSizeSelector currentSize={pageSize} on:change={onPageSizeSelectChange} />
        </div>
      {/if}
    {/if}
  </div>
{/if}
