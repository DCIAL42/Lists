<script>
    let { children } = $props();
</script>

<div
    style="display: flex; flex-direction: column; gap: 10px; align-items: center;"
>
    {@render children()}
</div>
