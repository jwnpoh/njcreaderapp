<script>
  import { Icon } from "svelte-awesome";
  import book from "svelte-awesome/icons/book";
  import users from "svelte-awesome/icons/users";
  import feed from "svelte-awesome/icons/feed";

  import Container from "$lib/Container.svelte";
  import NotesContainer from "$lib/NotesContainer.svelte";

  let active = "notes";

  export let data;
  const notes = data.notes;
  const following = data.following;
  const discover = data.discover;
  const API_URL = data.API_URL;
  let liked_notes = data.liked_notes.data;
</script>

<Container>
  {#if active === "notes"}
    <NotesContainer data={notes ?? ""} bind:liked_notes {API_URL} />
  {/if}
  {#if active === "following"}
    <NotesContainer data={following ?? ""} bind:liked_notes {API_URL} />
  {/if}
  {#if active === "discover"}
    <NotesContainer data={discover ?? ""} bind:liked_notes {API_URL} />
  {/if}
</Container>

<div class="btm-nav">
  <button class:active={active === "notes"} on:click={() => (active = "notes")}
    ><Icon data={book} scale={2} /></button
  >
  <button
    class:active={active === "following"}
    on:click={() => (active = "following")}
  >
    <Icon data={users} scale={2} />
  </button>
  <button
    class:active={active === "discover"}
    on:click={() => (active = "discover")}
  >
    <Icon data={feed} scale={2} />
  </button>
</div>
