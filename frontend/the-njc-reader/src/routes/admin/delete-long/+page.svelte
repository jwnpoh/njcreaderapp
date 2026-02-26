<script>
  import PageTitle from "$lib/PageTitle.svelte";

  export let data;

  const articles = data.articles;


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

  let selectedIds = new Set();
let deleteStatus = ""; // "success" | "error" | ""
let deleteError = "";

  const toggleSelection = (id, checked) => {
    if (checked) {
      selectedIds.add(id);
    } else {
      selectedIds.delete(id);
    }
    selectedIds = selectedIds; // trigger Svelte reactivity
  };

  // -------------------------------------------------------------------------
  // NEW: CUSTOM FORM SUBMISSION
  // -------------------------------------------------------------------------
  const handleDelete = async (event) => {
    event.preventDefault();

    if (selectedIds.size === 0) return;

    const formData = new FormData();
    for (const id of selectedIds) {
      formData.append("selection", id);
    }

    const res = await fetch("?/delete", {
      method: "POST",
      body: formData,
    });

    if (res.ok) {
  selectedIds = new Set();
  searchQuery = "";
  deleteStatus = "success";
  window.location.reload();
} else {
  deleteStatus = "error";
  deleteError = "Deletion failed. Please try again.";
}
  };
</script>

<PageTitle>Delete long articles</PageTitle>
<div class="mx-auto px-5 pt-5">
  <div class="pt-5">
    <a href="/admin" class="btn-link">
      <p class="align-middle">&#8678; Back to admin dashboard</p>
    </a>
  </div>

  <div class="mx-auto py-5">
    <p>
      {#if selectedIds.size > 0}
        <span class="font-semibold text-error">{selectedIds.size} article(s) selected for deletion.</span>
      {:else}
        Select article(s) to delete. You can filter to find specific articles, then check them.
      {/if}
    </p>
  </div>

  <!-- ======================================================================
    NEW: Filter bar
    Client-side only — instantly filters the table as the admin types.
  ======================================================================= -->
  <div class="pb-3">
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
    Changes:
    - on:submit calls handleDelete (our custom JS handler)
    - {#each} loops over displayedArticles instead of articles
    - checkboxes use toggleSelection() and checked={selectedIds.has(...)}
      to stay in sync with the selectedIds Set across filter changes
  ======================================================================= -->
  <div>
    <form on:submit={handleDelete}>
      <button class="btn btn-sm btn-error mb-5" disabled={selectedIds.size === 0}>
        Delete {selectedIds.size > 0 ? `${selectedIds.size} selected` : "selected"} article(s)
      </button>
      {#if deleteStatus === "error"}
  <div class="alert alert-error max-w-lg mb-3">
    <span>{deleteError}</span>
  </div>
{/if}
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
              <td>{article.topic}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </form>
  </div>
</div>
