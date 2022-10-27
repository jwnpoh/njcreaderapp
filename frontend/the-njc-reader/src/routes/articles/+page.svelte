<script>
  import CardContainer from "$lib/CardContainer.svelte";
  export let data;
  const articles = data.articles;
  const page = 1;
</script>

<CardContainer {articles} {page} />
