<script>
  import { DateInput } from "date-picker-svelte";
  import { page } from "$app/stores";
  import { enhance, applyAction } from "$app/forms";

  export let url;
  export let title;
  export let tags;
  export let date = new Date();

  const session = $page.data.user.session;

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  const getTitle = async (url) => {
    const payload = { url: url };
    const res = await fetch("http://localhost:8080/api/articles/get-title", {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    title = response.data;
  };
</script>

<form
  method="POST"
  action="?/queue"
  use:enhance={({ form }) => {
    return async ({ result, update }) => {
      if (result.type === "success") {
        form.reset();
      }
      if (result.type === "invalid") {
        await applyAction(result);
      }
      update();
    };
  }}
>
  <div class="flex pt-5 px-5">
    <div class="flex-auto">
      <div class="pt-3">
        <input
          required
          name="url"
          type="text"
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

  <button class="btn btn-sm btn-secondary mx-7">Add to queue</button>
</form>
