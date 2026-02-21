// import "dotenv/config"
// let data;
// let q;

// export async function load() {
//   const error = data.error;
//   const articles = data.data;
//   const message = data.message;
//   const query = q;

//   return {
//     error: error,
//     articles: articles,
//     message: message,
//     query: query,
//   };
// }

// export const actions = {
//   default: async ({ request }) => {
//     const formData = await request.formData()
//     const query = await formData.get("query")

//     const queryURL = `${process.env.API_URL}/api/articles/find?term=${query}`;
//     const res = await fetch(queryURL);
//     const response = await res.json();

//     data = response;
//     q = query;
//   }
// }

// import "dotenv/config"
// let data;
// let q;

// export async function load({ url }) {
//   // Check if there's a URL param first
//   const urlQuery = url.searchParams.get('term');
  
//   if (urlQuery) {
//     const queryURL = `${process.env.API_URL}/api/articles/find?term=${urlQuery}`;
//     const res = await fetch(queryURL);
//     const response = await res.json();
    
//     // Set module variables just like the action does
//     data = response;
//     q = urlQuery;
//   }
  
//   // Return data (either from URL param or from previous form action)
//   const error = data?.error;
//   const articles = data?.data;
//   const message = data?.message;
//   const query = q;

//   return {
//     error: error,
//     articles: articles,
//     message: message,
//     query: query,
//   };
// }

// export const actions = {
//   default: async ({ request }) => {
//     const formData = await request.formData()
//     const query = await formData.get("query")

//     const queryURL = `${process.env.API_URL}/api/articles/find?term=${query}`;
//     const res = await fetch(queryURL);
//     const response = await res.json();

//     data = response;
//     q = query;
//   }
// }

import "dotenv/config"

export async function load({ url }) {
  const query = url.searchParams.get('query') || url.searchParams.get('term');
  
  if (!query) {
    return {
      error: false,
      articles: null,
      message: null,
      query: null,
    };
  }
  
  const queryURL = `${process.env.API_URL}/api/articles/find?term=${query}`;
  const res = await fetch(queryURL);
  const response = await res.json();
  
  return {
    error: response.error,
    articles: response.data,
    message: response.message,
    query: query,
  };
}
