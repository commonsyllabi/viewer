<template>
  <div class="container p-3">
    <div class="container pt-4 mb-3 border rounded">
      <form id="submit-form" action="/syllabi/" method="POST">
        <div>
          <label for="title">Title of the course</label>
          <input type="text" name="title" id="title" />
        </div>

        <div>
          <label for="description">Description of the course</label>
          <input type="text" name="description" id="description" />
        </div>

        <button
          id="course-submit"
          type="button"
          class="btn btn-primary mb-4"
          @click="submit()"
        >submit</button>
      </form>
    </div>
    <!-- upload form -->
    <div class="container pt-4 mb-3 border rounded">
      <form id="upload-form" action="/api/upload" method="post">
        <div class="form-group">
          <label
            for="cartridgeInput"
            class="d-block h5 mb-3"
          >Upload a common cartridge (.imscc) file</label>
          <input
            id="upload-file"
            type="file"
            name="cartridge"
            class="form-control-file d-block mb-2"
          />
          <button
            id="upload-submit"
            type="button"
            class="btn btn-primary mb-3"
            @click="upload()"
          >upload</button>
        </div>
      </form>
    </div>

    <!-- status log -->
    <div id="log" class="container py-4 mb-3 border rounded text-light bg-dark">
      <pre>{{ log }}</pre>
    </div>

    <!-- metadata viewer -->
    <div class="row mb-3">
      <div id="metadata-accordion" class="accordion">
        <div class="accordion-item">
          <h2 id="headingOne" class="accordion-header">
            <button
              class="accordion-button"
              type="button"
              data-bs-toggle="collapse"
              data-bs-target="#collapseMeta"
              aria-expanded="true"
              aria-controls="collapseOne"
            >Cartridge Metadata</button>
          </h2>
          <div
            id="collapseMeta"
            class="accordion-collapse collapse show"
            aria-labelledby="headingOne"
            data-bs-parent="#metadata-accordion"
          >
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
    <!-- <div v-if="isUploaded" class="container cartridge"> -->
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
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ManifestType, ItemType, ResourceType } from "./js/types";

import Resource from './components/Resource.vue'
import Item from './components/Item.vue'

import { stub } from './js/stub'

let cartridge = reactive({ name: "" });
let manifest = new Object() as ManifestType;
let items = new Array<ItemType>()
let resources = new Array<ResourceType>()

let log = ref("ready");
let isUploaded = ref(false);
const HOST = import.meta.env.DEV ? "http://localhost:3046" : ""

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

  fetch(`${HOST}/api/upload`, {
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

      Object.assign(manifest, JSON.parse(data.data))
      Object.assign(items, JSON.parse(data.items))

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

let submit = () => {
  const formElem = document.getElementById("submit-form") as HTMLFormElement;
  const formData = new FormData(formElem);

  if (validateSubmission(formData)) {
    console.warn("can't submit an empty title or description!");
    log.value = "can't submit an empty title or description!";
    return;
  }

  fetch(`${HOST}/syllabi/`, {
    method: "POST",
    body: formData,
  })
    .then(res => {
      res.json()
    })
    .then(data => {
      console.log(data)
      log.value = "submitted syllabus!"
    })
}

let validateSubmission = (_data: FormData) => {
  if (_data == null)
    return false
  else if (_data.get("title") != undefined || _data.get("description") != undefined)
    return false
  else if (_data.get("title") != null || _data.get("description") != null)
    return false

  const title = _data.get("title") as string
  const description = _data.get("description") as string
  if (title.length < 5 || description.length < 5)
    return false

  return true
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
</style>
