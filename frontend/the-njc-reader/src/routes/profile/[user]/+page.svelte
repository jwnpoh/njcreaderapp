<script>
  import { page } from "$app/stores";
  import Container from "$lib/Container.svelte";
  import NotesContainer from "$lib/NotesContainer.svelte";

  export let data;

  const user_id = $page.data.user.id;
  const session = $page.data.user.session;
  const userInfo = data.userInfo;
  const userNotes = data.userNotes;
  const friends = data.friends;
  const API_URL = data.API_URL;

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  let isFriend = friends.data.followed_by_ids.includes(user_id) ? true : false;

  let followersCount = friends.data.followed_by_ids
    ? friends.data.followed_by_ids.length
    : 0;
  let followingCount = friends.data.following_ids
    ? friends.data.following_ids.length
    : 0;
  const userNotesCount = userNotes.data ? userNotes.data.length : 0;

  const updateFollows = async (to_follow) => {
    let follow = !isFriend;
    const payload = { user_id: user_id, to_follow: to_follow, follow: follow };
    const res = await fetch(`${API_URL}/api/users/follow`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();

    if (!response.error) {
      if (follow) {
        friends.data.followed_by_ids.push(user_id);
        followersCount++;
      } else {
        friends.data.followed_by_ids.splice(
          friends.data.followed_by_ids.indexOf(user_id)
        );
        followersCount--;
      }
      isFriend = !isFriend;
    }
  };
</script>

<Container>
  <div class="px-10">
    <div class="grid pt-4 justify-center items-center">
      <h1 class="card-title pb-2 text-2xl">
        Viewing {userInfo.display_name}'s Public Profile
      </h1>
      <h2 class="pb-2 text-xl text-neutral opacity-80 italic font-semibold">
        {userInfo.class}
      </h2>
      {#if userInfo.id != user_id}
        <button
          class="btn btn-sm md:btn-xs btn-info text-info bg-opacity-30 rounded-full text-xs"
          class:btn-outline={!isFriend}
          on:click={() => {
            updateFollows(userInfo.id);
          }}>{isFriend ? "Following" : "Follow"}</button
        >
      {/if}
    </div>
    <div class="flex pt-3 justify-center items-center ">
      <div
        class="stats stats-vertical w-9/12 md:w-7/12 md:stats-horizontal shadow mt-4"
      >
        <div class="stat place-items-center ">
          <div class="stat-figure text-secondary" />
          <div class="stat-title text-sm">Public notes</div>
          <div class="stat-value ">
            {userNotesCount}
          </div>
        </div>

        <div class="stat place-items-center ">
          <div class="stat-figure text-secondary" />
          <div class="stat-title">Followers</div>
          <div class="stat-value ">{followersCount}</div>
        </div>

        <div class="stat place-items-center ">
          <div class="stat-figure text-secondary" />
          <div class="stat-title">Following</div>
          <div class="stat-value ">{followingCount}</div>
        </div>
      </div>
    </div>
    <div class="grid pt-20 justify-center items-center">
      <h1 class="card-title pb-2 text-2xl">
        {userInfo.display_name}'s Public Notes
      </h1>
    </div>
    <NotesContainer data={userNotes ?? ""} {API_URL} />
  </div>
</Container>
