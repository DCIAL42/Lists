<script lang="ts">
    import Button from "./Button.svelte";
    import Plus from "$lib/icons/Plus.svelte";

    const {
        src,
        alt = "",
        title,
        subtitle,
        add = false,
        onclick,
    }: {
        src?: string;
        alt?: string;
        title?: string;
        subtitle?: string;
        add?: boolean;
        onclick?: () => void;
    } = $props();
</script>

<div class="card">
    <div class="card-left">
        {#if src}
            <img {src} class="cover" {alt} />
        {/if}
        <h1 class="title">{title}</h1>
        <p class="subtitle">{subtitle}</p>
    </div>
    {#if add}
        <Button variant="icon" {onclick}><Plus /></Button>
    {/if}
</div>

<style>
    .card {
        display: flex;
        align-items: center;
        width: 300px;
        border: 1px, solid, var(--border);
        text-align: center;
        padding: 10px;
        border-radius: 10px;
        gap: 10px;
        transition: 0.5s ease;
        background-color: var(--background);
    }

    .card-left {
        display: flex;
        flex: 1;
        flex-direction: column;
        align-items: center;
        text-align: center;
    }

    .card:hover {
        box-shadow: gray 2px 2px;
        transform: scale(1.005);
    }

    .cover {
        width: 250px;
    }
</style>
