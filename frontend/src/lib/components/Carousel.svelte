<script lang="ts">
    import Button from "$lib/components/Button.svelte";
    import ListPreview from "$lib/ListPreview.svelte";
    import type { ListMeta } from "$lib/types";
    import { onMount } from "svelte";

    let { lists }: { lists: ListMeta[] } = $props();
    let container: HTMLDivElement;
    let canScrollLeft = $state(false);
    let canScrollRight = $state(false);
    function updateScrollButtons() {
        if (!container) return;

        canScrollLeft = container.scrollLeft > 0;

        canScrollRight =
            container.scrollLeft + container.clientWidth <
            container.scrollWidth - 1;
    }
    onMount(() => {
        updateScrollButtons();
    });
</script>

<div class="carousel">
    <div class="lists" bind:this={container} onscroll={updateScrollButtons}>
        {#each lists as list}
            <ListPreview {list} />
        {/each}
    </div>
    {#if canScrollRight}
        <div class="scroll right">
            <Button
                style="width: 100%; height: 100%; background-color: #00000077;"
                onclick={() =>
                    container.scrollBy({
                        left: 600,
                        behavior: "smooth",
                    })}
            >
                &#10095;
            </Button>
        </div>
    {/if}
    {#if canScrollLeft}
        <div class="scroll left">
            <Button
                style="width: 100%; height: 100%; background-color: #00000077;"
                onclick={() =>
                    container.scrollBy({
                        left: -600,
                        behavior: "smooth",
                    })}
            >
                &#10094;
            </Button>
        </div>
    {/if}
</div>

<style>
    .scroll {
        position: absolute;
        top: 0;
        width: 50px;
        height: 100%;
        background-color: #00000044;
        display: flex;
        align-items: center;
        justify-content: center;
        opacity: 0;
        transition: opacity 0.2s;
    }

    .scroll.right {
        left: 100%;
        transform: translateX(-100%);
    }

    .carousel {
        position: relative;

        &:hover .scroll {
            opacity: 1;
        }
    }

    .lists {
        display: flex;
        overflow: auto;
        scroll-behavior: smooth;
        scrollbar-width: none;
        gap: 5px;
    }
</style>
