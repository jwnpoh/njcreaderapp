<script>
  export let user;

  let role;
  if (user) {
    role = user.role;
  }
  console.log(role);

  let showMenu;
</script>

<div class="navbar bg-primary text-white fixed top-0 z-50 ">
  <div class="md:navbar-start">
    <div class="dropdown">
      <button
        class="btn btn-ghost btn-circle"
        on:click={() => (showMenu = !showMenu)}
      >
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
      </button>
      <div class={!showMenu ? "hidden" : ""}>
        <ul
          class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-primary rounded-box w-52"
        >
          <li>
            <div class="md:hidden text-xl font-medium">The NJC Reader</div>
          </li>
          <li><a data-sveltekit-reload href="/articles/1">News Feed</a></li>
          <li><a href="/columns/1">Longer Reads</a></li>
          <li><a href="/about">About</a></li>
          {#if user}
            <li><a href="/feedback">Feedback</a></li>
            <li>
              <a href="/preferences">Preferences</a>
            </li>
            {#if role === "admin"}
              <li>
                <a href="/admin">Admin dashboard</a>
              </li>
            {/if}
          {/if}
          <li>
            <a
              data-sveltekit-reload
              class="bg-secondary text-black"
              href={user ? "/logout" : "/login"}
              >{user ? "Log out" : "Log in"}</a
            >
          </li>
        </ul>
      </div>
    </div>
  </div>
  <div class="md:navbar-center invisible md:visible">
    <a class="btn btn-ghost normal-case text-xl" data-sveltekit-reload href="/"
      >The NJC Reader</a
    >
  </div>
  <div class="md:navbar-end">
    <form method="POST" action="/search">
      <div class="form-control px-3 fixed top-2 right-1 text-black">
        <input
          type="text"
          placeholder="Search"
          class="input input-bordered"
          name="query"
        />
      </div>
    </form>
  </div>
</div>
