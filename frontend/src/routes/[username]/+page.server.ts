import { env } from "$env/dynamic/private";
import type { UserResponse } from "$lib/types";
import type { PageServerLoad } from "./$types"

export const load: PageServerLoad = async ({ locals, params }) => {
    const { userId, ...auth } = locals.auth();
    const token = await auth.getToken()

    const backendURL = env.BACKEND_URL
    const userRes = await fetch(`${backendURL}/${params.username}`)
    const userData: UserResponse = await userRes.json()

    const listsRes = await fetch(`${backendURL}/${params.username}/lists?order=desc`)

    if (!listsRes.ok) {
        throw new Error(`Failed to fetch lists`)
    }

    const listsData = await listsRes.json()

    if (userId === userData.id) {
        const trackingRes = await fetch(`${backendURL}/tracking?type=album|movie|game&status=backlog`, {
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
