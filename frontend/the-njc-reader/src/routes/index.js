export async function GET() {
  const res = await fetch("http://localhost:8080/api/articles/1");
  const data = await res.json();
  const articles = data.data;

  return {
    body: {
      articles,
    }
  };
}
