import "dotenv/config"
// import { redirect } from "@sveltejs/kit"

export async function load({ fetch, locals }) {
  // if (!locals.user.loggedIn) {
  //   throw redirect(302, "/login")
  // }

  const queryURL = `${process.env.API_URL}/api/long/all`;
  const res = await fetch(queryURL);
  const data = await res.json();

  const articles = data.data;

  // Group articles by topic
  const groupedArticles = articles.reduce((acc, article) => {
    if (!acc[article.topic]) {
      acc[article.topic] = [];
    }
    acc[article.topic].push(article);
    return acc;
  }, {});

  return {
    groupedArticles
  };
}
