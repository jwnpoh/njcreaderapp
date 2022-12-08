<script>
  import { page } from "$app/stores";

  import dayjs from "dayjs";
  import relativeTime from "dayjs/plugin/relativeTime";

  import Icon from "svelte-awesome";
  import heart from "svelte-awesome/icons/heart";
  import heartO from "svelte-awesome/icons/heartO";

  export let note;
  export let liked_notes;
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
  };

  let liked;
  $: if (liked_notes.includes(note.id)) {
    liked = true;
  }

  let note_likes = $page.data.user.note_likes;
  const updateLikes = async (like) => {
    const payload = { user_id: user_id, post_id: note.id, like: like };
    await fetch(`${API_URL}/api/posts/like`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    if (like) {
      note_likes[note.id]++;
      liked_notes.push(note.id);
    } else {
      note_likes[note.id]--;
      liked_notes.splice(liked_notes.indexOf(note.id));
    }

    liked = like;
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
          {#if !note.public}
            <p class="text-neutral text-sm italic">Private note</p>
          {:else if note.public}
            {#if note.user_id == user_id}
              <p class="text-neutral text-sm italic">
                {#if note_likes[note.id] > 0}
                  <Icon class="text-primary" data={heart} />
                {:else}
                  <Icon class="text-primary" data={heartO} />
                {/if}
                {note_likes[note.id]}
                {note_likes[note.id] > 1 ? "likes" : "like"} received
              </p>
            {:else if liked}
              <button
                class="btn btn-xs text-primary border-none bg-transparent hover:bg-transparent "
                on:click={() => {
                  updateLikes(false);
                }}
              >
                <Icon data={heart} />
              </button>
              <p class="text-sm translate-y-0.5">{note_likes[note.id]}</p>
            {:else}
              <button
                class="btn btn-xs text-primary border-none bg-transparent hover:bg-transparent "
                on:click={() => {
                  updateLikes(true);
                }}
              >
                <Icon data={heartO} />
              </button>
              <p class="text-sm translate-y-0.5">{note_likes[note.id]}</p>
            {/if}
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

<style>
</style>
