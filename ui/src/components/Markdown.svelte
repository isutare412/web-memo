<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import { Marked } from '@ts-stack/markdown'
  import DOMPurify from 'isomorphic-dompurify'
  import ImageLightbox from '$components/ImageLightbox.svelte'

  export let content: string
  export let editable: boolean = false

  const dispatch = createEventDispatcher<{ checkboxToggle: { index: number } }>()

  const forbiddenTags = ['form', 'button']

  let lightboxUrl: string | null = null

  function processCheckboxes(html: string, isEditable: boolean): string {
    let index = 0
    return html.replace(/<li>(\s*)\[([ xX])\]/g, (_, space, checked) => {
      const isChecked = checked.toLowerCase() === 'x'
      const checkbox = `<li><input type="checkbox" data-checkbox-index="${index}" ${
        isChecked ? 'checked' : ''
      } ${isEditable ? '' : 'disabled'} class="checkbox checkbox-sm mr-2 align-middle" />`
      index++
      return checkbox
    })
  }

  function handleClick(event: MouseEvent) {
    const target = event.target as HTMLElement

    if (target.tagName === 'INPUT' && target.getAttribute('type') === 'checkbox') {
      const indexStr = target.getAttribute('data-checkbox-index')
      if (indexStr !== null && editable) {
        event.preventDefault()
        dispatch('checkboxToggle', { index: parseInt(indexStr, 10) })
      }
      return
    }

    if (target.tagName !== 'IMG') return

    const link = target.closest('a')
    if (!link) return

    // Check if this is an image link (link wrapping an image)
    event.preventDefault()
    lightboxUrl = link.href
  }

  function closeLightbox() {
    lightboxUrl = null
  }

  $: sanitizedHtml = processCheckboxes(
    DOMPurify.sanitize(Marked.parse(content, { breaks: true }), {
      FORBID_TAGS: forbiddenTags,
      ADD_TAGS: ['input'],
      ADD_ATTR: ['type', 'checked', 'disabled', 'data-checkbox-index'],
    }),
    editable
  )
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions a11y-no-noninteractive-element-interactions -->
<article class="prose max-w-none break-words" on:click={handleClick}>
  {@html sanitizedHtml}
</article>

{#if lightboxUrl}
  <ImageLightbox src={lightboxUrl} on:close={closeLightbox} />
{/if}
