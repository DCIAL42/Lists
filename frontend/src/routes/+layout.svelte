<script lang="ts">
    import favicon from "$lib/assets/favicon.svg";

    let { children } = $props();
</script>

<svelte:head>
    <link rel="icon" href={favicon} />
</svelte:head>

<nav>
    <a href="/" class="page-title">Home</a>
    <div class="page-links">
        <a href="/music">Music</a>
        <a href="/movies">Movies</a>
    </div>
    <div>
        <a href="/profile">Profile</a>
    </div>
</nav>

{@render children()}

<style>
    nav {
        border-bottom: 1px solid grey;
        margin-bottom: 5px;
        padding: 10px;
        display: flex;
        justify-content: space-between;
    }

    .page-title {
        text-decoration: none;
        color: inherit;
    }

    .page-links {
        position: absolute;
        left: 50%;
        transform: translateX(-50%);
    }

    .page-links a {
        text-decoration: none;
        color: inherit;
    }

    .page-links a:hover {
        text-decoration: underline solid white 2px;
        text-underline-offset: 6px;
    }
</style>
