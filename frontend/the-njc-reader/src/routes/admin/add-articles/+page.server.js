import "dotenv/config";
import { invalid, redirect } from "@sveltejs/kit";

let queue = [];

export const load = ({ locals }) => {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }
  if (locals.user.role != "admin") {
    throw redirect(302, "/profile")
  }
  return {
    queue: queue,
    API_URL: `${process.env.API_URL}`
  }
}

export const actions = {
  queue: async ({ request }) => {
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
        queue
      })
    }

    const input = {
      index: Math.floor(Math.random() * 100) + 1,
      title: title,
      url: url,
      tags: tags,
      date: date,
      must_read: must_read
    };
    queue.push(input);
    return {
      success: true,
    }
  },
  remove: async ({ request }) => {
    const formData = await request.formData()
    const index = formData.get("index")

    queue = queue.filter(input => input.index != index)
    return {
      success: true,
    }
  },
  send: async ({ cookies }) => {
    if (queue.length < 1) {
      return invalid(400, { failed: true, message: "No articles queued. Add one or more articles before adding to the database." })

    }
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
      return invalid(400, { failed: true, message: response.data })
    }
    if (!response.error) {
      queue = [];
    }
    return {
      success: true,
      sent: true
    }
  },
}
