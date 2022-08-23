<template>
  <footer class="container bottom">
    <div class="copyright">&#169; 2022 - GPLv3</div>
    <div class="authors">
      <a
        href="http://commonsyllabi.org"
        target="_blank"
        rel="noopener noreferrer"
        >Common Syllabi</a
      >
    </div>
  </footer>
</template>

<script setup lang="ts"></script>

<style lang="scss" scoped>
  footer {
    position: relative;
    bottom: 0;
    width: 100%;
    display: flex;
    justify-content: center;
    border-top: 1px solid lightgray;

    a {
      color: black;
      text-decoration: none;
    }

    a:hover {
      color: black;
      text-decoration: underline;
    }
  }

  .copyright,
  .authors {
    margin: 1rem;
  }

  .authors {
    display: flex;
    flex-direction: row;
  }

  .author {
    margin-right: 1rem;
  }
</style>
