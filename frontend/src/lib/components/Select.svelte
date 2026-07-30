<script lang="ts">
    import type { SvelteHTMLElements } from "svelte/elements";
    import Field from "./Field.svelte";

    let {
        options,
        label,
        value = $bindable(),
        ...props
    }: SvelteHTMLElements["select"] & {
        options: any[];
        label: string;
    } = $props();

    let focused = $state(false);
    let filled = $derived(value !== "");
</script>

<Field {label} {focused} {filled}>
    <select
        bind:value
        {...props}
        onfocus={() => (focused = true)}
        onblur={() => (focused = false)}
    >
        <selectedcontent></selectedcontent>
        {#each options as option}
            <option value={option}>{option}</option>
        {/each}
    </select>
</Field>

<style>
    option,
    select {
        text-transform: capitalize;
        text-align: left;
    }

    select::picker-icon {
        transform: scale(0.75);
    }

    select {
        border: 0;
        background-color: var(--background);
        width: 100%;
        padding: 16.5px 14px;
        font-size: inherit;
        appearance: base-select;
    }

    select:focus {
        outline: none;
    }
</style>
