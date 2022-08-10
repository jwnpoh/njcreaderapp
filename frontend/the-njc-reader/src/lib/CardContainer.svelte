<script>
  import Card from "$lib/Card.svelte";
  export let articles;
</script>

<div class="md:px-80 ">
  {#each articles as article}
    <div class=" px-5 py-6">
      <Card
        title={article.title}
        url={article.url}
        topics={article.topics}
        questions={article.questions}
        date={article.date}
      />
    </div>
  {/each}
</div>
