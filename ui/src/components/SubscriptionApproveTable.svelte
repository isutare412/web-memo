<script lang="ts">
  import type { Subscriber } from '$lib/apis/backend/memo'
  import { createEventDispatcher } from 'svelte'

  export let subscribers: Subscriber[]

  const dispatch = createEventDispatcher()

  function onCheckboxClick(userId: string): (
    event: MouseEvent & {
      currentTarget: EventTarget & HTMLInputElement
    }
  ) => void {
    const handler = (
      event: MouseEvent & {
        currentTarget: EventTarget & HTMLInputElement
      }
    ) => {
      dispatch('change', { userId, checked: event.currentTarget.checked })
    }

    return handler
  }
</script>

<div class="overflow-x-auto">
  <table class="table">
    <thead>
      <tr>
        <th>Photo</th>
        <th>Name</th>
        <th>Approve</th>
      </tr>
    </thead>
    <tbody>
      {#each subscribers as subscriber (subscriber.id)}
        <tr>
          <td>
            <div class="h-10 w-10">
              <img src={subscriber.photoUrl} alt="profile" class="mask mask-circle" />
            </div>
          </td>
          <td>
            <span class="text-lg">{subscriber.userName}</span>
          </td>
          <td>
            <input
              type="checkbox"
              class="toggle toggle-primary focus:outline-none"
              checked={subscriber.approved}
              on:click={onCheckboxClick(subscriber.id)}
            />
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
