<script>
  import "../app.css";
  import Navbar from "$lib/Navbar.svelte";
  import Container from "$lib/Container.svelte";
</script>

<main>
  <header>
    <Navbar />
  </header>
  <Container>
    <slot />
  </Container>
</main>

<footer class="footer footer-center p-4 bg-base-300 text-base-content">
  <div>
    <p>Copyright © 2022 National Junior College</p>
  </div>
</footer>
