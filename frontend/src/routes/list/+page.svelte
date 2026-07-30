<script lang="ts">
    import Link from "$lib/components/Link.svelte";
    import ListPreview from "$lib/ListPreview.svelte";

    let { data } = $props();
    let lists = $derived(data.lists);
    $inspect(lists);
</script>

<Link href="/list/create">Create new list</Link>

{#each lists as list}
    <ListPreview {list} />
{/each}
