import { redirect } from "@sveltejs/kit"

export const load = async ({ locals }) => {
  if (!locals.user.loggedIn) {
    throw redirect(302, "/login")
  }

  return {
  }
}
