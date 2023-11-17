<script lang="ts">
  import { page } from '$app/stores'
  import LoadingSpinner from '$components/LoadingSpinner.svelte'
  import MemoList from '$components/MemoList.svelte'
  import SignInStack from '$components/SignInStack.svelte'
  import TagFilter from '$components/TagFilter.svelte'
  import Plus from '$components/icons/Plus.svelte'
  import Refresh from '$components/icons/Refresh.svelte'
  import { authStore } from '$lib/auth'
  import { clearTagFilter, fetchPagedMemos, memoStore } from '$lib/memo'
  import { addToast } from '$lib/toast'
  import { getErrorMessage } from '$lib/utils/error'
  import { clone, isEqual } from 'lodash-es'
  import { onDestroy, onMount } from 'svelte'

  $: user = $authStore.user
  $: pagedMemos = $memoStore.pagedMemos
  $: currentPageStr = $page.url.searchParams.get('page') ?? '1'
  $: pageSizeStr = $page.url.searchParams.get('pageSize') ?? '10'

  let selectedTags: string[] = []
  let isFetchingMemo: Promise<void>

  const unsubscribe = memoStore.subscribe((state) => {
    if (isEqual(state.selectedTags, selectedTags)) return

    selectedTags = clone(state.selectedTags)
    isFetchingMemo = fetchMemos()
  })

  onMount(async () => {
    isFetchingMemo = fetchMemos()
  })

  onDestroy(() => unsubscribe)

  function onRefreshButtonClick() {
    isFetchingMemo = fetchMemos()
    clearTagFilter()
  }

  async function fetchMemos() {
    if (user === undefined) {
      return
    }

    try {
      const currentPage = Number(currentPageStr)
      if (isNaN(currentPage)) {
        addToast(`invalid page "${currentPageStr}"`, 'error')
        return
      }

      const pageSize = Number(pageSizeStr)
      if (isNaN(currentPage)) {
        addToast(`invalid page size "${pageSizeStr}"`, 'error')
        return
      }

      await fetchPagedMemos(currentPage, pageSize, selectedTags)
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
    <LoadingSpinner />
  {:then}
    {#if pagedMemos === null}
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
        <MemoList {...pagedMemos} />
      </div>
    {/if}
  {/await}
{/if}
