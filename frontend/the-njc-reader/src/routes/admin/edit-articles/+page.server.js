import { redirect } from "@sveltejs/kit"
import "dotenv/config"

export async function load({ fetch, cookies, locals }) {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }

  const session = cookies.get("session")

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  const queryURL = `${process.env.API_URL}/api/admin/articles/update`;
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
    const tags = formData.get("tags")
    const date = formData.get("date")
    const must_read = formData.get("must_read")
    const id = formData.get("id")

    if (url.length < 1 || title.length < 1 || tags.length < 1) {
      return invalid(400, {
        error: true,
        message: "All fields must be filled.",
        id,
        url,
        title,
        tags,
        date,
        must_read,
      })
    }
    const payload = [{
      id: id,
      title: title,
      url: url,
      tags: tags,
      date: date,
      must_read: must_read
    }];

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    const res = await fetch(`${process.env.API_URL}/api/admin/articles/update`, {
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
