<script lang="ts">
    import favicon from "$lib/assets/favicon.svg";
    import Link from "$lib/Link.svelte";
    import "../app.css";

    let { children } = $props();
</script>

<svelte:head>
    <link rel="icon" href={favicon} />
</svelte:head>

<nav>
    <Link href="/" class="page-title">Home</Link>
    <div class="page-links">
        <Link href="/music">Music</Link>
        <Link href="/movies">Movies</Link>
    </div>
    <div>
        <Link href="/profile">Profile</Link>
    </div>
</nav>

{@render children()}
