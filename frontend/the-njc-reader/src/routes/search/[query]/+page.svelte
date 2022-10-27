<script>
  import CardContainer from "$lib/CardContainer.svelte";
  export let data;
  const error = data.error;
  const message = data.message;
  console.log(message);
  const articles = data.articles;
  const query = data.query;
</script>

<p class="py-4 italic font-medium text-xl">Showing results for: "{query}"</p>

{#if !error}
  <CardContainer {articles} />
{:else}
  {message}
{/if}
