<script>
  export let title;
  export let url;
  export let topics;
  export let questions;
  export let date;
</script>

<div class="card bg-base-100 shadow-md min-h-full ">
  <div class="card-body">
    <h2 class="card-title"><a href={url} target="_blank">{title}</a></h2>
    <p>{date}</p>
    <div class="inline">
      {#each topics as topic}
        <a href="/search"><div class="badge badge-outline mx-1">{topic}</div></a
        >
      {/each}
    </div>
    {#each questions as question}
      {question}
      <br />
    {/each}
  </div>
</div>
