<script>
  import dayjs from "dayjs";
  import relativeTime from "dayjs/plugin/relativeTime";

  import Icon from "svelte-awesome";
  import heart from "svelte-awesome/icons/heart";
  import heartO from "svelte-awesome/icons/heartO";

  export let note;

  dayjs().format();
  dayjs.extend(relativeTime);
</script>

<div class="flex place-content-center">
  <div class="card w-full md:w-3/5 bg-secondary bg-opacity-10 shadow-md">
    <div class="card-body py-5">
      <div class="chat chat-start relative">
        <div class="chat-header">
          {note.author ?? "anonymous"}
          <time class="text-xs opacity-50">{dayjs(note.date).fromNow()}</time>
        </div>
        <div class="chat-bubble ">
          <div>
            {note.tldr}
          </div>
          <div class="pt-5 px-4">
            <blockquote class="italic link">
              <a href={note.article.url} target="blank" rel="noreferrer"
                >{note.article.title}</a
              >
            </blockquote>
          </div>
        </div>
        <div class="absolute top-1 right-0 flex">
          <button
            class="btn btn-xs text-primary border-none bg-transparent hover:bg-transparent "
          >
            <Icon data={heartO} />
          </button>
          <p class="text-sm translate-y-0.5">38</p>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
</style>
