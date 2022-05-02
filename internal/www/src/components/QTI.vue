<script setup lang="ts">
import { ResourceType } from "../js/types";

const props = defineProps<{
  qti: ResourceType
}>();
</script>

<template>
  <h4>QTI</h4>
  <div>{{ props.qti.Assessment.Title }}</div>
  <div>{{ props.qti.Assessment.Text }}</div>
</template>

<style scoped></style>
