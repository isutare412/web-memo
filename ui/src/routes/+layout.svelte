<script lang="ts">
  import Avatar from '$components/Avatar.svelte'
  import ThemeToggle from '$components/ThemeToggle.svelte'
  import { setDocumentDataTheme, setPreferredTheme, type ThemeMode } from '$lib/theme'
  import '@fontsource-variable/inter'
  import '../app.css'
  import type { LayoutData } from './$types'

  export let data: LayoutData

  function onThemeToggle(theme: ThemeMode) {
    setDocumentDataTheme(theme)
    setPreferredTheme(theme)
  }
</script>

<div class="flex min-h-screen flex-col">
  <nav class="border-base-300 border-b shadow md:mb-4">
    <div class="mx-auto flex max-w-3xl items-center justify-between px-4 py-3">
      <a class="text-3xl" href="/">Web Memo</a>
      <div class="flex items-center gap-x-3">
        <ThemeToggle {onThemeToggle} initialMode={data.preferredTheme} />
        <Avatar />
      </div>
    </div>
  </nav>

  <main
    class="md:border-base-300 mx-auto mb-6 w-full max-w-3xl p-6 md:rounded-xl md:border md:shadow-md"
  >
    <slot />
  </main>
</div>
