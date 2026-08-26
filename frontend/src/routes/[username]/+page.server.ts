import type { UserResponse } from "$lib/types";
import type { PageServerLoad } from "./$types"

export const load: PageServerLoad = async ({ locals, params }) => {
    const { userId, ...auth } = locals.auth();
    const token = await auth.getToken()

    const userRes = await fetch(`http://localhost:8080/api/${params.username}`)
    const userData: UserResponse = await userRes.json()

    const listsRes = await fetch(`http://localhost:8080/api/${params.username}/lists?order=desc`)

    if (!listsRes.ok) {
        throw new Error(`Failed to fetch lists`)
    }

    const listsData = await listsRes.json()

    if (userId === userData.id) {
        const trackingRes = await fetch(`http://localhost:8080/api/tracking?type=album|movie|game&status=backlog`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${token}`
            }
        })

        if (!trackingRes.ok) {
            throw new Error(`Failed to fetch tracking items`)
        }

        const trackingData = await trackingRes.json()

        return { userData, trackingData, listsData }
    }

    return { userData, listsData }
}
