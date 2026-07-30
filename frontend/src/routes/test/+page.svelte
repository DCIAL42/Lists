<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import ListItem from "$lib/ListItem.svelte";
    import {
        removeItem,
        SortableList,
        sortItems,
    } from "@rodrigodagostino/svelte-sortable-list";
    let val = $state("");
    let val2 = $state("Test");
    const ops = ["op1", "op2"];
    let val3 = $state("");
    let things = [];
    for (let i = 0; i < 10; i++) {
        things.push(`thing${i}`);
    }
    let cols = $state(1);
    let items: SortableList.ItemData[] = $state([
        {
            id: "1",
            title: "Either/Or",
            artist: "Elliott Smith",
            cover: "https://placehold.co/250",
        },
        {
            id: "2",
            title: "In Rainbows",
            artist: "Radiohead",
            cover: "https://placehold.co/250",
        },
        {
            id: "3",
            title: "Currents",
            artist: "Tame Impala",
            cover: "https://placehold.co/250",
        },
        {
            id: "4",
            title: "Random Access Memories",
            artist: "Daft Punk",
            cover: "https://placehold.co/250",
        },
    ]);

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
        items = removeItem(items, i);
    }

    let dragIdx = $state<number | null>(null);
    let loading = $state(false);
</script>

<!-- <div class="outer"> -->
<!--     <Input label="Name" bind:value={val} /> -->
<!--     <Input label="Name" bind:value={val2} /> -->
<!--     <Select options={ops} bind:value={val3} label="Name" /> -->
<!--     <div style="display: block;"> -->
<!--         <button onclick={() => (cols = Math.max(1, cols - 1))}>-</button> -->
<!--         <button onclick={() => (cols = Math.min(things.length, cols + 1))} -->
<!--             >+</button -->
<!--         > -->
<!--     </div> -->
<!--     <div class="grid" style="--n: {cols}"> -->
<!--         {#each things as thing} -->
<!--             <div>{thing}</div> -->
<!--         {/each} -->
<!--     </div> -->
<!-- </div> -->
<!-- <Button variant="icon"><Plus /></Button> -->
<!-- <Button><Plus /></Button> -->

<!-- <Skeleton size={250} /> -->
<Button onclick={() => (loading = !loading)}>Toggle loading</Button>
<main>
    <div class="sidebar">
        <h1 class="title">List Name</h1>
        <p class="subtitle">List creator</p>
        <p>
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Enim
            beatae, minima explicabo repudiandae obcaecati nihil totam facilis
            iure deleniti a vero alias possimus provident, eveniet nobis
            quibusdam delectus fugiat optio?
        </p>
    </div>
    <div class="list">
        <SortableList.Root
            ondragend={handleDragEnd}
            ondrag={(e) => (dragIdx = e.draggedItemIndex)}
            ondrop={() => (dragIdx = null)}
            isLocked={loading}
        >
            {#each items as item, index (item.id)}
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
    </div>
</main>

<style>
    .outer {
        padding: 5px;
        display: flex;
        flex-direction: column;
        gap: 25px;
    }

    .grid {
        display: grid;
        grid-template-columns: repeat(var(--n), auto);
        gap: 10px;
        padding: 10px;
    }

    .grid div {
        border: 1px solid var(--border);
    }

    .sidebar {
        padding: 5px;
        border: 1px solid var(--border);
        border-radius: 4px;
        height: 100%;
    }

    main {
        display: grid;
        grid-template-columns: 250px 1fr;
        gap: 5px;
        margin: 5px;
    }

    .list {
        display: flex;
        flex-direction: column;
        gap: 10px;
        width: 50%;
        margin-inline: auto;
    }
</style>
