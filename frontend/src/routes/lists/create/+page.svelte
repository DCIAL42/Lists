<script lang="ts">
    import Editable from "$lib/components/Editable.svelte";
    import Button from "$lib/components/Button.svelte";
    import ArrowLeft from "$lib/icons/ArrowLeft.svelte";
    import { goto } from "$app/navigation";
    import List from "$lib/List.svelte";
    import SearchBar from "$lib/SearchBar.svelte";
    import type { MediaItem } from "$lib/types";

    let items = $state([]);

    async function handleSave() {
        let list = $state.snapshot(listForm);

        const response = await fetch("/api/lists", {
            method: "POST",
            body: JSON.stringify(list),
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (!response.ok) {
            throw new Error("Unable to save list");
        }

        const data = await response.json();

        goto(`/lists/${data.id}`);
    }

    let loading = $state(false);
    let showBackConfirm = $state(false);

    let listForm: { title: string; items: MediaItem[] } = $state({
        title: "",
        items: [],
    });
</script>

<main>
    <div class="topbar">
        <Button variant="ghost" onclick={() => (showBackConfirm = true)}>
            <ArrowLeft />
        </Button>
        <div class="list-details">
            <Editable
                bind:content={listForm.title}
                placeholder="Enter title"
                class="title"
            />
        </div>
        <Button onclick={handleSave}>Save</Button>
    </div>
    <div class="search-panel">
        <SearchBar
            bind:items
            small
            handleAdd={(item) => (listForm.items = [...listForm.items, item])}
        />
    </div>
    <List bind:items={listForm.items} {loading} editing />
</main>

{#if showBackConfirm}
    <div class="back-confirm">
        <p>Are you sure you want to go back, changes will be lost</p>
        <Button onclick={() => (showBackConfirm = false)}>Cancel</Button>
        <Button variant="warning" onclick={() => goto("/lists")}>
            Confirm
        </Button>
    </div>
{/if}

<style>
    main {
        display: flex;
        flex-direction: column;
        gap: 10px;
        margin: 5px;
    }

    .topbar {
        display: flex;
        position: relative;
        border: 1px solid var(--border);
        padding: 5px;
        justify-content: space-between;
        text-align: center;
        align-items: center;
        height: 10vh;
    }

    .search-panel {
        padding: 5px;
        height: 100%;
        width: 15%;
        margin-inline: auto;
    }

    .list-details {
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
    }

    .back-confirm {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        border: 1px solid var(--border);
        background-color: var(--background);
        padding: 20px;
        z-index: 1;
    }
</style>
