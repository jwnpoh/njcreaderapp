<script>
  import PageTitle from "$lib/PageTitle.svelte";
  export let data;
  $: articles = data.articles;
  $: topic = data.topic;
</script>

<PageTitle>Longer Reads - {topic}</PageTitle>
<div class="mt-6 px-10 md:px-20 lg:px-60 text-justify">
<ul class="list-decimal list-outside">
  {#each articles as article}
    <li class="min-w-full text-lg py-2">
      <a href="{article.url}" class="link" rel="noreferrer" target="_blank">{article.title}</a>
    </li>
  {/each}
</ul></div>

