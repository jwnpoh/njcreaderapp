import "dotenv/config"
import { redirect } from "@sveltejs/kit"

export async function load({ fetch, locals, params }) {
  const queryURL = `${process.env.API_URL}/api/long/${params.topic}`;
  const res = await fetch(queryURL);
  const data = await res.json();

  const articles = data.data;
  const topic = params.topic;

  return {
    articles: articles,
    topic: topic,
  };
}
