<script lang="ts">
    import SearchBar from "$lib/SearchBar.svelte";
    import type { MediaItem, MediaType } from "$lib/types";

    let {
        data = $bindable(),
    }: { data: { items: MediaItem[]; query: string; tab?: MediaType } } =
        $props();
    let items = $state.raw(data.items);
    let query = $state.raw(data.query);
    let tab = $state.raw(data.tab ?? "movie");
</script>

<SearchBar bind:items bind:query bind:tab floating={false} />

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
