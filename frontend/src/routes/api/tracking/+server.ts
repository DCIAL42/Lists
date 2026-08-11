import type { TrackingPayload } from "$lib/types";
import { json, type RequestHandler } from "@sveltejs/kit";

export const POST: RequestHandler = async ({ request, locals }) => {
    const auth = locals.auth();
    const token = await auth.getToken()
    const url = `http://localhost:8080/api/tracking`

    const item: TrackingPayload = await request.json()
    console.log(item)

    const res = await fetch(url, {
        method: 'POST',
        body: JSON.stringify(item),
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
        }
    })

    const data = await res.json()

    return json(data)
}
