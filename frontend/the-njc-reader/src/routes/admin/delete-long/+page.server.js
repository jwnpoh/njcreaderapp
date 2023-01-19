import { redirect } from "@sveltejs/kit"
import "dotenv/config"

export async function load({ fetch, cookies, locals }) {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }
  if (locals.user.role != "admin") {
    throw redirect(302, "/profile")
  }

  const session = cookies.get("session")

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  const queryURL = `${process.env.API_URL}/api/admin/long/update`;
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

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    const res = await fetch(`${process.env.API_URL}/api/admin/long/delete`, {
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
