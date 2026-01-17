<script lang="ts">
  import Toast from '$components/Toast.svelte'
  import { toastStore } from '$lib/toast'

  $: toasts = $toastStore.toasts
</script>

<div class="toast toast-center w-full min-w-0 max-w-[500px] sm:toast-end">
  {#each toasts as toast (toast.id)}
    <Toast {...toast} />
  {/each}
</div>
