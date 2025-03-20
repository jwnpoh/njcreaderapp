<script>
  var date = new Date();


  let greeting = "Hello";
  const h = date.getHours();
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
  export let showMenu;

  let role;
  if (user.loggedIn) {
    role = user.role;
  }

  const toggleMenu = () => {
    showMenu = !showMenu;
  };
</script>

<div class="navbar bg-primary-focus text-white fixed top-0 z-50 ">
  <div class="md:navbar-start">
    <div class="dropdown">
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
      <!-- svelte-ignore a11y-label-has-associated-control -->
      <label tabindex="0" class="btn btn-ghost btn-circle" on:click={toggleMenu}>
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
      <div class={!showMenu ? "hidden" : ""}>
        <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
        <ul tabindex="0" class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-primary-focus rounded-box w-64" >
          <li>
            <div
              class="md:hidden py-1 text-xl font-medium  hover:cursor-default"
            >
              <a href="/"> The NJC Reader </a>
            </div>
          </li>
            <li>
              <div
                class="md:py-1 px-4 font-semibold text-lg  hover:cursor-default"
              >
                {greeting}!
              </div>
            </li>
            <li><a data-sveltekit-reload href="/articles">News Feed</a></li>
            <li><a href="/long">Long Reads</a></li>
          {#if user.loggedIn}
            <div class="divider px-2 before:bg-white opacity-30 after:bg-white opacity-30" />
            {#if role === "admin"}
              <li>
                <a href="/admin">Admin dashboard</a>
              </li>
            {/if}
          {/if}
          <li>
            <a
              data-sveltekit-reload
              class="text-secondary"
              href={user.loggedIn ? "/logout" : "/login"}
              >{user.loggedIn ? "Log out" : "Admin log in"}</a
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
