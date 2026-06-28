<script>
    import AlbumCard from "$lib/AlbumCard.svelte";

    let { data } = $props();
    $inspect(data);
</script>

<div
    style="display: flex; flex-direction: column; gap: 10px; align-items: center;"
>
    {#each data.data as result}
        <AlbumCard album={result} />
    {/each}
</div>
