import { env } from "$env/dynamic/private";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url }) => {
    const query = url.searchParams.get("query")
    const tab = url.searchParams.get("type")
    if (tab === null || query === null) {
        return {
            items: [], query: "", tab: ""
        }
    }

    const backendURL = env.BACKEND_URL
    const res = await fetch(`${backendURL}/search?query=${query}&type=${tab}`)

    if (!res.ok) {
        throw new Error(`Failed to fetch tracking items`)
    }

    const items = await res.json()

    return { items, query, tab }
}
