<script setup lang="ts">
import { ref } from 'vue'
import { ResourceType } from '../js/types'

const props = defineProps<{
  resource: ResourceType,
  cartridge: string
}>()

const previewPath = ref("")

function getFile(_evt: Event, _id: string) {
  fetch(`/api/file/${_id}?cartridge=${props.cartridge}`, {
    method: 'GET'
  })
    .then(res => { return res.json() })
    .then(body => {
      previewPath.value = body.path
    })
    .catch(err => console.error(err));
}
</script>

<template>
  <div>
    <div class="meta">
      <div>
        xml:
        <span class="resource-value">{{ props.resource.XMLName.Local }}</span>
      </div>
      <div>
        title:
        <span class="resource-value">{{ props.resource.Title }}</span>
      </div>
      <div>
        type:
        <span class="resource-value">{{ props.resource.Type }}</span>
      </div>
      <div>
        id:
        <span class="resource-value">{{ props.resource.Identifier }}</span>
      </div>
    </div>

    <ol class="resource-files" v-if="props.resource.XMLName.Local == 'resource'">
      <li v-for="f in props.resource.File">
        <div
          class="resource-file"
          @click="getFile($event, props.resource.Identifier)"
        >{{ f.Href }}</div>
        <div class="preview" v-if="previewPath != ''">
          <iframe :src="previewPath" frameborder="0"></iframe>
        </div>
      </li>
    </ol>
  </div>
</template>

<style scoped>

.meta {
  display: flex;
  flex-direction: column;
  font-size: 0.9em;
}

.resource-files {
  font-size: 0.85em;
}

.resource-file:hover {
  text-decoration: underline;
  cursor: pointer;
}

.preview {
  margin: 5px;
  padding: 5px;
  border: 2px solid black;
}

.preview iframe {
  width: 100%;
}
</style>
