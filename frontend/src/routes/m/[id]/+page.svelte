<script lang="ts">
    import Carousel from "$lib/components/Carousel.svelte";
    import ItemCard from "$lib/ItemCard.svelte";
    import { isAlbum, isArtist, toDuration } from "$lib/utils";
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
            {media.name}
        </h1>
        {#if isAlbum(media)}
            <a href={`/m/${media.data.artist.id}`} class="subtitle">
                {media.data.artist.name}
            </a>
            <hr />
            {#each media.data.tracks as track}
                <p>
                    {track.name} - {toDuration(track.duration)}
                </p>
            {/each}
        {:else if isArtist(media)}
            <h2>Albums</h2>
            <Carousel>
                {#each media.data.albums as a}
                    <div>
                        <ItemCard item={a} small style="width: 250px;" />
                    </div>
                {/each}
            </Carousel>
        {/if}
    </div>
</main>

<style>
    main.grid {
        display: grid;
        grid-template-columns: minmax(0, 1fr) minmax(0, 2fr) minmax(0, 1fr);
        gap: 50px;
        align-items: start;
    }
</style>
