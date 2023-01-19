import "dotenv/config"

export async function load({ fetch }) {
  const queryURL = `${process.env.API_URL}/api/long`;
  const res = await fetch(queryURL);
  const data = await res.json();

  const topics = data.data;

  return {
    topics: topics
  };
}