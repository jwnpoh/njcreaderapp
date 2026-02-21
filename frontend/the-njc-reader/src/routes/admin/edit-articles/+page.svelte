<script>
  import PageTitle from "$lib/PageTitle.svelte";
  import { page } from "$app/stores";
  import { DateInput } from "date-picker-svelte";

  export let data;
  export let form;

  const articles = data.articles; // the 100 most recent articles, loaded on page open

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

  // -------------------------------------------------------------------------
  // NEW: SEARCH STATE
  // These variables manage the search bar and its results.
  // -------------------------------------------------------------------------

  // What the admin has typed in the search box
  let searchQuery = "";

  // The list of articles currently shown in the table.
  // Starts as the default 100 recent articles, and gets replaced with
  // search results when the admin performs a search.
  let displayedArticles = articles;

  // Controls the status message shown below the search bar
  // e.g. "Searching...", "No results found", "Showing results for: X"
  let searchStatus = "";

  // True while waiting for the API to respond, so we can show a loading state
  let isSearching = false;

  // -------------------------------------------------------------------------
  // NEW: DEBOUNCE TIMER
  // This holds a reference to a pending timer. Each time the user types,
  // we cancel the previous timer and start a new one. The API call only
  // fires after the user has stopped typing for 400ms.
  // Without this, every single keystroke would trigger a separate API call.
  // -------------------------------------------------------------------------
  let debounceTimer;

  // -------------------------------------------------------------------------
  // NEW: SEARCH FUNCTION
  // -------------------------------------------------------------------------
  const handleSearch = async (query) => {
    // If the search box is empty, just reset back to the default 100 articles
    if (!query.trim() || query.trim().length < 3) {
      displayedArticles = articles;
      searchStatus = "";
      return;
    }

    isSearching = true;
    searchStatus = "Searching...";

    try {
      // This is the same endpoint the public search page uses.
      // No auth header needed since this endpoint is public.
      const res = await fetch(`${data.API_URL}/api/articles/find?term=${encodeURIComponent(query)}`);
      const response = await res.json();

      if (response.error || !response.data || response.data.length === 0) {
        // No results found — clear the table and show a message
        displayedArticles = [];
        searchStatus = `No results found for "${query}".`;
      } else {
        // Results found — replace the table contents with search results
        displayedArticles = response.data;
        searchStatus = `Showing ${response.data.length} result(s) for "${query}".`;
      }
    } catch (e) {
      // Something went wrong with the network request
      displayedArticles = [];
      searchStatus = "Search failed. Please try again.";
    }

    isSearching = false;
  };

  // -------------------------------------------------------------------------
  // NEW: REACTIVE DEBOUNCED TRIGGER
  // -------------------------------------------------------------------------
  $: {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => handleSearch(searchQuery), 400);
  }

  // -------------------------------------------------------------------------
  // NEW: CLEAR SEARCH
  // Resets everything back to the default state (100 most recent articles).
  // Setting searchQuery to "" will trigger the reactive block above,
  // which will call handleSearch("") and reset displayedArticles.
  // -------------------------------------------------------------------------
  const clearSearch = () => {
    searchQuery = "";
  };
</script>

<PageTitle>Edit articles</PageTitle>
<div class="mx-auto px-5">
  <div class="px-5 pt-5">
    <a href="/admin" class="btn-link">
      <p class="align-middle">&#8678; Back to admin dashboard</p>
    </a>
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

  <!-- ======================================================================
    NEW: Search bar
  ======================================================================= -->
  <div class="px-5 pb-3">
    <p class="pb-3 text-sm opacity-70">
      Search across all articles by title, topic, or question tag. Or scroll
      through the 100 most recent articles below.
    </p>

    <div class="flex gap-2 max-w-lg">
      <input
        type="text"
        placeholder="e.g. climate change, 2022-Q3, globalisation..."
        class="input input-bordered flex-1"
        bind:value={searchQuery}
      />
      <!-- Only show the Clear button once the admin has typed something -->
      {#if searchQuery}
        <button type="button" class="btn btn-ghost" on:click={clearSearch}>
          Clear
        </button>
      {/if}
    </div>

    <!-- Status message: shows result count, "no results", loading, or errors -->
    {#if searchStatus}
      <p class="pt-2 text-sm opacity-70">{searchStatus}</p>
    {/if}
  </div>

  <div class="mx-auto pb-5 px-5">
    {#if !searchStatus}
      <!-- Default state: remind admin this is just the 100 most recent -->
      <p class="text-sm opacity-70">
        Displaying 100 most recent articles. Select article to edit.
      </p>
    {/if}
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
        {#each displayedArticles as article}
          <tr>
            <th>
              <input
                type="radio"
                name="selection"
                value={article.id}
                on:input={() => {
                  // When a row is selected, populate the edit form fields above
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
              />
            </th>
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
