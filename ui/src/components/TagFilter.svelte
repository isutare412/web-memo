<script lang="ts">
  import Autocomplete from '$components/Autocomplete.svelte'
  import Tag from '$components/Tag.svelte'
  import Funnel from '$components/icons/Funnel.svelte'
  import { insertTagFilter, memoStore, removeTagFilter } from '$lib/memo'
  import { partition, sortBy, sortedUniq } from 'lodash-es'

  $: selectedTags = $memoStore.selectedTags

  let showAutocomplete = false
  let tagCandidates: string[] = []
  let selectedCandidate: string | undefined

  let inputValue: string = ''
  let tagInput: HTMLInputElement
  let tagInputContainer: HTMLDivElement

  function updateTagCandidates() {
    const uniqueTags = sortedUniq(
      sortBy(
        $memoStore.memos.flatMap((memo) => memo.tags),
        (tag) => tag.toLowerCase()
      )
    )

    if (inputValue === '') {
      tagCandidates = uniqueTags
      return
    }

    const inputValueLowered = inputValue.toLowerCase()
    const includeInput = uniqueTags.filter((tag) => tag.toLowerCase().includes(inputValueLowered))
    tagCandidates = partition(includeInput, (tag) =>
      tag.toLowerCase().startsWith(inputValueLowered)
    ).flat()
  }

  function onTagInput() {
    updateTagCandidates()
    showAutocomplete = true
  }

  function onTagInputFocus() {
    updateTagCandidates()
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
    removeTagFilter(event.detail.name)
  }

  function addTag() {
    const trimmedInput = inputValue.trim()
    if (trimmedInput.trim() === '') return

    insertTagFilter(trimmedInput)
    inputValue = ''
    tagInput.blur()
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
          on:input={onTagInput}
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
