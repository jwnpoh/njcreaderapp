import "dotenv/config";
import { invalid, redirect } from "@sveltejs/kit";


export async function load({ fetch, locals }) {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }
  if (locals.user.role != "admin") {
    throw redirect(302, "/profile")
  }

  const queryURL = `${process.env.API_URL}/api/long`;
  const res = await fetch(queryURL);
  const data = await res.json();

  const topics = data.data;

  return {
    API_URL: `${process.env.API_URL}`,
    topics: topics
  }
}

export const actions = {
  add: async ({ request, cookies }) => {
    const formData = await request.formData()
    const input = formData.get("input")
    console.log(input)

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    let payload = input

    const res = await fetch(`${process.env.API_URL}/api/admin/long/insert`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    if (response.error) {
      return invalid(400, { failed: true, message: response.data })
    }

    return {
      success: true,
      sent: true
    }
  },
}
