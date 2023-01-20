import { invalid, redirect } from '@sveltejs/kit';
import "dotenv/config"

export const load = async () => {
  return {
  }
}

export const actions = {
  default: async ({ cookies, request }) => {
    const formData = await request.formData()
    const email = formData.get('email')

    if (
      typeof email != "string" || !email) {
      return invalid(400, { invalid: true })
    }


    const response = await fetch(`${process.env.API_URL}/api/users/reset-password`, {
      method: 'POST',
      body: JSON.stringify(email),
      headers: {
        'content-type': 'application/json'
      },
    });
    const res = await response.json()

    if (res.error) {
      return invalid(400, { message: res.data })
    }

    return {
      success: true
    }
  }
}
