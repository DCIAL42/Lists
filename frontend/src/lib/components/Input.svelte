<script lang="ts">
    import type { SvelteHTMLElements } from "svelte/elements";
    import Field from "./Field.svelte";

    type Props = Omit<SvelteHTMLElements["input"], "placeholder"> & {
        label?: string;
        value?: string;
        variant?: "text" | "search";
        focused?: boolean;
    };

    let {
        label = "Label",
        value = $bindable(),
        variant = "text",
        focused = $bindable<boolean>(),
        ...props
    }: Props = $props();

    let filled = $derived(value);
</script>

<Field {label} {focused} {filled}>
    <input
        placeholder=""
        bind:value
        onfocus={() => (focused = true)}
        onblur={() => (focused = false)}
        {...props}
    />
</Field>

<style>
    input {
        background-color: var(--background);
        width: 100%;
        border: 0;
        padding: 16.5px 14px;
        height: 1.5em;
        font-size: inherit;
    }

    input:focus {
        outline: none;
    }
</style>
