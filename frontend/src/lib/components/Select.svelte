<script lang="ts">
    import type { SvelteHTMLElements } from "svelte/elements";

    let {
        options = [],
        label,
        value = $bindable(),
        displayText = (t) => t,
        ...props
    }: SvelteHTMLElements["select"] & {
        options?: string[];
        label: string;
        displayText?: (text: string) => string;
    } = $props();

    let empty = $derived(value === "");
</script>

<select bind:value {...props} placeholder={label} class:empty>
    {#each options as option}
        <option value={option}>{displayText(option)}</option>
    {/each}
</select>

<style>
    select {
        border: 1px solid var(--border);
        padding: 16.5px 14px;
        appearance: base-select;
        border-radius: 4px;
        gap: 0;

        &::picker-icon {
            transform: scale(0.75);
            margin-inline-start: 10px;
        }

        &:focus {
            outline: 1px solid black;
        }
    }

    .empty::after {
        content: attr(placeholder);
    }
</style>
