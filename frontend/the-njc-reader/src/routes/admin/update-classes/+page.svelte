<script>
  import PageTitle from "$lib/PageTitle.svelte"
  export let form;
</script>

<PageTitle>Update classes</PageTitle>
<div class="px-5 pt-5 ">
  <a href="/admin" class="btn-link"
    ><p class="align-middle">&#8678; Back to admin dashboard</p></a
  >
</div>
<div class="pt-4 px-10">
  <h2 class="text-lg font-semibold">Update classes</h2>
  <div class="py-5">
    <p class=" text-md  text-justify opacity-90">
      Upload the CSV template. 
    </p>
      <p class="text-md text-justify opacity-90">
        As this function takes a little longer to process, confirmation of success will be sent to your email instead. Check your email in a while to see if the function was successfully completed.
      </p>
  </div>
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
      <span>{form.message}</span>
    </div>
  {/if}
  {#if form?.failed}
    <div class="alert alert-error max-w-fit place-self-center">
      <span class="text-center">{form?.message}</span>
    </div>
  {/if}
</div>