<script lang="ts">
    import Card from "./Card.svelte";

    const {
        movie,
        add = false,
        onclick,
    }: { movie: any; add?: boolean; onclick?: () => void } = $props();
</script>

<Card
    src={movie.poster}
    alt={`${movie.title} movie poster`}
    title={movie.title}
    {add}
    {onclick}
/>

<style>
</style>
