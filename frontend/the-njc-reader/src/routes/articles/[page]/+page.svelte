<script>
  import CardContainer from "$lib/CardContainer.svelte";

  export let data;
  $: console.log(data);
  $: articles = data.articles;
  $: page = data.page;
</script>

<CardContainer {articles} {page} />
