import "dotenv/config"
import { redirect } from "@sveltejs/kit"

export async function load({ fetch,locals}) {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }

  const queryURL = `${process.env.API_URL}/api/long`;
  const res = await fetch(queryURL);
  const data = await res.json();

  const topics = data.data;

  return {
    topics: topics
  };
}