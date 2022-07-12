<script setup lang="ts">
  import { ResourceType } from '../js/types'

  const props = defineProps<{
    topic: ResourceType
  }>()
</script>

<template>
  <div class="discussion-topic">
    <div class="discussion-topic-name py-1 px-2 ms-2 rounded-bottom">
      Discussion Topic
    </div>
    <div class="meta p-3">
      <div class="title">{{ props.topic.Title }}</div>
      <div class="text" v-html="props.topic.Text.Text"></div>
    </div>
  </div>
</template>

<style scoped>
  .discussion-topic {
    border: 2px solid #cd794a;
  }

  .discussion-topic-name {
    background-color: #cd794a;
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
