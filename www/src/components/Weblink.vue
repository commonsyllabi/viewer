<script setup lang="ts">
import { ResourceType } from "../js/types";

const props = defineProps<{
  weblink: ResourceType
}>();
</script>

<template>
     <div class="weblink">
    <div class="weblink-name p-1 rounded-bottom">Weblink</div>
    <div class="meta p-1">
      <div class="title">{{ props.weblink.Title }}</div>
      <div class="text">{{ props.weblink.URL.Text }}</div>
      <div class="url">{{ props.weblink.URL.Href }}</div>
    </div>
  </div>
</template>

<style scoped>
.weblink {
  border: 2px solid #9bcbf6;
}

.weblink-name {
  background-color: #9bcbf6;
  color: white;
  width: max-content;
  font-size: 0.8em;
}

.title{
  font-size: 0.9em;
}

.text{
  font-size: 0.8em;
}
</style>
