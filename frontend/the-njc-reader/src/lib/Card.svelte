<script>
  export let title;
  export let url;
  export let topics;
  export let question_display;
  export let date;
  export let mustRead;

  let checkBox;
  $: collapseTitle = checkBox
    ? "Past year questions"
    : "View past year questions";
</script>

<div class="card bg-base-100 shadow-lg ">
  <div class="card-body pb-5">
    {#if mustRead}
      <div class="badge badge-secondary">Must read!</div>
    {/if}
    <h2 class="card-title">
      <a href={url} rel="noreferrer" target="_blank">{title}</a>
    </h2>
    <p>{date}</p>
    <div />
    <div class="inline">
      {#each topics as topic}
        <form class="inline" action="/search" method="POST">
          <input type="hidden" value={topic} name="query" />
          <button class="badge badge-outline ">{topic}</button>
        </form>
      {/each}
    </div>
    <div class="divider" />

    <div class="collapse collapse-arrow">
      <input type="checkbox" bind:checked={checkBox} />
      <div class="collapse-title px-0 font-medium">
        {collapseTitle}
      </div>
      <div class="collapse-content px-0 py-0">
        <ul class="list-none">
          {#each question_display as question}
            <form action="/search" method="POST">
              <input type="hidden" value={question} name="query" />
              <button class="submit-btn">{question}</button>
            </form>
            <br />
          {/each}
        </ul>
      </div>
    </div>
  </div>
</div>

<style>
  .submit-btn {
    width: 100%;
    position: relative;
    background-color: none;
    text-align: start;
    cursor: default;
  }

  .submit-btn:hover {
    cursor: pointer;
  }
</style>
