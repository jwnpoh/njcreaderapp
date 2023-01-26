import "dotenv/config"

export const handle = async ({ event, resolve }) => {
  const session = event.cookies.get("session")

  if (!session) {
    event.locals.user = {
      email: "",
      role: "",
      display_name: "",
      class: "",
      loggedIn: false,
      session: "",
    }
    return await resolve(event)
  }

  const user = await getUser(session)

  if (!user.error) {
    const note_likes = await getDiscover(session)
    const liked_notes = await getLikedNotes(session, user.data.id)
    event.locals.user = {
      id: user.data.id,
      email: user.data.email,
      role: user.data.role,
      display_name: user.data.display_name,
      class: user.data.class,
      loggedIn: true,
      session: session,
      note_likes: note_likes,
      liked_notes: liked_notes,
    }
  } else {
    event.locals.user = {
      id: "",
      email: "",
      role: "",
      display_name: "",
      class: "",
      loggedIn: false,
      session: "",
    }
  }

  return await resolve(event)
}

const getUser = async (session) => {
  const myHeaders = new Headers();
  myHeaders.append('Content-Type', 'application/json');
  myHeaders.append('Authorization', 'Bearer ' + session);

  const requestOptions = {
    method: "POST",
    headers: myHeaders,
  }

  const res = await fetch(`${process.env.API_URL}/api/users/get-user`, requestOptions)
  const user = await res.json()
  return user
}

const getDiscover = async (session) => {
  const myHeaders = new Headers();
  myHeaders.append('Content-Type', 'application/json');
  myHeaders.append('Authorization', 'Bearer ' + session);

  const res = await fetch(`${process.env.API_URL}/api/posts/public?user=all`, {
    method: "GET",
    headers: myHeaders,
  })

  const notes = await res.json()

  let note_likes = {}
  if (notes.data) {
    for (var note of notes.data) {
      note_likes[note.id] = note.likes
    }
  }
  return note_likes
}

const getLikedNotes = async (session, userID) => {
  const myHeaders = new Headers();
  myHeaders.append('Content-Type', 'application/json');
  myHeaders.append('Authorization', 'Bearer ' + session);

  const res = await fetch(`${process.env.API_URL}/api/posts/get-liked-posts?user=${userID}`, {
    method: "GET",
    headers: myHeaders,
  })

  const likes = await res.json()

  return likes.data
}
