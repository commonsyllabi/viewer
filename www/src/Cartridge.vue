<template>
  <Header></Header>
  <main class="container p-3">

    <!-- upload form -->
    <div class="container pt-4 mb-3 border rounded">
      <form id="upload-form" action="/api/upload" method="post">
        <div class="form-group">
          <label for="cartridgeInput" class="d-block h5 mb-3">Upload a common cartridge (.imscc) file</label>
          <input id="upload-file" type="file" name="cartridge" class="form-control-file d-block mb-2" />
          <!-- <button class="btn btn-primary mb-3" ">upload</button> -->
          <uiButton @click="upload()" text="Upload"  classes="btn btn-primary mb-3" />
        </div>
      </form>
    </div>

    <!-- upload modal -->
    <div class="modal" @keydown.esc="showUpload = false" tabindex="-1" v-if="showUpload">
      <Upload :title="manifest.Metadata.Lom.General.Title.String.Text"
        :description="manifest.Metadata.Lom.General.Description.String.Text" @close="showUpload = false"></Upload>
    </div>


    <!-- status log -->
    <div id="log" class="container py-1 mb-3">
      <pre>{{ log }}</pre>
    </div>


    <div v-if="isUploaded" class="container cartridge">
      <!-- actions -->
      <div class="row mb-2 actions">
        <button type="button" class="action" @click="showMetadata = !showMetadata">{{
            showMetadata ? "hide metadata & index" : "show metadata & index"
        }}</button>
        <button type="button" class="action" @click="reset()">load another file</button>
        <button type="button" class="action" @click="showModal()">upload this file</button>
      </div>
      <!-- metadata viewer -->
      <div class="row mb-3" v-if="showMetadata">
        <div id="metadata-accordion" class="accordion">
          <div class="accordion-item">
            <h2 id="headingOne" class="accordion-header">
              <button class="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseMeta"
                aria-expanded="true" aria-controls="collapseOne">Cartridge Metadata</button>
            </h2>
            <div id="collapseMeta" class="accordion-collapse collapse show" aria-labelledby="headingOne"
              data-bs-parent="#metadata-accordion">
              <div class="accordion-body">
                <!-- file metadata placeholder -->
                <div v-if="!isUploaded" class="metadata-placeholder">
                  <span class="text-muted">
                    <em>metadata goes here</em>
                  </span>
                </div>

                <!-- file metadata -->
                <div v-if="isUploaded" class="metadata">
                  <div class="title">{{ manifest.Metadata.Lom.General.Title.String.Text }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- viewer -->
      <div class="container border rounded mb-5 cartridge-viewer">
        <!-- file navigator -->
        <div class="row">
          <!-- items panel -->
          <div class="col-6 overflow-scroll items-panel">
            <h6 class="my-2">Items Who Knows Index</h6>

            <!-- items listing -->
            <!-- TODO: dont use index, bad prac -->
            <div v-for="(i, index) in items" :key="index">
              <Item :item="i" :cartridge="cartridge.name" />
              <hr />
            </div>
          </div>

          <!-- resources panel -->
          <div class="col-6 overflow-scroll resources-panel">
            <h6 class="my-2">Resources Index</h6>

            <!-- resources listing -->
            <!-- TODO: dont use index, bad prac -->
            <div v-for="(r, index) in resources" :key="index" class>
              <Resource :resource="r" :cartridge="cartridge.name" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
  <Footer></Footer>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ManifestType, ItemType, ResourceType } from "./js/types";

import Resource from './components/Resource.vue'
import Item from './components/Item.vue'
import Upload from './components/Upload.vue'
import uiButton from './components/ui/ui-button.vue'

import Footer from './components/Footer.vue';
import Header from './components/Header.vue';

import { stub } from './js/stub'

//-- form fields
const syllabus = reactive({
  title: "",
  description: "",
  email: ""
})

let cartridge = reactive({ name: "" });
let manifest = new Object() as ManifestType;
let items = new Array<ItemType>()
let resources = new Array<ResourceType>()

let showMetadata = ref(false);
let showUpload = ref(false);

let log = ref("ready");
let isUploaded = ref(false);
const HOST = import.meta.env.DEV ? "http://localhost:3046" : ""

let upload = function () {
  const formElem = document.getElementById("upload-form") as HTMLFormElement;
  const formData = new FormData(formElem);

  let cc = formData.get("cartridge") as File
  if (cc.name == "" || cc.size  == 0) {
    console.warn("can't submit an empty cartridge!");
    log.value = "can't submit an empty cartridge!";
    return;
  }

  cartridge.name = (formData.get("cartridge") as File).name;
  log.value = `uploading ${cartridge.name}`;

  fetch(`${HOST}/api/upload`, {
    method: "POST",
    body: formData,
  })
    .then((res) => {
      return res.json();
    })
    .then((data) => {
      log.value = `uploaded ${cartridge.name}`;
      isUploaded.value = true;

      Object.assign(manifest, JSON.parse(data.data))
      Object.assign(items, JSON.parse(data.items))

      //-- todo, here we have to get rid of the Item field of the returned struct... what to do?
      for (let r of JSON.parse(data.resources)) {
        resources.push(r.Resource);
      }

      console.log(manifest, items, resources)

      syllabus.title = manifest.Metadata.Lom.General.Title.String.Text
      syllabus.description = manifest.Metadata.Lom.General.Description.String.Text
    })
    .catch((err) => {
      log = ref(err);
      console.error(err);
    });
};

let reset = () => {
  isUploaded.value = false
  showMetadata.value = false

  Object.assign(manifest, {})
  Object.assign(items, {})
  Object.assign(resources, {})

  log.value = "reset cartridge"
}

let showModal = () => {
  showUpload.value = true;
  setTimeout(() => {
    let md = document.getElementById("modal-dialog") as HTMLElement
    md.focus()
  }, 100);
}

// onMounted(() => {
//   isUploaded.value = true;

//   Object.assign(manifest, JSON.parse(stub.data))
//   Object.assign(items, JSON.parse(stub.items))

//   for (let r of JSON.parse(stub.resources)) {
//     resources.push(r.Resource);
//   }
// });


</script>

<style lang="scss">
@import "./css/global-vars.scss";

// page typograph

.section-title-small {
  font-size: 0.6rem;
  text-transform: uppercase;
  font-weight: 300;
}

.cartridge-viewer {
  // border: 1px solid orange;
  height: 75vh;

  .items-panel {
    // border: 1px solid blue;
    height: 75vh;
  }

  .resources-panel {
    // border: 1px solid red;
    height: 75vh;
  }
}

.actions {
  display: flex;

  .action {
    width: max-content;
    border: none;
    background-color: white;
    color: grey;
    text-decoration: underline;
  }

  .action:hover {
    color: black;
  }
}

.modal {
  display: block;
}
</style>
