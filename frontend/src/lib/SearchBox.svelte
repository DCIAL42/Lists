<script lang="ts">
    import Input from "$lib/components/Input.svelte";
    import { title } from "$lib/utils";
    import type { MediaItem } from "$lib/types";
    import Select from "./components/Select.svelte";

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

    $effect(() => search());

    $effect(() => {
        void resultType;

        searchQuery = "";
        results = null;
    });

    $effect(() => {
        items = results?.items || [];
    });
</script>

<form action="" class="list-form">
    <Select
        label="Select a type to search for"
        bind:value={resultType}
        options={resultTypes}
        displayText={(t) => title(t)}
    />
    <Input
        label="Search..."
        bind:value={searchQuery}
        disabled={resultType === ""}
    />
</form>

<style>
    .list-form {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        margin: 10px;
    }
</style>
