export async function load({ fetch, params }) {
  const queryURL = `http://localhost:8080/api/articles/${params.page}`;
  const res = await fetch(queryURL);
  const data = await res.json();

  const articles = data.data;
  const page = params.page;
  console.log(page)

  return {
    articles: articles,
    page: page,
  };
}
