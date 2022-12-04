<script>
  import { page } from "$app/stores";

  import dayjs from "dayjs";
  import relativeTime from "dayjs/plugin/relativeTime";

  import Icon from "svelte-awesome";
  import heart from "svelte-awesome/icons/heart";
  import heartO from "svelte-awesome/icons/heartO";

  export let note;
  export let API_URL;

  dayjs().format();
  dayjs.extend(relativeTime);

  const user_id = $page.data.user.id;
  const session = $page.data.user.session;
  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  const addFollow = async (to_follow) => {
    const payload = { user_id: user_id, to_follow: to_follow };
    const res = await fetch(`${API_URL}/api/users/follow`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    const msg = await response.message;
    console.log(msg);
  };
</script>

<div class="flex place-content-center">
  <div class="card w-full md:w-3/5 bg-secondary bg-opacity-10 shadow-md">
    <div class="card-body py-5">
      <div class="chat chat-start relative">
        <div class="chat-header">
          <button
            on:click={() => {
              addFollow(note.user_id);
            }}
          >
            {note.author ?? "anonymous"}
            {note.user_id ?? "anonymous"}
          </button>
          <time class="text-xs opacity-50">{dayjs(note.date).fromNow()}</time>
        </div>
        <div class="chat-bubble ">
          <div>
            {note.tldr}
          </div>
          <div class="pt-5 px-4">
            <blockquote class="italic link">
              <a href={note.article_url} target="blank" rel="noreferrer"
                >{note.article_title}</a
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
