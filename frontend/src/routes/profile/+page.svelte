<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import ListItemCard from "$lib/ListItemCard.svelte";
    import type { MediaType, TrackingStatus, MediaItem } from "$lib/types.js";

    let { data }: { data: { items: MediaItem[] } } = $props();
    let items = $derived(data.items);
    let tab: TrackingStatus = $state("backlog");
    let mediaTypes: MediaType[] = $state(["album", "movie", "game"]);
    let showToast = $state(false);
    let loading = $state(true);

    async function getTrackingItems(
        mediaTypes: MediaType[],
        status: TrackingStatus,
    ) {
        loading = true;
        const res = await fetch(
            `/api/tracking?type=${mediaTypes.join("|")}&status=${status}`,
        );

        const data = await res.json();

        items = data.items;
        loading = false;
    }

    $effect(() => {
        getTrackingItems(mediaTypes, tab);
    });
</script>

<main>
    <div class="tabs">
        <Button
            variant="ghost"
            selected={tab === "backlog"}
            onclick={() => (tab = "backlog")}>Backlog</Button
        >
        <hr />
        <Button
            variant="ghost"
            selected={tab === "paused"}
            onclick={() => (tab = "paused")}>Paused</Button
        >
        <hr />
        <Button
            variant="ghost"
            selected={tab === "done"}
            onclick={() => (tab = "done")}>Done</Button
        >
    </div>

    <div class="types">
        <Button
            variant="ghost"
            selected={mediaTypes.includes("album")}
            onclick={() => {
                if (mediaTypes.includes("album")) {
                    mediaTypes = mediaTypes.filter((e) => e !== "album");
                } else {
                    mediaTypes.push("album");
                }
            }}
        >
            Album
        </Button>
        <hr />
        <Button
            variant="ghost"
            selected={mediaTypes.includes("movie")}
            onclick={() => {
                if (mediaTypes.includes("movie")) {
                    mediaTypes = mediaTypes.filter((e) => e !== "movie");
                } else {
                    mediaTypes.push("movie");
                }
            }}
        >
            Movie
        </Button>
        <hr />
        <Button
            variant="ghost"
            selected={mediaTypes.includes("game")}
            onclick={() => {
                if (mediaTypes.includes("game")) {
                    mediaTypes = mediaTypes.filter((e) => e !== "game");
                } else {
                    mediaTypes.push("game");
                }
            }}
        >
            Game
        </Button>
    </div>

    <div class="items">
        {#each items as item, index}
            <ListItemCard
                {item}
                {index}
                onRemoveClick={(i: number) => {
                    items = items.toSpliced(i, 1);
                    showToast = true;
                    setTimeout(() => {
                        showToast = false;
                    }, 3000);
                }}
            />
        {/each}
    </div>

    {#if showToast}
        <div class="toast">
            <p>Item deleted</p>
            <Button>Undo</Button>
        </div>
    {/if}
</main>

<style>
    .toast {
        position: absolute;
        border: 1px solid var(--border);
        border-radius: 4px;
        top: 100%;
        left: 100%;
        width: 10%;
        padding: 5px;
        transform: translate(-105%, -110%);
    }

    main {
        display: flex;
        flex-direction: column;
        gap: 5px;
        padding: 5px;
    }

    .tabs,
    .types {
        display: grid;
        margin-inline: auto;
        grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
    }

    .items {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr;
        gap: 5px;
    }
</style>
