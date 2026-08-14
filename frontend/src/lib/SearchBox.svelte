<script lang="ts">
    import Input from "$lib/components/Input.svelte";
    import { title } from "$lib/utils";
    import type { MediaItem } from "$lib/types";
    import { useClerkContext } from "svelte-clerk";

    type SearchResponse = {
        next: string;
        items: MediaItem[];
    };

    let { items = $bindable() }: { items: MediaItem[] } = $props();

    const resultTypes = ["album", "movie"];
    let resultType = $state("");
    let searchQuery = $state("");
    let results: SearchResponse | null = $state(null);
    let timer: ReturnType<typeof setTimeout>;

    const ctx = useClerkContext();

    const search = (url: string = "") => {
        const q = searchQuery;

        timer = setTimeout(async () => {
            if (!searchQuery || !resultType) {
                results = null;
                return;
            }

            const token = await ctx.session?.getToken();

            if (url === "") {
                url = `/api/search?query=${q}&type=${resultType}`;
            }

            const res = await fetch(url, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
            });
            if (!res.ok) throw new Error("Search failed");
            results = await res.json();
        }, 500);

        return () => clearTimeout(timer);
    };

    $effect(() => search());
    $effect(() => {
        void resultType;

        searchQuery = "";
        results = null;
    });

    $effect(() => {
        if (results !== null) {
            items = results.items;
        }
    });
</script>

<form action="" class="list-form">
    <select bind:value={resultType}>
        <option value="" selected>Select a type to search for</option>
        {#each resultTypes as t}
            <option value={t}>{title(t)}</option>
        {/each}
    </select>
    <div class="search-container">
        <Input
            label="Search..."
            bind:value={searchQuery}
            disabled={resultType === ""}
        />
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
