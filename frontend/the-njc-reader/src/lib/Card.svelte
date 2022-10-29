<script>
  export let title;
  export let url;
  export let topics;
  export let questions;
  export let date;

  let checkBox;
  $: collapseTitle = checkBox
    ? "Past year questions"
    : "View past year questions";
</script>

<div class="card bg-base-100 shadow-lg ">
  <div class="card-body pb-5">
    <h2 class="card-title">
      <a href={url} rel="noreferrer" target="_blank">{title}</a>
    </h2>
    <p>{date}</p>
    <div />
    <div class="inline">
      {#each topics as topic}
        <a data-sveltekit-reload href="/search/{topic}"
          ><div class="badge badge-outline mx-1">{topic}</div></a
        >
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
          {#each questions as question}
            <li>{question}</li>
            <br />
          {/each}
        </ul>
      </div>
    </div>
  </div>
</div>
