<script lang="ts">
    import AlbumCard from "$lib/AlbumCard.svelte";
    import SearchResults from "$lib/SearchResults.svelte";

    const resultTypes = ["album", "movie"];
    let resultType: string = $state("");
    let searchQuery = $state("");
    let results: Promise<any> | null = $state(null);
    let timer: ReturnType<typeof setTimeout>;
    $effect(() => {
        const q = searchQuery;

        timer = setTimeout(async () => {
            if (!searchQuery) {
                results = null;
                return;
            }

            const res = await fetch(`/api/search?query=${q}`);
            if (!res.ok) throw new Error("Search failed");
            results = await res.json();
        }, 250);

        return () => clearTimeout(timer);
    });
    // $inspect(resultType, searchQuery, results);
</script>

<form action="" class="list-form">
    <select bind:value={resultType}>
        <option value="" selected>Select a type to search for</option>
        {#each resultTypes as t}
            <option value={t}>{t}</option>
        {/each}
    </select>
    <input type="text" bind:value={searchQuery} />
</form>

{#await results}
    <p>Loading...</p>
{:then data}
    <SearchResults>
        {#each data as item}
            <AlbumCard album={item} />
        {/each}
    </SearchResults>
{:catch error}
    <p>Error: {error.message}</p>
{/await}

<style>
    .list-form {
        display: flex;
        flex-direction: column;
        align-items: center;
        margin: 10px;
    }
</style>
