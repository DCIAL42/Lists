import type { List } from "$lib/types";
import type { RequestHandler } from "./$types";
import { json } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ locals }) => {
    const auth = locals.auth();
    const userId = auth.userId

    let u = `http://localhost:8080/api/lists/user/${userId}`

    const res = await fetch(u)

    const data: List = await res.json()

    return json(data)
}
