<template>
  <div class="container p-3">
    <!-- upload -->
    <div class="container py-4 mb-3 border rounded"> 
      <form action="/api/upload" method="post" id="upload-form">
        <div class="form-group">
          <label for="cartridgeInput" class="d-block  h5 mb-1">Upload a common cartridge (.imscc) file</label>
          <input type="file" name="cartridge" class="form-control-file d-block  mb-2" id="upload-file" />
          <button id="upload-submit" type="button" class="btn btn-primary  mb-3" @click="upload()">upload</button>
        </div>
      </form>
    </div>

    <!-- status log -->
    <div id="log" class="container py-4 mb-3 border rounded text-light bg-dark">
      <pre>{{ log }}</pre>
    </div>

    <!-- viewer -->
    <!-- <div v-if="isUploaded" class="container cartridge"> -->
    <div class="container cartridge">
      <!-- file metadata -->
      <div class="row">
        <div class="metadata">
          <h2 class="section-title-small">Metadata</h2>
          <div class="title">
            {{ manifest.Metadata.Lom.General.Title.String.Text }}
          </div>
        </div>
      </div>

      <!-- file navigator -->
      <div class="row">
        <div class="col items">
          <h2 class="section-title-small">Items</h2>
          <!-- TODO: dont use index, bad prac -->
          <div class="item" v-for="(i, index) in items" :key="index">
            <Item :item="i" :cartridge="cartridge.name" />
            <hr />
          </div>
        </div>
        <div class="col resources">
          <h2 class="section-title-small">Resources</h2>
          <!-- TODO: dont use index, bad prac -->
          <div class="resource" v-for="(r, index) in resources" :key="index">
            <Resource :resource="r" :cartridge="cartridge.name" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { ManifestType, ItemType, ResourceType } from "./js/types";

import Resource from './components/Resource.vue'
import Item from './components/Item.vue'

let cartridge = reactive({ name: "" });
let manifest = reactive<ManifestType>({
  Metadata: { Lom: { General: { Title: { String: { Text: "" } } } } },
});
let items = reactive<Array<ItemType>>([
  { Item: { Identifier: "", Title: "" }, Children: [], Resources: [] },
]);
let resources = reactive<Array<ResourceType>>([
  {
    XMLName: { Local: "" },
    Title: "",
    Type: "",
    Identifier: "",
    File: [],
    Text: { Text: "" },
    Attachments: {
      Text: "",
      Attachment: [
        {
          Text: "",
          Href: "",
        },
      ],
    },
    Gradable: {
      Text: "",
      PointsPossible: "",
    },
    SubmissionFormats: {
      Text: "",
      Format: [
        {
          Text: "",
          Type: "",
        },
      ],
    },

    Description: "",
    LaunchURL: "",
    SecureLaunchURL: "",
    Vendor: {
      Text: "",
      Name: "",
      Description: "",
      URL: "",
    },

    Assessment: {
      Title: "",
      Text: "",
    },

    URL: {
      Text: "",
      Href: "",
    },
  },
]);

let log = ref("ready");
let isUploaded = ref(false);

let upload = function () {
  const formElem = document.getElementById("upload-form") as HTMLFormElement;
  const formData = new FormData(formElem);

  if (formData.get("cartridge") == null) {
    console.warn("can't submit an empty cartridge!");
    log.value = "can't submit an empty cartridge!";
    return;
  }

  cartridge.name = (formData.get("cartridge") as File).name;
  log.value = `uploading ${cartridge.name}`;

  fetch("http://localhost:2046/api/upload", {
    method: "POST",
    body: formData,
  })
    .then((res) => {
      return res.json();
    })
    .then((data) => {
      console.log(data);
      log.value = `uploaded ${cartridge.name}`;
      isUploaded.value = true;

      manifest = JSON.parse(data.data);
      items = JSON.parse(data.items);

      //-- todo, here we have to get rid of the Item field of the returned struct... what to do?
      for (let r of JSON.parse(data.resources)) {
        resources.push(r.Resource);
      }
    })
    .catch((err) => {
      log = ref(err);
      console.error(err);
    });
};
</script>

<style lang="scss">
@import "./css/global-vars.scss";

// page typograph

.section-title-small {
  font-size: 0.6rem;
  text-transform: uppercase;
  font-weight: 300;
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
