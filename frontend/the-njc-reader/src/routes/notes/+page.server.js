export const load = async ({ fetch, locals, cookies }) => {
  const userID = locals.user.id;

  const session = cookies.get("session");

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + session);

  const getNotes = async () => {
    const res = await fetch(`http://localhost:8080/api/posts/notebook?user=${userID}`, {
      method: "GET",
      headers: myHeaders,
    })

    const notes = await res.json()
    return notes
  }

  // getFollowing
  const getFollowing = async () => {
    const res = await fetch(`http://localhost:8080/api/posts/following?user=${userID}`, {
      method: "GET",
      headers: myHeaders,
    })

    const notes = await res.json()
    return notes
  }

  // getPublic
  const getDiscover = async () => {
    const res = await fetch(`http://localhost:8080/api/posts/public?user=all`, {
      method: "GET",
      headers: myHeaders,
    })

    const notes = await res.json()
    return notes

  }

  return {
    notes: getNotes(),
    following: getFollowing(),
    discover: getDiscover()
  }
}

