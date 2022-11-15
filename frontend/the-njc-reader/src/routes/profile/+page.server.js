import { invalid, redirect } from "@sveltejs/kit"

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
    const email = locals.user.email

    console.log(email, display_name, new_password, old_password)
    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    let payload = {
      email: email,
      old_password: old_password,
      new_password: new_password,
      display_name: display_name,
    }

    console.log("payload: ", payload)
    const res = await fetch("http://localhost:8080/api/users/update-user", {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    if (response.error) {
      return invalid(400, { failed: true, message: response.message })
    }
    locals.user.display_name = display_name
    console.log(locals.user)
    return {
      success: true,
      sent: true
    }
  }
}
