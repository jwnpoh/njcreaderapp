import "dotenv/config";
import { invalid, redirect } from "@sveltejs/kit";

export const load = ({ locals }) => {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }
  if (locals.user.role != "admin") {
    throw redirect(302, "/profile")
  }
  return {
    API_URL: `${process.env.API_URL}`
  }
}

export const actions = {
  send: async ({ cookies, request }) => {
    let queue = [];
    const formData = await request.formData()
    const url = formData.get("url")
    const title = formData.get("title")
    const tags = formData.get("tags")
    const date = formData.get("date")
    const must_read = formData.get("must_read")

    if (url.length < 1 || title.length < 1 || tags.length < 1) {
      return invalid(400, {
        error: true,
        message: "All fields must be filled.",
        url,
        title,
        tags,
        date,
        must_read,
      })
    }

    const input = {
      title: title,
      url: url,
      tags: tags,
      date: date,
      must_read: must_read
    };
    queue.push(input);

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    let payload = queue

    const res = await fetch(`${process.env.API_URL}/api/admin/articles/insert`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    if (response.error) {
      return invalid(400, { 
        failed: true, 
        message: response.data,
        url,
        title,
        tags,
        date,
        must_read,
     })
    }
    return {
      success: true,
      sent: true
    }
  },
}
