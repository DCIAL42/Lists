<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import type {
        MediaItem,
        TrackingItem,
        TrackingPayload,
        TrackingStatus,
    } from "$lib/types";
    import Bookmark from "./icons/Bookmark.svelte";
    import Check from "./icons/Check.svelte";
    import Pause from "./icons/Pause.svelte";
    import { useClerkContext } from "svelte-clerk";
    import { error } from "@sveltejs/kit";
    import Plus from "./icons/Plus.svelte";
    import type { HTMLAttributes } from "svelte/elements";

    let {
        item = $bindable(),
        add = false,
        onAddClick = () => {},
        onTrackingChange = () => {},
        ...rest
    }: {
        item: MediaItem;
        add?: boolean;
        onAddClick?: (item: MediaItem) => void;
        onTrackingChange?: (from: TrackingStatus, to: TrackingStatus) => void;
    } & HTMLAttributes<HTMLDivElement> = $props();

    const ctx = useClerkContext();
    const userId = $derived(ctx.auth.userId);

    let tracking = $derived(item.tracking?.status);

    async function removeTracking() {
        if (item.tracking === undefined) {
            return;
        }

        const res = await fetch(`/api/tracking/${item.tracking.id}`, {
            method: "DELETE",
        });

        if (!res.ok) {
            throw error(500, "failed to delete tracking item");
        }

        const data: TrackingItem = await res.json();
        item.tracking = data;
    }

    async function updateTracking(newStatus: TrackingStatus) {
        if (item.tracking === undefined) {
            return;
        }

        const res = await fetch(`/api/tracking/${item.tracking.id}`, {
            method: "PATCH",
            body: JSON.stringify({
                status: newStatus,
            }),
        });

        if (!res.ok) {
            throw error(500, "failed to udpate tracking item");
        }

        const data: TrackingItem = await res.json();
        item.tracking = data;
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

        const data: TrackingItem = await res.json();
        item.tracking = data;
    }

    function handleTrackingChange(status: TrackingStatus) {
        if (userId === null) return;

        if (tracking === status) {
            removeTracking();
            onTrackingChange(tracking, "none");
        } else if (tracking === undefined) {
            newTracking(status);
            onTrackingChange("none", status);
        } else {
            updateTracking(status);
            onTrackingChange(tracking, status);
        }
        tracking = tracking === status ? undefined : status;
    }
</script>

<div class="cover" {...rest}>
    <img src={item.cover || "https://placehold.co/250"} alt="cover" />

    <div class="hover details">
        <p>{item.data.title}</p>
    </div>

    <div class="hover actions">
        {#if add}
            <Button
                variant="ghost"
                text="light"
                onclick={() => onAddClick(item)}
            >
                <Plus />
            </Button>
        {/if}
        <Button
            variant="ghost"
            text="light"
            selected={tracking === "backlog"}
            onclick={() => handleTrackingChange("backlog")}
        >
            <Bookmark />
        </Button>
        <Button
            variant="ghost"
            text="light"
            selected={tracking === "paused"}
            onclick={() => handleTrackingChange("paused")}
        >
            <Pause />
        </Button>
        <Button
            variant="ghost"
            text="light"
            selected={tracking === "done"}
            onclick={() => handleTrackingChange("done")}
        >
            <Check />
        </Button>
    </div>
</div>

<style>
    .cover {
        position: relative;
        width: fit-content;
        height: fit-content;

        img {
            display: block;
            width: 100%;
        }

        &::after {
            content: "";
            position: absolute;
            inset: 0;
            background: rgba(0, 0, 0, 0);
            transition: background 0.2s ease;
            pointer-events: none;
        }

        &:hover {
            &::after {
                background: rgba(0, 0, 0, 0.2);
            }

            .hover {
                display: flex;
                flex-direction: column;
            }
        }
    }

    .hover {
        position: absolute;
        display: none;
        background-color: rgb(from var(--primary) r g b / 40%);
    }

    .actions {
        top: 0;
        left: 100%;
        gap: 5px;
        transform: translateX(-100%);
    }

    .details {
        padding: 5px;
        p {
            margin: 0;
        }
        top: 100%;
        left: 0;
        transform: translateY(-100%);
        color: var(--primary-foreground);
    }
</style>
