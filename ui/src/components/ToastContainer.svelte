<script lang="ts">
  import Toast from '$components/Toast.svelte'
  import { toastStore } from '$lib/toast'

  $: toasts = $toastStore.toasts
</script>

<div class="toast-center toast sm:toast-end w-full min-w-0 max-w-[500px]">
  {#each toasts as toast (toast.id)}
    <Toast {...toast} />
  {/each}
</div>
