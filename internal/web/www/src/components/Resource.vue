<script setup>
import { ref } from 'vue'

const props = defineProps({
  resource: {},
  cartridge: String
})

const previewPath = ref("")

function getFile(_evt, _id,) {
  fetch(`/api/file/${_id}?cartridge=${props.cartridge}`, {
    method: 'GET'
  })
    .then(res => { return res.json() })
    .then(body => {
      previewPath = body.path
    })
    .catch(err => console.error(err));
}
</script>

<template>
  <div>
    <div class="meta">
      <div>
        xml:
        <span class="resource-value">{{ props.resource.Resource.XMLName.Local }}</span>
      </div>
      <div>
        title:
        <span class="resource-value">{{ props.resource.Resource.Title }}</span>
      </div>
      <div>
        type:
        <span class="resource-value">{{ props.resource.Resource.Type }}</span>
      </div>
      <div>
        id:
        <span class="resource-value">{{ props.resource.Resource.Identifier }}</span>
      </div>
    </div>

    <ol class="resource-files" v-if="props.resource.Resource.XMLName.Local == 'resource'">
      <li v-for="f in props.resource.Resource.File">
        <div
          class="resource-file"
          @click="getFile($event, props.resource.Resource.Identifier)"
        >{{ f.Href }}</div>
        <div class="preview">
          <iframe :src="previewPath" frameborder="0"></iframe>
        </div>
      </li>
    </ol>
  </div>

  <hr />
</template>

<style scoped>
.resource {
  margin: 5px;
  padding: 5px;
  border: 2px solid black;
}

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
