<script lang="ts">
  import { Marked } from '@ts-stack/markdown'
  import DOMPurify from 'isomorphic-dompurify'
  import ImageLightbox from '$components/ImageLightbox.svelte'

  export let content: string

  const forbiddenTags = ['form', 'button']

  let lightboxUrl: string | null = null

  function handleClick(event: MouseEvent) {
    const target = event.target as HTMLElement
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
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions a11y-no-noninteractive-element-interactions -->
<article class="prose max-w-none break-words" on:click={handleClick}>
  {@html DOMPurify.sanitize(Marked.parse(content, { breaks: true }), {
    FORBID_TAGS: forbiddenTags,
  })}
</article>

{#if lightboxUrl}
  <ImageLightbox src={lightboxUrl} on:close={closeLightbox} />
{/if}
