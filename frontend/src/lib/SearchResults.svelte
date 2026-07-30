<script lang="ts">
    import type { SvelteHTMLElements } from "svelte/elements";
    import type { Snippet } from "svelte";
    let {
        children,
        ...props
    }: { children: Snippet<[]> } & SvelteHTMLElements["div"] = $props();
</script>

<div {...props} class="scrollable-element">
    {@render children()}
</div>

<style>
    div {
        display: flex;
        flex-direction: column;
        gap: 10px;
        align-items: center;
        position: absolute;
        left: 50%;
        transform: translateX(-50%);
        top: 100%;
        overflow-y: auto;
        max-height: 500px;
        padding: 10px;
        border: 2px solid var(--border);
        border-radius: 4px;
        background-color: var(--background);
    }
</style>
