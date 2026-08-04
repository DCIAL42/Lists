<script lang="ts">
    import favicon from "$lib/assets/favicon.svg";
    import Link from "$lib/components/Link.svelte";
    import { ClerkProvider, Show, SignInButton } from "svelte-clerk";
    import "../app.css";

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
            <!-- <Link href="/movies">Movies</Link> -->
        </div>
        <div>
            <Link href="/profile">Profile</Link>
            <Show when="signed-out">
                <SignInButton />
            </Show>
        </div>
    </nav>

    {@render children()}
</ClerkProvider>
