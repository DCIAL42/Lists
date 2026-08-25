<script lang="ts">
    import type { Album } from "$lib/types";
    import Card from "$lib/components/Card.svelte";

    const {
        album,
        cover,
        add = false,
        onclick,
    }: {
        album: Album;
        cover: string;
        add?: boolean;
        onclick?: () => void;
    } = $props();
</script>

<Card
    src={cover}
    alt={`${album.title} by ${album.artist} album cover`}
    title={album.title}
    subtitle={album.artist}
    {add}
    {onclick}
/>
