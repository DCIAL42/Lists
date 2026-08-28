<script lang="ts">
    import type { Snippet } from "svelte";
    import type { SvelteHTMLElements } from "svelte/elements";

    let {
        variant = "primary",
        text = "default",
        selected = false,
        children,
        ...props
    }: {
        variant?: "primary" | "icon" | "ghost" | "warning";
        text?: "default" | "dark" | "light";
        selected?: boolean;
        children?: Snippet<[]>;
    } & SvelteHTMLElements["button"] = $props();
</script>

<button type="button" {...props} class="{variant} {text}" class:selected>
    {@render children?.()}
</button>

<style>
    button {
        text-decoration: none;
        color: inherit;
        background-color: var(--primary);
        padding: 8px;
        border: none;
        /* border: 1px solid var(--border); */
        transition: background-color 0.2s ease;
        cursor: pointer;
    }

    .default,
    .light {
        color: var(--primary-foreground);
    }

    .dark {
        color: var(--primary);
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
        border: none;
    }

    .ghost.default {
        color: var(--primary);
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
        color: var(--focused) !important;
    }
</style>
