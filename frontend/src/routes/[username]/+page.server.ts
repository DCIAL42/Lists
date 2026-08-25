import type { PageServerLoad } from "./$types"

export const load: PageServerLoad = async ({ locals, params }) => {
    const { userId, ...auth } = locals.auth();
    const token = await auth.getToken()

    const trackingRes = await fetch(`http://localhost:8080/api/tracking?type=album|movie|game&status=backlog`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
        }
    })

    const listsRes = await fetch(`http://localhost:8080/api/${params.username}/lists`)

    if (!trackingRes.ok) {
        throw new Error(`Failed to fetch tracking items`)
    }

    if (!listsRes.ok) {
        throw new Error(`Failed to fetch lists`)
    }

    const items = await trackingRes.json()
    const listsData = await listsRes.json()

    return { items, listsData }
}
