<script setup lang="ts">
  import { ref, computed } from "vue";
  import Resource from "./Resource.vue";
  import { ItemType } from "../js/types";

  const props = defineProps<{
    item: ItemType;
    cartridge: string;
  }>();

  const showResources = ref(false);
</script>

<template>
  <h6 class="fs-6 text-uppercase">Title:</h6>
  <h3>{{ props.item.Item.Title }}</h3>

  <div v-if="props.item.Children">
    <h6 class="fs-6 text-uppercase">Children:</h6>
    <div class="sub-item" v-for="child in props.item.Children">
      <div>{{ child.Item.Identifier }} - {{ child.Item.Title }}</div>
    </div>
  </div>

  <h4 @click="showResources = !showResources" class="resources">Resources</h4>
  <ul
    v-if="showResources"
    class="sub-resource"
    v-for="res in props.item.Resources"
  >
    <li>
      <Resource :resource="res" :cartridge="cartridge" />
    </li>
  </ul>
</template>

<style scoped>
.resources {
  cursor: pointer;
}
</style>
