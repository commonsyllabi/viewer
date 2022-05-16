<template>
    <footer class="container bottom">
        <div class="copyright">
 	&#169; GPL3
        </div>
        <div class="authors">
            <div class="author">@periode</div>
            <div class="author">@patshiu</div>
        </div>
    </footer>
</template>

<script  setup lang="ts">
    
</script>

<style>
    footer {
        position: relative;
        bottom: 0;
        width: 100%;
        display: flex;
        justify-content: center;
    }

    .copyright, .authors{
        margin: 1rem;
    }

    .authors {
        display: flex;
        flex-direction: row;
    }

    .author{
        margin-right: 1rem;
    }

    
</style>