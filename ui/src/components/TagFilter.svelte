<script lang="ts">
  import Tag from '$components/Tag.svelte'
  import Funnel from '$components/icons/Funnel.svelte'
  import { insertTagFilter, memoStore, removeTagFilter } from '$lib/memo'

  let value: string = ''
  $: tags = $memoStore.selectedTags

  function addTag(event: KeyboardEvent & { currentTarget: EventTarget & HTMLInputElement }) {
    if (value === '' || event.key !== 'Enter') return

    insertTagFilter(value)
    value = ''
    event.currentTarget.blur()
  }

  function removeTag(event: CustomEvent<{ name: string }>) {
    removeTagFilter(event.detail.name)
  }
</script>

<div>
  <div class="flex items-center gap-x-2">
    <div class="w-4">
      <Funnel />
    </div>
    <input
      type="text"
      placeholder="Tag"
      maxlength="20"
      bind:value
      on:keyup={addTag}
      class="input input-sm input-bordered focus:border-primary h-7 w-full max-w-[200px] focus:outline-none"
    />
  </div>
  <div class="mb-3 mt-2 flex flex-1 flex-wrap gap-1">
    {#each tags as tag (tag)}
      <Tag value={tag} color={'secondary'} isClose={true} on:click={removeTag} />
    {/each}
  </div>
</div>
