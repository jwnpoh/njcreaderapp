<script>
  import { Icon } from "svelte-awesome";
  import book from "svelte-awesome/icons/book";
  import users from "svelte-awesome/icons/users";
  import feed from "svelte-awesome/icons/feed";
  import globe from "svelte-awesome/icons/globe";

  import Container from "$lib/Container.svelte";
  import NotesContainer from "$lib/NotesContainer.svelte";
  import PageTitle from "$lib/PageTitle.svelte";

  let active = "following";

  export let data;
  const notes = data.notes;
  const following = data.following;
  const discover = data.discover;
  const API_URL = data.API_URL;
</script>

<Container>
  <PageTitle>Noteworthy news</PageTitle>
  {#if active === "following"}
    <NotesContainer
      data={following ?? ""}
      {API_URL}
      section={"Notes from people I follow"}
    />
  {/if}
  {#if active === "discover"}
    <NotesContainer
      data={discover ?? ""}
      {API_URL}
      section={"Notes from everyone"}
    />
  {/if}
  {#if active === "notes"}
    <NotesContainer data={notes ?? ""} {API_URL} section={"My notes"} />
  {/if}
</Container>

<div class="btm-nav">
  <button
    class:active={active === "following"}
    on:click={() => (active = "following")}
  >
    <Icon data={feed} scale={1.6} />
    <p class="text-xs md:text-sm">Following</p>
  </button>
  <button
    class:active={active === "discover"}
    on:click={() => (active = "discover")}
  >
    <Icon data={globe} scale={1.6} />
    <p class="text-xs md:text-sm">Discover</p>
  </button>
  <button class:active={active === "notes"} on:click={() => (active = "notes")}
    ><Icon data={book} scale={1.6} />
    <p class="text-xs md:text-sm">My Notes</p></button
  >
  <button
    class:active={active === "activity"}
    on:click={() => (active = "activity")}
  >
    <div class="indicator">
      <span
        class="indicator-item badge badge-xs badge-accent w-fit translate-x-8 md:translate-x-10"
        >new</span
      >
      <Icon data={users} scale={1.6} />
    </div>
    <p class="text-xs md:text-sm">Activity</p>
  </button>
</div>
