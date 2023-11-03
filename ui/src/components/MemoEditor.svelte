<script lang="ts">
  import Tag from '$components/Tag.svelte'
  import { createEventDispatcher } from 'svelte'

  export let tags: string[] = []
  export let title: string = ''
  export let content: string = ''

  const dispatch = createEventDispatcher()

  let tagInputValue = ''
  let warnTitle = false

  function onTagInputKeyUp(
    event: KeyboardEvent & { currentTarget: EventTarget & HTMLInputElement }
  ) {
    switch (event.key) {
      case 'Enter':
        addTag(tagInputValue)
        break
    }
  }

  function onTagInputButtonClick() {
    addTag(tagInputValue)
  }

  function onTitleInput() {
    if (title.trim() !== '') {
      warnTitle = false
    }
  }

  function onTagClick(event: CustomEvent<{ name: string }>) {
    tags = tags.filter((tag) => tag !== event.detail.name)
  }

  async function onSubmit() {
    if (title.trim() === '') {
      warnTitle = true
      return
    }

    if (tagInputValue !== '') {
      addTag(tagInputValue)
    }

    dispatch('submit', {
      title,
      content,
      tags,
    })
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
    {#if warnTitle}
      <label for="title" class="text-error text-xs">Need title</label>
    {/if}
    <input
      type="text"
      placeholder="Title"
      id="title"
      bind:value={title}
      on:input={onTitleInput}
      class="input input-bordered focus:border-primary w-full focus:outline-none"
      class:border-error={warnTitle}
    />
  </div>
  <div>
    <div class="flex">
      <input
        type="text"
        placeholder="Tag"
        maxlength="20"
        bind:value={tagInputValue}
        on:keyup={onTagInputKeyUp}
        class="input input-bordered focus:border-primary w-full max-w-xs rounded-r-none border-r-0 focus:outline-none"
      />
      <button on:click={onTagInputButtonClick} class="btn btn-primary btn-outline rounded-l-none"
        >Add</button
      >
    </div>
    {#if tags.length > 0}
      <div class="mt-2 flex flex-wrap gap-1">
        {#each tags as tag (tag)}
          <Tag value={tag} outline={true} isClose={true} on:click={onTagClick} />
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
<div class="mt-4 flex justify-end">
  <button on:click={onSubmit} class="btn btn-primary btn-outline">Submit</button>
</div>
