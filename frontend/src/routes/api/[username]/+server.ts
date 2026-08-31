import { env } from "$env/dynamic/private";
import type { Follow, ListsPreviewData, TrackingListData, UserResponse } from "$lib/types";
import { json, type RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ locals, params }) => {
    const { userId, ...auth } = locals.auth();
    const token = await auth.getToken()

    const backendURL = env.BACKEND_URL
    const userRes = await fetch(`${backendURL}/${params.username}`)
    const userData: UserResponse = await userRes.json()

    const listsRes = await fetch(`${backendURL}/${params.username}/lists?order=desc`)

    if (!listsRes.ok) {
        throw new Error(`Failed to fetch lists`)
    }

    const listsData: ListsPreviewData = await listsRes.json()

    if (userId === userData.id) {
        const trackingRes = await fetch(`${backendURL}/tracking?type=album|movie|game&status=backlog`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${token}`,
            }
        })

        if (!trackingRes.ok) {
            throw new Error(`Failed to fetch tracking items`)
        }

        const trackingData: TrackingListData = await trackingRes.json()

        return json({ self: true, userData, trackingData, listsData })
    }

    const followRes = await fetch(`${backendURL}/${params.username}/follow`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        },
    })
    const followData: Follow = await followRes.json()

    return json({ self: false, userData, listsData, followData })
}
