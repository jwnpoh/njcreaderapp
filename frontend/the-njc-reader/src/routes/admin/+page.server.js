import "dotenv/config";
import { redirect } from "@sveltejs/kit"

export async function load({cookies, fetch, locals }) {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }
  if (locals.user.role != "admin") {
    throw redirect(302, "/profile")
  }

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

  const queryURL = `${process.env.API_URL}/api/admin/stats/get-stats`;
  const res = await fetch(queryURL, {
      headers: myHeaders,
  });

  const data = await res.json();

  const stats = data.data;

  console.log(stats)

  return {
    API_URL: `${process.env.API_URL}`,
    user: locals.user,
    stats: stats
  }
}

// export const load = async ({ locals }) => {
//   if (!locals.user.loggedIn) {
//     throw redirect(302, "/login")
//   }

//   if (locals.user.role != "admin") {
//     throw redirect(302, "/profile")
//   }
//   return {
//     user: locals.user
//   }
// }
