<script lang="ts">
    import type { Snippet } from "svelte";
    import type { SvelteHTMLElements } from "svelte/elements";
    import { v4 as uuidv4 } from "uuid";

    let {
        children,
        focused,
        filled,
        label,
        ...props
    }: SvelteHTMLElements["div"] & {
        children: Snippet<[id: string]>;
        focused: boolean;
        filled: boolean;
        label: string;
    } = $props();

    const id = uuidv4();
</script>

<div class="outer" class:focused class:filled {...props}>
    <label for={id}>{label}</label>
    {@render children?.(id)}
    <fieldset>
        <legend>
            <span>{label}</span>
        </legend>
    </fieldset>
</div>

<style>
    .outer {
        position: relative;
        display: flex;
    }

    fieldset {
        position: absolute;
        inset: -9px 0 0 0;
        border: 1px solid var(--border);
        pointer-events: none;
        border-radius: 4px;
        padding: 0 8px;
        margin: 0;
    }

    span {
        opacity: 0;
        padding: 0 5px;
    }

    .outer:not(:has(*:disabled), .focused):hover fieldset {
        border-color: black;
    }

    .outer.focused fieldset {
        border-color: var(--focused);
    }

    .outer.focused label {
        color: var(--focused-light);
    }

    legend {
        max-width: 0;
        padding: 0;
        font-size: 0.75em;
    }

    .outer label {
        position: absolute;
        transform: translate(14px, 16px);
        transition: transform 0.2s ease;
        color: var(--text-light);
        pointer-events: none;
    }

    .outer.focused legend,
    .outer.filled legend {
        max-width: 100%;
    }

    .outer.focused label,
    .outer.filled label {
        transform: translate(8px, -12px) scale(0.75);
    }
</style>
