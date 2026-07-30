import type { PageServerLoad } from "./$types"

export const load: PageServerLoad = async ({ params }) => {
    const response = await fetch(`http://localhost:8080/api/list/${params.id}`)

    if (!response.ok) {
        throw new Error(`Failed to fetch list ${params.id}`)
    }

    const list = await response.json()

    return {
        list: list
    }
}
