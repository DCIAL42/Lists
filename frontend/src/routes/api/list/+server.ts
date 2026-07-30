import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

interface ListItem {
    external_id: string
    type: string
}

interface List {
    title: string
    created_by: string
    items: ListItem[]
}

export const POST: RequestHandler = async ({ request }) => {
    let list: List = await request.json()

    list.items = list.items.map(item => ({ ...item, external_id: String(item.external_id) }))

    const url = `http://localhost:8080/api/list`
    const res = await fetch(url, {
        method: 'POST',
        body: JSON.stringify(list),
        headers: {
            'Content-Type': 'application/json',
        },
    })
    const data = await res.json()
    return json(data)
}
