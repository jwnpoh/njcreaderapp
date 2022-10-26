<script>
  import Card from "$lib/Card.svelte";
  export let articles;
</script>

<div class="px-10 grid md:grid-cols-2 gap-10">
  {#each articles as article}
    <div>
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
