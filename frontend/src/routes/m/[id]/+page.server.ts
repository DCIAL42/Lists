import { env } from "$env/dynamic/private";
import type { MediaItem } from "$lib/types";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, locals }) => {
    const token = await locals.auth().getToken()
    const res = await fetch(`${env.BACKEND_URL}/media/${params.id}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
        }
    })
    const media: MediaItem = await res.json()
    return { media }
}
