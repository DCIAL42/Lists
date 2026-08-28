import { env } from "$env/dynamic/private";
import { json, type RequestHandler } from "@sveltejs/kit";

export const DELETE: RequestHandler = async ({ params, locals }) => {
    const token = await locals.auth().getToken();
    const url = `${env.BACKEND_URL}/lists/${params.id}`

    const res = await fetch(url, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
        }
    })

    const data = await res.json()

    return json(data)
}

export const PATCH: RequestHandler = async ({ request, params, locals }) => {
    const token = await locals.auth().getToken();
    const url = `${env.BACKEND_URL}/lists/${params.id}`

    const payload = await request.json()

    const res = await fetch(url, {
        method: 'PATCH',
        body: JSON.stringify(payload),
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        }
    })

    const data = await res.json()

    return json(data)
}
