import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ url, locals }) => {
    const token = await locals.auth().getToken()
    const query = url.searchParams.get("query")
    const type = url.searchParams.get("type")
    const page = url.searchParams.get("page")

    if (!query || !type) {
        return json(
            { error: "Missing query parameter" },
            { status: 400 },
        )
    }

    let u = `http://localhost:8080/api/search?query=${encodeURIComponent(query)}&type=${encodeURIComponent(type)}`

    if (page !== null) {
        u += `&page=${page}`
    }

    const res = await fetch(u, {
        method: "GET",
        headers: {
            Authorization: `Bearer ${token}`,
        },
    })
    const data = await res.json()

    return json(data)
}
