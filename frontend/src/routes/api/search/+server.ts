import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ url }) => {
    const query = url.searchParams.get("query")
    const type = url.searchParams.get("type")

    if (!query || !type) {
        return json(
            { error: "Missing query parameter" },
            { status: 400 },
        )
    }

    const res = await fetch(`http://localhost:8080/search?query=${encodeURIComponent(query)}&type=${encodeURIComponent(type)}`)
    const data = await res.json()

    return json(data)
}
