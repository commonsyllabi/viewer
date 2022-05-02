<script setup lang="ts">
import { ResourceType } from "../js/types";

const props = defineProps<{
  assignment: ResourceType
}>();
</script>

<template>
  <h4>Assignment</h4>
  <div>{{ props.assignment.Title }}</div>
  <div>{{ props.assignment.Text.Text }}</div>
</template>

<style scoped></style>
