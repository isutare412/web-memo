<script lang="ts">
  import { authStore, signInGoogle, signOut } from '$lib/auth'
  import userImage from '$media/user.png'

  $: user = $authStore.user
  $: photoUrl = user?.photoUrl ?? userImage
</script>

<div class="dropdown dropdown-end">
  <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
  <!-- svelte-ignore a11y-label-has-associated-control -->
  <label tabindex="0" class="btn btn-ghost btn-circle avatar">
    <div class="h-10 w-10 rounded-full">
      <img src={photoUrl} alt="profile" />
    </div>
  </label>
  <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
  <ul
    tabindex="0"
    class="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-4 w-36 p-2 shadow"
  >
    {#if user}
      <li><button on:click={signOut}>Sign Out</button></li>
    {:else}
      <li><button on:click={signInGoogle}>Sign In</button></li>
    {/if}
  </ul>
</div>
