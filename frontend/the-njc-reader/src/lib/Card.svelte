<script>
  import Icon from "svelte-awesome";
  import bookmark from "svelte-awesome/icons/bookmark";

  export let id;
  export let title;
  export let url;
  export let topics;
  export let question_display;
  export let date;
  export let mustRead;
  export let loggedIn;

  let checkBox;
  $: collapseTitle = checkBox
    ? "Past year questions"
    : "View past year questions";
</script>

<div class="relative">
  <div class="card bg-base-100 box-shadow ">
    {#if loggedIn}
      <div class="bookmark bg-base-100">
        <form method="POST" action="/notes/add-note?/newnote">
          <input name="article_id" type="hidden" hidden value={id} />
          <button>
            <Icon data={bookmark} scale={2} class="bookmark text-primary" />
          </button>
        </form>
      </div>
    {/if}
    <div class="card-body pb-5">
      {#if mustRead}
        <div class="badge badge-secondary py-3">Must read!</div>
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
</div>

<style>
  .box-shadow {
    box-shadow: 1px 0 8px 0 rgb(0 0 0 / 0.1);
  }
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

  .bookmark {
    position: absolute;
    top: -1%;
    right: 1.5em;
    width: 25px;
    height: 30px;
  }

  .bookmark:hover {
    transform: scaleY(1.5);
    transition: 300ms cubic-bezier(0.19, 1, 0.22, 1);
  }
</style>
