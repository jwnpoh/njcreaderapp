<script>
  import { invalidateAll } from "$app/navigation";
  import { Icon } from "svelte-awesome";
  import book from "svelte-awesome/icons/book";
  import feed from "svelte-awesome/icons/feed";
  import globe from "svelte-awesome/icons/globe";

  import Container from "$lib/Container.svelte";
  import NotesContainer from "$lib/NotesContainer.svelte";
  import PageTitle from "$lib/PageTitle.svelte";

  let active = "notes";

  export let data;
  $: notes = data.notes;
  $: following = data.following;
  $: discover = data.discover;

  const API_URL = data.API_URL;

  const refresh = () => {
    invalidateAll();
  }
</script>

<Container>
  <PageTitle>The Social Notebook</PageTitle>
  {#if active === "notes"}
    <NotesContainer data={notes ?? ""} {API_URL} section={"My notes"} />
  {/if}
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
</Container>

<button class="btm-nav" on:click={refresh}>
  <button class:active={active === "notes"} on:click={() => (active = "notes")}
    ><Icon data={book} scale={1.6} />
    <p class="text-xs md:text-sm">My Notes</p></button
  >
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
</button>
