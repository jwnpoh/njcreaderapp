import { redirect } from "@sveltejs/kit";
import "dotenv/config"

export const load = async ({ cookies }) => {
  // delete session token
  const session = cookies.get("session");

  const res = await logOut(session)
  console.log(res)

  // eat the cookie
  cookies.set("session", "", {
    path: "/",
    Expires: new Date(0),
  })

  throw redirect(302, "/login")
}

const logOut = async (session) => {
  const myHeaders = new Headers();
  myHeaders.append('Content-Type', 'application/json');
  myHeaders.append('Authorization', 'Bearer ' + session);

  const requestOptions = {
    method: "POST",
    headers: myHeaders,
  }

  const res = await fetch(`${process.env.API_URL}/api/users/logout`, requestOptions)
  const response = await res.json()
  return response
}
