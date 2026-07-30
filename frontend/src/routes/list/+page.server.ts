import type { PageServerLoad } from "./$types"

export const load: PageServerLoad = async () => {
    const response = await fetch(`http://localhost:8080/api/list/`)

    if (!response.ok) {
        throw new Error(`Failed to fetch lists`)
    }

    const lists = await response.json()

    return {
        lists: lists
    }
}
