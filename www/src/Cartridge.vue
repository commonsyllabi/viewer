<template>
  <Header></Header>
  <main class="container p-3">

    <!-- upload form -->
    <div class="container pt-4 mb-3 w-50 rounded upload-box" v-show="!isUploaded">
      <form id="upload-form" action="/api/upload" method="post">
        <div class="form-group">
          <div class="upload-icon w-25">
            <img src="./assets/upload.svg" alt="upload icon" class="w-100">
          </div>
          <label for="cartridgeInput" class="d-block h5 mb-3">Upload a common cartridge (.imscc) file</label>
          <input id="upload-file" type="file" name="cartridge" class="form-control-file d-block mb-2 visually-hidden" />
          <div id="upload-file-name"></div>
          <button id="upload-submit" @click.prevent="selectFile()" class="btn btn-primary mb-3 cc-btn">select a
            file</button>
        </div>
      </form>
    </div>
    
      <div class="examples-dropdown m-auto" v-show="!isUploaded">
        Or select an example from the list:
        <select name="examples" id="examples" @change.prevent="loadExample($event)" @select.prevent="loadExample($event)">
            <option value="" default>---</option>
            <option value="0">Test File</option>
        </select>
      </div>

    <!-- upload modal -->
    <div class="modal" @keydown.esc="showUpload = false" tabindex="-1" v-if="showUpload">
      <Upload :title="manifest.Metadata.Lom.General.Title.String.Text"
        :description="manifest.Metadata.Lom.General.Description.String.Text" @close="showUpload = false"></Upload>
    </div>

    <!-- status log -->
    <!-- <div id="log" class="container py-1 mb-3">
      <pre>{{ log }}</pre>
    </div> -->

    <div v-if="isUploaded" class="container cartridge">
      <div class="metadata">
        <div class="title">
          {{ syllabus.title }}
        </div>
        <div class="file-name">
          {{ cartridge.name }}
        </div>
        <div class="description">
          {{ syllabus.description }}
        </div>
      </div>

      <!-- actions -->
      <div class="row mb-2 actions">
        <button type="button" class="action" @click="showMetadata = !showMetadata">{{
            showMetadata ? "hide metadata & index" : "show metadata & index"
        }}</button>
        <button id="reset-upload" type="button" class="action" @click="reset()">load another file</button>
        <button id="show-upload" type="button" class="action" @click="showModal()" v-if="!isExample">upload this file</button>
      </div>

      <div class="row mb-5 cartridge-viewer">

        <!-- metadata and items -->
        <div class="col-4 h-100 mb-1" v-if="showMetadata">

          <!-- metadata viewer -->
          <div class="metadata-panel h-50 d-flex flex-column">
            <h6 class="my-2">Cartridge metadata</h6>
            <!-- file metadata placeholder -->
            <div class="metadata-container border rounded p-2 overflow-scroll flex-grow-1">
              <div class="metadata-element">
                <div class="legend">file name</div>
                <div>{{ cartridge.name }}</div>
              </div>

              <div class="metadata-element">
                <div class="legend">schema</div>
                <div>{{ manifest.Metadata.Schema }}</div>
              </div>

              <div class="metadata-element">
                <div class="legend">schema version</div>
                <div>{{ manifest.Metadata.Schemaversion }}</div>
              </div>

              <div class="metadata-element">
                <div class="legend">role</div>
                <div>{{ manifest.Metadata.Lom.LifeCycle.Contribute.Entity }} - {{ manifest.Metadata.Lom.LifeCycle.Contribute.Role }}</div>
              </div>

              <div class="metadata-element">
                <div class="legend">copyright</div>
                <div>{{ manifest.Metadata.Lom.Rights.Description }}</div>
              </div>

              <div class="metadata-element">
                <div class="legend">language</div>
                <div>{{ manifest.Metadata.Lom.General.Language }}</div>
              </div>
            </div>
          </div>
          <!-- items panel -->
          <div class="items-panel h-50 d-flex flex-column">
            <h6 class="my-2">Items in Index</h6>

            <!-- items listing -->
            <div class="metadata-container border rounded p-2 overflow-scroll flex-grow-1">
              <div v-for="i in items" :key="i.Item.Identifier">
                <Item :item="i" :cartridge="cartridge.name" />
                <hr />
              </div>
            </div>

          </div>
        </div>

        <!-- resources panel -->
        <div class="col resources-panel h-100 d-flex flex-column">
          <h6 class="my-2">Resources in Cartridge</h6>
          <div class="resources-container border rounded overflow-scroll flex-grow-1">
            <!-- resources listing -->
            <div v-for="r in resources" :key="r.Identifier" class>
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
let manifest = reactive(new Object() as ManifestType)
let items = reactive(new Array<ItemType>())
let resources = reactive(new Array<ResourceType>())

let showMetadata = ref(true);
let showUpload = ref(false);

let log = ref("ready");
let isUploaded = ref(false);
let isExample = ref(false);
const HOST = import.meta.env.DEV ? "http://localhost:3046" : ""

let selectFile = () => {
  document.getElementById("upload-file")?.click()
}

let upload = function () {
  const formElem = document.getElementById("upload-form") as HTMLFormElement;
  const formData = new FormData(formElem);

  let cc = formData.get("cartridge") as File
  if (cc.name == "" || cc.size == 0) {
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

  isExample.value = false
  let examples = document.getElementById("examples") as HTMLSelectElement
  examples.selectedIndex = 0

  items.length = 0
  resources.length = 0

  log.value = "reset cartridge"
}

let showModal = () => {
  showUpload.value = true;
  setTimeout(() => {
    let md = document.getElementById("modal-dialog") as HTMLElement
    md.focus()
  }, 100);
}

let loadExample = (_evt : Event) => {
  let v = parseInt((<HTMLInputElement>_evt.target).value)
  console.log(v);
  
  if(isNaN(v) || v >= stub.length){
    console.log(`attempting to get #${v} from examples, ${stub.length} available.`);
    
    return
  }
  
  isUploaded.value = true;
  isExample.value = true

  Object.assign(manifest, JSON.parse(stub[v].data))
  Object.assign(items, JSON.parse(stub[v].items))

  syllabus.title = manifest.Metadata.Lom.General.Title.String.Text
  syllabus.description = manifest.Metadata.Lom.General.Description.String.Text
  cartridge.name = "example cartridge"

  for (let r of JSON.parse(stub[v].resources)) {
    resources.push(r.Resource);
  }
}

onMounted(() => {
  let uploadForm = document.getElementById("upload-form") as HTMLElement
  uploadForm.ondrop = (e: Event) => {
    e.preventDefault()
    console.log("dropped", e);

  }
  let uploadField = document.getElementById("upload-file") as HTMLElement

  uploadField.onchange = (e: Event) => {
    let t = e.target as HTMLInputElement
    let f = t.files as FileList
    let uploadedFile = document.getElementById("upload-file-name") as HTMLElement
    uploadedFile.innerText = f[0].name

    upload()
  }
});


</script>

<style lang="scss">
@import "./css/global-vars.scss";

// page typograph

.upload-box {
  display: flex;
  flex-direction: column;

  background-color: #eaeaea;
  border: 3px dashed grey;
}

.form-group {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.examples-dropdown{
  width: fit-content;

  select{
    background-color: white;
    border: 1px solid black;
    border-radius: 25px;
    padding: 5px;
  }
}

.cc-btn, .cc-btn:hover {
  border-radius: 25px;
  background-color: white;
  color: black;
  border-color: black;
}

.cc-btn:hover {
  text-decoration: underline;
}

.file-name {
  color: lightgray;
}

.metadata-element {
  padding-bottom: 10px;
  font-size: 0.9em;

  .legend {
    color: grey;
  }
}

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
