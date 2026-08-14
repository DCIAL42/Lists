<script lang="ts">
    import ListItemCard from "$lib/ListItemCard.svelte";
    import SearchBox from "$lib/SearchBox.svelte";
    import type { MediaItem } from "$lib/types";

    let items: MediaItem[] = $state([]);
</script>

<SearchBox bind:items />

<div class="results">
    {#if items}
        {#each items as item}
            <ListItemCard {item} />
        {/each}
    {/if}
</div>

<style>
    .results {
        padding: 5px;
    }
</style>
