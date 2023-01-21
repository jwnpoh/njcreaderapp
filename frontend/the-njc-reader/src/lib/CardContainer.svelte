<script>
  import Card from "$lib/Card.svelte";
  export let articles;
  export let page;
  export let query;
export let loggedIn;

  const navigate = (p) => {
    window.location.href = `/articles/${p}`;
  };
</script>

{#if query}
  <p class="px-10 py-5 italic font-medium text-xl">
    Showing results for: "{query}"
  </p>
{/if}

<div class="container mx-auto">
  <div class="py-5 px-10 grid md:grid-cols-2 gap-10">
    {#each articles as article}
      <div>
        <Card
          id={article.id}
          title={article.title}
          url={article.url}
          topics={article.topics}
          question_display={article.question_display}
          date={article.date}
          mustRead={article.must_read}
          {loggedIn}
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
</div>
