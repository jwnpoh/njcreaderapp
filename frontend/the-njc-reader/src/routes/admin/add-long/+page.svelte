<script>
  import PageTitle from "$lib/PageTitle.svelte";
  export let form;
  export let data;
  $: topics = data.topics;
</script>

<PageTitle>Add long articles</PageTitle>
<div class="px-5 pt-5 ">
  <a href="/admin" class="btn-link"
    ><p class="align-middle">&#8678; Back to admin dashboard</p></a
  >
</div>

<form method="POST" action="?/add">
  <div class="pt-4 px-10">
    <h2 class="text-lg font-semibold">Add long reads</h2>
    <p class="py-3 text-md  text-justify opacity-80">
      Enter each article on a new line, separating the url and the topic with a
      semicolon.
    </p>
    <textarea
      required
      name="input"
      type="text"
      placeholder="https://example.com/article123-456 ; Topic
https://example2.com/article2; topic 2"
      class="textarea px-2 py-2 w-screen max-w-full bg-secondary bg-opacity-5"
    />
    <label for="tldr" class="label">
      <span class="label-text-alt" />
    </label>
    <button class="btn btn-md btn-primary">Add to database</button>
    {#if form?.sent}
      <div class="alert alert-success max-w-fit place-self-center">
        <span>Articles added successfully.</span>
      </div>
    {/if}
    {#if form?.failed}
      <div class="alert alert-error max-w-fit place-self-center">
        <span class="text-center">{form?.message}</span>
      </div>
    {/if}
  </div>
</form>

<br/>

<details class="dropdown mt-0">
  <summary class="collapse-title text-lg px-5">Show existing topics</summary>
  <ul class="mx-5 py-2 text-md shadow menu dropdown-content z-[1] bg-base-100 rounded-box w-60">
  {#each topics as topic}
    <li class="px-5">{topic}</li>
  {/each}
  </ul>
</details>
