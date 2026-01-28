<script lang="ts">
  import { navigating } from '$app/stores'
  import { loading } from '$lib/stores/loading'
</script>

{#if $navigating || $loading}
  <progress class="progress progress-primary fixed left-0 top-0 z-50 h-1 w-full" />
{/if}
