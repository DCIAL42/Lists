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

<div>
    <input
        placeholder={label}
        bind:value
        onfocus={() => (focused = true)}
        onblur={() => (focused = false)}
        {...props}
        type="search"
    />
</div>

<style>
    input {
        border: 1px solid var(--border);
        border-radius: 4px;
        background-color: var(--background);
        padding: 16.5px 14px;
        font-size: inherit;
    }

    input[type="search"]::-webkit-search-cancel-button {
        -webkit-appearance: none;
        height: 1.2em;
        width: 1.2em;
        background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9ImN1cnJlbnRDb2xvciIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiIGNsYXNzPSJsdWNpZGUgbHVjaWRlLXgtaWNvbiBsdWNpZGUteCI+PHBhdGggZD0iTTE4IDYgNiAxOCIvPjxwYXRoIGQ9Im02IDYgMTIgMTIiLz48L3N2Zz4=)
            no-repeat 50% 50%;
        background-size: contain;
        opacity: 0.5;
        cursor: pointer;
    }

    input[type="search"]::-webkit-search-cancel-button:hover {
        opacity: 1;
    }

    input:focus {
        outline: 1px solid black;
    }
</style>
