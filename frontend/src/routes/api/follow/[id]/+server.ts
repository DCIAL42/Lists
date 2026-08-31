import { env } from "$env/dynamic/private";
import type { Follow } from "$lib/types";
import { json, type RequestHandler } from "@sveltejs/kit";

export const DELETE: RequestHandler = async ({ locals, params }) => {
    const token = await locals.auth().getToken()
    const res = await fetch(`${env.BACKEND_URL}/follow/${params.id}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        },
    })
    const follow: Follow = await res.json()
    return json(follow)
}
