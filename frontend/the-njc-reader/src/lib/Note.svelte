<script>
  import { page } from "$app/stores";

  import dayjs from "dayjs";
  import relativeTime from "dayjs/plugin/relativeTime";

  import Icon from "svelte-awesome";
  import heart from "svelte-awesome/icons/heart";
  import heartO from "svelte-awesome/icons/heartO";
  import edit from "svelte-awesome/icons/edit";
  import trashO from "svelte-awesome/icons/trashO";

  export let note;
  export let API_URL;

  dayjs().format();
  dayjs.extend(relativeTime);

  const user_id = $page.data.user.id;
  const session = $page.data.user.session;
  const liked_notes = $page.data.user.liked_notes;
  let note_likes = $page.data.user.note_likes;

  $: liked = liked_notes.includes(note.id) ? true : false;

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

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

  const deleteNote = async (noteID) => {
    const payload = noteID;
    await fetch(`${API_URL}/api/posts/delete`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    location.reload();
  };

  let seeMore;
</script>

<div class="flex place-content-center">
  <div class="card w-full md:w-3/5 bg-secondary bg-opacity-10 shadow-md">
    <div class="card-body pt-5 pb-3">
      <div class="chat chat-start relative">
        <div class="chat-header">
          <a href="/profile/{note.user_id}">
            {note.author ?? "anonymous"}
          </a>
          <time class="block text-xs opacity-50"
            >{dayjs(note.date).fromNow()}</time
          >
        </div>
        <div class="chat-bubble pt-2">
          <div class="pt-2">
            <blockquote class="italic link text-info">
              <a href={note.article_url} target="blank" rel="noreferrer"
                >{note.article_title}</a
              >
            </blockquote>
          </div>
          <div class="py-2">
            <h3 class="italic font-semibold">TL;DR</h3>
            {note.tldr}
          </div>
          <div class="inline">
            {#if note.tags[0] !== ""}
              {#each note.tags as tag}
                <button class="mr-1 badge badge-outline">{tag}</button>
              {/each}
            {/if}
          </div>
        </div>
        <div class="absolute top-1 right-0 flex">
          {#if !note.public}
            <p class="text-neutral text-sm italic pr-1">Private note</p>
            <form method="POST" action="/notes/edit-note?/edit">
              <input name="note_id" type="hidden" hidden value={note.id} />
              <button
                class="btn btn-xs text-neutral border-none bg-transparent hover:bg-transparent"
                ><Icon data={edit} /></button
              >
            </form>
            <button
              class="btn btn-xs text-neutral border-none bg-transparent hover:bg-transparent"
              on:click={deleteNote(note.id)}><Icon data={trashO} /></button
            >
          {:else if note.public}
            {#if note.user_id == user_id}
              <p class="text-neutral text-sm italic pr-1">
                {#if note_likes[note.id] > 0}
                  <Icon class="text-primary " data={heart} />
                {:else}
                  <Icon class="text-primary " data={heartO} />
                {/if}
                {#if !note_likes[note.id]}
                  0
                {:else}
                  {note_likes[note.id]}
                {/if}
                {note_likes[note.id] > 1 ? "likes" : "like"} received
              </p>
              <form method="POST" action="/notes/edit-note?/edit">
                <input name="note_id" type="hidden" hidden value={note.id} />
                <button
                  class="btn btn-xs text-neutral border-none bg-transparent hover:bg-transparent"
                  ><Icon data={edit} /></button
                >
              </form>
              <button
                class="btn btn-xs text-neutral border-none bg-transparent hover:bg-transparent"
                on:click={deleteNote(note.id)}><Icon data={trashO} /></button
              >
            {:else if liked}
              <button
                class="btn btn-xs text-primary border-none bg-transparent hover:bg-transparent heart "
                on:click={() => {
                  updateLikes(false);
                }}
              >
                <Icon data={heart} />
              </button>
              <p class="text-sm translate-y-0.5">{note_likes[note.id]}</p>
            {:else}
              <button
                class="btn btn-xs text-primary border-none bg-transparent hover:bg-transparent heart "
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
      <div class="collapse right-4">
        <input type="checkbox" bind:checked={seeMore} />
        <div class="collapse-title text-sm font-normal">
          <p class="text-sm text-info">
            {seeMore ? "See less" : "See more..."}
          </p>
        </div>
        <div class="collapse-content">
          <div class="py-1">
            <h3 class="italic font-semibold">Examples from article</h3>
            <div class="whitespace-pre-line">{note.examples}</div>
          </div>
          <div class="py-2">
            <h3 class="italic font-semibold">Further reflection</h3>
            {note.notes}
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .heart:hover {
    transform: scale(1.15);
    transition: 300ms cubic-bezier(0.19, 1, 0.22, 1);
  }
</style>
