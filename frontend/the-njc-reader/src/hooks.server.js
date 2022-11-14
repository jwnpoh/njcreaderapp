
export const handle = async ({ event, resolve }) => {
  const session = event.cookies.get("session")

  if (!session) {
    event.locals.user = {
      email: "",
      role: "",
      name: "",
      loggedIn: false,
      session: "",
    }
    return await resolve(event)
  }

  const user = await getUser(session)

  if (!user.error) {
    event.locals.user = {
      email: user.data.email,
      role: user.data.role,
      name: user.data.name,
      loggedIn: true,
      session: session
    }
  } else {
    event.locals.user = {
      email: "",
      role: "",
      name: "",
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

  const res = await fetch("http://localhost:8080/api/users/get-user", requestOptions)
  const user = await res.json()
  return user
}
