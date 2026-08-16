<script lang="ts">
    import type { Snippet } from "svelte";
    import type { SvelteHTMLElements } from "svelte/elements";

    type Props = Omit<SvelteHTMLElements["input"], "placeholder"> & {
        label?: string;
        value?: string;
        variant?: "text" | "search";
        focused?: boolean;
        left_icon?: Snippet<[]>;
    };

    let {
        label = "Label",
        value = $bindable(),
        variant = "text",
        focused = $bindable<boolean>(),
        left_icon,
        ...props
    }: Props = $props();
</script>

<div class="input-outer">
    <div class="input-left-icon">
        {@render left_icon?.()}
    </div>
    <input
        placeholder={label}
        bind:value
        onfocus={() => (focused = true)}
        onblur={() => (focused = false)}
        {...props}
        type="search"
        class:left={left_icon !== undefined}
    />
</div>

<style>
    .input-outer {
        background-color: var(--background);
        border: 1px solid var(--border);
        border-radius: 4px;
        display: flex;
        align-items: center;
    }

    .input-outer:focus-within {
        outline: 1px solid var(--primary);
    }

    input {
        font-size: inherit;
        border: none;
        padding: 16.5px 14px;
        width: 100%;
        flex: 1;
        background-color: var(--background);
        color: var(--primary);
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
        outline: none;
    }

    .input-left-icon {
        padding-left: 10px;
    }

    .input-left-icon :global(svg) {
        height: 16px;
    }
</style>
