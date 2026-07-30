<script lang="ts">
    import Card from "./components/Card.svelte";

    let { list } = $props();

    $inspect(list);
</script>

<Card title={list.title} />
