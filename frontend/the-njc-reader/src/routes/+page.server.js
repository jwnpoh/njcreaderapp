export async function load({ fetch }) {
  const queryURL = `http://localhost:8080/api/articles/1`;
  const res = await fetch(queryURL);
  const data = await res.json();

  const articles = data.data;

  return {
    articles: articles,
  };
}
