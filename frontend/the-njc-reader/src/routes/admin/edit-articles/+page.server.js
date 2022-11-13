export async function load({ fetch, cookies }) {
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
  const articles = data.data;

  return {
    articles: articles,
  };
}

export const actions = {
  edit: async ({ request, cookies }) => {
    const formData = await request.formData()
    const url = formData.get("url")
    const title = formData.get("title")
    const tags = formData.get("tags")
    const date = formData.get("date")
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
      })
    }
    const payload = [{
      id: id,
      title: title,
      url: url,
      tags: tags,
      date: date
    }];

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    const res = await fetch("http://localhost:8080/api/admin/articles/update", {
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
