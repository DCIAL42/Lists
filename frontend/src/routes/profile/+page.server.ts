import type { PageServerLoad } from "./$types"

export const load: PageServerLoad = async ({ locals }) => {
    const auth = locals.auth();
    const token = await auth.getToken()
    const response = await fetch(`http://localhost:8080/api/tracking?type=album&status=done`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
        }
    })

    if (!response.ok) {
        throw new Error(`Failed to fetch tracking items`)
    }

    const items = await response.json()

    return items
}
