<script>
  import { page } from "$app/stores";
  import Note from "$lib/Note.svelte";

  export let data;
  export let API_URL;
  export let section;
  const message = data.message;
  const notes = data.data;
  const user_id = $page.data.user.id;
</script>

<h2 class="py-2 text-lg font-semibold italic text-center">
  {section ? section : ""}
</h2>
<div class="px-5 md:px-10 py-5 grid gap-5 mb-20">
  {#if notes}
    {#each notes as note}
      {#if section !== "My notes"}
        {#if note.user_id !== user_id}
          <Note {note} {API_URL} />
        {/if}
      {:else}
          <Note {note} {API_URL} />
      {/if}
    {/each}
  {:else if message}
    <p class="text-center">
      {message}
    </p>
  {/if}
</div>
