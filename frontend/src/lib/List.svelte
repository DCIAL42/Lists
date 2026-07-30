<script lang="ts">
    import ListItem from "$lib/ListItem.svelte";
    import {
        SortableList,
        sortItems,
    } from "@rodrigodagostino/svelte-sortable-list";

    let {
        items = $bindable(),
        loading = false,
    }: { items: any[]; loading?: boolean } = $props();

    function tag_items(items: any[]) {
        return items.map((item, i) => ({ id: String(i), ...item }));
    }

    function handleDragEnd(e: SortableList.RootEvents["ondragend"]) {
        const { draggedItemIndex, targetItemIndex, isCanceled } = e;
        if (
            !isCanceled &&
            typeof targetItemIndex === "number" &&
            draggedItemIndex !== targetItemIndex
        )
            items = sortItems(items, draggedItemIndex, targetItemIndex);
    }

    function onRemoveClick(_: MouseEvent, i: number) {
        items = items.toSpliced(i, 1);
    }

    let dragIdx = $state<number | null>(null);
</script>

<div class="list">
    {#if items.length === 0}
        <div class="empty-list">
            <p>Search for something to add to the list</p>
        </div>
    {:else}
        <SortableList.Root
            ondragend={handleDragEnd}
            ondrag={(e) => (dragIdx = e.draggedItemIndex)}
            ondrop={() => (dragIdx = null)}
            isLocked={loading}
        >
            {#each tag_items(items) as item, index (item.id)}
                <SortableList.Item {...item} {index}>
                    <ListItem
                        {item}
                        {index}
                        dragging={dragIdx === index}
                        bind:loading
                        {onRemoveClick}
                    />
                </SortableList.Item>
            {/each}
        </SortableList.Root>
    {/if}
</div>

<style>
    .list {
        display: flex;
        flex-direction: column;
        gap: 10px;
        width: 50%;
        margin-inline: auto;
    }

    .empty-list {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        color: var(--primary-muted);
    }
</style>
