<script lang="ts">
    import Card from "./Card.svelte";
    import type { Movie } from "$lib/types";

    const {
        movie,
        add = false,
        onclick,
    }: { movie: Movie; add?: boolean; onclick?: () => void } = $props();
</script>

<Card
    src={movie.cover}
    alt={`${movie.title} movie poster`}
    title={movie.title}
    {add}
    {onclick}
/>

<style>
</style>
