<template>
  <Header></Header>
  <main class="container p-3">
    <!-- upload form -->
    <div class="container d-flex flex-column align-items-center">
      <div class="col-lg-6 mb-3 p-5 rounded upload-box" v-show="!isUploaded">
        <form id="upload-form" action="/api/upload" method="post">
          <div class="form-group d-flex flex-column align-items-center">
            <div class="upload-icon w-25 mb-4">
              <img src="./assets/upload.svg" alt="upload icon" class="w-100" />
            </div>
            <label for="cartridgeInput" class="d-block h5 mb-3"
              >Upload a common cartridge (.imscc) file</label
            >
            <input
              id="upload-file"
              type="file"
              name="cartridge"
              class="form-control-file d-block mb-2 visually-hidden"
            />
            <div id="upload-file-name" class="text-center mb-3"></div>
            <button
              id="upload-submit"
              @click.prevent="selectFile()"
              class="btn btn-lg btn-dark"
            >
              select a file
            </button>
          </div>
        </form>
      </div>

      <!-- status log -->
      <div id="log" class="container p-0 w-50 mb-5 text-center">
        {{ log }}
      </div>

      <div v-show="!isUploaded" class="w-50">
        <label class="mb-3">Or select an example file from the list:</label>
        <select
          class="form-select"
          aria-label="Select a demo file"
          id="examples"
          @change.prevent="loadExample($event)"
          @select.prevent="loadExample($event)"
        >
          <option value="" default>---</option>
          <option value="0">IMSCC - Test File</option>
          <option value="1">
            Liz Falconer - Computers, Canvas and Community
          </option>
          <option value="2">Openmed - Intro to Chemistry</option>
        </select>
      </div>
    </div>

    <!-- upload modal -->
    <div
      class="modal"
      @keydown.esc="showUpload = false"
      tabindex="-1"
      v-if="showUpload"
    >
      <Upload
        :title="manifest.Metadata.Lom.General.Title.String.Text"
        :description="manifest.Metadata.Lom.General.Description.String.Text"
        @close="showUpload = false"
      ></Upload>
    </div>

    <div v-if="isUploaded" class="container cartridge">
      <div>
        <h2 class="title fs-4">
          {{ syllabus.title }}
        </h2>
        <div class="file-name">
          {{ cartridge.name }}
        </div>
        <div class="description">
          {{ syllabus.description }}
        </div>
      </div>

      <!-- actions -->
      <div class="row mb-2 actions">
        <button
          type="button"
          class="action"
          @click="showMetadata = !showMetadata"
        >
          {{ showMetadata ? 'hide metadata & index' : 'show metadata & index' }}
        </button>
        <button id="reset-upload" type="button" class="action" @click="reset()">
          load another file
        </button>
        <button
          id="show-upload"
          type="button"
          class="action"
          @click="showModal()"
          v-if="!isExample"
        >
          upload this file
        </button>
      </div>

      <div class="row mb-5 cartridge-viewer">
        <!-- metadata and items -->
        <div class="col-4 h-100 mb-1" v-if="showMetadata">
          <!-- metadata viewer -->
          <div class="metadata-panel h-50 d-flex flex-column">
            <h3 class="my-2 fs-6">Cartridge metadata</h3>
            <!-- file metadata placeholder -->
            <div
              class="metadata-container border rounded p-2 overflow-scroll flex-grow-1"
            >
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
                <div>
                  {{ manifest.Metadata.Lom.LifeCycle.Contribute.Entity.String }}
                  - {{ manifest.Metadata.Lom.LifeCycle.Contribute.Role.String }}
                </div>
              </div>

              <div class="metadata-element">
                <div class="legend">copyright</div>
                <div>{{ manifest.Metadata.Lom.Rights.Description.String }}</div>
              </div>

              <div class="metadata-element">
                <div class="legend">language</div>
                <div>{{ manifest.Metadata.Lom.General.Language }}</div>
              </div>
            </div>
          </div>
          <!-- items panel -->
          <div class="items-panel h-50 d-flex flex-column">
            <h3 class="my-2 fs-6">Items in Index</h3>

            <!-- items listing -->
            <div
              class="metadata-container border rounded p-2 overflow-scroll flex-grow-1"
            >
              <div v-for="i in items" :key="i.Item.Identifier">
                <Item :item="i" :cartridge="cartridge.name" />
                <hr />
              </div>
            </div>
          </div>
        </div>

        <!-- resources panel -->
        <div class="col-8 resources-panel h-100 d-flex flex-column">
          <h3 class="fs-6 my-2">Resources in Cartridge</h3>
          <div
            class="resources-container border rounded overflow-scroll flex-grow-1"
          >
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
  import { ref, reactive, onMounted } from 'vue'
  import { ManifestType, ItemType, ResourceType } from './js/types'

  import Resource from './components/Resource.vue'
  import Item from './components/Item.vue'
  import Upload from './components/Upload.vue'

  import Footer from './components/Footer.vue'
  import Header from './components/Header.vue'

  import { stub } from './js/stub'

  //-- form fields
  const syllabus = reactive({
    title: '',
    description: '',
    email: '',
  })

  let cartridge = reactive({ name: '' })
  let manifest = reactive(new Object() as ManifestType)
  let items = reactive(new Array<ItemType>())
  let resources = reactive(new Array<ResourceType>())

  let showMetadata = ref(true)
  let showUpload = ref(false)

  let log = ref('')
  let isUploaded = ref(false)
  let isExample = ref(false)
  const HOST = import.meta.env.DEV ? 'http://localhost:3046' : ''
  const DEBUG = import.meta.env.DEV ? true : false

  let selectFile = () => {
    document.getElementById('upload-file')?.click()
  }

  let upload = function () {
    const formElem = document.getElementById('upload-form') as HTMLFormElement
    const formData = new FormData(formElem)

    let cc = formData.get('cartridge') as File
    if (cc.name == '' || cc.size == 0) {
      if (DEBUG) console.warn("can't submit an empty cartridge!")
      log.value = "can't submit an empty cartridge!"
      return
    }

    cartridge.name = (formData.get('cartridge') as File).name
    log.value = `Uploading ${cartridge.name}...`

    const uploadBtn = document.getElementById('upload-file') as HTMLInputElement
    uploadBtn.disabled = true

    fetch(`${HOST}/api/upload`, {
      method: 'POST',
      body: formData,
    })
      .then((res) => {
        return res.json()
      })
      .then((payload) => {
        log.value = `uploaded ${cartridge.name}`
        isUploaded.value = true

        Object.assign(manifest, payload.data)
        Object.assign(items, payload.items)

        //-- todo, here we have to get rid of the Item field of the returned struct... what to do?
        for (let r of payload.resources) {
          resources.push(r.Resource)
        }

        syllabus.title = manifest.Metadata.Lom.General.Title.String.Text
        syllabus.description =
          manifest.Metadata.Lom.General.Description.String.Text
      })
      .catch((err) => {
        log = ref(err)
        console.error(err)
      })
      .finally(() => {
        uploadBtn.disabled = false
      })
  }

  let reset = () => {
    isUploaded.value = false
    showMetadata.value = false

    isExample.value = false
    let examples = document.getElementById('examples') as HTMLSelectElement
    examples.selectedIndex = 0

    items.length = 0
    resources.length = 0

    log.value = 'reset cartridge'
  }

  let showModal = () => {
    showUpload.value = true
    setTimeout(() => {
      let md = document.getElementById('modal-dialog') as HTMLElement
      md.focus()
    }, 100)
  }

  let loadExample = (_evt: Event) => {
    let v = parseInt((<HTMLInputElement>_evt.target).value)

    if (isNaN(v) || v >= stub.length) {
      console.log(
        `attempting to get #${v} from examples, ${stub.length} available.`
      )

      return
    }

    isUploaded.value = true
    isExample.value = true

    Object.assign(manifest, stub[v].data)
    Object.assign(items, stub[v].items)

    syllabus.title = manifest.Metadata.Lom.General.Title.String.Text
    syllabus.description = manifest.Metadata.Lom.General.Description.String.Text
    cartridge.name = stub[v].name

    for (let r of stub[v].resources) {
      resources.push(r.Resource as unknown as ResourceType)
    }
  }

  onMounted(() => {
    let uploadForm = document.getElementById('upload-form') as HTMLElement
    uploadForm.ondrop = (e: Event) => {
      e.preventDefault()
      console.log('dropped', e)
    }
    let uploadField = document.getElementById('upload-file') as HTMLElement

    uploadField.onchange = (e: Event) => {
      let t = e.target as HTMLInputElement
      let f = t.files as FileList
      let uploadedFile = document.getElementById(
        'upload-file-name'
      ) as HTMLElement
      uploadedFile.innerText = f[0].name

      upload()
    }
  })
</script>

<style lang="scss">
  @import './css/global-vars.scss';

  .upload-box {
    background-color: $light;
    border: 1px solid $secondary;
  }

  .file-name {
    color: $secondary;
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
