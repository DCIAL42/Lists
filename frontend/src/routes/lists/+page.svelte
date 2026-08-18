<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import Link from "$lib/components/Link.svelte";
    import ListPreview from "$lib/ListPreview.svelte";
    import type { List } from "$lib/types";

    let { data }: { data: { lists: List[]; next: string } } = $props();
    let lists = $derived(data.lists);

    async function nextPage(url: string) {
        const res = await fetch(url);
        if (!res.ok) {
            throw new Error("Error fetching next page");
        }
        data = await res.json();
    }

    $inspect(data);
</script>

<Link href="/lists/create">Create new list</Link>

{#each lists as list}
    <ListPreview {list} />
{/each}

<Button onclick={() => nextPage(data.next)}>Next</Button>
