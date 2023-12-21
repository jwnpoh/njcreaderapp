<script>
  import PageTitle from "$lib/PageTitle.svelte"
  export let data;
  $: topics = data.topics;
</script>

<PageTitle>Longer Reads</PageTitle>
<div
  class="grid gap-4 md:grid-cols-4 justify-center place-content-center mt-6 md:px-24 lg:px-60"
>
  {#each topics as topic}
    <div class="card card-compact max-w-max bg-base-100 box-shadow text-center">
      <a href="/long/{topic}">
        <figure><img src="https://source.unsplash.com/random/300x200/?{topic}" alt="{topic}" /></figure>
        <div class="card-body text-lg">
          <p class="text-xl">{topic}</p>
        </div>
      </a>
    </div>
  {/each}
</div>

<style>
  .box-shadow {
    box-shadow: 1px 0 8px 0 rgb(0 0 0 / 0.1);
  }
</style>
