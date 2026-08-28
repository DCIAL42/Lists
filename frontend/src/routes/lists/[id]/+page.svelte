<script lang="ts">
    import { goto } from "$app/navigation";
    import Button from "$lib/components/Button.svelte";
    import ArrowLeft from "$lib/icons/ArrowLeft.svelte";
    import Edit from "$lib/icons/Edit.svelte";
    import Trash from "$lib/icons/Trash.svelte";
    import List from "$lib/List.svelte";
    import { useClerkContext } from "svelte-clerk";
    import type { PageData } from "./$types";
    import Check from "$lib/icons/Check.svelte";
    import SearchBar from "$lib/SearchBar.svelte";
    import Editable from "$lib/components/Editable.svelte";
    import Cross from "$lib/icons/Cross.svelte";

    let { data = $bindable() }: { data: PageData } = $props();

    let showDeleteConfirm = $state(false);
    let editing = $state(false);
    let listForm = $state({
        title: data.list.title,
        items: data.list.items,
    });
    $inspect(data.list.items[0]);

    async function deleteList() {
        const res = await fetch(`/api/lists/${data.list.id}`, {
            method: "DELETE",
        });
        const body = await res.json();
        console.log(body);
        goto("..");
    }

    async function updateList() {
        const res = await fetch(`/api/lists/${data.list.id}`, {
            method: "PATCH",
            body: JSON.stringify(listForm),
        });
        const body = await res.json();
        console.log(body);
        editing = false;
    }

    function cancelUpdate() {
        editing = false;
        listForm = {
            title: data.list.title,
            items: data.list.items,
        };
        window.location.reload();
    }
</script>

<main>
    <div class="topbar">
        <Button variant="ghost" onclick={() => goto("/lists")}>
            <ArrowLeft />
        </Button>
        <div class="list-details">
            {#if editing}
                <Editable
                    bind:content={listForm.title}
                    placeholder="Enter title"
                    class="title"
                />
            {:else}
                <h1 class="title">{data.list.title}</h1>
            {/if}
            <a href={`/${data.list.created_by}`} class="subtitle"
                >{data.list.created_by}</a
            >
        </div>
        <div class="list-actions">
            {#if data.list.created_by === useClerkContext().user?.username}
                {#if editing}
                    <Button variant="ghost" onclick={cancelUpdate}>
                        <Cross />
                    </Button>
                    <Button variant="ghost" onclick={updateList}>
                        <Check />
                    </Button>
                {:else}
                    <Button variant="ghost" onclick={() => (editing = true)}>
                        <Edit />
                    </Button>
                {/if}
                <Button
                    variant="ghost"
                    style="color: var(--warning);"
                    onclick={() => (showDeleteConfirm = true)}
                >
                    <Trash />
                </Button>
            {/if}
        </div>
    </div>
    {#if editing}
        <div class="search-panel">
            <SearchBar
                small
                floating
                add
                handleAdd={(item) => listForm.items.push(item)}
            />
        </div>
    {/if}
    <List bind:items={listForm.items} {editing} />
</main>

{#if showDeleteConfirm}
    <div class="popover">
        <p>Are you sure you want to delete this list</p>
        <Button onclick={() => (showDeleteConfirm = false)}>Cancel</Button>
        <Button variant="warning" onclick={deleteList}>Confirm</Button>
    </div>
{/if}

<style>
    main {
        display: flex;
        flex-direction: column;
        gap: 10px;
        margin: 5px;
    }

    .subtitle {
        text-decoration: none;
        &:hover {
            color: var(--focused);
        }
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

    .list-details {
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
    }

    .popover {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        border: 1px solid var(--border);
        background-color: var(--background);
        padding: 20px;
        z-index: 1;
    }

    .search-panel {
        padding: 5px;
        height: 100%;
        width: 15%;
        margin-inline: auto;
    }
</style>
