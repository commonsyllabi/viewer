<script setup lang="ts">
  import { ResourceType } from '../js/types'

  const props = defineProps<{
    lti: ResourceType
  }>()
</script>

<template>
  <div class="lti">
    <div class="lti-name py-1 px-2 ms-2 rounded-bottom">
      LTI - learning tools interop
    </div>
    <div class="meta p-3">
      <div class="title">{{ props.lti.Title }}</div>
      <div class="text">{{ props.lti.Description }}</div>
      <div class="launch-url">{{ props.lti.LaunchURL }}</div>
      <div class="launch-url-secure">{{ props.lti.SecureLaunchURL }}</div>
    </div>
  </div>
</template>

<style scoped>
  .lti {
    border: 2px solid lightgray;
  }

  .lti-name {
    background-color: lightgray;
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
