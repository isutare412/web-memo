<script lang="ts">
  import { Marked } from '@ts-stack/markdown'

  export let content: string
</script>

<article class="prose break-words">
  {@html Marked.parse(content, { breaks: true })}
</article>
