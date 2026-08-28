<script lang="ts">
    import type { HTMLAttributes } from "svelte/elements";

    let {
        content = $bindable(""),
        class: className = "",
        ...rest
    }: {
        content?: string;
    } & HTMLAttributes<HTMLDivElement> = $props();

    function metaEditKeyDown(e: KeyboardEvent) {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            (e.target as HTMLElement).blur();
        }
    }
</script>

<div
    class="editable-field {className}"
    class:empty={!content}
    contenteditable
    bind:textContent={content}
    onkeydown={metaEditKeyDown}
    role="textbox"
    tabindex="0"
    {...rest}
></div>

<style>
    .editable-field {
        margin: 0;
        cursor: text;

        &.empty[placeholder] {
            &::before {
                content: attr(placeholder);
                color: var(--primary);
                font-style: italic;
            }

            &:focus::before {
                content: "";
            }
        }
    }
</style>
