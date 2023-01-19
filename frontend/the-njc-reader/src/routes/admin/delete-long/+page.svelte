<script>
  import PageTitle from "$lib/PageTitle.svelte";
  export let data;
  const articles = data.articles;
</script>

<PageTitle>Delete long articles</PageTitle>
<div class="mx-auto px-5 pt-5">
  <div class=" pt-5 ">
    <a href="/admin" class="btn-link"
      ><p class="align-middle">&#8678; Back to admin dashboard</p></a
    >
  </div>
  <div class="mx-auto py-5">
    <p>Select article(s) to delete.</p>
  </div>
  <div>
    <form method="POST" action="?/delete">
      <button class="btn btn-sm btn-primary mb-5"
        >Delete selected articles</button
      >
      <table class="table table-compact w-full">
        <thead>
          <tr>
            <th>check</th>
            <th>Title</th>
            <th>Topic</th>
          </tr>
        </thead>
        <tbody>
          {#each articles as article}
            <tr>
              <th
                ><input
                  type="checkbox"
                  name="selection"
                  value={article.id}
                /></th
              >
              <td>{article.title}</td>
              <td>{article.topic}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </form>
  </div>
</div>
