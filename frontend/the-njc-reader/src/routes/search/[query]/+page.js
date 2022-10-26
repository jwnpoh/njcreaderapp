export async function load({ fetch, params }) {
  const queryURL = `http://localhost:8080/api/articles/find?term=${params.query}`;
  const res = await fetch(queryURL);
  const data = await res.json();
  const articles = data.data;

  console.log(articles);
  return {
    articles
  };
}
