import { json, type RequestHandler } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";

export const POST: RequestHandler = async ({ params, locals }) => {
    const token = await locals.auth().getToken()

    const backendURL = env.BACKEND_URL
    const res = await fetch(`${backendURL}/like`, {
        method: 'POST',
        body: JSON.stringify({
            list_id: Number(params.id)
        }),
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        },
    })

    const data = await res.json()

    if (!res.ok) {
        return json(data, {
            status: res.status
        })
    }

    return json(data)
}

export const DELETE: RequestHandler = async ({ params, locals }) => {
    const token = await locals.auth().getToken();
    const url = `${env.BACKEND_URL}/like/${params.id}`

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
