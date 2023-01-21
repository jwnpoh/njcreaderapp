<script>
  import { page } from "$app/stores";
  import CardContainer from "$lib/CardContainer.svelte";
  import PageTitle from "$lib/PageTitle.svelte";

  export let data;
  $: articles = data.articles;
  $: pageNo = data.page;

  let user = $page.data.user;
</script>

<PageTitle>News Feed</PageTitle>
{#if !user.loggedIn}
<div class="grid gap-6  my-6 md:px-52">
  <h1 class="font-bold text-2xl">Welcome!</h1>
  <p class="text-xl">
    Too many articles and sources but not enough time to read them all?
  </p>
  <p class="text-lg">
    Your teachers have curated articles which are relevant to our course for
    your daily reading. Each article is also appended with relevant past year 'A' Level essay
    questions so that you can practise making connections
    between the news you read and the issues discussed.
  </p>

  <p class="text-lg">Log in to enjoy more features!</p>

  <p class="text-lg">We hope that you find this resource helpful!</p>
</div>
<div class="divider" />
{/if}
<CardContainer {articles} page={pageNo} loggedIn={user.loggedIn}/>
