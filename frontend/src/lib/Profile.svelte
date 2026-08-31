<script lang="ts">
    import Skeleton from "./components/Skeleton.svelte";
    import ListPreview from "./ListPreview.svelte";
    import type { ProfileData } from "./types";

    const {
        data,
        loading = false,
    }: {
        data?: ProfileData;
        loading?: boolean;
    } = $props();
</script>

<div>
    <h2>Recent lists</h2>
    <hr />
    <div class="recent-lists">
        {#if loading || !data}
            {#each Array(4) as _}
                <Skeleton width="25%" height={170} />
            {/each}
        {:else}
            {#each data.listsData.lists.slice(0, 4) as list}
                <ListPreview {list} />
            {/each}
        {/if}
    </div>
</div>

<style>
    .recent-lists {
        display: flex;
        gap: 5px;
    }
</style>
