<script>
  import Card from "$lib/Card.svelte";
  export let articles;
  export let page;

  const navigate = (p) => {
    window.location.href = `/articles/${p}`;
  };
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
{#if page}
  <div class="flex justify-center pt-10">
    {#if page > 1}
      <div class="flex-initial btn-group grid grid-cols-2 ">
        <button
          class="btn btn-outline"
          on:click={() => {
            page--;
            navigate(page);
          }}>Previous page</button
        >
        <button
          class="btn btn-outline"
          on:click={() => {
            page++;
            navigate(page);
          }}>Next page</button
        >
      </div>
    {:else}
      <div class="flex-initial btn-group grid ">
        <button
          class="btn btn-outline"
          on:click={() => {
            page++;
            navigate(page);
          }}>Next page</button
        >
      </div>
    {/if}
  </div>
{/if}
