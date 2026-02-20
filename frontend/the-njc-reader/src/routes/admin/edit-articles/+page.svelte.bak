<script>
  import PageTitle from "$lib/PageTitle.svelte";
  import { page } from "$app/stores";
  import { DateInput } from "date-picker-svelte";

  export let data;
  export let form;
  const articles = data.articles;

  let url;
  let title;
  let tags;
  let date;
  let must_read;

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

<PageTitle>Edit articles</PageTitle>
<div class="mx-auto px-5">
  <div class="px-5 pt-5 ">
    <a href="/admin" class="btn-link"
      ><p class="align-middle">&#8678; Back to admin dashboard</p></a
    >
  </div>
  <form method="POST" action="?/edit">
    <input name="id" type="text" class="hidden" hidden bind:value={id} />
    <div class="flex pt-5 px-5">
      <div class="flex-auto">
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
      <div class="flex-auto">
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
    <div class="flex pt-5 px-5">
      <input
        name="title"
        type="text"
        placeholder="Article title (will be auto-populated after the URL is entered)"
        class="input w-full max-w-100"
        bind:value={title}
      />
    </div>
    <div class="flex py-5 px-5">
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
      <input
        name="must_read"
        bind:checked={must_read}
        type="checkbox"
        class="checkbox"
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
    <p>
      Displaying 100 most recent articles. Select article to edit. To edit older
      articles, contact database admin.
    </p>
  </div>
  <div>
    <table class="table table-compact w-full">
      <thead>
        <tr>
          <th>check</th>
          <th>Date</th>
          <th>Title</th>
          <th>Topics</th>
          <th>Questions</th>
          <th>Must read?</th>
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
                  tags =
                    article.topics.join(";") +
                    "; " +
                    article.questions.join(";");
                  must_read = article.must_read;
                  date = new Date(article.date);
                  form = {};
                }}
              /></th
            >
            <td>{article.date}</td>
            <td>{article.title}</td>
            <td>{article.topics}</td>
            <td>{article.questions}</td>
            <td>{article.must_read}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
