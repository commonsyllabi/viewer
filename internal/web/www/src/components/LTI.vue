<script setup lang="ts">
import { ResourceType } from '../js/types';

const props = defineProps<{
    lti: ResourceType
}>()
</script>

<template>
    <h4>LTI</h4>
    <div>{{ props.lti.Title }}</div>
    <div>{{ props.lti.Description }}</div>
    <div>{{ props.lti.LaunchURL }}</div>
    <div>{{ props.lti.SecureLaunchURL }}</div>
</template>

<style scoped>

</style>