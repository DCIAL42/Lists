<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import Skeleton from "$lib/components/Skeleton.svelte";
    import Tabs from "$lib/components/Tabs.svelte";
    import ListItemCard from "$lib/ListItemCard.svelte";
    import ListPreview from "$lib/ListPreview.svelte";
    import type {
        MediaType,
        TrackingStatus,
        MediaItem,
        ListMeta,
        UserResponse,
    } from "$lib/types.js";
    import { untrack } from "svelte";
    import { page } from "$app/state";

    interface ListsPreviewData {
        lists: ListMeta[];
        next: string;
        page: number;
        count: number;
    }

    interface TrackingListData {
        items: MediaItem[];
        count: number;
    }

    let {
        data = $bindable(),
    }: {
        data: {
            trackingData?: TrackingListData;
            listsData: ListsPreviewData;
            userData: UserResponse;
        };
    } = $props();

    let trackingData = $state.raw(data.trackingData);
    let listsData = $state.raw(data.listsData);
    let tab: "lists" | "tracking" = $state("tracking");
    let trackingTab: TrackingStatus = $state("backlog");
    let mediaTypes: MediaType[] = $state(["album", "movie", "game"]);
    let toastData: { show: boolean; message: string } = $state({
        show: false,
        message: "",
    });
    let loading = $state(true);

    async function getTrackingItems() {
        if (trackingData === undefined) return;
        loading = true;
        trackingData = { ...trackingData, items: [] };
        const res = await fetch(
            `/api/tracking?type=${mediaTypes.join("|")}&status=${trackingTab}`,
        );

        const body: TrackingListData = await res.json();

        trackingData = body;
        loading = false;
    }

    async function getMoreTrackingItems() {
        if (trackingData === undefined) return;
        loading = true;
        const page = Math.floor(trackingData.items.length / 9) + 1;
        const res = await fetch(
            `/api/tracking?type=${mediaTypes.join("|")}&status=${trackingTab}&page=${page}`,
        );

        const body: { items: MediaItem[]; total: number } = await res.json();

        trackingData = {
            ...trackingData,
            items: [...trackingData.items, ...body.items],
        };
        loading = false;
    }

    function onTrackingChange(from: TrackingStatus, to: TrackingStatus) {
        if (to === "none") {
            toastData = { show: true, message: "Item removed" };
        } else if (from === "none") {
            toastData = { show: true, message: "Item added" };
        } else {
            toastData = { show: true, message: "Item updated" };
        }
        setTimeout(() => {
            toastData = { show: false, message: "" };
        }, 3000);
    }

    $effect(() => {
        void trackingTab;
        if (tab === "tracking") {
            untrack(getTrackingItems);
        }
    });
</script>

<main>
    <div class="user-details">
        <h1>{page.params.username}</h1>
    </div>
    <hr style="width: 100%;" />
    {#if trackingData !== undefined}
        <Tabs tabs={["tracking", "lists"]} bind:selected={tab} />

        {#if tab === "tracking"}
            <Tabs
                tabs={["backlog", "paused", "done"]}
                bind:selected={trackingTab}
            />

            <Tabs
                tabs={["album", "movie", "game"]}
                bind:selected={mediaTypes}
            />
        {/if}

        {#if tab === "tracking"}
            <div class="items">
                {#each trackingData.items as item, index}
                    <ListItemCard {item} {index} {onTrackingChange} />
                {/each}
                {#if loading}
                    {#each Array(9) as _}
                        <Skeleton height={250} />
                    {/each}
                {/if}
            </div>
            {#if trackingData.count > trackingData.items.length}
                <Button onclick={getMoreTrackingItems}>more</Button>
            {/if}
        {:else if tab === "lists"}
            <div class="lists">
                {#each listsData.lists as list}
                    <ListPreview {list} />
                {/each}
            </div>
        {/if}
    {:else}
        <div class="lists">
            {#each listsData.lists as list}
                <ListPreview {list} />
            {/each}
        </div>
    {/if}

    {#if toastData.show}
        <div class="toast">
            <p>{toastData.message}</p>
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

    .items {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr;
        gap: 5px;
    }

    .lists {
        display: grid;
        grid-template-columns: repeat(5, 1fr);
        gap: 5px;
        justify-items: center;
    }
</style>
