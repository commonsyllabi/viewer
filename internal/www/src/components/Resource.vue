<script setup lang="ts">
import { ref } from 'vue'
import { ResourceType } from '../js/types'

import Assignment from './Assignment.vue'
import DiscussionTopic from './DiscussionTopic.vue'
import LTI from './LTI.vue'
import QTI from './QTI.vue'
import Weblink from './Weblink.vue'

const props = defineProps<{
  resource: ResourceType,
  cartridge: string
}>()

const previewPath = ref("")

function getFile(_evt: Event, _id: string) {
  fetch(`/api/file/${_id}?cartridge=${props.cartridge}`, {
    method: 'GET'
  })
    .then(res => { 
      return res.blob().then(blob => {
        return {
          contentType: res.headers.get("Content-Type"),
          raw: blob
        }
      })
      // previewPath.value = createFileURL(res) 
    })
    .then(data => {
      console.log(data);
      
      previewPath.value = URL.createObjectURL(data.raw)
    })
    .catch(err => console.error(err));
}
</script>

<template>
  <div>
    <div v-if="props.resource.XMLName.Local === 'assignment'">
      <Assignment :assignment="resource" />
    </div>

    <div v-if="props.resource.XMLName.Local === 'topic'">
      <DiscussionTopic :topic="resource" />
    </div>

    <div v-if="props.resource.XMLName.Local === 'webLink'">
      <Weblink :weblink="resource" />
    </div>

    <div v-if="props.resource.XMLName.Local === 'questestinterop'">
      <QTI :qti="resource" />
    </div>

    <div v-if="props.resource.XMLName.Local === 'cartridge_basiclti_link'">
      <LTI :lti="resource" />
    </div>

    <div v-if="props.resource.XMLName.Local === 'resource'">
      <div class="meta">
        <div>
          xml:
          <span class="resource-value">{{ props.resource.XMLName.Local }}</span>
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

      <ol class="resource-files">
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
