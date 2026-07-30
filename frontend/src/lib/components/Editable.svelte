<script lang="ts">
    import type { HTMLAttributes } from "svelte/elements";

    let {
        content = $bindable(""),
        class: className = "",
        ...rest
    }: {
        content?: string;
        className?: string;
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
    }

    .editable-field[placeholder]:empty::before {
        content: attr(placeholder);
        color: #00000044;
        font-style: italic;
    }

    .editable-field[placeholder]:empty:focus::before {
        content: "";
    }
</style>
