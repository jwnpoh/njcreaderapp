<script>
  import { page } from "$app/stores";
  import "../app.css";
  import Navbar from "$lib/Navbar.svelte";
  import Container from "$lib/Container.svelte";

  let user = $page.data.user;

  let showMenu;
</script>

<header>
  <Navbar {user} bind:showMenu />
</header>
<main>
  <Container>
    <slot />
  </Container>
</main>
