import "dotenv/config"

export async function load({ fetch }) {
  const queryURL = `${process.env.API_URL}/api/articles/1`;
  const res = await fetch(queryURL);
  const data = await res.json();

  const articles = data.data;

  return {
    articles: articles,
    page: 1,
  };
}