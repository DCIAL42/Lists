<script lang="ts">
    import type { ListMeta } from "$lib/types";

    let { list }: { list: ListMeta } = $props();
</script>

<a href={`/lists/${list.id}`}>
    <div
        class="preview"
        style:background-image={`linear-gradient(to right, rgba(0, 0, 0, 0.8), rgba(0, 0, 0, 0.1)), url("${list.cover || "https://placehold.co/250"}")`}
    >
        <h1>{list.title}</h1>
        <p>{list.created_by}</p>
    </div>
</a>

<style>
    .preview {
        border: 1px solid var(--border);
        border-radius: 4px;
        padding: 10px;
        width: 300px;
        display: flex;
        flex-direction: column;
        color: var(--primary-foreground-muted);

        background-position: center;
        background-size: cover;
        background-repeat: no-repeat;

        h1 {
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }

        &:hover {
            border: 1px solid var(--primary);
        }
    }

    a {
        color: inherit;
        text-decoration: none;
    }
</style>
