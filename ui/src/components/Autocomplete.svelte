<script lang="ts">
  import { mod } from '$lib/utils/math'
  import { createEventDispatcher } from 'svelte'

  export let items: string[]
  export let selectedItem: string | undefined = undefined
  export let selectedIndex: number | undefined = undefined
  let itemElements: HTMLLIElement[] = []

  const dispatch = createEventDispatcher()

  $: {
    items
    selectedItem = undefined
    selectedIndex = undefined
  }

  function onWindowKeyDown(
    event: KeyboardEvent & {
      currentTarget: EventTarget & Window
    }
  ) {
    switch (event.key) {
      case 'ArrowUp':
        event.preventDefault()
        increaseSelectedIndex()
        break
      case 'ArrowDown':
        event.preventDefault()
        decreaseSelectedIndex()
        break
    }
  }

  function onWindowKeyUp(
    event: KeyboardEvent & {
      currentTarget: EventTarget & Window
    }
  ) {
    switch (event.key) {
      case 'Enter':
        dispatchSelectedItem()
        break
    }
  }

  function dispatchSelectedItem() {
    if (selectedItem === undefined) return

    dispatch('select', {
      item: selectedItem,
    })
  }

  function increaseSelectedIndex() {
    selectedIndex =
      selectedIndex === undefined ? items.length - 1 : mod(selectedIndex - 1, items.length)
    scrollSelectedIntoView()
    selectedItem = items[selectedIndex]
  }

  function decreaseSelectedIndex() {
    selectedIndex = selectedIndex === undefined ? 0 : mod(selectedIndex + 1, items.length)
    scrollSelectedIntoView()
  }

  function scrollSelectedIntoView() {
    if (selectedIndex === undefined || !itemElements[selectedIndex]) return
    itemElements[selectedIndex].scrollIntoView({ behavior: 'instant', block: 'nearest' })
    selectedItem = items[selectedIndex]
  }
</script>

<svelte:window on:keydown={onWindowKeyDown} on:keyup={onWindowKeyUp} />

{#if items.length > 0}
  <div class="relative">
    <ul class="bg-base-100 absolute top-1 max-h-60 overflow-y-scroll rounded-md p-2 shadow-lg">
      {#each items as option, i (option)}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
        <li
          bind:this={itemElements[i]}
          on:mouseenter={() => {
            selectedIndex = i
            selectedItem = option
          }}
          on:click={dispatchSelectedItem}
          class="cursor-pointer rounded px-2 py-1"
          class:bg-primary={i === selectedIndex}
          class:text-primary-content={i === selectedIndex}
        >
          {option}
        </li>
      {/each}
    </ul>
  </div>
{/if}
