<script>
  import PageTitle from "$lib/PageTitle.svelte"
  export let data;
  $: topics = data.topics;
</script>

<PageTitle>Longer Reads</PageTitle>
<div
  class="grid gap-6 md:grid-cols-2 justify-center place-content-center mt-6 md:px-24 lg:px-60"
>
  {#each topics as topic}
    <div class="card bg-base-100 box-shadow min-w-max text-center">
      <a href="/long/{topic}">
        <div class="card-body text-lg">
          {topic}
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
