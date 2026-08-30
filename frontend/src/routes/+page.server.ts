import { error } from "@sveltejs/kit"
import type { PageServerLoad } from "./$types"
import type { ListsPreviewData } from "$lib/types"
import { env } from "$env/dynamic/private"

export const load: PageServerLoad = async () => {
    const backendURL = env.BACKEND_URL
    const res = await fetch(`${backendURL}/lists?order=desc&order_by=likes`)

    if (!res.ok) {
        const errorData = await res.json().catch(() => ({}))

        throw error(res.status, errorData.error || `Failed to home page`)
    }

    const top: ListsPreviewData = await res.json()

    return { top }
}
