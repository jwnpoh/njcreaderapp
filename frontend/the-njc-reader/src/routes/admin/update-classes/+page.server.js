import "dotenv/config";
import Papa from "papaparse";
import { invalid, redirect } from "@sveltejs/kit";


export const load = ({ locals }) => {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }
  if (locals.user.role != "admin") {
    throw redirect(302, "/profile")
  }
  return {
  }
}

export const actions = {
  add: async ({ request, cookies }) => {
    const formData = await request.formData()
    const input = await formData.get("input")
    const inputString = await input.text()
    const parsed = Papa.parse(inputString , {header: true})

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    let payload = parsed.data

    const res = await fetch(`${process.env.API_URL}/api/admin/users/update-classes`, {
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
      sent: true,
      message: response.message
    }
  },
}