<script lang="ts">
    import Input from "$lib/components/Input.svelte";
    import Select from "$lib/components/Select.svelte";

    let testOption = $state("");
</script>

<main>
    <Input />

    <Select
        options={["op 1", "op 2"]}
        label="Test select"
        bind:value={testOption}
    />
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        gap: 5px;
        padding: 5px;
    }
</style>
