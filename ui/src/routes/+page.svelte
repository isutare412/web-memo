<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import LoadingSpinner from '$components/LoadingSpinner.svelte'
  import MemoList from '$components/MemoList.svelte'
  import PageNavigator from '$components/PageNavigator.svelte'
  import SignInStack from '$components/SignInStack.svelte'
  import TagFilter from '$components/TagFilter.svelte'
  import Plus from '$components/icons/Plus.svelte'
  import Refresh from '$components/icons/Refresh.svelte'
  import { listMemos } from '$lib/apis/backend/memo'
  import { authStore } from '$lib/auth'
  import { defaultPageSize, mapToMemo, type MemoListPageData } from '$lib/memo'
  import { setPageOfSearchParams } from '$lib/searchParams'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'
  import { onMount } from 'svelte'
  import { get } from 'svelte/store'

  $: user = $authStore.user
  $: tags = $page.url.searchParams.getAll('tag')

  let currentPage: number = 1
  let listData: MemoListPageData | undefined

  $: {
    const rawPageStr = $page.url.searchParams.get('p') ?? '1'
    const rawPage = Number(rawPageStr)
    if (isNaN(rawPage)) {
      addToast(`page '${rawPageStr}' is invalid`, 'error')
      goto('/')
    } else {
      currentPage = rawPage
    }
  }

  onMount(() => {
    fetchMemos()
  })

  function onRefreshButtonClick() {
    fetchMemos()
  }

  function onNavigateEvent(event: CustomEvent<{ page: number }>) {
    const searchParams = get(page).url.searchParams
    if (!setPageOfSearchParams(searchParams, event.detail.page)) return

    goto(`/?${searchParams.toString()}`)
  }

  async function fetchMemos() {
    if (user === undefined) {
      return
    }

    listData = undefined

    try {
      const response = await listMemos(currentPage, defaultPageSize, tags)
      listData = {
        page: response.page,
        pageSize: response.pageSize,
        lastPage: response.lastPage,
        totalMemoCount: response.totalMemoCount,
        memos: response.memos.map(mapToMemo),
      }
    } catch (error) {
      addToast(getErrorMessage(error), 'error')
      goto('/')
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
