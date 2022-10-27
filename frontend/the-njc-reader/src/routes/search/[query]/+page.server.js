export async function load({ fetch, params }) {
  const queryURL = `http://localhost:8080/api/articles/find?term=${params.query}`;
  const res = await fetch(queryURL);
  const data = await res.json();
  const error = data.error;
  const articles = data.data;
  const message = data.message;
  const query = params.query;

  console.log(articles);
  return {
    error: error,
    articles: articles,
    message: message,
    query: query,
  };
}
