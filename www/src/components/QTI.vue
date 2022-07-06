<script setup lang="ts">
import { ResourceType } from "../js/types";

const props = defineProps<{
  qti: ResourceType
}>();
</script>

<template>
  <div class="qti">
    <div class="qti-name p-1 rounded-bottom">QTI - Test</div>
    <div class="meta p-1">
      <div class="title">{{ props.qti.Assessment.Title }}</div>
      <div class="text">{{ props.qti.Assessment.Text }}</div>
    </div>
  </div>
</template>

<style scoped>
.qti {
  border: 2px solid #e9bf53;
}

.qti-name {
  background-color: #e9bf53;
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
