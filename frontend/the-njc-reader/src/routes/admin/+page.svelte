<script>
  import { page } from "$app/stores";
  import PageTitle from "$lib/PageTitle.svelte";

  const userName = $page.data.user.display_name;
  const stats = $page.data.stats
  const questions_most_articles = $page.data.stats.questions_with_most_articles
  const questions_fewest_articles = $page.data.stats.questions_with_fewest_articles
  const topics_most_articles = $page.data.stats.topics_with_most_articles
  const topics_fewest_articles = $page.data.stats.topics_with_fewest_articles

  console.log(stats)
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

<!--   
  <div class="card bg-base-100 box-shadow min-w-max">
      <div class="card-body pb-5">
        <div class="grid justify-items-center">
          <h3 class="px-5 text-xl md:text-2xl font-medium">Manage users</h3>
          <button class="btn btn-link text-xl"
            ><a href="/admin/add-users">Add users</a></button
          >
          <button class="btn btn-link text-xl"
            ><a href="/admin/update-classes">Update classes</a></button
          >
        </div>
      </div>
    </div>
-->
</div>
<div
    class="grid justify-center place-content-center mt-4 mb-4 "
  >
          <button class="btn btn-link text-xl"
            ><a href="/admin/telegram">Send Telegram digest</a></button
          >
          </div>

<div
    class="grid justify-center place-content-center mt-10 "
  >
  <div class="stats stats-vertical lg:stats-horizontal shadow text-lg">
    <div class="stat">
      <div class="stat-title">Number of articles to date</div>
      <div class="stat-value">{stats.number_of_articles}</div>
    </div>
  </div>
</div>

<div
    class="grid justify-center place-content-center mt-10 "
  >
  <div class="stats stats-vertical lg:stats-horizontal shadow ">

  <div class="stat">
    <div class="stat-title">Topics with most articles</div>
    {#each topics_most_articles as topic}
      <div class="stat-value text-2xl pt-2">{topic.K}</div>
      <div class="stat-desc">{topic.V} articles</div>
    {/each}
  </div>

  <div class="stat">
    <div class="stat-title">Topics with fewest articles</div>
    {#each topics_fewest_articles as topic}
      <div class="stat-value text-2xl pt-2">{topic.K}</div>
      <div class="stat-desc">{topic.V} {topic.V > 1 ? "articles" : "article"}</div>
    {/each}
  </div>
  <div class="stat">
    <div class="stat-title">Questions with most articles</div>
    {#each questions_most_articles as question}
      <div class="stat-value text-2xl pt-2">{question.K}</div>
      <div class="stat-desc">{question.V} articles</div>
    {/each}
  </div>
  <div class="stat">
    <div class="stat-title">Questions with fewest articles</div>
    {#each questions_fewest_articles as question}
      <div class="stat-value text-2xl pt-2">{question.K}</div>
      <div class="stat-desc">{question.V} {question.V > 1 ? "articles" : "article"}</div>
    {/each}
  </div>
  </div>
</div>
<style>
  .box-shadow {
    box-shadow: 1px 0 8px 0 rgb(0 0 0 / 0.1);
  }
</style>
