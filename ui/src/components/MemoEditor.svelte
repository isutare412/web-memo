<script lang="ts">
  import Tag from '$components/Tag.svelte'
  import { reservedTags } from '$lib/memo'
  import { createEventDispatcher, tick } from 'svelte'

  export let tags: string[] = []
  export let title: string = ''
  export let content: string = ''

  const dispatch = createEventDispatcher()

  let tagInputValue = ''
  let textareaElement: HTMLTextAreaElement
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

  function onTextareaKeydown(
    event: KeyboardEvent & {
      currentTarget: EventTarget & HTMLTextAreaElement
    }
  ) {
    if (event.key === 'Enter' && !event.isComposing) {
      insertListSymbolTextarea(event)
      return
    }

    if (event.key === '~') {
      strikeThroughTextarea(event)
      return
    }

    if (event.key === '`') {
      codeBlockTextarea(event)
      return
    }

    if (event.key === '(' || event.key === ')') {
      parenthesisTextarea(event)
      return
    }

    if (event.key === '[' || event.key === ']') {
      squareBraketTextarea(event)
      return
    }
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

  async function insertListSymbolTextarea(
    event: KeyboardEvent & {
      currentTarget: EventTarget & HTMLTextAreaElement
    }
  ) {
    const cursorPos = textareaElement.selectionStart
    const textBeforeCursor = textareaElement.value.substring(0, cursorPos)
    const textAfterCursor = textareaElement.value.substring(cursorPos)
    const linesBeforeCursor = textBeforeCursor.split('\n')
    const lastLine = linesBeforeCursor[linesBeforeCursor.length - 1]

    const match = lastLine.match(/^(\s*)([-*+]|[0-9]+[.)])\s(.*)/)
    if (!match) return

    const indent = match[1]
    const listSymbol = match[2]
    const listContents = match[3]

    if (listContents === '') {
      event.preventDefault()
      content = textBeforeCursor.slice(0, -lastLine.length) + textAfterCursor
      await tick()
      textareaElement.selectionStart = cursorPos - lastLine.length
      textareaElement.selectionEnd = textareaElement.selectionStart
      return
    }

    let newLineText = ''
    let newListSymbol = ''
    if (listSymbol.match(/[0-9]+[.)]/)) {
      const currentNumber = parseInt(listSymbol, 10)
      newLineText = '\n' + indent + (currentNumber + 1) + '. '
      newListSymbol = (currentNumber + 1).toString() + '.'
    } else {
      newLineText = '\n' + indent + listSymbol + ' '
      newListSymbol = listSymbol
    }

    // Korean in iOS Safari does not fire composition event. As we cannot check
    // event.isComposing, we just delay the modification after composition
    // terminates by iOS.
    // https://discussionskorea.apple.com/thread/251376323?sortBy=best
    setTimeout(async () => {
      content = textBeforeCursor + newLineText + textAfterCursor
      await tick()
      textareaElement.selectionStart = cursorPos + indent.length + newListSymbol.length + 2
      textareaElement.selectionEnd = textareaElement.selectionStart
    }, 50)
  }

  async function strikeThroughTextarea(
    event: KeyboardEvent & {
      currentTarget: EventTarget & HTMLTextAreaElement
    }
  ) {
    const selectionStart = textareaElement.selectionStart
    const selectionEnd = textareaElement.selectionEnd
    if (selectionStart >= selectionEnd) return

    event.preventDefault()

    const textBeforeSelection = textareaElement.value.substring(0, selectionStart)
    const textInsideSelection = textareaElement.value.substring(selectionStart, selectionEnd)
    const textAfterSelection = textareaElement.value.substring(selectionEnd)

    if (
      textInsideSelection.length > 4 &&
      textInsideSelection.startsWith('~~') &&
      textInsideSelection.endsWith('~~')
    ) {
      content =
        textBeforeSelection +
        textInsideSelection.slice(2, textInsideSelection.length - 2) +
        textAfterSelection
      await tick()
      textareaElement.selectionStart = selectionStart
      textareaElement.selectionEnd = selectionEnd - 4
    } else {
      content = textBeforeSelection + '~~' + textInsideSelection + `~~` + textAfterSelection
      await tick()
      textareaElement.selectionStart = selectionStart
      textareaElement.selectionEnd = selectionEnd + 4
    }
  }

  async function codeBlockTextarea(
    event: KeyboardEvent & {
      currentTarget: EventTarget & HTMLTextAreaElement
    }
  ) {
    const selectionStart = textareaElement.selectionStart
    const selectionEnd = textareaElement.selectionEnd
    if (selectionStart >= selectionEnd) return

    event.preventDefault()

    const textBeforeSelection = textareaElement.value.substring(0, selectionStart)
    const textInsideSelection = textareaElement.value.substring(selectionStart, selectionEnd)
    const textAfterSelection = textareaElement.value.substring(selectionEnd)

    if (
      textInsideSelection.length > 2 &&
      textInsideSelection.startsWith('`') &&
      textInsideSelection.endsWith('`')
    ) {
      content =
        textBeforeSelection +
        textInsideSelection.slice(1, textInsideSelection.length - 1) +
        textAfterSelection
      await tick()
      textareaElement.selectionStart = selectionStart
      textareaElement.selectionEnd = selectionEnd - 2
    } else {
      content = textBeforeSelection + '`' + textInsideSelection + '`' + textAfterSelection
      await tick()
      textareaElement.selectionStart = selectionStart
      textareaElement.selectionEnd = selectionEnd + 2
    }
  }

  async function parenthesisTextarea(
    event: KeyboardEvent & {
      currentTarget: EventTarget & HTMLTextAreaElement
    }
  ) {
    const selectionStart = textareaElement.selectionStart
    const selectionEnd = textareaElement.selectionEnd
    if (selectionStart >= selectionEnd) return

    event.preventDefault()

    const textBeforeSelection = textareaElement.value.substring(0, selectionStart)
    const textInsideSelection = textareaElement.value.substring(selectionStart, selectionEnd)
    const textAfterSelection = textareaElement.value.substring(selectionEnd)

    if (
      textInsideSelection.length > 2 &&
      textInsideSelection.startsWith('(') &&
      textInsideSelection.endsWith(')')
    ) {
      content =
        textBeforeSelection +
        textInsideSelection.slice(1, textInsideSelection.length - 1) +
        textAfterSelection
      await tick()
      textareaElement.selectionStart = selectionStart
      textareaElement.selectionEnd = selectionEnd - 2
    } else {
      content = textBeforeSelection + '(' + textInsideSelection + ')' + textAfterSelection
      await tick()
      textareaElement.selectionStart = selectionStart
      textareaElement.selectionEnd = selectionEnd + 2
    }
  }

  async function squareBraketTextarea(
    event: KeyboardEvent & {
      currentTarget: EventTarget & HTMLTextAreaElement
    }
  ) {
    const selectionStart = textareaElement.selectionStart
    const selectionEnd = textareaElement.selectionEnd
    if (selectionStart >= selectionEnd) return

    event.preventDefault()

    const textBeforeSelection = textareaElement.value.substring(0, selectionStart)
    const textInsideSelection = textareaElement.value.substring(selectionStart, selectionEnd)
    const textAfterSelection = textareaElement.value.substring(selectionEnd)

    if (
      textInsideSelection.length > 2 &&
      textInsideSelection.startsWith('[') &&
      textInsideSelection.endsWith(']')
    ) {
      content =
        textBeforeSelection +
        textInsideSelection.slice(1, textInsideSelection.length - 1) +
        textAfterSelection
      await tick()
      textareaElement.selectionStart = selectionStart
      textareaElement.selectionEnd = selectionEnd - 2
    } else {
      content = textBeforeSelection + '[' + textInsideSelection + ']' + textAfterSelection
      await tick()
      textareaElement.selectionStart = selectionStart
      textareaElement.selectionEnd = selectionEnd + 2
    }
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
    bind:this={textareaElement}
    bind:value={content}
    on:keydown={onTextareaKeydown}
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
