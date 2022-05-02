<script setup lang="ts">
import { ResourceType } from "../js/types";

const props = defineProps<{
  topic: ResourceType;
}>();
</script>

<template>
  <h4>Discussion Topic</h4>
  <div>{{ props.topic.Title }}</div>
  <div>{{ props.topic.Text.Text }}</div>
</template>

<style scoped></style>
