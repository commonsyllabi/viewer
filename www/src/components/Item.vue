<script setup lang="ts">
  import { ref, computed } from 'vue'
  import { ItemType } from '../js/types'

  const props = defineProps<{
    item: ItemType
    cartridge: string
  }>()

  let scrollo = (_id : string) => {
    console.warn("please implement me! I'm trying to do something with:",_id)
  }
</script>

<template>
  <h4 class="item-title fs-6 fw-bold">
    {{ props.item.Item.Title }}
  </h4>
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
          <div
            v-if="child.Item.Identifier"
            class="item-ref d-flex align-items-baseline"
          >
            Refers to:
            <button
              class="btn btn-link btn-sm p-0"
              @click="scrollo(child.Item.Identifier)"
            >
              <span class="text-break">
                {{ child.Item.Identifier }}
              </span>
            </button>
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
