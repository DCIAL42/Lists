<script lang="ts">
    import AlbumCard from "$lib/components/AlbumCard.svelte";
    import Input from "$lib/components/Input.svelte";
    import Link from "$lib/components/Link.svelte";
    import MovieCard from "$lib/components/MovieCard.svelte";
    import SearchResults from "$lib/SearchResults.svelte";
    import { isAlbum, isMovie, title } from "$lib/utils";
    import type { MediaItem } from "$lib/types";

    type SearchResponse = {
        next: string;
        items: MediaItem[];
    };

    let { items = $bindable() }: { items: MediaItem[] } = $props();

    const resultTypes = ["album", "movie"];
    let resultType = $state("");
    let searchQuery = $state("");
    let results: Promise<SearchResponse> | null = $state(null);
    let timer: ReturnType<typeof setTimeout>;

    const search = (url: string = "") => {
        const q = searchQuery;

        timer = setTimeout(async () => {
            if (!searchQuery || !resultType) {
                results = null;
                return;
            }

            if (url === "") {
                url = `/api/search?query=${q}&type=${resultType}`;
            }

            const res = await fetch(url);
            if (!res.ok) throw new Error("Search failed");
            results = await res.json();
        }, 500);

        return () => clearTimeout(timer);
    };

    let focused = $state(false);
    let container: HTMLDivElement;

    $effect(() => search());
    $effect(() => {
        void resultType;

        searchQuery = "";
        results = null;
    });
</script>

<form action="" class="list-form">
    <select bind:value={resultType}>
        <option value="" selected>Select a type to search for</option>
        {#each resultTypes as t}
            <option value={t}>{title(t)}</option>
        {/each}
    </select>
    <div
        class="search-container"
        bind:this={container}
        onfocusin={() => (focused = true)}
        onfocusout={(e) => {
            const next = e.relatedTarget as Node | null;

            if (!next || !container.contains(next)) {
                focused = false;
            }
        }}
    >
        <Input
            label="Search..."
            bind:value={searchQuery}
            disabled={resultType === ""}
        />
        {#await results}
            <p>Loading...</p>
        {:then data}
            {#if data !== null && data.items?.length > 0 && focused}
                <SearchResults
                    onmousedown={(e) => e.preventDefault()}
                    style="z-index: 10000;"
                >
                    {#if data.next !== ""}
                        <Link
                            href="#top"
                            onclick={() => search("api" + data?.next)}
                        >
                            Next
                        </Link>
                    {/if}
                    {#each data.items as item}
                        {#if isAlbum(item)}
                            <AlbumCard
                                album={item.data}
                                add
                                onclick={() => items.push(item)}
                            />
                        {:else if isMovie(item)}
                            <MovieCard
                                movie={item.data}
                                add
                                onclick={() => {
                                    items.push(item);
                                }}
                            />
                        {/if}
                    {/each}
                </SearchResults>
            {/if}
        {:catch error}
            <p>Error: {error.message}</p>
        {/await}
    </div>
</form>

<style>
    .list-form {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        margin: 10px;
    }

    .search-container {
        position: relative;
        display: flex;
        flex-direction: column;
        align-items: center;
    }
</style>
