import "dotenv/config"

export async function load({ fetch, params }) {
  const queryURL = `${process.env.API_URL}/api/articles/${params.page}`;
  const res = await fetch(queryURL);
  const data = await res.json();

  const articles = data.data;
  const page = params.page;

  return {
    articles: articles,
    page: page,
  };
}
