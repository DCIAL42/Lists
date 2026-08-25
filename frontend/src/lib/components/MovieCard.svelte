<script lang="ts">
    import Card from "./Card.svelte";
    import type { Movie } from "$lib/types";

    const {
        movie,
        cover,
        add = false,
        onclick,
    }: {
        movie: Movie;
        cover: string;
        add?: boolean;
        onclick?: () => void;
    } = $props();
</script>

<Card
    src={cover}
    alt={`${movie.title} movie poster`}
    title={movie.title}
    {add}
    {onclick}
/>

<style>
</style>
