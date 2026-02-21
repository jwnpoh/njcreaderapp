import { redirect, invalid } from "@sveltejs/kit"
import "dotenv/config"

export async function load({ fetch, cookies, locals }) {
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

  // Fetch recent articles (reusing the existing endpoint but limiting to 25)
  const queryURL = `${process.env.API_URL}/api/admin/articles/update`;
  const res = await fetch(queryURL, {
    method: "GET",
    headers: myHeaders
  });

  const data = await res.json();
  if (data.error) {
    throw redirect(302, "/login")
  }

  // Limit to 25 most recent
  const articles = data.data.slice(0, 25);
  
  // Get today's timestamp (start of day)
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const todayTimestamp = Math.floor(today.getTime() / 1000);

  return {
    articles: articles,
    todayTimestamp: todayTimestamp
  };
}

export const actions = {
  send: async ({ request, cookies }) => {
    const formData = await request.formData()
    let payload = [];
    
    for (var pair of formData.entries()) {
      if (pair[0] === "selection") {
        payload.push(JSON.parse(pair[1]))
      }
    }

    if (payload.length === 0) {
      return invalid(400, { 
        failed: true, 
        message: "No articles selected" 
      })
    }

    const session = cookies.get("session")

    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + session);

    const res = await fetch(`${process.env.API_URL}/api/admin/articles/telegram-send`, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: myHeaders,
    });

    const response = await res.json();
    if (response.error) {
      return invalid(400, { 
        failed: true, 
        message: response.message 
      })
    }
    
    return {
      success: true
    }
  }
}
