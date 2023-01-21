<script>
  import PageTitle from "$lib/PageTitle.svelte";
  export let data;
  $: articles = data.articles;
  $: topic = data.topic;
</script>

<PageTitle>Longer Reads - {topic}</PageTitle>
<div class="grid gap-6  mt-6 md:px-60">
  {#each articles as article, index}
    <div class="min-w-full text-lg">
      {index+1}. <a href="{article.url}" class="link" rel="noreferrer" target="_blank">{article.title}</a>
    </div>
  {/each}
</div>

