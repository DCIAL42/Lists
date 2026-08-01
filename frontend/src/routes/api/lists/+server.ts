import type { RequestHandler } from "./$types";
import type { List, ListPayload } from "$lib/types";
import { json } from "@sveltejs/kit";

export const POST: RequestHandler = async ({ request }) => {
    let list: List = await request.json()

    let body: ListPayload = {
        title: list.title,
        created_by: list.created_by,
        items: list.items.map(item => ({ ...item, external_id: String(item.external_id) }))
    }

    const url = `http://localhost:8080/api/lists`

    const res = await fetch(url, {
        method: 'POST',
        body: JSON.stringify(body),
        headers: {
            'Content-Type': 'application/json',
        },
    })

    const data = await res.json()

    return json(data)
}

export const GET: RequestHandler = async ({ url }) => {
    let page = url.searchParams.get("page")

    let u = `http://localhost:8080/api/lists`
    if (page !== null && page !== undefined) {
        u += `?page=${page}`
    }

    const res = await fetch(u)

    const data = await res.json()

    return json(data)
}
