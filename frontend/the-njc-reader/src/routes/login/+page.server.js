import { invalid, redirect } from '@sveltejs/kit';
import "dotenv/config"

export const load = async ({ url }) => {
  const returnURL = url.searchParams.get('redirect')

  return {
    returnURL
  }
}

export const actions = {
  default: async ({ cookies, request }) => {
    const formData = await request.formData()
    const email = formData.get('email')
    const password = formData.get('password')
    const returnURL = formData.get('returnURL')

    if (
      typeof email != "string" ||
      typeof password != "string" ||
      !email || !password
    ) {
      return invalid(400, { invalid: true })
    }

    const input = {
      email: email,
      password: password,
    }

    const authResponse = await authenticate(JSON.stringify(input))
    const res = await authResponse.json()

    if (res.error) {
      return invalid(400, { credentials: true })
    }

    let token = res.data.token
    let expiry = res.data.expiry

    cookies.set("session", token, {
      path: "/",
      HttpOnly: true,
      Expires: expiry,
      secure: process.env.NODE_ENV == "production",
    })

    if (returnURL) {
      throw redirect(302, `${returnURL}`)
    } else {
      throw redirect(302, "/articles")
    }
  }
}

const authenticate = async (userObjString) => {
  const response = await fetch(`${process.env.API_URL}/api/auth`, {
    method: 'POST',
    body: userObjString,
    headers: {
      'content-type': 'application/json'
    },
  });

  return response;
}
