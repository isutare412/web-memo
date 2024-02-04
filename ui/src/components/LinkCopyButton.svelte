<script lang="ts">
  import CopyIcon from '$components/icons/CopyIcon.svelte'
  import { addToast } from '$lib/toast'

  export let link: string

  function onClickCopyButton() {
    // https://forums.developer.apple.com/forums/thread/691873
    // Apple WebKit does not allow async.
    if (/^((?!chrome|android).)*safari/i.test(navigator.userAgent)) {
      navigator.clipboard
        .write([
          new ClipboardItem({
            'text/plain': Promise.resolve(link),
          }),
        ])
        .then(showCopyResult)
    } else {
      navigator.clipboard.writeText(link).then(showCopyResult)
    }
  }

  function showCopyResult() {
    addToast('Copied memo URL!', 'info', { timeout: 1500 })
  }
</script>

<div class="opacity-70">
  <button on:click={onClickCopyButton} class="btn btn-sm btn-outline rounded-full">
    <div class="w-[15px]">
      <CopyIcon />
    </div>
  </button>
</div>
