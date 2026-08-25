<script lang="ts">
    import SearchBar from "$lib/SearchBar.svelte";
    import type { MediaItem } from "$lib/types";

    let items: MediaItem[] = $state([]);
</script>

<main>
    <!-- <SearchBar bind:items /> -->
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        gap: 5px;
        padding: 5px;
    }
</style>
