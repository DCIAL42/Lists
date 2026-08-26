<script lang="ts">
    import type { HTMLAttributes } from "svelte/elements";

    let {
        size,
        width,
        height,
        borderRadius,
        ...rest
    }: {
        size?: number | string;
        width?: number | string;
        height?: number | string;
        borderRadius?: number | string;
    } & HTMLAttributes<HTMLDivElement> = $props();

    const css = (v: number | string | undefined) =>
        typeof v === "number" ? `${v}px` : v;
</script>

<div
    class="skeleton"
    style:--width={css(width) || css(size)}
    style:--height={css(height) || css(size)}
    style:--border-radius={css(borderRadius) || "4px"}
    {...rest}
></div>

<style>
    /* Skeleton Styles */
    .skeleton {
        background-color: var(--skeleton-base);
        animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
        width: var(--width);
        height: var(--height);
        border-radius: var(--border-radius);
    }

    /* Animation Keyframes */
    @keyframes pulse {
        50% {
            opacity: 0.5;
        }
    }

    /* Reduced Motion Support */
    @media (prefers-reduced-motion: reduce) {
        .skeleton {
            animation: none;
            background: var(--skeleton-base);
        }
    }
</style>
