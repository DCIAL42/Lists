import type { RequestHandler } from "./$types";
import type { List, ListPayload } from "$lib/types";
import { json } from "@sveltejs/kit";

export const POST: RequestHandler = async ({ request, locals }) => {
    const token = await locals.auth().getToken()
    let list: List = await request.json()

    let body: ListPayload = {
        title: list.title,
        items: list.items.map(item => ({ ...item, external_id: String(item.external_id) }))
    }

    const res = await fetch('http://localhost:8080/api/lists', {
        method: 'POST',
        body: JSON.stringify(body),
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        },
    })

    const data = await res.json()

    return json(data)
}

export const GET: RequestHandler = async ({ url }) => {
    let page = url.searchParams.get("page")

    let u = `http://localhost:8080/api/lists?order=desc`

    if (page !== null) {
        u += `&page=${page}`
    }

    const res = await fetch(u)

    const data = await res.json()

    return json(data)
}
