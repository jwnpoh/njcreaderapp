<script>
  import { page } from "$app/stores";
  import { onDestroy } from "svelte";
  import Note from "$lib/Note.svelte";
  import SearchStore  from "$lib/stores/notesSearch"

  import Icon from "svelte-awesome";
  import close from "svelte-awesome/icons/close";

  export let data;
  export let API_URL;
  export let section;
  const message = data.message;
  let notes = data.data;
  const user_id = $page.data.user.id;

  let searchTerm = "";

  const unsubscribe = SearchStore.subscribe(term => {
    searchTerm = term;
    console.log(searchTerm);
  });

  const removeSearchTerm = () => {
    SearchStore.set("");
  }

  onDestroy(() => {
    unsubscribe();
  })

</script>

<h2 class="py-2 text-lg font-semibold italic text-center">
  {section ? section : ""}
</h2>
{#if searchTerm}
      <p class="text-center text-lg">
      Showing notes tagged <badge class="mr-1 badge badge-outline hover:cursor-default" >{searchTerm} <button on:click={removeSearchTerm}><Icon data={close} style="padding-left: 3px;" /></button></badge>
      </p>
      {/if}
<div class="px-5 md:px-10 py-5 grid gap-5 mb-20">
  {#if notes}
    {#each notes as note}
      {#if searchTerm}
        {#if note.tags.includes(searchTerm)}
          {#if section !== "My notes"}
            {#if note.user_id !== user_id}
              <Note {note} {API_URL} />
            {/if}
          {:else}
            <Note {note} {API_URL} />
          {/if}
        {/if}
      {:else}
          {#if section !== "My notes"}
            {#if note.user_id !== user_id}
              <Note {note} {API_URL} />
            {/if}
          {:else}
            <Note {note} {API_URL} />
          {/if}
      {/if}
    {/each}
  {:else if message}
    <p class="text-center">
      {message}
    </p>
  {/if}
</div>
