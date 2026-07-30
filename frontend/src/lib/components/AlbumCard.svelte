<script lang="ts">
    import Card from "./Card.svelte";

    const {
        album,
        add = false,
        onclick,
    }: { album: any; add?: boolean; onclick?: () => void } = $props();
</script>

<Card
    src={album.cover}
    alt={`${album.title} by ${album.artist} album cover`}
    title={album.title}
    subtitle={album.artist}
    {add}
    {onclick}
/>
