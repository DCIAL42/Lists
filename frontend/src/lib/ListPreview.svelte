<script lang="ts">
    import Card from "./components/Card.svelte";
    import type { List } from "$lib/types";

    let { list }: { list: List } = $props();
</script>

<Card title={list.title} />
