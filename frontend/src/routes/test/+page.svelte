<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import Card from "$lib/components/Card.svelte";
    import ListItemCard from "$lib/ListItemCard.svelte";
    import type { MediaItem } from "$lib/types";
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

    let loading = $state(false);

    const item: MediaItem = {
        type: "album",
        external_id: "",
        data: {
            title: "test",
            artist: "other",
            cover: "https://placehold.co/50",
        },
    };
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

<ListItemCard {item} />

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
