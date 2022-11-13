<script>
  export let data;
  const articles = data.articles;
</script>

<div class=" pt-5 ">
  <a href="/admin" class="btn-link"
    ><p class="align-middle">&#8678; Back to admin dashboard</p></a
  >
</div>
<div class="mx-auto py-5">
  <p>
    Displaying 100 most recent articles. Select articles to delete. To delete
    older articles, contact database admin.
  </p>
</div>
<div>
  <form method="POST" action="?/delete">
    <button class="btn btn-sm btn-primary mb-5">Delete selected articles</button
    >
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
        {#each articles as article}
          <tr>
            <th
              ><input type="checkbox" name="selection" value={article.id} /></th
            >
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
