<script>
  import "../app.css";
  import Navbar from "$lib/Navbar.svelte";
  import Container from "$lib/Container.svelte";

  let signedIn;
</script>

<header>
  <Navbar {signedIn} />
</header>
<main>
  <Container>
    <slot />
  </Container>
</main>

<footer class="footer footer-center p-4 bg-primary text-white">
  <div>
    <p>Copyright © 2022 National Junior College</p>
  </div>
</footer>
