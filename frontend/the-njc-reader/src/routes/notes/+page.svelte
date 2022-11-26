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
</script>

<Container>
  {#if active === "notes"}
    <NotesContainer data={notes ?? ""} />
  {/if}
  {#if active === "following"}
    <NotesContainer data={following ?? ""} />
  {/if}
  {#if active === "discover"}
    <NotesContainer data={discover ?? ""} />
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
