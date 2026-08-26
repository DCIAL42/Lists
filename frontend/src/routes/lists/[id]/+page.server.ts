import { error } from "@sveltejs/kit"
import type { PageServerLoad } from "./$types"
import type { List } from "$lib/types"

export const load: PageServerLoad = async ({ params, locals }) => {
    const token = await locals.auth().getToken()

    const res = await fetch(`http://localhost:8080/api/lists/${params.id}`, {
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
