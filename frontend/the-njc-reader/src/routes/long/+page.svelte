<script>
  import PageTitle from "$lib/PageTitle.svelte";
  export let data;

  $: groupedArticles = data.groupedArticles;
  $: sortedTopics = Object.keys(groupedArticles).sort();
</script>

<PageTitle>Longer Reads by Topic</PageTitle>

<div class="max-w-4xl mx-auto mt-6 px-4">
  {#each sortedTopics as topic}
    <div
      class="collapse collapse-arrow mb-2 border border-base-300 bg-base-100 rounded-box"
    >
      <input type="checkbox" name="my-accordion-2" />
      <div class="collapse-title text-xl font-bold text-black">
        {topic}
      </div>
      <div class="collapse-content">
        <ul class="list-disc list-inside ml-4">
          {#each groupedArticles[topic] as article}
            <li
              class="text-lg py-1 hover:text-blue-600 transition duration-200"
            >
              <a href={article.url} target="_blank" rel="noreferrer"
                >{article.title}</a
              >
            </li>
          {/each}
        </ul>
      </div>
    </div>
  {/each}
</div>

<style>
  .box-shadow {
    box-shadow: 1px 0 8px 0 rgb(0 0 0 / 0.1);
  }
</style>
