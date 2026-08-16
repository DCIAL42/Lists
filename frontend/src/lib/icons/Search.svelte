<script lang="ts">
    import type { SvelteHTMLElements } from "svelte/elements";

    let props: SvelteHTMLElements["svg"] = $props();
</script>

<svg
    xmlns="http://www.w3.org/2000/svg"
    width="24"
    height="24"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
    class="lucide lucide-search-icon lucide-search"
    {...props}
>
    <path d="m21 21-4.34-4.34" />
    <circle cx="11" cy="11" r="8" />
</svg>

<style>
    svg {
        display: flex;
    }
</style>
