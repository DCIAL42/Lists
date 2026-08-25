<script lang="ts">
    import favicon from "$lib/assets/favicon.svg";
    import { ClerkProvider } from "svelte-clerk";
    import "../app.css";
    import { ModeWatcher } from "mode-watcher";
    import Nav from "$lib/Nav.svelte";

    let { children } = $props();
</script>

<svelte:head>
    <link rel="icon" href={favicon} />
</svelte:head>

<ClerkProvider>
    <Nav />

    <ModeWatcher />
    {@render children()}
</ClerkProvider>
