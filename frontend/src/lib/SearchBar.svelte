<script lang="ts">
    import Input from "$lib/components/Input.svelte";
    import Search from "$lib/icons/Search.svelte";
    import type { MediaItem, MediaType, SearchResponse } from "$lib/types";
    import type { HTMLAttributes } from "svelte/elements";
    import SearchResults from "./SearchResults.svelte";

    let {
        items = $bindable([]),
        tab = $bindable("movie"),
        query = $bindable(""),
        small = false,
        add = false,
        floating = true,
        handleAdd = () => {},
        ...rest
    }: {
        items?: MediaItem[];
        tab?: MediaType;
        query?: string;
        small?: boolean;
        add?: boolean;
        floating?: boolean;
        handleAdd?: (item: MediaItem) => void;
    } & HTMLAttributes<HTMLInputElement> = $props();

    let results: SearchResponse | null = $state(null);
    let timer: ReturnType<typeof setTimeout>;
    let focused = $state(false);
    let container: HTMLDivElement;

    const search = (url: string = "") => {
        const q = query;

        timer = setTimeout(async () => {
            if (!query) {
                results = null;
                return;
            }

            if (url === "") {
                url = `/api/search?query=${q}&type=all`;
            }

            const res = await fetch(url);
            if (!res.ok) throw new Error("Search failed");
            results = await res.json();
        }, 500);

        return () => clearTimeout(timer);
    };

    $effect(() => search());

    $effect(() => {
        items = results?.[tab]?.items || [];
    });
</script>

<div
    class="outer"
    bind:this={container}
    onfocusin={() => (focused = true)}
    onfocusout={(e) => {
        const next = e.relatedTarget as Node | null;

        if (!next || !container?.contains(next)) {
            focused = false;
        }
    }}
>
    <Input bind:value={query} style="height: 10px;" label="search..." {...rest}>
        {#snippet left_icon()}
            <Search />
        {/snippet}
    </Input>
    {#if query !== "" && (focused || !floating)}
        <SearchResults bind:items {small} {add} bind:tab {handleAdd} />
    {/if}
</div>

<style>
    .outer {
        position: relative;
        width: 100%;
    }
</style>
