export async function GET({ params }) {
  const queryURL = `http://localhost:8080/api/articles/${params.page}`;
  const res = await fetch(queryURL);
  const data = await res.json();
  const articles = data.data;

  return {
    body: {
      articles,
    }
  };
}

