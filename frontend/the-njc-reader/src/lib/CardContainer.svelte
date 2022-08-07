<script>
  import Card from "$lib/Card.svelte";
  export let articles;
</script>

<div class="grid grid-cols-2 gap-y-8">
  {#each articles as article}
    <div class="basis-1/2 px-5">
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
