import "dotenv/config"
let data;
let q;

export async function load() {
  const error = data.error;
  const articles = data.data;
  const message = data.message;
  const query = q;

  return {
    error: error,
    articles: articles,
    message: message,
    query: query,
  };
}

export const actions = {
  default: async ({ request }) => {
    const formData = await request.formData()
    const query = await formData.get("query")

    const queryURL = `${process.env.API_URL}/api/articles/find?term=${query}`;
    const res = await fetch(queryURL);
    const response = await res.json();

    data = response;
    q = query;
  }
}
