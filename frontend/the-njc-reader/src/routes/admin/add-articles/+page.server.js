import { invalid } from "@sveltejs/kit"

let queue = [];

export const load = () => {
  return {
    queue: queue
  }
}

export const actions = {
  queue: async ({ request }) => {
    const formData = await request.formData()
    const url = formData.get("url")
    const title = formData.get("title")
    const tags = formData.get("tags")
    const date = formData.get("date")

    if (url.length < 1 || title.length < 1 || tags.length < 1) {
      return invalid(400, {
        error: true,
        message: "All fields must be filled.",
        url,
        title,
        tags,
        date,
        queue
      })
    }

    const input = {
      index: Math.floor(Math.random() * 100) + 1,
      title: title,
      url: url,
      tags: tags,
      date: date
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

    const res = await fetch("http://localhost:8080/api/admin/articles/insert", {
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
