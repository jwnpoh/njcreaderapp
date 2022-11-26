<script>
  import Note from "$lib/Note.svelte";

  export let data;
  const notes = data.data;
  const message = data.message;
</script>

{#if notes}
  <div class="px-5 md:px-10 py-5 grid gap-5 mb-20">
    {#each notes as note}
      <Note {note} />
    {/each}
  </div>
{:else if message}
  {message}
{/if}
