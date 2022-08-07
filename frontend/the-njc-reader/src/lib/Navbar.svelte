<script>
  import Icon from "svelte-awesome";
  import user from "svelte-awesome/icons/user";

  let searchTerm = "";

  const submitSearch = () => {
    console.log("submitted search for", searchTerm);
    searchTerm = "";
  };
</script>

<div class="navbar bg-base-100">
  <div class="flex-1">
    <a class="btn btn-ghost normal-case text-xl">The NJC Reader</a>
  </div>
  <div class="flex-none gap-2">
    <form on:submit|preventDefault={submitSearch}>
      <div class="form-control  ">
        <input
          type="text"
          placeholder="Search"
          class="input input-bordered"
          bind:value={searchTerm}
        />
      </div>
    </form>

    <div class="dropdown dropdown-end">
      <div tabindex="0" class="btn btn-ghost btn-circle avatar">
        <div class="w-10 rounded-full shadow-md py-2">
          <Icon data={user} scale={1.3} />
        </div>
      </div>
      <ul
        tabindex="0"
        class="mt-3 p-2 shadow menu menu-compact dropdown-content bg-base-100 rounded-box w-52"
      >
        <li>
          <a class="justify-between">
            Profile
            <span class="badge">New</span>
          </a>
        </li>
        <li><a>Settings</a></li>
        <li><a>Logout</a></li>
      </ul>
    </div>
  </div>
</div>
