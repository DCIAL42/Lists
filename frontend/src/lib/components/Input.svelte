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
        focused = $bindable(),
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

        &:focus-within {
            outline: 1px solid var(--primary);
        }
    }

    input {
        font-size: inherit;
        border: none;
        padding: 16.5px 14px;
        width: 100%;
        flex: 1;
        background-color: var(--background);
        color: var(--primary);

        &:focus {
            outline: none;
        }

        &[type="search"]::-webkit-search-cancel-button {
            -webkit-appearance: none;

            width: 1.2em;
            height: 1.2em;

            background: none;

            -webkit-mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M18 6 6 18M6 6l12 12' fill='none' stroke='black' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E")
                center / contain no-repeat;
            mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M18 6 6 18M6 6l12 12' fill='none' stroke='black' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E")
                center / contain no-repeat;

            background-color: var(--primary);
            opacity: 0.5;
            cursor: pointer;

            &:hover {
                opacity: 1;
            }
        }
    }

    .input-left-icon {
        padding-left: 10px;

        :global(svg) {
            height: 16px;
        }
    }
</style>
