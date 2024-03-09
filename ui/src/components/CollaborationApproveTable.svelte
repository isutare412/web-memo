<script lang="ts">
  import type { Collaborator } from '$lib/apis/backend/memo'
  import { createEventDispatcher } from 'svelte'

  export let collaborators: Collaborator[]

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
      {#each collaborators as collaborator (collaborator.id)}
        <tr>
          <td>
            <div class="h-10 w-10">
              <img src={collaborator.photoUrl} alt="profile" class="mask mask-circle" />
            </div>
          </td>
          <td>
            <span class="text-lg">{collaborator.userName}</span>
          </td>
          <td>
            <input
              type="checkbox"
              class="toggle toggle-primary focus:outline-none"
              checked={collaborator.isApproved}
              on:click={onCheckboxClick(collaborator.id)}
            />
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
