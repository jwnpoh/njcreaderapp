<script>
  export let data;
  $: articles = data.articles;
  $: topic = data.topic;
</script>

{#each articles as article}
  {article.title}
  {article.url}
{/each}
{topic}
