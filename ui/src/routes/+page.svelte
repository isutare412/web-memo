<script lang="ts">
  import SignInStack from '$components/SignInStack.svelte'
  import { authStore } from '$lib/auth'

  $: user = $authStore.user
</script>

{#if !user}
  <SignInStack />
{:else}
  <h1>Welcome to SvelteKit</h1>
  <p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>
{/if}
