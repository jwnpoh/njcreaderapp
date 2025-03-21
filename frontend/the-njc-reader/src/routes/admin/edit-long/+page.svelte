<script>
  import PageTitle from "$lib/PageTitle.svelte";
  import { page } from "$app/stores";

  export let data;
  export let form;
  const articles = data.articles;

  let url;
  let title;
  let topic;
  let id;

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
</script>

<PageTitle>Edit long articles</PageTitle>
<div class="mx-auto px-5">
  <div class="px-5 pt-5 ">
    <a href="/admin" class="btn-link"
      ><p class="align-middle">&#8678; Back to admin dashboard</p></a
    >
  </div>
  <form method="POST" action="?/edit">
    <input name="id" type="text" class="hidden" hidden bind:value={id} />
    <div class="pt-5 px-5">
      <div class="pt-3">
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
    <div class="flex pt-5 px-5">
      <input
        name="title"
        type="text"
        placeholder="Article title (will be auto-populated after the URL is entered)"
        class="input w-full max-w-lg"
        bind:value={title}
      />
    </div>
    <div class="flex py-5 px-5">
      <input
        name="topic"
        type="text"
        placeholder="Topic"
        class="input w-full max-w-lg"
        bind:value={topic}
      />
    </div>

    <button class="btn btn-sm btn-primary mx-7">Save changes to article</button>
  </form>

  {#if form?.sent}
    <div class="mx-5 my-4 alert alert-success max-w-fit place-self-center">
      <span>Changes saved successfully.</span>
    </div>
  {/if}
  {#if form?.error}
    <div class="mx-5 my-4 alert alert-error max-w-fit place-self-center">
      <span class="text-center">{form?.message}</span>
    </div>
  {/if}

  <div class="divider py-3" />
  <div class="mx-auto pb-5">
    <p>Displaying all long articles. Select article to edit.</p>
  </div>
  <div>
    <table class="table table-compact w-full">
      <thead>
        <tr>
          <th>check</th>
          <th>Title</th>
          <th>Topic</th>
        </tr>
      </thead>
      <tbody>
        {#each articles as article}
          <tr>
            <th
              ><input
                type="radio"
                name="selection"
                value={article.id}
                on:input={() => {
                  id = article.id;
                  url = article.url;
                  title = article.title;
                  topic = article.topic;
                  form = {};
                  console.log(id);
                }}
              /></th
            >
            <td>{article.title}</td>
            <td>{article.topic}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
