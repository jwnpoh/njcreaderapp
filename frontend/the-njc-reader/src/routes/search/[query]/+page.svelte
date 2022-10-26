<script>
  import CardContainer from "$lib/CardContainer.svelte";
  export let data;
  const { articles } = data;
</script>

<CardContainer {articles} />
