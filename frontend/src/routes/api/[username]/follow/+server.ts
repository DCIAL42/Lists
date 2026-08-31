import { env } from "$env/dynamic/private";
import type { Follow } from "$lib/types";
import { json, type RequestHandler } from "@sveltejs/kit";

export const POST: RequestHandler = async ({ params, locals }) => {
    const token = await locals.auth().getToken()
    const res = await fetch(`${env.BACKEND_URL}/${params.username}/follow`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        },
    })
    const follow: Follow = await res.json()

    return json(follow)
}

export const DELETE: RequestHandler = async ({ params, locals }) => {
    const token = await locals.auth().getToken()
    const res = await fetch(`${env.BACKEND_URL}/${params.username}/follow`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        },
    })
    const follow: Follow = await res.json()

    return json(follow)
}
