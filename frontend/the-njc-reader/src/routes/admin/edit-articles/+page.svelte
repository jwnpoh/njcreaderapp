<script>
  import { DateInput } from "date-picker-svelte";

  export let data;
  export let form;
  const articles = data.articles;

  let url;
  let title;
  let tags;
  let date = new Date();

  let id;

  const getTitle = async (url) => {
    const payload = { url: url };
    const res = await fetch(
      "http://localhost:8080/api/admin/articles/get-title",
      {
        method: "POST",
        body: JSON.stringify(payload),
        headers: myHeaders,
      }
    );

    const response = await res.json();
    title = response.data;
  };
</script>

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

  <button class="btn btn-sm btn-primary mx-7">Save changes to article</button>
</form>

{#if form?.error}
  <p class="mx-7 pt-7 text-primary">{form?.message}</p>
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
        <th>Title</th>
        <th>Topics</th>
        <th>Questions</th>
        <th>Date</th>
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
                  article.topics.join(";") + "; " + article.questions.join(";");
              }}
            /></th
          >
          <td>{article.title}</td>
          <td>{article.topics}</td>
          <td>{article.questions}</td>
          <td>{article.date}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
