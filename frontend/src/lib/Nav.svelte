<script lang="ts">
    import { goto } from "$app/navigation";
    import { toggleMode, mode } from "mode-watcher";
    import type { MediaItem, MediaType } from "$lib/types";
    import {
        Show,
        SignInButton,
        useClerkContext,
        UserButton,
    } from "svelte-clerk";
    import Link from "./components/Link.svelte";
    import Moon from "./icons/Moon.svelte";
    import Sun from "./icons/Sun.svelte";
    import SearchBar from "./SearchBar.svelte";
    import Button from "./components/Button.svelte";

    let items: MediaItem[] = $state([]);
    let tab: MediaType = $state("movie");
    let query = $state("");
    const ctx = useClerkContext();
    const username = $derived(ctx.user?.username);
</script>

<nav>
    <Link href="/" class="page-title">Home</Link>
    <div class="page-links">
        <Link href="/lists">Lists</Link>
        <div class="vr"></div>
        <Link href="/search">Search</Link>
    </div>
    <div class="right-group">
        <SearchBar
            bind:items
            bind:tab
            bind:query
            small
            floating
            onkeydown={(e) => {
                if (e.key === "Enter") {
                    e.preventDefault();
                    goto(`/search?query=${query}&type=${tab}`);
                    query = "";
                }
            }}
        />
        <Button variant="ghost" onclick={toggleMode}>
            {#if mode.current === "light"}
                <Moon />
            {:else}
                <Sun />
            {/if}
        </Button>
        <Show when="signed-in">
            <UserButton>
                <UserButton.MenuItems>
                    <UserButton.Link label="profile" href={`/${username}`}>
                        {#snippet labelIcon()}
                            <span>👤</span>
                        {/snippet}
                    </UserButton.Link>
                </UserButton.MenuItems>
            </UserButton>
        </Show>
        <Show when="signed-out">
            <SignInButton />
        </Show>
    </div>
</nav>

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
