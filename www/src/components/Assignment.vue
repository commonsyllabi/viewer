<script setup lang="ts">
import { ResourceType } from "../js/types";

const props = defineProps<{
  assignment: ResourceType
}>();
</script>

<template>
  <div class="assignment">
    <div class="assignment-name p-1 rounded-bottom">Assignment</div>
    <div class="meta p-1">
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

.title{
  font-size: 0.9em;
}

.text{
  font-size: 0.8em;
}
</style>
