<script>
  import { page } from "$app/stores";
  import CardContainer from "$lib/CardContainer.svelte";
  import PageTitle from "$lib/PageTitle.svelte";

  export let data;
  $: articles = data.articles;
  $: pageNo = data.page;

  let user = $page.data.user;
</script>

<PageTitle>News Feed</PageTitle>
<CardContainer {articles} page={pageNo} loggedIn={user.loggedIn} />
