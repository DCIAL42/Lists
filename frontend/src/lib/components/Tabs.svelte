<script lang="ts">
    import { title } from "$lib/utils";
    import Button from "./Button.svelte";

    let {
        tabs,
        selected = $bindable(),
    }: {
        tabs: string[];
        selected: string | string[];
    } = $props();
</script>

<div class="tabs" style:--n={tabs.length * 2 - 1}>
    {#each tabs as tab, i}
        <Button
            variant="ghost"
            selected={Array.isArray(selected)
                ? selected.includes(tab)
                : selected === tab}
            onclick={() => {
                if (Array.isArray(selected)) {
                    if (selected.includes(tab)) {
                        selected = selected.filter((e) => e !== tab);
                    } else {
                        selected.push(tab);
                    }
                } else {
                    selected = tab;
                }
            }}
        >
            {title(tab)}
        </Button>
        {#if i < tabs.length - 1}
            <hr />
        {/if}
    {/each}
</div>

<style>
    .tabs {
        display: grid;
        margin-inline: auto;
        grid-template-columns: repeat(var(--n), 1fr);
    }
</style>
