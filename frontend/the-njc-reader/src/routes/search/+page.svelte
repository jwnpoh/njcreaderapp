<script>
  import { page } from "$app/stores";
  import CardContainer from "$lib/CardContainer.svelte";
  import PageTitle from "$lib/PageTitle.svelte";

  let user = $page.data.user;

  export let data;
  const error = data.error;
  const message = data.message;
  const articles = data.articles;
  const query = data.query;
</script>

<PageTitle>Search</PageTitle>
{#if !error}
  <CardContainer {articles} {query} loggedIn={user.loggedIn}/>
{:else}
  {message}
{/if}
