<script>
  import CardContainer from "$lib/CardContainer.svelte";
  import PageTitle from "$lib/PageTitle.svelte";
  export let data;
  const error = data.error;
  const message = data.message;
  const articles = data.articles;
  const query = data.query;
</script>

<PageTitle>Search</PageTitle>
{#if !error}
  <CardContainer {articles} {query} />
{:else}
  {message}
{/if}
