<script lang="ts">
    import Editable from "$lib/components/Editable.svelte";
    import Button from "$lib/components/Button.svelte";
    import ArrowLeft from "$lib/icons/ArrowLeft.svelte";
    import { goto } from "$app/navigation";
    import List from "$lib/List.svelte";
    import SearchBar from "$lib/SearchBar.svelte";
    import type { MediaItem } from "$lib/types";

    let items = $state([]);
    let errors = $state({
        title: "",
        items: "",
    });

    async function handleSave() {
        let list = $state.snapshot(listForm);
        let err = false;
        if (list.title === "") {
            errors.title = "Empty title";
            err = true;
        }
        if (list.items.length === 0) {
            errors.items = "No items";
            err = true;
        }
        if (err) return;

        const response = await fetch("/api/lists", {
            method: "POST",
            body: JSON.stringify(list),
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (!response.ok) {
            const data: { error: string } = await response.json();
            throw new Error(data.error);
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

    $effect(() => {
        if (listForm.title !== "") {
            errors.title = "";
        }
        if (listForm.items.length !== 0) {
            errors.items = "";
        }
    });
</script>

<main>
    <div class="topbar">
        <Button variant="ghost" onclick={() => (showBackConfirm = true)}>
            <ArrowLeft />
        </Button>
        <div class="list-details">
            {#if errors.title}
                <p class="error-text">
                    {errors.title}
                </p>
            {/if}
            <Editable
                bind:content={listForm.title}
                placeholder="Enter title"
                class="title"
            />
        </div>
        <Button onclick={handleSave}>Save</Button>
    </div>
    {#if errors.items}
        <p class="error-text">
            {errors.items}
        </p>
    {/if}
    <div class="search-panel">
        <SearchBar
            bind:items
            small
            floating
            add
            handleAdd={(item) => listForm.items.push(item)}
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

    .error-text {
        color: var(--warning);
        margin-inline: auto;
    }
</style>
