<!--
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
-->

<script>
  import PageTitle from "$lib/PageTitle.svelte";
  import { page } from "$app/stores";
  import { DateInput } from "date-picker-svelte";

  export let data;

  const session = $page.data.user.session;
  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  // -------------------------------------------------------------------------
  // FORM STATE
  // One set of fields — identical to the original page.
  // -------------------------------------------------------------------------
  let url = "";
  let title = "";
  let tags = "";
  let date = new Date();
  let must_read = false;
  let fetchingTitle = false;
  let formError = "";

  const getTitle = async (url) => {
    if (!url) {
      title = "";
      return;
    }
    fetchingTitle = true;
    const res = await fetch(`${data.API_URL}/api/admin/articles/get-title`, {
      method: "POST",
      body: JSON.stringify({ url }),
      headers: myHeaders,
    });
    const response = await res.json();
    title = response.data ?? "";
    fetchingTitle = false;
  };

  const resetForm = () => {
    url = "";
    title = "";
    tags = "";
    date = new Date();
    must_read = false;
    formError = "";
  };

  // -------------------------------------------------------------------------
  // QUEUE STATE
  // Each entry matches the ArticlePayload struct the API expects:
  //   { title, url, tags, date, must_read }
  // date is stored as date.toString() — the broker's formatDate() strips
  // the leading day-of-week token and extracts "Mon DD YYYY", which
  // ParseUnixTime then parses with the "Jan 02 2006" layout.
  // This is exactly what the original form action was sending.
  // -------------------------------------------------------------------------
  let queue = [];

  const addToQueue = () => {
    if (!url.trim() || !title.trim() || !tags.trim()) {
      formError = "Please fill in URL, title, and tags before adding to the queue.";
      return;
    }
    queue = [
      ...queue,
      {
        url: url.trim(),
        title: title.trim(),
        tags: tags.trim(),
        date: date.toString(),
        must_read: must_read ? "on" : "",
      },
    ];
    resetForm();
  };

  const removeFromQueue = (index) => {
    queue = queue.filter((_, i) => i !== index);
  };

  // -------------------------------------------------------------------------
  // SUBMISSION
  // Sends the entire queue in one POST to the existing insert endpoint.
  // On success, clears the queue. On failure, preserves it so no work is lost.
  // -------------------------------------------------------------------------
  let isSubmitting = false;
  let submitStatus = ""; // "success" | "error" | ""
  let submitError = "";

  const sendToDatabase = async () => {
    if (queue.length === 0) return;
    isSubmitting = true;
    submitStatus = "";

    try {
      const res = await fetch(`${data.API_URL}/api/admin/articles/insert`, {
        method: "POST",
        body: JSON.stringify(queue),
        headers: myHeaders,
      });
      const response = await res.json();

      if (response.error) {
        submitStatus = "error";
        submitError = response.message ?? "Something went wrong. Your queue has been preserved — please try again.";
      } else {
        submitStatus = "success";
        queue = [];
      }
    } catch (e) {
      submitStatus = "error";
      submitError = "Request failed. Please check your connection. Your queue has been preserved.";
    }

    isSubmitting = false;
  };
</script>

<PageTitle>Add articles</PageTitle>
<div class="px-5 pt-5">
  <a href="/admin" class="btn-link">
    <p class="align-middle">&#8678; Back to admin dashboard</p>
  </a>
</div>

<div class="px-5 pt-5">
  <p class="text-sm opacity-60 pb-4">
    Fill in one article at a time and add it to the queue. When you're done, send all queued articles to the database at once.
  </p>

  <!-- ======================================================================
    ENTRY FORM
    Row 1: URL (wide) | Date | Must read
    Row 2: Title | Tags
  ======================================================================= -->
  <div class="card bg-base-100 shadow-sm border border-base-300 p-5 max-w-3xl mb-4">

    <!-- Row 1 -->
    <div class="grid grid-cols-12 gap-3 items-end mb-3">
      <div class="col-span-6">
        <label class="label py-0 pb-1">
          <span class="label-text text-xs opacity-60">Article URL</span>
        </label>
        <input
          type="url"
          placeholder="https://..."
          class="input input-bordered input-sm w-full"
          bind:value={url}
          on:change={() => getTitle(url)}
        />
      </div>

      <div class="col-span-4">
        <label class="label py-0 pb-1">
          <span class="label-text text-xs opacity-60">Published on</span>
        </label>
        <DateInput
          placeholder="dd-MM-yyyy"
          format="dd-MM-yyyy"
          closeOnSelection
          bind:value={date}
        />
      </div>

      <div class="col-span-2 flex items-center gap-2 pb-1">
        <input type="checkbox" class="checkbox checkbox-sm" bind:checked={must_read} />
        <span class="text-sm">Must read</span>
      </div>
    </div>

    <!-- Row 2 -->
    <div class="grid grid-cols-2 gap-3 mb-4">
      <div>
        <label class="label py-0 pb-1">
          <span class="label-text text-xs opacity-60">
            Title{#if fetchingTitle} <span class="italic opacity-60">(fetching...)</span>{/if}
          </span>
        </label>
        <input
          type="text"
          placeholder="Auto-filled from URL"
          class="input input-bordered input-sm w-full"
          bind:value={title}
        />
      </div>

      <div>
        <label class="label py-0 pb-1">
          <span class="label-text text-xs opacity-60">Tags</span>
        </label>
        <input
          type="text"
          placeholder="e.g. 2019-Q6; leadership"
          class="input input-bordered input-sm w-full"
          bind:value={tags}
        />
      </div>
    </div>

    <div class="flex items-center gap-4">
      <button type="button" class="btn btn-primary btn-sm" on:click={addToQueue}>
        Add to queue
      </button>
      {#if formError}
        <p class="text-error text-sm">{formError}</p>
      {/if}
    </div>
  </div>

  <!-- ======================================================================
    QUEUE
    Only rendered once at least one article has been staged.
    "Send to database" lives here, clearly separated from the entry form.
  ======================================================================= -->
  {#if queue.length > 0 || submitStatus === "success"}
    <div class="card bg-base-100 shadow-sm border border-base-300 p-5 max-w-3xl">

      {#if submitStatus === "success" && queue.length === 0}
        <div class="alert alert-success max-w-fit">
          <span>All articles added successfully.</span>
        </div>
      {:else}
        <!-- Queue header -->
        <div class="flex items-center justify-between mb-3">
          <p class="text-sm font-medium opacity-70">
            {queue.length} article{queue.length !== 1 ? "s" : ""} queued
          </p>
          <button
            type="button"
            class="btn btn-primary btn-sm"
            on:click={sendToDatabase}
            disabled={isSubmitting}
          >
            {isSubmitting
              ? "Sending..."
              : `Send ${queue.length} article${queue.length !== 1 ? "s" : ""} to database`}
          </button>
        </div>

        <!-- Queue table -->
        <table class="table table-compact w-full text-sm">
          <thead>
            <tr>
              <th>Title</th>
              <th>Tags</th>
              <th>Date</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {#each queue as entry, i}
              <tr>
                <td class="max-w-xs"><span class="block truncate">{entry.title}</span></td>
                <td>{entry.tags}</td>
                <!-- Slice the first 15 chars of the toString() date for display,
                     e.g. "Wed Feb 12 2026" — the full string is what gets sent -->
                <td>{entry.date.slice(0, 15)}</td>
                <td>
                  <button
                    type="button"
                    class="btn btn-ghost btn-xs text-error"
                    on:click={() => removeFromQueue(i)}
                  >
                    Remove
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>

        {#if submitStatus === "error"}
          <div class="alert alert-error max-w-fit mt-4">
            <span>{submitError}</span>
          </div>
        {/if}
      {/if}

    </div>
  {/if}
</div>
