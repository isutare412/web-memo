<script lang="ts">
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import Autocomplete from '$components/Autocomplete.svelte'
  import Tag from '$components/Tag.svelte'
  import Funnel from '$components/icons/Funnel.svelte'
  import { listTags } from '$lib/apis/backend/memo'
  import {
      addTagToSearchParams,
      deleteTagFromSearchParams,
      setPageOfSearchParams,
  } from '$lib/searchParams'
  import { debounce, partition } from 'lodash-es'
  import { get } from 'svelte/store'

  $: selectedTags = $page.url.searchParams.getAll('tag')

  let showAutocomplete = false
  let tagCandidates: string[] = []
  let selectedCandidate: string | undefined

  let inputValue: string = ''
  let tagInput: HTMLInputElement
  let tagInputContainer: HTMLDivElement

  async function updateTagCandidates() {
    const tags = await listTags(inputValue.trim())

    tagCandidates = partition(tags, (tag) =>
      tag.toLowerCase().startsWith(inputValue.toLowerCase())
    ).flat()
  }

  async function onTagInput() {
    await updateTagCandidates()
    showAutocomplete = true
  }

  async function onTagInputFocus() {
    await updateTagCandidates()
    showAutocomplete = true
  }

  function onAutocompleteSelect(event: CustomEvent<{ item: string }>) {
    inputValue = event.detail.item
    addTag()
    showAutocomplete = false
  }

  function onTagInputKeyUp(
    event: KeyboardEvent & { currentTarget: EventTarget & HTMLInputElement }
  ) {
    if (showAutocomplete && selectedCandidate !== undefined) return

    switch (event.key) {
      case 'Enter':
        addTag()
        showAutocomplete = false
    }
  }

  function onTagFilterClick(event: CustomEvent<{ name: string }>) {
    const searchParams = get(page).url.searchParams
    if (!deleteTagFromSearchParams(searchParams, event.detail.name)) return

    setPageOfSearchParams(searchParams, 1)
    goto(`/?${searchParams.toString()}`)
  }

  function addTag() {
    const trimmedInput = inputValue.trim()
    if (trimmedInput.trim() === '') return

    const searchParams = get(page).url.searchParams
    if (!addTagToSearchParams(searchParams, trimmedInput)) return

    setPageOfSearchParams(searchParams, 1)
    goto(`/?${searchParams.toString()}`)
  }
</script>

<svelte:window
  on:keyup={(event) => {
    if (event.key === 'Escape') {
      showAutocomplete = false
      tagInput.blur()
    }
  }}
  on:click={(event) => {
    if (!(event.target instanceof Element)) return

    if (!tagInputContainer.contains(event.target)) {
      showAutocomplete = false
    }
  }}
/>

<div>
  <div class="flex justify-between gap-3">
    <div class="flex flex-1 items-center gap-2">
      <div class="w-4">
        <Funnel />
      </div>
      <div bind:this={tagInputContainer} class="w-full max-w-xs">
        <input
          type="text"
          placeholder="Tag"
          maxlength="20"
          bind:this={tagInput}
          bind:value={inputValue}
          on:keyup={onTagInputKeyUp}
          on:input={debounce(onTagInput, 500)}
          on:focus={onTagInputFocus}
          class="input input-sm input-bordered focus:border-primary w-full text-base focus:outline-none"
        />
        {#if showAutocomplete}
          <Autocomplete
            items={tagCandidates}
            bind:selectedItem={selectedCandidate}
            on:select={onAutocompleteSelect}
          />
        {/if}
      </div>
    </div>
    <slot />
  </div>
  <div class="mb-3 mt-2 flex flex-1 flex-wrap gap-1">
    {#each selectedTags as tag (tag)}
      <Tag value={tag} isClose={true} on:click={onTagFilterClick} />
    {/each}
  </div>
</div>
