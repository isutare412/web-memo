<script lang="ts">
  import { navigating } from '$app/stores'
  import Avatar from '$components/Avatar.svelte'
  import LoadingSpinner from '$components/LoadingSpinner.svelte'
  import ThemeToggle from '$components/ThemeToggle.svelte'
  import ToastContainer from '$components/ToastContainer.svelte'
  import {
      getPreferredTheme,
      setDocumentDataTheme,
      setPreferredTheme,
      type ThemeMode,
  } from '$lib/theme'
  import '@fontsource-variable/inter'
  import { onMount } from 'svelte'
  import '../app.css'

  let preferredTheme: ThemeMode | undefined = undefined

  onMount(async () => {
    preferredTheme = getPreferredTheme()
    setDocumentDataTheme(preferredTheme)
  })

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
        <ThemeToggle {onThemeToggle} initialMode={preferredTheme} />
        <Avatar />
      </div>
    </div>
  </nav>

  <main
    class="md:border-base-300 mx-auto mb-6 w-full max-w-3xl p-6 md:rounded-xl md:border md:shadow-md"
  >
    {#if $navigating !== null}
      <LoadingSpinner />
    {:else}
      <slot />
    {/if}
  </main>
  <ToastContainer />
</div>
