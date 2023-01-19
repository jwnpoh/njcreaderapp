<script>
  import CardContainer from "$lib/CardContainer.svelte";
  import PageTitle from "$lib/PageTitle.svelte";

  export let data;
  $: articles = data.articles;
  $: page = data.page;
</script>

<PageTitle>News Feed</PageTitle>
<CardContainer {articles} {page} />
