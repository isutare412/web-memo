<script lang="ts">
  import Tag from '$components/Tag.svelte'
  import { reservedTags } from '$lib/memo'
  import { createEventDispatcher } from 'svelte'

  export let tags: string[] = []
  export let title: string = ''
  export let content: string = ''

  const dispatch = createEventDispatcher()

  let tagInputValue = ''
  let titleWarning = false
  let tagWarning: string | undefined = undefined
  let isSubmitting = false

  function onTagInputKeyUp(
    event: KeyboardEvent & { currentTarget: EventTarget & HTMLInputElement }
  ) {
    switch (event.key) {
      case 'Enter':
        if (!validateTag(tagInputValue)) return

        addTag(tagInputValue)
        break
    }
  }

  function onTagInputButtonClick() {
    if (!validateTag(tagInputValue)) return

    addTag(tagInputValue)
  }

  function onTitleInput() {
    if (title.trim() !== '') {
      titleWarning = false
    }
  }

  function onTagInput() {
    if (tagInputValue.trim() !== '') {
      tagWarning = undefined
    }
  }

  function onTagClick(event: CustomEvent<{ name: string }>) {
    tags = tags.filter((tag) => tag !== event.detail.name)
  }

  async function onSubmit() {
    if (title.trim() === '') {
      titleWarning = true
      return
    }

    if (tagInputValue !== '') {
      if (!validateTag(tagInputValue)) return

      addTag(tagInputValue)
    }

    isSubmitting = true
    dispatch('submit', {
      title,
      content,
      tags,
    })
  }

  function onCancel() {
    dispatch('cancel')
  }

  function validateTag(value: string): boolean {
    value = value.trim()
    if (!reservedTags.includes(value)) return true

    tagWarning = `"${value}" is a reserved tag`
    return false
  }

  function addTag(value: string) {
    value = value.trim()
    if (value === '') return

    if (tags.find((tag) => tag === value) !== undefined) {
      return
    }

    tags.push(value)
    tags = tags.toSorted()
    tagInputValue = ''
  }
</script>

<div class="flex flex-col gap-y-3">
  <div>
    {#if titleWarning}
      <label for="title" class="text-error text-xs">Need title</label>
    {/if}
    <input
      type="text"
      placeholder="Title"
      id="title"
      bind:value={title}
      on:input={onTitleInput}
      class="input input-bordered focus:border-primary w-full focus:outline-none"
      class:border-error={titleWarning}
      class:focus:border-error={titleWarning}
    />
  </div>
  <div>
    {#if tagWarning !== undefined}
      <label for="tag-input" class="text-error text-xs">{tagWarning}</label>
    {/if}
    <div class="flex">
      <input
        type="text"
        placeholder="Tag"
        id="tag-input"
        maxlength="20"
        bind:value={tagInputValue}
        on:input={onTagInput}
        on:keyup={onTagInputKeyUp}
        class="input input-bordered focus:border-primary w-full max-w-xs rounded-r-none border-r-0 focus:outline-none"
        class:border-error={tagWarning !== undefined}
        class:focus:border-error={tagWarning !== undefined}
      />
      <button
        on:click={onTagInputButtonClick}
        class="btn btn-primary btn-outline rounded-l-none"
        class:btn-error={tagWarning !== undefined}>Add</button
      >
    </div>
    {#if tags.length > 0}
      <div class="mt-2 flex flex-wrap gap-1">
        {#each tags as tag (tag)}
          {#if reservedTags.includes(tag)}
            <Tag value={tag} isButton={false} outline={true} />
          {:else}
            <Tag value={tag} outline={true} isClose={true} on:click={onTagClick} />
          {/if}
        {/each}
      </div>
    {/if}
  </div>
  <textarea
    placeholder="Content"
    bind:value={content}
    class="textarea textarea-bordered focus:border-primary h-[360px] text-base focus:outline-none"
  />
</div>
<div class="mt-4 flex justify-end gap-x-2">
  <button on:click={onCancel} class="btn btn-primary btn-outline">Cancel</button>
  <button on:click={onSubmit} disabled={isSubmitting} class="btn btn-primary">
    {#if isSubmitting}
      <span class="loading loading-spinner" />
    {:else}
      Submit
    {/if}
  </button>
</div>
