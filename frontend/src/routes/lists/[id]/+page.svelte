<script lang="ts">
    import { goto } from "$app/navigation";
    import Button from "$lib/components/Button.svelte";
    import ArrowLeft from "$lib/icons/ArrowLeft.svelte";
    import List from "$lib/List.svelte";

    let { data } = $props();

    let list = $derived(data.list);
    let items = $derived(list.items);
</script>

<main>
    <div class="topbar">
        <Button variant="ghost" onclick={() => goto("/lists")}>
            <ArrowLeft />
        </Button>
        <div class="list-details">
            <h1 class="title">{list.title}</h1>
            <p class="subtitle">{list.created_by}</p>
        </div>
    </div>
    <List bind:items />
</main>

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
</style>
