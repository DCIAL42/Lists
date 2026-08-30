<script lang="ts">
    import type { PageData } from "./$types";
    import Carousel from "$lib/components/Carousel.svelte";

    let { data }: { data: PageData } = $props();
</script>

<main>
    <h2>Top Lists</h2>
    <hr />
    <Carousel lists={data.top.lists} />
</main>

<style>
    main {
        padding: 10px;
    }
</style>
