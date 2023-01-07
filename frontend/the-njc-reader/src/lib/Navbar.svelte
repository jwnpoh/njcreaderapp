<script>
  import dayjs from "dayjs";
  dayjs().format();

  let greeting = "Hello";
  const h = dayjs().hour();
  if (h >= 18 && h < 24) {
    greeting = "Good evening";
  }
  if (h >= 0 && h < 12) {
    greeting = "Good morning";
  }
  if (h >= 12 && h < 18) {
    greeting = "Good afternoon";
  }

  export let user;

  let role;
  if (user.loggedIn) {
    role = user.role;
  }

  export let showMenu;
  const toggleMenu = () => {
    showMenu = !showMenu;
  };
</script>

<div class="navbar bg-primary text-white fixed top-0 z-50 ">
  <div class="md:navbar-start">
    <div class="dropdown">
      <button class="btn btn-ghost btn-circle" on:click={toggleMenu}>
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
          class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-primary rounded-box w-64"
        >
          <li>
            <div
              class="md:hidden py-1 text-xl font-medium  hover:cursor-default"
            >
              The NJC Reader
            </div>
          </li>
          <li class="md:hidden" />
          {#if user.loggedIn}
            <li>
              <div
                class="md:py-1 px-2 font-semibold text-lg  hover:cursor-default"
              >
                {greeting}, {user.display_name}!
              </div>
            </li>
            <li />
            <li>
              <a href="/notes">Notebook</a>
            </li>
          {/if}
          <li><a data-sveltekit-reload href="/articles/1">News Feed</a></li>
          <li><a href="/columns/1">Long Reads</a></li>
          <li><a href="/about">About</a></li>
          {#if user.loggedIn}
            <li><a href="/feedback">Feedback</a></li>
            <li />
            {#if role === "admin"}
              <li>
                <a href="/admin">Admin dashboard</a>
              </li>
            {/if}
            <li>
              <a href="/user-profile">Profile</a>
            </li>
          {/if}
          <li>
            <a
              data-sveltekit-reload
              class="bg-secondary text-black"
              href={user.loggedIn ? "/logout" : "/login"}
              >{user.loggedIn ? "Log out" : "Log in"}</a
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
          placeholder="Search articles"
          class="input input-bordered "
          name="query"
        />
      </div>
    </form>
  </div>
</div>
