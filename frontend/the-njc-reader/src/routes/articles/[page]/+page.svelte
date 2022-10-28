<script>
  import CardContainer from "$lib/CardContainer.svelte";

  export let data;
  const articles = data.articles;
  const page = data.page;
</script>

<CardContainer {articles} {page} />
