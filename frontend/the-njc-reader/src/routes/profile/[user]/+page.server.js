import { redirect } from "@sveltejs/kit"
import "dotenv/config"

export const load = async ({ locals, cookies, fetch, params }) => {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }

  const userID = locals.user.id;
  const session = cookies.get("session");

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  const getUserInfo = async () => {
    const queryURL = `${process.env.API_URL}/api/users/${params.user}`;
    const res = await fetch(queryURL, {
      method: "GET",
      headers: myHeaders,
    });
    const data = await res.json();
    const userInfo = data.data;

    return userInfo
  }

  const getUserNotes = async () => {
    const res = await fetch(`${process.env.API_URL}/api/posts/public?user=${params.user}`, {
      method: "GET",
      headers: myHeaders,
    })

    const notes = await res.json()
    if (!notes.data) {
      notes.message = "This user has not created any note yet."
      return notes
    }
    let note_likes = {}
    for (var note of notes.data) {
      note_likes[note.id] = note.likes
    }
    locals.user.note_likes = note_likes
    return notes
  }

  const getLikedNotes = async () => {
    const res = await fetch(`${process.env.API_URL}/api/posts/get-liked-posts?user=${userID}`, {
      method: "GET",
      headers: myHeaders,
    })

    const likes = await res.json()
    locals.user.liked_notes = likes.data
    return likes
  }

  const getFriends = async () => {
    const res = await fetch(`${process.env.API_URL}/api/users/friends?user=${params.user}`, {
      method: "GET",
      headers: myHeaders,
    })

    const friends = await res.json()
    return friends
  }

  return {
    userInfo: getUserInfo(),
    userNotes: getUserNotes(),
    liked_notes: getLikedNotes(),
    friends: getFriends(),
    API_URL: `${process.env.API_URL}`
  }
}

