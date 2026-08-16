<script lang="ts">
    import favicon from "$lib/assets/favicon.svg";
    import Link from "$lib/components/Link.svelte";
    import { ClerkProvider, Show, SignInButton } from "svelte-clerk";
    import "../app.css";
    import Input from "$lib/components/Input.svelte";
    import Search from "$lib/icons/Search.svelte";
    import Button from "$lib/components/Button.svelte";
    import Moon from "$lib/icons/Moon.svelte";
    import Sun from "$lib/icons/Sun.svelte";
    import { ModeWatcher, toggleMode, mode } from "mode-watcher";

    let { children } = $props();
</script>

<svelte:head>
    <link rel="icon" href={favicon} />
</svelte:head>

<ClerkProvider>
    <nav>
        <Link href="/" class="page-title">Home</Link>
        <div class="page-links">
            <Link href="/lists">Lists</Link>
            <div class="vr"></div>
            <Link href="/search">Search</Link>
        </div>
        <div class="right-group">
            <div>
                <Input style="height: 10px;" label="search...">
                    {#snippet left_icon()}
                        <Search />
                    {/snippet}
                </Input>
            </div>
            <Button variant="ghost" onclick={toggleMode}>
                {#if mode.current === "light"}
                    <Moon />
                {:else}
                    <Sun />
                {/if}
            </Button>
            <Show when="signed-in">
                <Link href="/profile">Profile</Link>
            </Show>
            <Show when="signed-out">
                <SignInButton />
            </Show>
        </div>
    </nav>

    <ModeWatcher />
    {@render children()}
</ClerkProvider>

<style>
    nav {
        align-items: center;
    }

    .page-links {
        display: flex;
        gap: 20px;
        align-items: center;
    }

    .vr {
        border-left: 1px solid var(--primary);
        height: 1em;
    }

    .right-group {
        display: flex;
        gap: 10px;
        align-items: center;
    }
</style>
