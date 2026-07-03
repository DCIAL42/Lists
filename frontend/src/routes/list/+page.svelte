<script lang="ts">
    import AlbumCard from "$lib/AlbumCard.svelte";
    import Input from "$lib/Input.svelte";
    import SearchResults from "$lib/SearchResults.svelte";

    const resultTypes = ["album", "movie"];
    let resultType = $state("");
    let searchQuery = $state("");
    let results: Promise<any> | null = $state(null);
    let timer: ReturnType<typeof setTimeout>;
    $effect(() => {
        const q = searchQuery;

        timer = setTimeout(async () => {
            if (!searchQuery || !resultType) {
                results = null;
                return;
            }

            const url = `/api/search?query=${q}&type=${resultType}`;
            const res = await fetch(url);
            if (!res.ok) throw new Error("Search failed");
            results = await res.json();
        }, 500);

        return () => clearTimeout(timer);
    });
    $inspect(resultType, searchQuery, results);
</script>

<form action="" class="list-form">
    <select bind:value={resultType}>
        <option value="" selected>Select a type to search for</option>
        {#each resultTypes as t}
            <option value={t}>{t}</option>
        {/each}
    </select>
    <Input
        label="Search..."
        bind:value={searchQuery}
        disabled={resultType === ""}
    />
    <!-- <div class="search"> -->
    <!--     <input -->
    <!--         type="text" -->
    <!--         bind:value={searchQuery} -->
    <!--         disabled={resultType === ""} -->
    <!--     /> -->
    <!-- </div> -->
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

    .search {
        position: relative;
        background-color: inherit;
        margin: 15px;
    }

    .search input {
        background-color: inherit;
        border: 2px solid black;
        padding: 5px;
    }

    .search::after {
        content: "test";
        position: absolute;
        padding: 5px;
        left: 5px;
        transition: 0.2s ease;
        pointer-events: none;
    }

    .search:focus-within::after {
        background-color: black;
        transform: translateY(-50%);
    }
</style>
