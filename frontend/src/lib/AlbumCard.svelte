<script>
    const { album } = $props();
</script>

<div class="album-card">
    <img
        src={album.cover}
        class="album-cover"
        alt={`${album.title} by ${album.artist} album cover`}
    />
    <h1 class="album-title">{album.title}</h1>
    <p class="album-subtitle">{album.artist}</p>
</div>

<style>
    .album-card {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 300px;
        border: 1px, solid, gray;
        text-align: center;
        padding: 10px;
        border-radius: 10px;
        transition: 0.5s ease;
    }

    .album-card:hover {
        box-shadow: gray 2px 2px;
        transform: scale(1.005);
    }

    .album-cover {
        width: 250px;
    }

    .album-subtitle {
        color: lightgray;
    }
</style>
