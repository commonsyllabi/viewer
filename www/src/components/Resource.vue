<template>
  <div class="single-resource m-1 p-1">
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

    <div v-if="props.resource.XMLName.Local === 'resource'" class="resource">
      <div class="resource-name py-1 px-2 ms-2 rounded-bottom">
        {{ props.resource.XMLName.Local }}
      </div>
      <div class="meta p-3">
        <!-- {{ props.resource.Type }} -->
        <div class="resource-id" :id="props.resource.Identifier">
          {{ props.resource.Identifier }}
        </div>
      </div>

      <div class="p-3" v-if="props.resource.File.length > 0">
        <div>Files</div>
        <ol class="resource-files">
          <li v-for="f in props.resource.File" :key="f.Href">
            <div
              class="resource-file"
              @click="getFile($event, props.resource.Identifier)"
            >
              {{ f.Href }}
            </div>
            <div v-if="previewPath != ''" class="preview">
              <iframe :src="previewPath" frameborder="0" />
            </div>
            <a v-if="previewPath != ''" :href="previewPath" download
              >download</a
            >
          </li>
        </ol>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { ResourceType } from '../js/types'

  import Assignment from './Assignment.vue'
  import DiscussionTopic from './DiscussionTopic.vue'
  import LTI from './LTI.vue'
  import QTI from './QTI.vue'
  import Weblink from './Weblink.vue'

  const props = defineProps<{
    resource: ResourceType
    cartridge: string
  }>()

  const previewPath = ref('')

  function getFile(_evt: Event, _id: string) {
    fetch(`/api/file/${_id}?cartridge=${props.cartridge}`, {
      method: 'GET',
    })
      .then((res) => {
        return res.blob().then((blob) => {
          return {
            contentType: res.headers.get('Content-Type'),
            raw: blob,
          }
        })
        // previewPath.value = createFileURL(res)
      })
      .then((data) => {
        console.log(data)

        previewPath.value = URL.createObjectURL(data.raw)
      })
      .catch((err) => console.error(err))
  }
</script>

<style lang="scss" scoped>
  @import '../css/global-vars.scss';
  .resource {
    border: 2px solid #b2d88d;
    border-radius: $border-radius-sm;
  }

  .resource-name {
    background-color: #b2d88d;
    color: white;
    width: max-content;
    font-size: 0.8em;
  }

  .meta {
    display: flex;
    flex-direction: column;
    font-size: 0.9em;
  }

  .resource-id {
    color: lightgrey;
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
