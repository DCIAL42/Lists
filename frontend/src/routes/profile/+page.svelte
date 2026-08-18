<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import ListItemCard from "$lib/ListItemCard.svelte";
    import ListPreview from "$lib/ListPreview.svelte";
    import type {
        MediaType,
        TrackingStatus,
        MediaItem,
        List,
    } from "$lib/types.js";

    let { data }: { data: { items: MediaItem[]; lists: List[] } } = $props();
    let items = $derived(data.items);
    let tab: "lists" | "tracking" = $state("tracking");
    let trackingTab: TrackingStatus = $state("backlog");
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

        const data: { items: MediaItem[] } = await res.json();

        items = data.items;
        loading = false;
    }

    $effect(() => {
        if (tab === "tracking") {
            getTrackingItems(mediaTypes, trackingTab);
        }
    });
</script>

<main>
    <div class="tabs">
        <Button
            variant="ghost"
            selected={tab === "tracking"}
            onclick={() => (tab = "tracking")}>Tracking</Button
        >
        <hr />
        <Button
            variant="ghost"
            selected={tab === "lists"}
            onclick={() => (tab = "lists")}>Lists</Button
        >
    </div>

    {#if tab === "tracking"}
        <div class="tracking-tabs">
            <Button
                variant="ghost"
                selected={trackingTab === "backlog"}
                onclick={() => (trackingTab = "backlog")}>Backlog</Button
            >
            <hr />
            <Button
                variant="ghost"
                selected={trackingTab === "paused"}
                onclick={() => (trackingTab = "paused")}>Paused</Button
            >
            <hr />
            <Button
                variant="ghost"
                selected={trackingTab === "done"}
                onclick={() => (trackingTab = "done")}>Done</Button
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
    {/if}

    {#if tab === "tracking"}
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
    {:else if tab === "lists"}
        <div class="lists">
            {#each data.lists as list}
                <ListPreview {list} />
            {/each}
        </div>
    {/if}

    {#if showToast}
        <div class="toast">
            <p>Item deleted</p>
            <Button>Undo</Button>
        </div>
    {/if}
</main>

<style>
    .toast {
        position: fixed;
        border: 1px solid var(--border);
        border-radius: 4px;
        top: 100%;
        left: 100%;
        width: 10%;
        padding: 5px;
        transform: translate(-105%, -105%);
    }

    main {
        display: flex;
        flex-direction: column;
        gap: 5px;
        padding: 5px;
    }

    .tabs {
        display: grid;
        margin-inline: auto;
        grid-template-columns: repeat(3, 1fr);
    }

    .tracking-tabs,
    .types {
        display: grid;
        margin-inline: auto;
        grid-template-columns: repeat(5, 1fr);
    }

    .items {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr;
        gap: 5px;
    }

    .lists {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr;
        gap: 5px;
    }
</style>
