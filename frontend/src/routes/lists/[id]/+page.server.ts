import { error } from "@sveltejs/kit"
import type { PageServerLoad } from "./$types"
import type { List } from "$lib/types"
import { env } from "$env/dynamic/private"

export const load: PageServerLoad = async ({ params, locals }) => {
    const token = await locals.auth().getToken()

    const backendURL = env.BACKEND_URL
    const res = await fetch(`${backendURL}/lists/${params.id}`, {
        method: 'GET',
        headers: {
            Authorization: `Bearer ${token}`
        }
    })

    if (!res.ok) {
        const errorData = await res.json().catch(() => ({}))

        throw error(res.status, errorData.error || `Failed to fetch list ${params.id}`)
    }

    const list: List = await res.json()

    return { list }
}
