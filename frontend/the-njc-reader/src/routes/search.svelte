<script context="module">
  export async function load({ fetch }) {
    const res = await fetch(
      "http://localhost:8080/api/articles/find?term=" + term
    );
    const data = await res.json();
    const articles = data.data;

    return {
      props: articles,
    };
  }
</script>

<script>
  import CardContainer from "$lib/CardContainer.svelte";
  export let articles;
</script>

<CardContainer {articles} />
