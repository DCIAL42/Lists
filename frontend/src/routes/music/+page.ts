import { json } from "@sveltejs/kit";

export async function load({ fetch, url }) {
    const query = url.searchParams.get('query') ?? '';

    const res = await fetch(`http://localhost:8080/search?query=${query}`)

    if (!res.ok) {
        throw new Error('Failed to fetch')
    }

    return {
        data: await res.json()
    }
}
