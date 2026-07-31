<script lang="ts">
    import type { Album } from "$lib/types";
    import Card from "$lib/components/Card.svelte";

    const {
        album,
        add = false,
        onclick,
    }: { album: Album; add?: boolean; onclick?: () => void } = $props();
</script>

<Card
    src={album.cover}
    alt={`${album.title} by ${album.artist} album cover`}
    title={album.title}
    subtitle={album.artist}
    {add}
    {onclick}
/>
