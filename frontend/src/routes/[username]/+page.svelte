<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import Skeleton from "$lib/components/Skeleton.svelte";
    import Tabs from "$lib/components/Tabs.svelte";
    import ItemCard from "$lib/ItemCard.svelte";
    import ListPreview from "$lib/ListPreview.svelte";
    import type {
        MediaType,
        TrackingStatus,
        MediaItem,
        ProfileData,
        Follow,
    } from "$lib/types.js";
    import { untrack } from "svelte";
    import { page } from "$app/state";
    import Profile from "$lib/Profile.svelte";

    let data: ProfileData | null = $state(null);

    let tab: "profile" | "lists" | "tracking" = $state("profile");
    let trackingTab: TrackingStatus = $state("backlog");
    let mediaTypes: MediaType[] = $state(["album", "movie", "game"]);
    let toastData: { show: boolean; message: string } = $state({
        show: false,
        message: "",
    });
    let loading = $state(true);

    async function getProfile() {
        loading = true;
        const username = page.params.username;
        if (!username) return;
        const res = await fetch(`/api/${username}`);
        data = await res.json();
        loading = false;
    }

    async function getTrackingItems() {
        if (!data || !data.self) return;
        loading = true;
        data = { ...data, trackingData: { ...data.trackingData, items: [] } };
        const res = await fetch(
            `/api/tracking?type=${mediaTypes.join("|")}&status=${trackingTab}`,
        );

        data = { ...data, trackingData: await res.json() };

        loading = false;
    }

    async function getMoreTrackingItems() {
        if (!data || !data.self) return;
        loading = true;
        const page = Math.floor(data.trackingData.items.length / 10) + 1;
        const res = await fetch(
            `/api/tracking?type=${mediaTypes.join("|")}&status=${trackingTab}&page=${page}`,
        );

        const body: { items: MediaItem[]; total: number } = await res.json();

        data = {
            ...data,
            trackingData: {
                ...data.trackingData,
                items: [...data.trackingData.items, ...body.items],
            },
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

    async function handleFollow() {
        const username = page.params.username;
        if (!username || !data || data.self) return;
        if (data.followData.followed) {
            if (!data.followData.id) return;
            const res = await fetch(`api/follow/${data.followData.id}`, {
                method: "DELETE",
            });
            const body: Follow = await res.json();
            data = { ...data, followData: body };
            return;
        }
        const res = await fetch(`/api/${username}/follow`, {
            method: "POST",
        });
        if (!res.ok) return;
        const body: Follow = await res.json();
        data = { ...data, followData: body };
    }

    $effect(() => {
        void page.params.username;
        getProfile();
    });

    $effect(() => {
        void trackingTab, mediaTypes;
        if (tab === "tracking") {
            untrack(getTrackingItems);
        }
    });
</script>

<main>
    <div class="user-details">
        <h1>{page.params.username}</h1>
        {#if !loading && data && !data.self}
            <Button onclick={handleFollow}>
                {data.followData.followed ? "followed" : "follow"}
            </Button>
        {/if}
    </div>
    <hr style="width: 100%;" />
    {#if data?.self}
        <Tabs tabs={["profile", "lists", "tracking"]} bind:selected={tab} />
    {:else}
        <Tabs tabs={["profile", "lists"]} bind:selected={tab} />
    {/if}

    {#if tab === "tracking"}
        <Tabs
            tabs={["backlog", "paused", "done"]}
            bind:selected={trackingTab}
        />

        <Tabs tabs={["album", "movie", "game"]} bind:selected={mediaTypes} />
    {/if}

    {#if tab === "profile"}
        {#if loading || !data}
            <Profile loading />
        {:else}
            <Profile {data} />
        {/if}
    {:else if tab === "lists"}
        {#if loading || !data}{:else}
            <div class="lists">
                {#each data.listsData.lists as list}
                    <ListPreview {list} />
                {/each}
            </div>
        {/if}
    {:else if tab === "tracking" && data?.self}
        <div class="items">
            {#each data.trackingData.items as _, index}
                <ItemCard
                    bind:item={data.trackingData.items[index]}
                    {index}
                    small
                    {onTrackingChange}
                />
            {/each}
            {#if loading}
                {#each Array(10) as _}
                    <Skeleton height={250} />
                {/each}
            {/if}
        </div>
        {#if data.trackingData.count > data.trackingData.items.length}
            <Button onclick={getMoreTrackingItems}>more</Button>
        {/if}
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
        margin-inline: auto;
        width: 71vw;
    }

    .items {
        display: grid;
        grid-template-columns: repeat(5, 1fr);
        gap: 5px;
        align-items: center;
    }

    .lists {
        display: grid;
        grid-template-columns: repeat(5, 1fr);
        gap: 5px;
        justify-items: center;
    }
</style>
