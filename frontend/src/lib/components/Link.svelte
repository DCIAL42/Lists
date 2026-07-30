<script lang="ts">
    let { class: className = "", children, ...props } = $props();
</script>

<a {...props} class={className}>{@render children?.()}</a>

<style>
    a {
        text-decoration: none;
        color: inherit;
    }

    a:hover {
        text-decoration: underline solid var(--primary) 2px;
        text-underline-offset: 6px;
    }
</style>
