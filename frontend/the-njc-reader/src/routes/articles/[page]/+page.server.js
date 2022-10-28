export async function load({ fetch, params }) {
  const queryURL = `http://localhost:8080/api/articles/${params.page}`;
  // export async function load({ fetch }) {
  // const queryURL = `http://localhost:8080/api/articles/1`;
  const res = await fetch(queryURL);
  const data = await res.json();

  const articles = data.data;
  const page = params.page;

  return {
    articles: articles,
    page: page,
  };
}
