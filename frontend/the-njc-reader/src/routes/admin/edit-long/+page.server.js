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
  const articles = data.data;

  return {
    articles: articles,
    API_URL: `${process.env.API_URL}`
  };
}

export const actions = {
  edit: async ({ request, cookies }) => {
    const formData = await request.formData()
    const url = formData.get("url")
    const title = formData.get("title")
    const topic = formData.get("topic")
    const id = parseInt(formData.get("id"))

    if (url.length < 1 || title.length < 1) {
      return invalid(400, {
        error: true,
        message: "All fields must be filled.",
        id,
        url,
        title,
        topic,
      })
    }
    const payload = {
      id: id,
      title: title,
      url: url,
      topic: topic,
    };

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    const res = await fetch(`${process.env.API_URL}/api/admin/long/update`, {
      method: "PUT",
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