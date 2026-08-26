<script lang="ts">
    import type { Snippet } from "svelte";
    import type { SvelteHTMLElements } from "svelte/elements";

    let {
        variant = "primary",
        selected = false,
        children,
        ...props
    }: {
        variant?: "primary" | "icon" | "ghost" | "warning";
        selected?: boolean;
        children?: Snippet<[]>;
    } & SvelteHTMLElements["button"] = $props();
</script>

<button type="button" {...props} class={variant} class:selected>
    {@render children?.()}
</button>

<style>
    button {
        text-decoration: none;
        color: inherit;
        background-color: var(--primary);
        color: var(--primary-foreground);
        padding: 8px;
        border: none;
        /* border: 1px solid var(--border); */
        transition: background-color 0.2s ease;
        cursor: pointer;
    }

    .primary,
    .icon {
        background-color: var(--primary);
    }

    .primary:hover,
    .icon:hover {
        background-color: color-mix(in srgb, var(--primary) 80%, white);
    }

    .ghost {
        background-color: transparent;
        color: var(--primary);
        border: none;
    }

    .ghost:hover {
        background-color: color-mix(in srgb, transparent 80%, black);
    }

    .warning {
        background-color: var(--warning);
        border: none;
    }

    .warning:hover {
        /* background-color: rgba(var(--warning), 0.9); */
        background-color: color-mix(in srgb, var(--warning) 70%, white);
    }

    .selected {
        color: var(--focused);
    }
</style>
