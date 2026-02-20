<script>
  import PageTitle from "$lib/PageTitle.svelte";

  export let data;
  export let form;

  const articles = data.articles;

  // -------------------------------------------------------------------------
  // NEW: SEARCH STATE
  // -------------------------------------------------------------------------

  // What the admin has typed in the search box
  let searchQuery = "";

  // The list of articles currently shown in the table. Starts as the default
  // 100 recent articles, replaced with search results when a search runs.
  let displayedArticles = articles;

  // Status message shown below the search bar
  let searchStatus = "";

  // True while waiting for the API to respond
  let isSearching = false;

  // -------------------------------------------------------------------------
  // IMPORTANT: TRACKING CHECKED ARTICLES ACROSS SEARCHES
  //
  // If we just relied on the checkboxes in the DOM, any checked items would
  // disappear from view when search results replace the table — and since
  // they're no longer rendered, they wouldn't be included in the form POST.
  //
  // The solution: maintain a Set of selected article IDs in JS. Checkboxes
  // update this Set when ticked/unticked. The form submission reads from
  // this Set, not from the DOM checkboxes directly.
  // -------------------------------------------------------------------------
  let selectedIds = new Set();

  const toggleSelection = (id, checked) => {
    // When a checkbox is ticked, add the article ID to our Set.
    // When unticked, remove it. Either way, reassign selectedIds so
    // Svelte knows the value changed and re-renders anything that depends on it.
    if (checked) {
      selectedIds.add(id);
    } else {
      selectedIds.delete(id);
    }
    selectedIds = selectedIds; // trigger Svelte reactivity
  };

  // -------------------------------------------------------------------------
  // NEW: DEBOUNCE TIMER
  // Prevents an API call on every single keystroke. The search only fires
  // after the user has stopped typing for 400ms.
  // -------------------------------------------------------------------------
  let debounceTimer;

  // -------------------------------------------------------------------------
  // NEW: SEARCH FUNCTION
  // Called automatically via the reactive block below whenever searchQuery changes.
  // -------------------------------------------------------------------------
  const handleSearch = async (query) => {
    if (!query.trim()) {
      // Empty search box — reset to default 100 articles
      displayedArticles = articles;
      searchStatus = "";
      return;
    }

    isSearching = true;
    searchStatus = "Searching...";

    try {
      const res = await fetch(`${data.API_URL}/api/articles/find?term=${encodeURIComponent(query)}`);
      const response = await res.json();

      if (response.error || !response.data || response.data.length === 0) {
        displayedArticles = [];
        searchStatus = `No results found for "${query}".`;
      } else {
        displayedArticles = response.data;
        searchStatus = `Showing ${response.data.length} result(s) for "${query}".`;
      }
    } catch (e) {
      displayedArticles = [];
      searchStatus = "Search failed. Please try again.";
    }

    isSearching = false;
  };

  // -------------------------------------------------------------------------
  // NEW: REACTIVE DEBOUNCED TRIGGER
  // Runs automatically whenever searchQuery changes (i.e. on each keystroke),
  // but waits 400ms of inactivity before actually calling handleSearch.
  // -------------------------------------------------------------------------
  $: {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => handleSearch(searchQuery), 400);
  }

  const clearSearch = () => {
    searchQuery = "";
  };

  // -------------------------------------------------------------------------
  // NEW: CUSTOM FORM SUBMISSION
  //
  // Because selected IDs are tracked in JS (not just DOM checkboxes), we
  // need to build the form payload ourselves rather than relying on the
  // browser's default form serialization.
  //
  // This function fires on form submit. It posts the selected IDs to the
  // existing ?/delete server action, which is completely unchanged.
  // -------------------------------------------------------------------------
  const handleDelete = async (event) => {
    event.preventDefault(); // stop the default browser form submission

    if (selectedIds.size === 0) {
      alert("No articles selected for deletion.");
      return;
    }

    // Build a FormData object manually, adding one entry per selected ID —
    // this matches exactly what the original form would have submitted
    const formData = new FormData();
    for (const id of selectedIds) {
      formData.append("selection", id);
    }

    // Post to the existing SvelteKit form action (unchanged on the server)
    const res = await fetch("?/delete", {
      method: "POST",
      body: formData,
    });

    if (res.ok) {
      // On success: clear selection and refresh the article list
      selectedIds = new Set();
      searchQuery = "";
      // Reload the page so the deleted articles disappear from the table
      window.location.reload();
    } else {
      alert("Deletion failed. Please try again.");
    }
  };
</script>

<PageTitle>Delete articles</PageTitle>
<div class="mx-auto px-5">
  <div class="pt-5">
    <a href="/admin" class="btn-link">
      <p class="align-middle">&#8678; Back to admin dashboard</p>
    </a>
  </div>

  <div class="mx-auto py-5">
    <!-- Show how many articles are currently selected, so the admin always
         knows what will be deleted even after searching and changing the view -->
    <p>
      {#if selectedIds.size > 0}
        <span class="font-semibold text-error">{selectedIds.size} article(s) selected for deletion.</span>
      {:else}
        Select articles to delete. You can search to find specific articles, then check them.
      {/if}
    </p>
  </div>

  <!-- ======================================================================
    NEW: Search bar (same pattern as edit-articles)
  ======================================================================= -->
  <div class="pb-3">
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
    Changes:
    - on:submit now calls handleDelete (our custom JS handler) instead of
      doing a standard POST, so selected IDs from previous searches are included
    - {#each} loops over displayedArticles instead of articles
    - Each checkbox calls toggleSelection() to update our selectedIds Set
    - checked={selectedIds.has(article.id)} keeps checkboxes visually in sync
      when you search back to an article you already checked
  ======================================================================= -->
  <div>
    <form on:submit={handleDelete}>
      <button class="btn btn-sm btn-error mb-5" disabled={selectedIds.size === 0}>
        Delete {selectedIds.size > 0 ? `${selectedIds.size} selected` : "selected"} article(s)
      </button>
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
          {#each displayedArticles as article}
            <tr class={selectedIds.has(article.id) ? "bg-error bg-opacity-10" : ""}>
              <th>
                <input
                  type="checkbox"
                  name="selection"
                  value={article.id}
                  checked={selectedIds.has(article.id)}
                  on:change={(e) => toggleSelection(article.id, e.target.checked)}
                />
              </th>
              <td>{article.title}</td>
              <td>{article.topics}</td>
              <td>{article.questions}</td>
              <td>{article.date}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </form>
  </div>
</div>
