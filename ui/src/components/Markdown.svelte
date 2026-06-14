<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import { Marked } from '@ts-stack/markdown'
  import DOMPurify from 'isomorphic-dompurify'
  import ImageLightbox from '$components/ImageLightbox.svelte'
  import { createHeadingRenderer } from '$lib/markdown'

  export let content: string
  export let editable: boolean = false
  export let disabled: boolean = false

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

    const heading = target.closest('h1, h2, h3, h4, h5, h6') as HTMLElement | null
    if (heading && heading.id && !target.closest('a')) {
      // Reflect the header in the URL without pushing a history entry, so the
      // back button still leaves the memo instead of walking prior clicks.
      history.replaceState(history.state, '', `#${heading.id}`)
      heading.scrollIntoView({ behavior: 'smooth' })
      return
    }

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
    DOMPurify.sanitize(Marked.parse(content, { breaks: true, renderer: createHeadingRenderer() }), {
      FORBID_TAGS: forbiddenTags,
      ADD_TAGS: ['input'],
      ADD_ATTR: ['type', 'checked', 'disabled', 'data-checkbox-index'],
    }),
    editable && !disabled
  )
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions a11y-no-noninteractive-element-interactions -->
<article class="prose max-w-none break-words prose-img:mx-auto" on:click={handleClick}>
  {@html sanitizedHtml}
</article>

{#if lightboxUrl}
  <ImageLightbox src={lightboxUrl} on:close={closeLightbox} />
{/if}

<style>
  article :global(h1),
  article :global(h2),
  article :global(h3),
  article :global(h4),
  article :global(h5),
  article :global(h6) {
    cursor: pointer;
  }

  article :global(h1:hover),
  article :global(h2:hover),
  article :global(h3:hover),
  article :global(h4:hover),
  article :global(h5:hover),
  article :global(h6:hover) {
    text-decoration: underline;
    text-decoration-thickness: 1px;
    text-underline-offset: 4px;
  }
</style>
