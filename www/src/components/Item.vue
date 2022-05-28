<script setup lang="ts">
import { ref, computed } from "vue";
import Resource from "./Resource.vue";
import { ItemType } from "../js/types";

const props = defineProps<{
  item: ItemType,
  cartridge: string
}>();

const showResources = ref(true);
</script>

<template>
  <div class="item-title">{{ props.item.Item.Title }}</div>
  <div class="item-id">{{ props.item.Item.Identifier }}</div>

  <div v-if="props.item.Children">
    <h6 class="fs-6 text-uppercase">
      Children:
    </h6>
    <div v-for="child in props.item.Children" class="sub-item">
      <div>{{ child.Item.Identifier }} - {{ child.Item.Title }}</div>
    </div>
  </div>

  <ul v-if="showResources" class="item-resources">
    <li class="border rounded p-1 m-1 res" v-for="res in props.item.Resources">
      <!-- <Resource :resource="res" :cartridge="cartridge" /> -->
      <div class="res-title"> {{ res.Title }}</div>
      <div class="res-id"> {{ res.Identifier }}</div>
      <div class="res-type"> {{ res.Type }}</div>
    </li>
  </ul>
</template>

<style scoped>
.resources {
  cursor: pointer;
}

.item-id{
  color: lightgrey;
  font-size: 0.8em;
}

.item-resources{
  padding-left: 5px;
  border-left: 1px solid black;
}

.res{
  list-style: none;
  font-size: 0.7em;
}

.res-id{
  color: lightgray;
}

</style>
