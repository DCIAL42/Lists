<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import Link from "$lib/components/Link.svelte";
    import Skeleton from "$lib/components/Skeleton.svelte";
    import ListPreview from "$lib/ListPreview.svelte";
    import type { ListMeta } from "$lib/types";
    import { onMount } from "svelte";

    let listsData: { lists: ListMeta[]; next: string; count: number } = $state({
        lists: [],
        next: "/api/lists",
        count: 0,
    });
    let lists = $derived(listsData.lists);
    let loading = $state(false);
    let page = $state(0);

    onMount(() => {
        nextPage("/api/lists");
    });

    async function nextPage(url: string) {
        loading = true;
        const res = await fetch(url);
        if (!res.ok) {
            throw new Error("Error fetching next page");
        }
        listsData = await res.json();
        page++;
        loading = false;
    }
</script>

<Link href="/lists/create">Create new list</Link>

<div class="lists">
    {#if loading}
        {#each Array(10) as _}
            <Skeleton width={270} height={170} />
        {/each}
    {:else}
        {#each lists as list}
            <ListPreview {list} />
        {/each}
    {/if}
</div>

{#if listsData.count - 10 * (page - 1) >= 10}
    <Button onclick={() => nextPage(listsData.next)}>Next</Button>
{/if}

<style>
    .lists {
        display: grid;
        grid-template-columns: repeat(5, 1fr);
        gap: 10px;
        justify-items: center;
        padding: 20px;
    }
</style>
