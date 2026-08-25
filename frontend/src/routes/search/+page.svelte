<script lang="ts">
    import SearchBar from "$lib/SearchBar.svelte";
    import type { MediaItem, MediaType } from "$lib/types";

    let {
        data,
    }: { data: { items: MediaItem[]; query: string; tab: MediaType } } =
        $props();
    let items = $state(data.items);
    let query = $state(data.query);
    let tab = $state(data.tab);
</script>

<SearchBar bind:items bind:query bind:tab />

<!-- <div class="results"> -->
<!--     {#if items} -->
<!--         {#each items as item} -->
<!--             <ListItemCard {item} /> -->
<!--         {/each} -->
<!--     {/if} -->
<!-- </div> -->

<!-- <style> -->
<!--     .results { -->
<!--         padding: 5px; -->
<!--     } -->
<!-- </style> -->
