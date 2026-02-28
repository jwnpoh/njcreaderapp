<script>
  import PageTitle from "$lib/PageTitle.svelte";
  import { page } from "$app/stores";
  import { DateInput } from "date-picker-svelte";

  export let data;

  const session = $page.data.user.session;
  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  // ---------------------------------------------------------------------------
  // FORM STATE
  // ---------------------------------------------------------------------------
  let url = "";
  let title = "";
  let tags = "";
  let date = new Date();
  let must_read = false;
  let fetchingTitle = false;
  let formError = "";
  let duplicateWarning = "";

  const getTitle = async (url) => {
    if (!url) { title = ""; return; }
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
    duplicateWarning = "";
  };

  // ---------------------------------------------------------------------------
  // TAG PARSING
  //
  // Replicates the backend's formatQuestionString + parseTags logic in JS.
  // The Go regex is: \s?\d{4}(\s?)+-?(\s?)+(q|Q)?\d{1,2}
  // Very permissive — matches YYYY followed by optional spacing/hyphens,
  // optional q/Q, then 1-2 digits.
  // ---------------------------------------------------------------------------
  const QUESTION_REGEX = /\d{4}\s*-?\s*(q|Q)?\d{1,2}/;

  // Normalises a matched question token to YYYY - QN, mirroring the backend's
  // fmt.Sprintf("%s - Q%s", year, qnNumber).
  const normaliseQuestion = (raw) => {
    const year = raw.match(/\d{4}/)?.[0] ?? "";
    const qn = raw.match(/(?:q|Q)?(\d{1,2})$/)?.[1] ?? "";
    return `${year} - Q${qn}`;
  };

  const parseTags = (tagString) => {
    // Mirror splitTags: trim, strip trailing semicolon, split on semicolon
    const raw = tagString.trim().replace(/;$/, "");
    const tokens = raw.split(";").map(t => t.trim()).filter(Boolean);

    const parsed = tokens.map(token => ({
      raw: token,
      isQuestion: QUESTION_REGEX.test(token),
      // display: normalised form for questions, raw for topics
      display: QUESTION_REGEX.test(token) ? normaliseQuestion(token) : token,
    }));

    // Topics first, then questions — easier to scan for misclassifications
    return [
      ...parsed.filter(t => !t.isQuestion),
      ...parsed.filter(t => t.isQuestion),
    ];
  };

  // ---------------------------------------------------------------------------
  // QUEUE STATE
  //
  // Each entry stores the Date object directly (not date.toString()) so that
  // round-trip re-hydration via editEntry is lossless. date.toString() is only
  // called at send time when building the API payload.
  // ---------------------------------------------------------------------------
  let queue = [];

  // ---------------------------------------------------------------------------
  // DUPLICATE DETECTION
  // Non-blocking — surfaces a warning but does not prevent adding to the queue.
  // ---------------------------------------------------------------------------
  const checkForDuplicates = (incomingUrl, incomingTitle) => {
    const urlMatch = queue.find(e => e.url === incomingUrl);
    const titleMatch = queue.find(e => e.title === incomingTitle);
    if (urlMatch && titleMatch) return "This URL and title are already in the queue.";
    if (urlMatch) return "This URL is already in the queue.";
    if (titleMatch) return "This title is already in the queue.";
    return "";
  };

  const addToQueue = () => {
    if (!url.trim() || !title.trim() || !tags.trim()) {
      formError = "Please fill in URL, title, and tags before adding to the queue.";
      return;
    }

    duplicateWarning = checkForDuplicates(url.trim(), title.trim());

    queue = [
      ...queue,
      {
        url: url.trim(),
        title: title.trim(),
        tags: tags.trim(),
        parsedTags: parseTags(tags.trim()),
        // Store the Date object — not a string — so editEntry can restore it
        // directly into the DateInput binding without any re-parsing.
        date: date,
        must_read: must_read ? "on" : "",
      },
    ];
    resetForm();
  };

  const removeFromQueue = (index) => {
    queue = queue.filter((_, i) => i !== index);
  };

  // ---------------------------------------------------------------------------
  // EDIT ENTRY
  //
  // Pulls an entry out of the queue and populates the form with its data so
  // the admin can correct it and re-add it. The entry is removed from the queue
  // immediately to prevent duplicates.
  //
  // Because date is stored as a Date object, it binds directly to DateInput
  // without any string parsing — no timezone risk.
  // ---------------------------------------------------------------------------
  const editEntry = (index) => {
    const entry = queue[index];
    url = entry.url;
    title = entry.title;
    // Reconstruct the tags string from parsedTags display values so the
    // field shows the normalised form (e.g. "2022 - Q11" not "2022q11")
    tags = entry.parsedTags.map(t => t.display).join("; ");
    date = entry.date;
    must_read = entry.must_read === "on";
    formError = "";
    duplicateWarning = "";
    queue = queue.filter((_, i) => i !== index);
  };

  // ---------------------------------------------------------------------------
  // SUBMISSION
  //
  // date.toString() is called here — once, at send time — on a Date object
  // that has never been serialised and re-parsed. formatDate() on the backend
  // strips the day-of-week and extracts "Mon DD YYYY", which ParseUnixTime
  // parses with the "Jan 02 2006" layout.
  // ---------------------------------------------------------------------------
  let isSubmitting = false;
  let submitStatus = "";
  let submitError = "";

  const sendToDatabase = async () => {
    if (queue.length === 0) return;
    isSubmitting = true;
    submitStatus = "";

    // Build the API payload: strip parsedTags (display only) and convert
    // the stored Date object to a string for the backend to parse.
    const payload = queue.map(({ parsedTags, date, ...rest }) => ({
      ...rest,
      date: date.toString(),
    }));

    try {
      const res = await fetch(`${data.API_URL}/api/admin/articles/insert`, {
        method: "POST",
        body: JSON.stringify(payload),
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

  <!-- ========================================================================
    ENTRY FORM
  ========================================================================= -->
  <div class="card bg-base-100 shadow-sm border border-base-300 p-5 max-w-3xl mb-4 overflow-visible">

    <!-- Row 1: URL | Date | Must read -->
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
          on:input={() => getTitle(url)}
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

    <!-- Row 2: Title | Tags -->
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

    {#if duplicateWarning}
      <div class="alert alert-warning mt-3 py-2">
        <span class="text-sm">⚠ {duplicateWarning} Check the queue below before sending.</span>
      </div>
    {/if}
  </div>

  <!-- ========================================================================
    QUEUE
  ========================================================================= -->
  {#if queue.length > 0 || submitStatus === "success"}
    <div class="card bg-base-100 shadow-sm border border-base-300 p-5 mx-5">

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
        <table class="table table-compact w-full table-fixed text-sm">
          <thead>
            <tr>
              <th class="w-1/3">Title</th>
              <th>
                Tags
                <span class="badge badge-outline badge-sm ml-2 font-normal normal-case">topic</span>
                <span class="badge badge-primary badge-sm ml-1 font-normal normal-case">question</span>
              </th>
              <th class="w-28">Date</th>
              <th class="w-24"></th>
            </tr>
          </thead>
          <tbody>
            {#each queue as entry, i}
              <tr>
                <td>
                  <span class="block truncate">{entry.title}</span>
                  {#if entry.must_read === "on"}
                    <span class="badge badge-secondary badge-sm mt-1">must read</span>
                  {/if}
                </td>

                <td>
                  <div class="flex flex-wrap gap-1">
                    {#each entry.parsedTags as tag}
                      {#if tag.isQuestion}
                        <span class="badge badge-primary badge-sm">{tag.display}</span>
                      {:else}
                        <span class="badge badge-outline badge-sm">{tag.display}</span>
                      {/if}
                    {/each}
                  </div>
                </td>

                <!-- Date is now a Date object — format it for display -->
                <td class="whitespace-nowrap">
                  {entry.date.toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" })}
                </td>

                <td>
                  <div class="flex gap-1">
                    <button
                      type="button"
                      class="btn btn-ghost btn-xs"
                      on:click={() => editEntry(i)}
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      class="btn btn-ghost btn-xs text-error"
                      on:click={() => removeFromQueue(i)}
                    >
                      Remove
                    </button>
                  </div>
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
