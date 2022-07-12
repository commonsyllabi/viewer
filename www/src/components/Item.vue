<script setup lang="ts">
  import { ref, computed } from 'vue'
  import Resource from './Resource.vue'
  import { ItemType } from '../js/types'

  const props = defineProps<{
    item: ItemType
    cartridge: string
  }>()

  const showResources = ref(true)
</script>

<template>
  <h6 class="item-title fs-6 fw-bold">
    {{ props.item.Item.Title }}
  </h6>
  <div class="item-id">{{ props.item.Item.Identifier }}</div>
  <div class="item-children-container">
    <div v-if="props.item.Children" class="ps-2">
      <div
        v-for="child in props.item.Children"
        :key="child.Item.Identifier"
        class="sub-item px-2 py-1 m-1"
      >
        <div>
          <div>
            {{ child.Item.Title }}
          </div>
          <div class="item-id">
            {{ child.Item.Identifier }}
          </div>
          <div v-if="child.Item.Identifierref" class="item-ref">
            Refers to:
            <span class="text-break">{{ child.Item.Identifierref }}</span>
          </div>
          <div v-if="child.Resources" class="item-type">
            resource type:
            <span>
              {{ child.Resources[0].Type }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>

  <ul v-if="showResources" class="item-resources">
    <li
      class="border rounded p-1 m-1 res"
      v-for="res in props.item.Resources"
      :key="res.Identifier"
    >
      <!-- <Resource :resource="res" :cartridge="cartridge" /> -->
      <div class="res-title">{{ res.Title }}</div>
      <div class="res-id">{{ res.Identifier }}</div>
      <div class="res-type">{{ res.Type }}</div>
    </li>
  </ul>
</template>

<style lang="scss" scoped>
  @import '../css/global-vars.scss';
  .resources {
    cursor: pointer;
  }

  .item-id,
  .item-ref,
  .item-type {
    color: $secondary;
    font-size: 0.8em;
  }

  .item-resources {
    padding-left: 5px;
    border-left: 1px solid black;
  }

  .item-children-container {
    border-left: 1px solid $secondary;
  }

  .sub-item {
    border: 2px solid $gray-200;
    border-radius: $border-radius-sm;
  }

  .res {
    list-style: none;
    font-size: 0.7em;
  }

  .res-id {
    color: lightgray;
  }
</style>
