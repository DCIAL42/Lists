<script lang="ts">
    import { SortableList } from "@rodrigodagostino/svelte-sortable-list";
    import Grip from "./icons/Grip.svelte";
    import Skeleton from "$lib/components/Skeleton.svelte";
    import Button from "$lib/components/Button.svelte";
    import Cross from "$lib/icons/Cross.svelte";
    import type { MediaItem, TrackingStatus } from "$lib/types";
    import { isAlbum } from "./utils";
    import type { HTMLAttributes } from "svelte/elements";
    import ItemCover from "./ItemCover.svelte";

    let {
        item = $bindable(),
        index = 0,
        dragging = false,
        editing = false,
        loading = $bindable(false),
        numbered = false,
        small = false,
        add = false,
        onRemoveClick = () => {},
        onAddClick = () => {},
        onTrackingChange = () => {},
        ...rest
    }: {
        item: MediaItem;
        index?: number;
        dragging?: boolean;
        editing?: boolean;
        loading?: boolean;
        numbered?: boolean;
        small?: boolean;
        add?: boolean;
        onRemoveClick?: (i: number, e?: MouseEvent) => void;
        onAddClick?: (item: MediaItem) => void;
        onTrackingChange?: (from: TrackingStatus, to: TrackingStatus) => void;
    } & HTMLAttributes<HTMLDivElement> = $props();
</script>

{#if loading}
    <div class="list-item" {...rest}>
        {#if editing}
            <div class="grip" style="margin-inline: 6px;">
                <Skeleton height={24} width={12} />
            </div>
        {/if}
        <Skeleton size="100%" />
        <div class="list-text">
            <Skeleton width={200} height={48} />
            <Skeleton width={150} height={24} />
        </div>
    </div>
{:else}
    <div class="list-item" class:small class:dragging {...rest}>
        {#if editing}
            <div class="grip">
                <SortableList.ItemHandle>
                    <Grip />
                </SortableList.ItemHandle>
            </div>
        {/if}

        <div class="cover" class:small>
            <ItemCover bind:item {add} {onAddClick} {onTrackingChange} />
        </div>

        {#if !small}
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
                        text="light"
                        style="color:var(--warning);"
                        onclick={(e) => onRemoveClick(index, e)}
                    >
                        <Cross />
                    </Button>
                {/if}
            </div>
        {/if}
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
        position: relative;
        border: 1px solid var(--border);
        border-radius: 4px;
        display: flex;
        gap: 5px;
        padding: 5px;
        align-items: center;

        &.small {
            padding: 0;
            width: fit-content;
        }
    }

    .list-text {
        margin-inline: auto;
        text-align: center;
        display: flex;
        flex-direction: column;
        gap: 5px;
        align-items: center;
    }

    .actions {
        display: flex;
        flex-direction: column;
    }

    .cover {
        max-width: 250px;

        &.small {
            max-width: unset;
            width: 100%;
        }
    }
</style>
