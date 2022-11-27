<script>
  import { page } from "$app/stores";
  import { DateInput } from "date-picker-svelte";

  export let data;
  export let form;

  let url = form?.url ?? "";
  let title = form?.title ?? "";
  let tags = form?.tags ?? "";
  let date = new Date();

  const session = $page.data.user.session;
  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  const getTitle = async (url) => {
    const payload = { url: url };
    const res = await fetch(`${data.API_URL}/api/admin/articles/get-title`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    title = response.data;
  };

  let queue = data.queue;
  let len = queue.length;
</script>

<div class="px-5 pt-5 ">
  <a href="/admin" class="btn-link"
    ><p class="align-middle">&#8678; Back to admin dashboard</p></a
  >
</div>
<form method="POST" action="?/queue">
  <div class="flex pt-5 px-5">
    <div class="flex-auto basis-7/12">
      <div class="">
        <input
          required
          name="url"
          type="url"
          placeholder="Article URL"
          class="input w-full max-w-lg"
          bind:value={url}
          on:input={getTitle(url)}
        />
      </div>
    </div>

    <div class="flex-auto  ">
      <div>
        <label for="date">Published on</label>
        <input type="text" name="date" bind:value={date} hidden />
        <DateInput
          placeholder="dd-MM-yyyy"
          format="dd-MM-yyyy"
          closeOnSelection
          bind:value={date}
        />
      </div>
    </div>
  </div>
  <div class="flex px-5">
    <input
      name="title"
      type="text"
      placeholder="Article title"
      class="input w-full max-w-100"
      bind:value={title}
    />
  </div>
  <div class="flex pt-5 px-5">
    <input
      name="tags"
      type="text"
      placeholder="Topic and question tags. Separate each tag with a semicolon (e.g. 2019-Q6; leadership)"
      class="input w-full max-w-100"
      bind:value={tags}
    />
  </div>

  <div class="flex py-5 px-5">
    <label class="px-2" for="must_read">Must read?</label>
    <input name="must_read" type="checkbox" class="checkbox" />
  </div>

  <button class="btn btn-sm btn-secondary mx-7">Add to queue</button>
</form>
{#if form?.error}
  <p class="mx-7 pt-7 text-primary">{form?.message}</p>
{/if}

<div class="divider py-5 font-bold">Queued articles ({len})</div>

<div class="px-2 grid gap-10 relative">
  <div class="max-w-sm place-self-center">
    <form method="POST" action="?/send">
      <button class="btn btn-md btn-primary">Add to database</button>
    </form>
  </div>
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
  {#each queue as item}
    <form method="POST" action="?/remove">
      <input type="hidden" name="index" value={item.index} />
      <div
        class="card bg-base-100 shadow-lg w-4/5 mx-auto"
        class:bg-secondary={item.must_read}
      >
        <div class="card-body pb-5">
          <h2 class="card-title">
            <a href={item.url} rel="noreferrer" target="_blank">{item.title}</a>
          </h2>
          <p>{item.date}</p>
          <div />
          <div class="inline">
            tags: {item.tags}
          </div>
        </div>
        <button class="btn btn-x btn-xs btn-circle submit">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            /></svg
          >
        </button>
      </div>
    </form>
  {/each}
</div>

<style>
  .btn-x {
    position: absolute;
    top: 0%;
    left: 0.1%;
  }
</style>
