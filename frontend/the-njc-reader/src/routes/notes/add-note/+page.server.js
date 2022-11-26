import { invalid, redirect } from "@sveltejs/kit"

let data;

export const load = async ({ locals }) => {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }

  return {
    data: data
  }
}

export const actions = {
  newnote: async ({ request, locals, cookies }) => {
    if (!locals.user.loggedIn) {
      throw redirect(302, "/login")
    }

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    const formData = await request.formData()
    const articleID = await formData.get("article_id")

    const queryURL = `http://localhost:8080/api/posts/get-article?id=${articleID}`;
    const res = await fetch(queryURL, {
      method: "GET",
      headers: myHeaders,
    });

    const response = await res.json();

    data = response.data;
  },
  add: async ({ request, cookies, locals }) => {
    const formData = await request.formData()
    const articleID = formData.get("article_id")
    const tldr = formData.get("tldr")
    const examples = formData.get("examples")
    const notes = formData.get("notes")
    const makePublic = formData.get("make_public")
    const date = new Date().toDateString()
    const userID = locals.user.id;
    const likes = 0;

    const tldrTags = formData.get("tldr_tags")
    const examplesTags = formData.get("examples_tags")
    const notesTags = formData.get("notes_tags")

    let tags = [];

    if (tldrTags.length > 0) {
      tags = tags.concat(tldrTags)
    }
    if (examplesTags.length > 0) {
      tags = tags.concat(examplesTags)
    }
    if (notesTags.length > 0) {
      tags = tags.concat(notesTags)
    }
    console.log(tags)

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    let payload = {
      article_id: articleID,
      tldr: tldr,
      examples: examples,
      notes: notes,
      tags: tags,
      make_public: makePublic,
      date: date,
      user_id: userID,
      likes: likes
    }

    const res = await fetch("http://localhost:8080/api/posts/insert", {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    if (response.error) {
      return invalid(400, {
        error: true, message: response.message,
        articleID,
        tldr,
        examples,
        notes,
        tags,
        makePublic,
      })
    }
    return {
      success: true,
      sent: true
    }
  }
}
