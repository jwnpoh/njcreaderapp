import { invalid, redirect } from "@sveltejs/kit"
import "dotenv/config"

let data = "";

export const load = async ({ url, locals, cookies }) => {
  let articleID = url.searchParams.get('articleID')
  let returnURL = `/notes/add-note?articleID=${articleID}`

  if (!locals.user.loggedIn) {
    throw redirect(302, `/login?redirect=${returnURL}`)
  }

  if (articleID) {
    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    const queryURL = `${process.env.API_URL}/api/posts/get-article?id=${articleID}`;
    const res = await fetch(queryURL, {
      method: "GET",
      headers: myHeaders,
    });

    const response = await res.json();

    data = response.data;
  }

  return {
    data
  }
}

export const actions = {
  add: async ({ request, cookies, locals }) => {
    const formData = await request.formData()
    const articleID = formData.get("article_id")
    const articleTitle = formData.get("article_title")
    const articleURL = formData.get("article_url")
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

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    let payload = {
      article_id: articleID,
      article_title: articleTitle,
      article_url: articleURL,
      tldr: tldr,
      examples: examples,
      notes: notes,
      tags: tags,
      public: makePublic,
      date: date,
      user_id: userID,
      likes: likes
    }

    const res = await fetch(`${process.env.API_URL}/api/posts/insert`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    if (response.error) {
      return invalid(400, {
        error: true, message: response.message,
        articleID,
        articleTitle,
        articleURL,
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
