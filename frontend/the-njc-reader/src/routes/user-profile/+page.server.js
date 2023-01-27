import { invalid, redirect } from "@sveltejs/kit"
import "dotenv/config"

export const load = async ({ locals }) => {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }

  const user = locals.user
  return {
    user
  }
}

export const actions = {
  updateUser: async ({ locals, cookies, request }) => {
    const formData = await request.formData()
    const new_password = formData.get("new_password")
    const old_password = formData.get("old_password")
    const display_name = formData.get("display_name")
    const userID = locals.user.id

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    let payload = {
      user_id: userID,
      old_password: old_password,
      new_password: new_password,
      display_name: display_name,
    }

    const res = await fetch(`${process.env.API_URL}/api/users/update-user`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    if (response.error) {
      return invalid(400, { failed: true, message: response.message })
    }
    locals.user.display_name = display_name
    return {
      success: true,
      sent: true
    }
  }
}
