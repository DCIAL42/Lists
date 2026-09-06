<script lang="ts">
    import ItemCard from "$lib/ItemCard.svelte";
    import { isAlbum, toDuration } from "$lib/utils";
    import type { PageData } from "./$types";

    let { data = $bindable() }: { data: PageData } = $props();
    let media = $state(data.media);
    $effect(() => {
        media = data.media;
    });
</script>

<main class="grid">
    <ItemCard bind:item={media} small />
    <div>
        <h1>
            {media.title}
        </h1>
        {#if isAlbum(media)}
            <p class="subtitle">
                {media.data.artist}
            </p>
            <hr />
            {#each media.data.tracks as track}
                <p>
                    {track.title} - {toDuration(track.duration)}
                </p>
            {/each}
        {/if}
    </div>
</main>

<style>
    main.grid {
        display: grid;
        grid-template-columns: 1fr 2fr 1fr;
        gap: 50px;
        align-items: start;
    }
</style>
