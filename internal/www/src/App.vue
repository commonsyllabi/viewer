<template>
  <h1>upload a common cartridge</h1>
  <div class="upload">
    <form action="/api/upload" method="post" id="upload-form">
      <input type="file" name="cartridge" id="upload-file" />
    </form>
    <button id="upload-submit" type="button" @click="upload()">upload</button>
  </div>
  <div id="log" class="log">{{ log }}</div>
  <div v-if="isUploaded" class="cartridge">
    <div class="metadata">
      <h2 class="section-title-small">Metadata</h2>
      <div class="title">{{ manifest.Metadata.Lom.General.Title.String.Text }}</div>
    </div>
    <div class="navigator-container">
      <div class="items">
          <h2 class="section-title-small">Items</h2>
          <!-- TODO: dont use index, bad prac -->
          <div class="item" v-for="(i, index) in items" :key="index">
            <Item :item="i" :cartridge="cartridge.name"/>
            <hr />
          </div>
        </div>
        <div class="resources">
          <h2 class="section-title-small">Resources</h2>
          <!-- TODO: dont use index, bad prac -->
          <div class="resource" v-for="(r, index) in resources" :key="index">
            <Resource :resource="r" :cartridge="cartridge.name" />
          </div>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ManifestType, ItemType, ResourceType } from './js/types'

import Resource from './components/Resource.vue'
import Item from './components/Item.vue'

let cartridge = reactive({ name: "" })
let manifest = reactive<ManifestType>({ Metadata: { Lom: { General: { Title: { String: { Text: "" } } } } } })
let items = reactive<Array<ItemType>>([{ Item: { Identifier: "", Title: "" }, Children: [], Resources: [] }])
let resources = reactive<Array<ResourceType>>([{
  XMLName: { Local: "" },
  Title: "",
  Type: "",
  Identifier: "",
  File: [],
  Text: { Text: "" },
  Attachments: {
    Text: "",
    Attachment: [{
      Text: "",
      Href: ""
    }]
  },
  Gradable: {
    Text: "",
    PointsPossible: ""
  },
  SubmissionFormats: {
    Text: "",
    Format: [{
      Text: "",
      Type: ""
    }]
  },
  
  Description: "",
  LaunchURL: "",
  SecureLaunchURL: "",
  Vendor: {
    Text: "",
    Name: "",
    Description: "",
    URL: ""
  },

  Assessment: {
    Title: "",
    Text: ""
  },

  URL: {
    Text: "",
    Href: ""
  }
}])

let log = ref("ready")
let isUploaded = ref(false)

let upload = function () {
  const formElem = document.getElementById("upload-form") as HTMLFormElement
  const formData = new FormData(formElem)

  if (formData.get("cartridge") == null) {
    console.warn("can't submit an empty cartridge!");
    log.value = "can't submit an empty cartridge!"
    return
  }

  cartridge.name = (formData.get("cartridge") as File).name
  log.value = `uploading ${cartridge.name}`

  fetch("http://localhost:2046/api/upload", {
    method: 'POST',
    body: formData
  }).then(res => {
    return res.json()
  }).then(data => {
    console.log(data);
    log.value = `uploaded ${cartridge.name}`
    isUploaded.value = true

    manifest = JSON.parse(data.data)
    items = JSON.parse(data.items)

    //-- todo, here we have to get rid of the Item field of the returned struct... what to do?
    for (let r of JSON.parse(data.resources)) {
      resources.push(r.Resource)
    }
  }).catch(err => {
    log = ref(err)
    console.error(err);
  })
}
</script>

<style lang="scss">
@import "./css/global-vars.scss";

// page typograph

.section-title-small {
  font-size: 0.6rem;
  text-transform: uppercase;
  font-weight: 300;
}

.upload {
  border-bottom: 1px solid $grey;
  box-sizing: border-box; 
  padding: 1rem 0 2rem 0; 
  form {
    width: 100%;
    #cartridge {
      width: 100%;
      border: 1px solid $grey;
      border-radius: 4px;
    }
  }
}
.log {
  border: 1px solid pink;
  padding: 1rem 0; 
}

.navigator-container {
  display: flex;
  flex-direction: row;
}

.cartridge {
  border: 1px solid orange;
  .metadata {
    border: 1px solid red;
  }
  .items {
    border: 1px solid blue;
  }
  .resources {
    border: 1px solid yellow;
  }
}
</style>
