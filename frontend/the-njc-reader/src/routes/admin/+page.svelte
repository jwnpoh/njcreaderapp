<script>
  import { page } from "$app/stores";
  import Container from "$lib/Container.svelte";
  import PageTitle from "$lib/PageTitle.svelte";

  const userName = $page.data.user.display_name;
</script>

<PageTitle>Admin dashboard</PageTitle>
<Container>
  <div class="pt-5 flex justify-center">
    <p class="px-5 md:text-2xl font-medium">
      Welcome to the admin dashboard, {userName}! What would you like to do
      today?
    </p>
  </div>
  <div class="divider" />
  <div class="grid md:grid-cols-3 justify-items-center">
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
</Container>
