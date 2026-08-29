<script lang="ts">
    import Tabs from "$lib/components/Tabs.svelte";
    import ItemCard from "$lib/ItemCard.svelte";
    import type { MediaItem, MediaType } from "$lib/types";
    import type { HTMLAttributes } from "svelte/elements";

    let {
        items = $bindable(),
        tab = $bindable("movie"),
        small = false,
        add = false,
        handleAdd = () => {},
        ...rest
    }: {
        items: MediaItem[];
        tab?: MediaType;
        small?: boolean;
        add?: boolean;
        handleAdd?: (item: MediaItem) => void;
    } & HTMLAttributes<HTMLDivElement> = $props();
</script>

<div class="results" class:small {...rest}>
    <Tabs tabs={["movie", "album"]} bind:selected={tab} />
    {#each items as _, i}
        <ItemCard
            bind:item={items[i]}
            {small}
            {add}
            onAddClick={handleAdd}
            onmousedown={(e) => e.preventDefault()}
        />
    {/each}
</div>

<style>
    .small {
        overflow: auto;
        height: 50vh;
        width: fit-content;
        position: absolute;
        left: 50%;
        transform: translateX(-50%);
    }

    .results {
        display: flex;
        flex-direction: column;
        border: 1px solid var(--border);
        padding: 5px;
        gap: 10px;
        margin-top: 10px;
        background-color: var(--background);
        z-index: 9999;
        width: 100%;
    }
</style>
