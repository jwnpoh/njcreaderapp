import { redirect } from "@sveltejs/kit"

export const load = async ({ locals }) => {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }

  if (locals.user.role != "admin") {
    throw redirect(302, "/profile")
  }
  return {
    user: locals.user
  }
}
