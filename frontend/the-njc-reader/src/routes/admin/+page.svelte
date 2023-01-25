<script>
  import { page } from "$app/stores";
  import PageTitle from "$lib/PageTitle.svelte";

  const userName = $page.data.user.display_name;
</script>

<PageTitle>Admin dashboard</PageTitle>
  <div class="flex justify-center pt-6">
    <p class="px-5 md:text-2xl font-medium text-center">
      Welcome to the admin dashboard, {userName}!<br> What would you like to do
      today?
    </p>
  </div>
  <div class="divider" />
  <div
    class="grid gap-6 lg:grid-cols-2 justify-center place-content-center md:px-60"
  >
    <div class="card bg-base-100 box-shadow min-w-max">
      <div class="card-body pb-5">
        <div class="grid justify-items-center">
          <h3 class="px-5 text-xl md:text-2xl font-medium">News Feed</h3>
          <button class="btn btn-link text-xl"
            ><a href="/admin/add-articles">Add article(s)</a></button
          >
          <button class="btn btn-link text-xl"
            ><a href="/admin/edit-articles">Edit article(s)</a></button
          >
          <button class="btn btn-link text-xl"
            ><a href="/admin/delete-articles">Delete article(s)</a></button
          >
        </div>
      </div>
    </div>
    <div class="card bg-base-100 box-shadow min-w-max">
      <div class="card-body pb-5">
        <div class="grid justify-items-center">
          <h3 class="px-5 text-xl md:text-2xl font-medium">Longer Reads</h3>
          <button class="btn btn-link text-xl"
            ><a href="/admin/add-long">Add article(s)</a></button
          >
          <button class="btn btn-link text-xl"
            ><a href="/admin/edit-long">Edit article(s)</a></button
          >
          <button class="btn btn-link text-xl"
            ><a href="/admin/delete-long">Delete article(s)</a></button
          >
        </div>
      </div>
    </div>
    <div class="card bg-base-100 box-shadow min-w-max">
      <div class="card-body pb-5">
        <div class="grid justify-items-center">
          <h3 class="px-5 text-xl md:text-2xl font-medium">Manage users</h3>
          <button class="btn btn-link text-xl"
            ><a href="/admin/add-users">Add users</a></button
          >
          <button class="btn btn-link text-xl"
            ><a href="/admin/delete-users">Delete users</a></button
          >
        </div>
      </div>
    </div>
  </div>

<style>
  .box-shadow {
    box-shadow: 1px 0 8px 0 rgb(0 0 0 / 0.1);
  }
</style>
