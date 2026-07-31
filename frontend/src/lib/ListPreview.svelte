<script lang="ts">
    import Card from "./components/Card.svelte";

    let { list } = $props();
</script>

<Card title={list.title} />
