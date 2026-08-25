import { error } from "@sveltejs/kit"
import type { PageServerLoad } from "./$types"

export const load: PageServerLoad = async ({ params }) => {
    const res = await fetch(`http://localhost:8080/api/lists/${params.id}`)

    if (!res.ok) {
        const errorData = await res.json().catch(() => ({}))

        throw error(res.status, errorData.error || `Failed to fetch list ${params.id}`)
    }

    const list = await res.json()

    return { list }
}
