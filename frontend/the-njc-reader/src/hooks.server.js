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
    event.locals.user = {
      id: user.data.id,
      email: user.data.email,
      role: user.data.role,
      display_name: user.data.display_name,
      class: user.data.class,
      loggedIn: true,
      session: session
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
