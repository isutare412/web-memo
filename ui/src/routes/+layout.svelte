<script lang="ts">
  import GitHub from '$components/icons/GitHub.svelte'
  import '@fontsource-variable/inter'
  import '../app.css'
</script>

<div class="flex min-h-screen flex-col">
  <nav class="border-base-300 border-b shadow md:mb-4">
    <div class="mx-auto flex max-w-3xl items-center justify-between px-4 py-3">
      <a class="text-2xl" href="/">Web Memo</a>
      <a href="https://github.com/isutare412">
        <div class="hover:text-primary h-7 w-7"><GitHub /></div>
      </a>
    </div>
  </nav>

  <main
    class="md:border-base-300 mx-auto mb-6 w-full max-w-3xl p-6 md:rounded-xl md:border md:shadow-md"
  >
    <slot />
  </main>
</div>
