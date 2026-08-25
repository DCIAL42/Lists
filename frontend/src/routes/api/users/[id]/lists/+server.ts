import type { List } from "$lib/types";
import { json, type RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ locals }) => {
    const { userId } = locals.auth();

    let u = `http://localhost:8080/api/users/${userId}/lists`

    const res = await fetch(u)

    const data: List = await res.json()

    return json(data)
}
