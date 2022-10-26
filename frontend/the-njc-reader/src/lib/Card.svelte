<script>
  export let title;
  export let url;
  export let topics;
  export let questions;
  export let date;
</script>

<div class="card bg-base-100 shadow-md min-h-full">
  <div class="card-body pb-5">
    <h2 class="card-title">
      <a href={url} rel="noreferrer" target="_blank">{title}</a>
    </h2>
    <p>{date}</p>
    <div class="inline">
      {#each topics as topic}
        <a data-sveltekit-reload href="/search/{topic}"
          ><div class="badge badge-outline mx-1">{topic}</div></a
        >
      {/each}
    </div>
    <div class="divider" />
    <div class="collapse">
      <input type="checkbox" />
      <div class="collapse-title px-0 pt-2">
        Click to view past year questions
      </div>
      <div class="collapse-content px-0">
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
