<template>
    <header class="top">
        <h1 class="title"><a href="/">Common Syllabi Viewer</a></h1>
        <div class="nav-link">
            <a href="/about.html">about</a>
        </div>
    </header>
</template>

<script setup>
</script>

<style lang="scss">
header {
    width: 100%;
    display: flex;
    justify-content: space-between;
    padding: 1rem;
    background-color: white;
    border-bottom: 1px solid black;

    a {
        color: black;
        text-decoration: none;
    }

    a:hover {
        text-decoration: underline;
        color: black;
    }
}

.title {
    font-weight: bold;
    font-size: 1.9em;

    a, a:hover {
        color: black;
        text-decoration: none;
    }
}
</style>