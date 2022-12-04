<script>
  import Note from "$lib/Note.svelte";

  export let data;
  export let API_URL;
  const notes = data.data;
  const message = data.message;
</script>

{#if notes}
  <div class="px-5 md:px-10 py-5 grid gap-5 mb-20">
    {#each notes as note}
      <Note {note} {API_URL} />
    {/each}
  </div>
{:else if message}
  {message}
{/if}
