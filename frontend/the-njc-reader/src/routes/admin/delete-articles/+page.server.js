import { redirect } from "@sveltejs/kit"

export async function load({ fetch, cookies, locals }) {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }

  const session = cookies.get("session")

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  const queryURL = `http://localhost:8080/api/admin/articles/update`;
  const res = await fetch(queryURL, {
    method: "GET",
    headers: myHeaders
  });

  const data = await res.json();
  if (data.error) {
    throw redirect(302, "/login")
  }

  const articles = data.data;

  return {
    articles: articles,
  };
}

export const actions = {
  delete: async ({ request, cookies }) => {
    const formData = await request.formData()
    let payload = [];
    for (var pair of formData.entries()) {
      payload.push(pair[1])
    }
    console.log(payload);

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    const res = await fetch("http://localhost:8080/api/admin/articles/delete", {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    if (response.error) {
      return invalid(400, { failed: true, message: response.data })
    }
    return {
      success: true,
      sent: true
    }
  }
}
