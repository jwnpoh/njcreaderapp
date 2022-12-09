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

  // getPublic
  const getDiscover = async () => {
    const res = await fetch(`${process.env.API_URL}/api/posts/public?user=all`, {
      method: "GET",
      headers: myHeaders,
    })

    const notes = await res.json()

    return notes
  }

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



  return {
    notes: getNotes(),
    following: getFollowing(),
    discover: getDiscover(),
    API_URL: `${process.env.API_URL}`
  }
}

