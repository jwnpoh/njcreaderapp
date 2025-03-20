<script>
  import PageTitle from "$lib/PageTitle.svelte"
  export let data;
  // $: topics = data.topics;

  let topics = data.topics;
  let grouped = topics.reduce((acc, topic) => {
    let firstLetter = topic[0].toUpperCase();
    if (!acc[firstLetter]) {
      acc[firstLetter] = [];
    }
    acc[firstLetter].push(topic);
    return acc;
  }, {});
</script>

<PageTitle>Longer Reads by Topic</PageTitle>


  <!--
<div
  class="grid gap-4 md:grid-cols-4 justify-center place-content-center mt-6 md:px-24 lg:px-60"
>
-->
<ul class="max-w-4xl mx-auto text-left grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10">
  {#each Object.keys(grouped).sort() as letter}
    <div>
    <li class="font-bold text-xl mt-8 border-b border-primary pb-1">{letter}</li>
    <ul class="ml-4 list-disc list-inside">
      {#each grouped[letter] as item}
        <li class="text-xl py-1 hover:text-blue-600 hover:font-semibold transition duration-200 cursor-pointer"><a href="/long/{item}">{item}</a></li>
      {/each}
    </ul>
    </div>
  {/each}
</ul>

<!--
  {#each topics as topic}
    <div class=" bg-base-100 text-center">
      <a href="/long/{topic}">
        <div class="card-body text-lg">
          <p class="text-xl">{topic}</p>
        </div>
      </a>
    </div>
  {/each}
  -->

<style>
  .box-shadow {
    box-shadow: 1px 0 8px 0 rgb(0 0 0 / 0.1);
  }
</style>
