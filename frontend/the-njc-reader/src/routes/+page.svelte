<script>
  import { page } from "$app/stores";
  import CardContainer from "$lib/CardContainer.svelte";
  import PageTitle from "$lib/PageTitle.svelte";
  import Welcome from "$lib/Welcome.svelte"

  export let data;
  $: articles = data.articles;
  $: pageNo = data.page;

  let user = $page.data.user;
</script>

<PageTitle>News Feed</PageTitle>
{#if !user.loggedIn}
<Welcome/>
{/if}
<CardContainer {articles} page={pageNo} loggedIn={user.loggedIn}/>
