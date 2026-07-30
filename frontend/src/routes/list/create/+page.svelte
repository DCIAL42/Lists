<script lang="ts">
    import Editable from "$lib/components/Editable.svelte";
    import SearchForm from "$lib/SearchForm.svelte";
    import Button from "$lib/components/Button.svelte";
    import ArrowLeft from "$lib/icons/ArrowLeft.svelte";
    import { goto } from "$app/navigation";
    import List from "$lib/List.svelte";

    let items = $state([]);

    async function handleSave() {
        let list = $state.snapshot(listForm);

        const response = await fetch("/api/list", {
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

        goto(`/list/${data.ID}`);
    }

    let loading = $state(false);
    let title = $state("");
    let created_by = $state("");
    let showBackConfirm = $state(false);

    let listForm = $derived({
        title: title,
        created_by: created_by,
        items: items,
    });
</script>

<main>
    <div class="topbar">
        <Button variant="ghost" onclick={() => (showBackConfirm = true)}>
            <ArrowLeft />
        </Button>
        <div class="list-details">
            <Editable
                bind:content={title}
                placeholder="Enter title"
                class="title"
            />
            <Editable bind:content={created_by} placeholder="Created by" />
        </div>
        <Button onclick={handleSave}>Save</Button>
    </div>
    <div class="search-panel">
        <SearchForm bind:items />
    </div>
    <List bind:items {loading} />
</main>

{#if showBackConfirm}
    <div class="back-confirm">
        <p>Are you sure you want to go back, changes will be lost</p>
        <Button onclick={() => (showBackConfirm = false)}>Cancel</Button>
        <Button variant="warning" onclick={() => goto("..")}>Confirm</Button>
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
        width: 15%;
        height: 100%;
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
