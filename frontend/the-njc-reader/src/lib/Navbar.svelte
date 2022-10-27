<script>
  import Icon from "svelte-awesome";
  import user from "svelte-awesome/icons/user";

  let searchTerm = "";
</script>

<div class="navbar bg-primary fixed top-0 z-10">
  <div class="navbar-start">
    <div class="dropdown">
      <!-- svelte-ignore a11y-label-has-associated-control -->
      <label tabindex="0" class="btn btn-ghost btn-circle">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 6h16M4 12h16M4 18h7"
          /></svg
        >
      </label>
      <ul
        tabindex="0"
        class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-52"
      >
        <li><a data-sveltekit-reload href="/articles/1">News Feed</a></li>
        <li><a href="/columns/1">Longer Reads</a></li>
        <li><a href="/about">About</a></li>
        <li><a href="/feedback">Feedback</a></li>
      </ul>
    </div>
  </div>
  <div class="navbar-center">
    <a
      data-sveltekit-reload
      href="/articles/1"
      class="btn btn-ghost normal-case text-xl">The NJC Reader</a
    >
  </div>
  <div class="navbar-end">
    <form action="/search/{searchTerm}">
      <div class="form-control px-3">
        <input
          type="text"
          placeholder="Search"
          class="input input-bordered"
          bind:value={searchTerm}
          name="term"
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
