<script lang="ts">
    import ListItemCard from "$lib/ListItemCard.svelte";

    let { data } = $props();
    let items = $derived(data.items);
    $inspect(data, items);
</script>

{#each items as item}
    <ListItemCard {item} />
{/each}
