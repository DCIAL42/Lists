<script lang="ts">
    import { goto } from "$app/navigation";
    import Button from "$lib/components/Button.svelte";
    import ArrowLeft from "$lib/icons/ArrowLeft.svelte";
    import Edit from "$lib/icons/Edit.svelte";
    import Trash from "$lib/icons/Trash.svelte";
    import List from "$lib/List.svelte";
    import { useClerkContext } from "svelte-clerk";
    import type { PageData } from "./$types";

    let { data = $bindable() }: { data: PageData } = $props();

    let showDeleteConfirm = $state(false);

    async function deleteList() {
        const res = await fetch(`/api/lists/${data.list.id}`, {
            method: "DELETE",
        });
        const body = await res.json();
        console.log(body);
        goto("..");
    }
</script>

<main>
    <div class="topbar">
        <Button variant="ghost" onclick={() => goto("/lists")}>
            <ArrowLeft />
        </Button>
        <div class="list-details">
            <h1 class="title">{data.list.title}</h1>
            <p class="subtitle">{data.list.created_by}</p>
        </div>
        <div class="list-actions">
            {#if data.list.created_by === useClerkContext().user?.username}
                <Button variant="ghost">
                    <Edit />
                </Button>
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
    <List bind:items={data.list.items} />
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
</style>
