<script lang="ts">
    import type { PageData } from "./$types";
    import Carousel from "$lib/components/Carousel.svelte";
    import ListPreview from "$lib/ListPreview.svelte";

    let { data }: { data: PageData } = $props();
</script>

<main>
    <h2>Top Lists</h2>
    <hr />
    <Carousel>
        {#each data.top.lists as list}
            <ListPreview {list} />
        {/each}
    </Carousel>
</main>

<style>
    main {
        padding: 10px;
    }
</style>
