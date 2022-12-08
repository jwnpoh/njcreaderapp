import "dotenv/config"
import { redirect } from "@sveltejs/kit"

export const load = async ({ fetch, locals, cookies }) => {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }

  const userID = locals.user.id;

  const session = cookies.get("session");

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  const getNotes = async () => {
    const res = await fetch(`${process.env.API_URL}/api/posts/notebook?user=${userID}`, {
      method: "GET",
      headers: myHeaders,
    })

    const notes = await res.json()
    return notes
  }

  // getFollowing
  const getFollowing = async () => {
    const res = await fetch(`${process.env.API_URL}/api/posts/following?user=${userID}`, {
      method: "GET",
      headers: myHeaders,
    })

    const notes = await res.json()
    return notes
  }

  // getPublic
  const getDiscover = async () => {
    const res = await fetch(`${process.env.API_URL}/api/posts/public?user=all`, {
      method: "GET",
      headers: myHeaders,
    })

    const notes = await res.json()
    let note_likes = {}
    for (var note of notes.data) {
      note_likes[note.id] = note.likes
    }
    locals.user.note_likes = note_likes
    return notes
  }

  // get user's liked posts
  const getLikedNotes = async () => {
    const res = await fetch(`${process.env.API_URL}/api/posts/get-liked-posts?user=${userID}`, {
      method: "GET",
      headers: myHeaders,
    })

    const likes = await res.json()
    locals.user.liked_articles = likes.data
    return likes
  }

  return {
    notes: getNotes(),
    following: getFollowing(),
    discover: getDiscover(),
    liked_notes: getLikedNotes(),
    API_URL: `${process.env.API_URL}`
  }
}

