<script>
  import PageTitle from "$lib/PageTitle.svelte";
  export let data;
  export let form;
  
  const articles = data.articles;
  const todayTimestamp = data.todayTimestamp;
</script>

<PageTitle>Send Telegram Digest</PageTitle>
<div class="mx-auto px-5">
  <div class="pt-5">
    <a href="/admin" class="btn-link"
      ><p class="align-middle">&#8678; Back to admin dashboard</p></a
    >
  </div>
  <div class="mx-auto py-5">
    <p>
      Displaying 25 most recent articles. Articles added today are checked by default. 
      Select articles to send to Telegram channel.
    </p>
  </div>
  
  {#if form?.success}
    <div class="alert alert-success mb-5">
      <div>
        <span>✓ Telegram message sent successfully!</span>
      </div>
    </div>
  {/if}
  
  {#if form?.failed}
    <div class="alert alert-error mb-5">
      <div>
        <span>✗ Failed to send: {form.message}</span>
      </div>
    </div>
  {/if}
  
  <div>
    <form method="POST" action="?/send">
      <button class="btn btn-sm btn-primary mb-5">
        Send to Telegram
      </button>
      <table class="table table-compact w-full">
        <thead>
          <tr>
            <th>Send</th>
            <th>Title</th>
            <th>Topics</th>
            <th>Date</th>
          </tr>
        </thead>
        <tbody>
          {#each articles as article}
            <tr>
              <th>
                <input
                  type="checkbox"
                  name="selection"
                  value={JSON.stringify({
                    title: article.title,
                    url: article.url,
                    topics: article.topics
                  })}
                  checked={article.published_on >= todayTimestamp}
                />
              </th>
              <td>{article.title} {#if article.must_read}<span class="badge badge-secondary">Must Read</span>{/if}</td>
              <td>{article.topics}</td>
              <td>{article.date}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </form>
  </div>
</div>
