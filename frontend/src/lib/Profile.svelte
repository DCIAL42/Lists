<script lang="ts">
    import type { PageData } from "../routes/[username]/$types";
    import ListPreview from "./ListPreview.svelte";

    const { data }: { data: PageData } = $props();
</script>

<div>
    <h2>Recent lists</h2>
    <hr />
    <div class="recent-lists">
        {#each data.listsData.lists.slice(0, 4) as list}
            <ListPreview {list} />
        {/each}
    </div>
</div>

<style>
    .recent-lists {
        display: flex;
        gap: 5px;
    }
</style>
