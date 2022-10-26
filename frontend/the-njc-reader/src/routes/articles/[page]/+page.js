export async function load({ params }) {
  const queryURL = `http://localhost:8080/api/articles/${params.page}`;
  const res = await fetch(queryURL);
  const data = await res.json();
  const articles = data.data;

  return {
    articles,
  };
}

// export const load = async ({ fetch, params }) => {

//   const queryURL = `http://localhost:8080/api/articles/${params.page}`;
//   const res = await fetch(queryURL);
//   const data = res.json();
//   const articles = data.data;

//   return {
//     articles,
//   }
// }
