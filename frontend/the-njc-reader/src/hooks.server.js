export const handle = async ({ event, resolve }) => {
  const session = event.cookies.get("session")
  console.log("session token: ", session)

  if (!session) {
    return await resolve(event)
  }

  const user = await getUser(session)

  if (user) {
    event.locals.user = {
      email: user.data.email,
      role: user.data.role,
      loggedIn: true
    }
  }

  console.log("event locals user: ", event.locals.user)

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
