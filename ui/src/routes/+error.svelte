<script lang="ts">
  import { page } from '$app/stores'
</script>

<div class="font-mono">
  <p class="mb-4 text-center text-4xl">
    {$page.status}
  </p>
  <p class="text-center text-xl">
    {$page.error?.message ?? 'Unknown Error'}
  </p>
</div>
