<script setup lang="ts">
  import { ResourceType } from '../js/types'

  const props = defineProps<{
    assignment: ResourceType
  }>()
</script>

<template>
  <div class="assignment">
    <div class="assignment-name py-1 px-2 ms-2 rounded-bottom">Assignment</div>
    <div class="meta p-3">
      <div class="title">{{ props.assignment.Title }}</div>
      <div class="text" v-html="props.assignment.Text.Text"></div>
    </div>
  </div>
</template>

<style scoped>
  .assignment {
    border: 2px solid lightgrey;
  }

  .assignment-name {
    background-color: lightgrey;
    color: white;
    width: max-content;
    font-size: 0.8em;
  }

  .title {
    font-size: 0.9em;
  }

  .text {
    font-size: 0.8em;
  }
</style>
