<script>
    import Input from "$lib/Input.svelte";
    import Select from "$lib/Select.svelte";
    let val = $state("");
    let val2 = $state("Test");
    const ops = ["op1", "op2"];
    let val3 = $state("");
    let things = [];
    for (let i = 0; i < 10; i++) {
        things.push(`thing${i}`);
    }
    $inspect(val3);
    let cols = $state(1);
    $inspect(cols);
</script>

<div class="outer">
    <Input label="Name" bind:value={val} />
    <Input label="Name" bind:value={val2} />
    <Select options={ops} bind:value={val3} label="Name" />
    <div style="display: block;">
        <button onclick={() => (cols = Math.max(1, cols - 1))}>-</button>
        <button onclick={() => (cols = Math.min(things.length, cols + 1))}
            >+</button
        >
    </div>
    <div class="grid" style="--n: {cols}">
        {#each things as thing}
            <div>{thing}</div>
        {/each}
    </div>
</div>

<style>
    .outer {
        padding: 5px;
        display: flex;
        flex-direction: column;
        gap: 25px;
    }

    .grid {
        display: grid;
        grid-template-columns: repeat(var(--n), auto);
        gap: 10px;
        padding: 10px;
    }

    .grid div {
        border: 1px solid var(--border);
    }
</style>
