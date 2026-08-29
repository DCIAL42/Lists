<script lang="ts">
    import favicon from "$lib/assets/favicon.ico";
    import { ClerkProvider } from "svelte-clerk";
    import "../app.css";
    import { ModeWatcher } from "mode-watcher";
    import Nav from "$lib/Nav.svelte";
    import { PUBLIC_CLERK_PUBLISHABLE_KEY } from "$env/static/public";

    let { children } = $props();
</script>

<svelte:head>
    <link rel="icon" href={favicon} />
</svelte:head>

<ClerkProvider publishableKey={PUBLIC_CLERK_PUBLISHABLE_KEY}>
    <Nav />

    <ModeWatcher />
    {@render children()}
</ClerkProvider>
