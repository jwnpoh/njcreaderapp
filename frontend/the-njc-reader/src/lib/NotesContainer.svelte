<script>
  import Note from "$lib/Note.svelte";

  export let data;
  export let API_URL;
  export let section;
  const message = data.message;
  const notes = data.data;
</script>

<h2 class="py-2 text-lg font-semibold italic text-center">
  {section ? section : ""}
</h2>
<div class="px-5 md:px-10 py-5 grid gap-5 mb-20">
  {#if notes}
    {#each notes as note}
      <Note {note} {API_URL} />
    {/each}
  {:else if message}
    <p class="text-center">
      {message}
    </p>
  {/if}
</div>
