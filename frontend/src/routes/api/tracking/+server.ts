import type { TrackingPayload } from "$lib/types";
import { json, type RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ url, locals }) => {
    const token = await locals.auth().getToken();
    const mediaType = url.searchParams.get("type")
    const status = url.searchParams.get("status")
    const page = url.searchParams.get("page")

    let targetUrl = `http://localhost:8080/api/tracking?type=${mediaType}&status=${status}`
    if (page !== null) {
        targetUrl += `&page=${page}`
    }

    const res = await fetch(targetUrl, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
        }
    })

    const data = await res.json()

    return json(data)
}

export const POST: RequestHandler = async ({ request, locals }) => {
    const token = await locals.auth().getToken();
    const url = `http://localhost:8080/api/tracking`

    const item: TrackingPayload = await request.json()

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
