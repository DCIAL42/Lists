<script lang="ts">
    import { SortableList } from "@rodrigodagostino/svelte-sortable-list";
    import Grip from "./icons/Grip.svelte";
    import Skeleton from "$lib/components/Skeleton.svelte";
    import Button from "$lib/components/Button.svelte";
    import Cross from "$lib/icons/Cross.svelte";
    import type { MediaItem } from "$lib/types";
    import { isAlbum } from "./utils";

    let {
        item,
        index,
        dragging,
        editing = false,
        loading = $bindable(false),
        onRemoveClick,
    }: {
        item: MediaItem;
        index: number;
        dragging: boolean;
        editing?: boolean;
        loading?: boolean;
        onRemoveClick: (e: MouseEvent, i: number) => void;
    } = $props();
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
            <h1 class="title">{index + 1}. {item.data.title}</h1>
            {#if isAlbum(item)}
                <p class="subtitle">{item.data.artist}</p>
            {/if}
        </div>
        {#if editing}
            <Button
                variant="ghost"
                style="color:red;"
                onclick={(e) => onRemoveClick(e, index)}
            >
                <Cross />
            </Button>
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
</style>
