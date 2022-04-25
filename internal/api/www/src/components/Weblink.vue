<script setup lang="ts">
import { ResourceType } from '../js/types';

const props = defineProps<{
    weblink: ResourceType
}>()
</script>

<template>
    <h4>Weblink</h4>
    <div>{{ props.weblink.Title }}</div>
    <div>{{ props.weblink.URL.Text }}</div>
    <div>{{ props.weblink.URL.Href }}</div>
</template>

<style scoped>

</style>