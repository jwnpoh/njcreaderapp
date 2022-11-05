<script>
  import CardContainer from "$lib/CardContainer.svelte";
  export let data;
  const error = data.error;
  const message = data.message;
  const articles = data.articles;
  const query = data.query;
</script>

{#if !error}
  <CardContainer {articles} {query} />
{:else}
  {message}
{/if}
