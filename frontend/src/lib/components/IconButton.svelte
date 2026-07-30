<script lang="ts">
    import type { Snippet } from "svelte";

    const { children }: { children?: Snippet<[]> } = $props();
</script>

<button aria-label="plus button">
    {@render children?.()}
</button>

<style>
    button {
        background: lightgrey;
        border: 1px solid var(--border);
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 4px;
        cursor: pointer;
    }
</style>
