<script>
  import PageTitle from "$lib/PageTitle.svelte";
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
  if (url === "") {
    title = "";
    return
  }
    form = {};
    const payload = { url: url };
    const res = await fetch(`${data.API_URL}/api/admin/articles/get-title`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    title = response.data;
  };

</script>

<PageTitle>Add articles</PageTitle>
<div class="px-5 pt-5 ">
  <a href="/admin" class="btn-link"
    ><p class="align-middle">&#8678; Back to admin dashboard</p></a
  >
</div>
<form method="POST" action="?/send">
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

  <button class="btn btn-md btn-primary mx-6">Save article to database</button>
</form>
  {#if form?.sent}
    <div class="alert alert-success max-w-fit mx-6">
      <span>Article added successfully.</span>
    </div>
  {/if}
{#if form?.error}
    <div class="alert alert-error max-w-fit mx-6">
      <span class="text-center">{form?.message}</span>
    </div>
{/if}
