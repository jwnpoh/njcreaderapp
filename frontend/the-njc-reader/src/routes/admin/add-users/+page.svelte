<script>
  import PageTitle from "$lib/PageTitle.svelte"
  export let form;
</script>

<PageTitle>Add users</PageTitle>
<div class="px-5 pt-5 ">
  <a href="/admin" class="btn-link"
    ><p class="align-middle">&#8678; Back to admin dashboard</p></a
  >
</div>
<div class="pt-4 px-10">
  <h2 class="text-lg font-semibold">Add users</h2>
  <p class="py-3 text-md  text-justify opacity-80">
    Upload the CSV template. It might take a while if the list is long.
  </p>
  <form method="POST" action="?/add" enctype="multipart/form-data">
    <input
      type="file"
      name="input"
      class="file-primary w-full max-w-xs"
      accept=".csv"
    />
    <div>
      <button class="btn btn-md btn-primary mt-3">Upload user list</button>
    </div>
  </form>
  {#if form?.sent}
    <div class="alert alert-success max-w-fit place-self-center">
      <span>Users added successfully.</span>
    </div>
  {/if}
  {#if form?.failed}
    <div class="alert alert-error max-w-fit place-self-center">
      <span class="text-center">{form?.message}</span>
    </div>
  {/if}
</div>
