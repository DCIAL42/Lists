<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import type {
        MediaItem,
        Rating,
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
    import Menu from "./icons/Menu.svelte";
    import { tick } from "svelte";
    import Cross from "./icons/Cross.svelte";

    let {
        item = $bindable(),
        add = false,
        rating,
        onAddClick = () => {},
        onTrackingChange = () => {},
        ...rest
    }: {
        item: MediaItem;
        add?: boolean;
        rating?: number;
        onAddClick?: (item: MediaItem) => void;
        onTrackingChange?: (from: TrackingStatus, to: TrackingStatus) => void;
    } & HTMLAttributes<HTMLDivElement> = $props();

    const ctx = useClerkContext();
    const userId = $derived(ctx.auth.userId);

    let tracking = $derived(item.tracking?.status);
    let showMenu = $state(false);

    let menuActions = $derived([
        {
            text:
                tracking === "backlog"
                    ? "Remove from backlog"
                    : "Add to backlog",
            onclick: () => handleTrackingChange("backlog"),
        },
        {
            text:
                tracking === "paused" ? "Remove from paused" : "Add to paused",
            onclick: () => handleTrackingChange("paused"),
        },
        {
            text: tracking === "done" ? "Remove from done" : "Add to done",
            onclick: () => handleTrackingChange("done"),
        },
    ]);

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
            media_id: item.id,
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

    let newRating = $state(item.rating.rating ?? 0);

    async function deleteRating() {
        const res = await fetch(`/api/rating/${item.rating.id}`, {
            method: "DELETE",
        });
        if (res.ok) {
            item.rating = {};
            newRating = 0;
        }
    }

    async function submitRating() {
        if (newRating === item.rating.rating) {
            await deleteRating();
            return;
        }
        if (item.rating.id !== undefined) {
            const payload = {
                rating: newRating,
            };
            const res = await fetch(`/api/rating/${item.rating.id}`, {
                method: "PATCH",
                body: JSON.stringify(payload),
            });

            if (!res.ok) return;

            const data: Rating = await res.json();
            item.rating = data;
            return;
        }
        const payload = {
            media_id: item.id,
            rating: newRating,
        };
        const res = await fetch("/api/rating", {
            method: "POST",
            body: JSON.stringify(payload),
        });

        if (!res.ok) return;

        const data: Rating = await res.json();
        item.rating = data;
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

    let menuContainer: HTMLDivElement;
</script>

<div class="cover" {...rest}>
    <img src={item.cover || "https://placehold.co/250"} alt="cover" />

    <div class="hover details">
        <p>{item.title}</p>
    </div>

    <div
        class="menu"
        bind:this={menuContainer}
        tabindex="-1"
        onfocusout={() => {
            setTimeout(() => {
                if (!menuContainer.contains(document.activeElement))
                    showMenu = false;
            });
        }}
    >
        {#if showMenu}
            <div class="rating">
                <div class="slider-wrapper">
                    <span
                        class="slider-value"
                        style={`left: ${newRating ?? 0 * 10}%`}
                    >
                        {newRating}
                    </span>

                    <input
                        type="range"
                        bind:value={newRating}
                        min={0}
                        max={10}
                        defaultvalue={0}
                        step={0.5}
                    />
                </div>
                <Button variant="ghost" onclick={submitRating}>
                    {#if newRating === item.rating.rating}
                        <Cross />
                    {:else}
                        <Check />
                    {/if}
                </Button>
            </div>
            {#each menuActions as button}
                <Button
                    variant="ghost"
                    style="border-top: 1px solid var(--primary-foreground); "
                    onclick={button.onclick}
                >
                    {button.text}
                </Button>
            {/each}
        {/if}
    </div>

    <div class="hover actions">
        <Button
            variant="ghost"
            text="light"
            onclick={async () => {
                showMenu = true;
                await tick();
                menuContainer.focus();
            }}
        >
            <Menu />
        </Button>
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

    .menu {
        position: absolute;
        display: flex;
        flex-direction: column;
        top: 0;
        left: 100%;
        width: 200px;
        background-color: var(--skeleton-base);
        z-index: 999;

        .rating {
            display: grid;
            grid-template-columns: 1fr auto;
            align-items: center;
            gap: 10px;
            padding: 10px;

            .slider-wrapper {
                position: relative;
                width: 100%;
                padding-top: 20px;
                display: flex;
                align-items: center;

                input {
                    width: 100%;
                    margin: 0;
                }

                .slider-value {
                    position: absolute;
                    top: 0;
                    transform: translateX(-50%);
                }
            }
        }
    }
</style>
