import type { RequestHandler } from "./$types";
import type { Rating } from "$lib/types";
import { json } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";

export const POST: RequestHandler = async ({ request, locals }) => {
    const token = await locals.auth().getToken()
    let rating: Rating = await request.json()

    const backendURL = env.BACKEND_URL
    const res = await fetch(`${backendURL}/rating`, {
        method: 'POST',
        body: JSON.stringify(rating),
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

// export const GET: RequestHandler = async ({ url }) => {
//     let page = url.searchParams.get("page")
//
//     const backendURL = env.BACKEND_URL
//     let u = `${backendURL}/lists?order=desc`
//
//     if (page !== null) {
//         u += `&page=${page}`
//     }
//
//     const res = await fetch(u)
//
//     if (!res.ok) {
//         const body = await res.json()
//         return json(body)
//     }
//
//     const data = await res.json()
//
//     return json(data)
// }
