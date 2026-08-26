import { env } from "$env/dynamic/private";
import type { List } from "$lib/types";
import { json, type RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ locals }) => {
    const { userId } = locals.auth();

    const backendURL = env.BACKEND_URL
    let u = `${backendURL}/users/${userId}/lists`

    const res = await fetch(u)

    const data: List = await res.json()

    return json(data)
}
