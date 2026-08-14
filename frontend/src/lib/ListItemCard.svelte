<script lang="ts">
    import { SortableList } from "@rodrigodagostino/svelte-sortable-list";
    import Grip from "./icons/Grip.svelte";
    import Skeleton from "$lib/components/Skeleton.svelte";
    import Button from "$lib/components/Button.svelte";
    import Cross from "$lib/icons/Cross.svelte";
    import type {
        MediaItem,
        TrackingPayload,
        TrackingStatus,
    } from "$lib/types";
    import { isAlbum } from "./utils";
    import Bookmark from "./icons/Bookmark.svelte";
    import Check from "./icons/Check.svelte";
    import Pause from "./icons/Pause.svelte";
    import { useClerkContext } from "svelte-clerk";
    import { error } from "@sveltejs/kit";

    let {
        item,
        index = 0,
        dragging = false,
        editing = false,
        loading = $bindable(false),
        numbered = false,
        onRemoveClick = () => {},
    }: {
        item: MediaItem;
        index?: number;
        dragging?: boolean;
        editing?: boolean;
        loading?: boolean;
        numbered?: boolean;
        onRemoveClick?: (i: number, e?: MouseEvent) => void;
    } = $props();

    const ctx = useClerkContext();
    const userId = $derived(ctx.auth.userId);

    let tracking = $derived(item.tracking.status);

    async function removeTracking() {
        const res = await fetch(`/api/tracking/${item.tracking.id}`, {
            method: "DELETE",
        });

        if (!res.ok) {
            throw error(500, "failed to delete tracking item");
        }
    }

    async function updateTracking(newStatus: TrackingStatus) {
        const res = await fetch(`/api/tracking/${item.tracking.id}`, {
            method: "PATCH",
            body: JSON.stringify({
                status: newStatus,
            }),
        });

        if (!res.ok) {
            throw error(500, "failed to udpate tracking item");
        }
    }

    async function newTracking(status: TrackingStatus) {
        const payload: TrackingPayload = {
            external_id: item.external_id,
            status: status,
            type: item.type,
        };

        const res = await fetch("/api/tracking", {
            method: "POST",
            body: JSON.stringify(payload),
        });

        if (!res.ok) {
            throw error(500, "failed to create tracking item");
        }
    }

    function handleTrackingChange(status: TrackingStatus) {
        if (userId === null) return;

        if (tracking === status) {
            removeTracking();
        } else if (tracking === undefined) {
            newTracking(status);
        } else {
            updateTracking(status);
        }
        tracking = tracking === status ? undefined : status;
        onRemoveClick(index);
    }
</script>

{#if loading}
    <div class="list-item">
        {#if editing}
            <div class="grip" style="margin-inline: 6px;">
                <Skeleton height={24} width={12} />
            </div>
        {/if}
        <Skeleton size={250} />
        <div class="list-text">
            <Skeleton width={200} height={48} />
            <Skeleton width={150} height={24} />
        </div>
    </div>
{:else}
    <div class="list-item" class:dragging>
        {#if editing}
            <div class="grip">
                <SortableList.ItemHandle>
                    <Grip />
                </SortableList.ItemHandle>
            </div>
        {/if}
        <img
            src={item.data.cover || "https://placehold.co/250"}
            alt={"cover"}
            class="item-cover"
        />
        <div class="list-text">
            <h1 class="title">
                {#if numbered}
                    {index + 1}.
                {/if}
                {item.data.title}
            </h1>
            {#if isAlbum(item)}
                <p class="subtitle">{item.data.artist}</p>
            {/if}
        </div>
        <div class="actions">
            {#if editing}
                <Button
                    variant="ghost"
                    style="color:red;"
                    onclick={(e) => onRemoveClick(index, e)}
                >
                    <Cross />
                </Button>
            {:else}
                <Button
                    variant="ghost"
                    selected={tracking === "backlog"}
                    onclick={() => handleTrackingChange("backlog")}
                >
                    <Bookmark />
                </Button>
                <Button
                    variant="ghost"
                    selected={tracking === "paused"}
                    onclick={() => handleTrackingChange("paused")}
                >
                    <Pause />
                </Button>
                <Button
                    variant="ghost"
                    selected={tracking === "done"}
                    onclick={() => handleTrackingChange("done")}
                >
                    <Check />
                </Button>
            {/if}
        </div>
    </div>
{/if}

<style>
    .dragging {
        visibility: hidden;
    }

    .grip {
        align-self: center;
    }

    .list-item {
        border: 1px solid var(--border);
        border-radius: 4px;
        display: flex;
        gap: 5px;
        padding: 5px;
        align-items: center;
    }

    .list-text {
        margin-inline: auto;
        text-align: center;
        display: flex;
        flex-direction: column;
        gap: 5px;
        align-items: center;
    }

    .item-cover {
        width: 250px;
    }

    .actions {
        display: flex;
        flex-direction: column;
    }
</style>
