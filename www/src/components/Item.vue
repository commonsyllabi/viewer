<script setup lang="ts">
import { ref, computed } from "vue";
import Resource from "./Resource.vue";
import { ItemType } from "../js/types";

const props = defineProps<{
  item: ItemType,
  cartridge: string
}>();

const showResources = ref(false);
</script>

<template>
  <h6 class="fs-6 text-uppercase">
    Title:
  </h6>
  <div>{{ props.item.Item.Title }}</div>

  <div v-if="props.item.Children">
    <h6 class="fs-6 text-uppercase">
      Children:
    </h6>
    <div v-for="child in props.item.Children" class="sub-item">
      <div>{{ child.Item.Identifier }} - {{ child.Item.Title }}</div>
    </div>
  </div>

  <h4 class="resources" @click="showResources = !showResources">
    Resources
  </h4>
  <ul v-for="res in props.item.Resources" v-if="showResources" class="sub-resource">
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
