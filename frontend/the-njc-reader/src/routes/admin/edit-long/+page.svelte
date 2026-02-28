<!--
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

  // -------------------------------------------------------------------------
  // NEW: SEARCH STATE
  // Since ALL long articles are loaded on page open, we can just filter the already-loaded list
  // entirely in the browser — no API calls needed, and no debouncing required
  // since we're not hitting a server.
  // -------------------------------------------------------------------------

  // What the admin has typed in the search box
  let searchQuery = "";

  // -------------------------------------------------------------------------
  // NEW: REACTIVE FILTER
  // -------------------------------------------------------------------------
  $: displayedArticles = searchQuery.trim()
    ? articles.filter(
        (a) =>
          a.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
          a.topic.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : articles;

  const clearSearch = () => {
    searchQuery = "";
  };
</script>

<PageTitle>Edit long articles</PageTitle>
<div class="mx-auto px-5">
  <div class="px-5 pt-5">
    <a href="/admin" class="btn-link">
      <p class="align-middle">&#8678; Back to admin dashboard</p>
    </a>
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

  <div class="px-5 pb-3">
    <p class="pb-3 text-sm opacity-70">
      Filter articles by title or topic.
    </p>
    <div class="flex gap-2 max-w-lg">
      <input
        type="text"
        placeholder="e.g. globalisation, technology..."
        class="input input-bordered flex-1"
        bind:value={searchQuery}
      />
      {#if searchQuery}
        <button type="button" class="btn btn-ghost" on:click={clearSearch}>
          Clear
        </button>
      {/if}
    </div>

    // Show result count when a filter is active 
    {#if searchQuery.trim()}
      <p class="pt-2 text-sm opacity-70">
        {#if displayedArticles.length > 0}
          Showing {displayedArticles.length} result(s) for "{searchQuery}".
        {:else}
          No results found for "{searchQuery}".
        {/if}
      </p>
    {/if}
  </div>
  <div class="mx-auto pb-5 px-5">
    {#if !searchQuery.trim()}
      <p class="text-sm opacity-70">Displaying all long articles. Select article to edit.</p>
    {/if}
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
        {#each displayedArticles as article}
          <tr>
            <th>
              <input
                type="radio"
                name="selection"
                value={article.id}
                on:input={() => {
                  id = article.id;
                  url = article.url;
                  title = article.title;
                  topic = article.topic;
                  form = {};
                }}
              />
            </th>
            <td>{article.title}</td>
            <td>{article.topic}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
-->

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

  let searchQuery = "";

  $: displayedArticles = searchQuery.trim()
    ? articles.filter(
        (a) =>
          a.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
          a.topic.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : articles;

  $: searchStatus = searchQuery.trim()
    ? `Showing ${displayedArticles.length} of ${articles.length} article(s).`
    : "";

  const clearSearch = () => {
    searchQuery = "";
  };
</script>

<PageTitle>Edit long articles</PageTitle>
<div class="mx-auto px-5">
  <div class="px-5 pt-5">
    <a href="/admin" class="btn-link">
      <p class="align-middle">&#8678; Back to admin dashboard</p>
    </a>
  </div>

  <!-- ======================================================================
    Edit form — restyled to match add-articles layout
    Row 1: URL | Topic (equal columns)
    Row 2: Title (full width)
  ======================================================================= -->
  <form method="POST" action="?/edit">
    <input name="id" type="text" class="hidden" hidden bind:value={id} />

    <div class="card bg-base-100 shadow-sm border border-base-300 p-5 max-w-3xl mt-5 mb-4">

      <!-- Row 1: URL + Topic -->
      <div class="grid grid-cols-2 gap-3 mb-3">
        <div>
          <label class="label py-0 pb-1">
            <span class="label-text text-xs opacity-60">Article URL</span>
          </label>
          <input
            required
            name="url"
            type="url"
            placeholder="https://..."
            class="input input-bordered input-sm w-full"
            bind:value={url}
            on:input={() => getTitle(url)}
          />
        </div>
        <div>
          <label class="label py-0 pb-1">
            <span class="label-text text-xs opacity-60">Topic</span>
          </label>
          <input
            name="topic"
            type="text"
            placeholder="Topic"
            class="input input-bordered input-sm w-full"
            bind:value={topic}
          />
        </div>
      </div>

      <!-- Row 2: Title full width -->
      <div class="mb-4">
        <label class="label py-0 pb-1">
          <span class="label-text text-xs opacity-60">Title</span>
        </label>
        <input
          name="title"
          type="text"
          placeholder="Auto-filled from URL"
          class="input input-bordered input-sm w-full"
          bind:value={title}
        />
      </div>

      <button class="btn btn-primary btn-sm">Save changes to article</button>
    </div>
  </form>

  {#if form?.sent}
    <div class="alert alert-success max-w-3xl mt-3">
      <span>Changes saved successfully.</span>
    </div>
  {/if}
  {#if form?.error}
    <div class="alert alert-error max-w-3xl mt-3">
      <span>{form?.message}</span>
    </div>
  {/if}

  <div class="divider py-3" />

  <!-- ======================================================================
    Search bar — client-side filter, no API call needed
  ======================================================================= -->
  <div class="px-5 pb-3">
    <p class="pb-3 text-sm opacity-70">Filter articles by title or topic.</p>
    <div class="flex gap-2 max-w-lg">
      <input
        type="text"
        placeholder="e.g. geopolitics, environment..."
        class="input input-bordered flex-1"
        bind:value={searchQuery}
      />
      {#if searchQuery}
        <button type="button" class="btn btn-ghost" on:click={clearSearch}>
          Clear
        </button>
      {/if}
    </div>
    {#if searchStatus}
      <p class="pt-2 text-sm opacity-70">{searchStatus}</p>
    {/if}
  </div>

  <!-- ======================================================================
    Articles table
  ======================================================================= -->
  <div class="mx-auto pb-5 px-5">
    {#if !searchQuery.trim()}
      <p class="text-sm opacity-70">Displaying all long articles. Select article to edit.</p>
    {/if}
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
        {#each displayedArticles as article}
          <tr>
            <th>
              <input
                type="radio"
                name="selection"
                value={article.id}
                on:input={() => {
                  id = article.id;
                  url = article.url;
                  title = article.title;
                  topic = article.topic;
                  form = {};
                }}
              />
            </th>
            <td>{article.title}</td>
            <td>{article.topic}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
