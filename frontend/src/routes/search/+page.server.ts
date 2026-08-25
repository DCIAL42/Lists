import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url }) => {
    const query = url.searchParams.get("query")
    const tab = url.searchParams.get("type")
    if (tab === null || query === null) {
        return {
            items: [], query: "", tab: ""
        }
    }

    const res = await fetch(`http://localhost:8080/api/search?query=${query}&type=${tab}`)

    if (!res.ok) {
        throw new Error(`Failed to fetch tracking items`)
    }

    const items = await res.json()

    return { items, query, tab }
}
