<script lang="ts">
    import { page } from "$app/state";
</script>

<div>
    <h1>{page.status}</h1>
    <p>{page.error?.message}</p>
</div>

<style>
    div {
        text-align: center;
    }
</style>
