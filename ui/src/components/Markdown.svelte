<script lang="ts">
  import { Marked } from '@ts-stack/markdown'
  import DOMPurify from 'isomorphic-dompurify'

  export let content: string

  const forbiddenTags = ['form', 'button']
</script>

<article class="prose break-words">
  {@html DOMPurify.sanitize(Marked.parse(content, { breaks: true }), {
    FORBID_TAGS: forbiddenTags,
  })}
</article>
