<template>
  <h1>upload a common cartridge</h1>
  <div class="upload">
    <form action="/api/upload" method="post" id="upload-form">
      <input type="file" name="cartridge" id="cartridge" />
    </form>
    <button type="button" @click="upload()">upload</button>
  </div>
  <div class="log">{{ log }}</div>
  <div v-if="manifest" class="cartridge">
    <div class="metadata">
      <h2>Metadata</h2>
      <div class="title">{{ manifest.Metadata.Lom.General.Title.String.Text }}</div>
    </div>
    <div class="items">
      <h2>Items</h2>
      <div class="item" v-for="i in items">
        <Item :item=i />
      </div>
    </div>
    <div class="resources">
      <h2>Resources</h2>
      <div class="resource" v-for="r in resources">
        <Resource :resource=r :cartridge="cartridge.name" />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import Resource from './components/Resource.vue'
import Item from './components/Item.vue'

// Check out https://vuejs.org/api/sfc-script-setup.html#script-setup for the newest way to use Vue3

export default {
  components: { Item, Resource },
  setup() {
    let cartridge
    let manifest
    let items
    let resources
    const log = ref("log ready")

    let upload = function () {
      let form = document.getElementById("upload-form")
      let formData = new FormData(form)

      if (formData.get("cartridge").name == "") {
        console.warn("can't submit an empty cartridge!");
        this.log = "can't submit an empty cartridge!"
        return
      }

      this.cartridge = formData.get("cartridge")
      this.log = `uploading ${this.cartridge.name}`

      fetch("/api/upload", {
        method: 'POST',
        body: formData
      }).then(res => {
        if (res.ok) {
          return res.json()
        } else {
          console.error(res.err)
          log = `internal server error on upload: ${res.err}`
        }
      }).then(data => {
        console.log(data);
        this.log = `loaded ${this.cartridge.name}`
        this.manifest = reactive(JSON.parse(data.data))
        this.items = reactive(JSON.parse(data.items))
        this.resources = reactive(JSON.parse(data.resources))
      }).catch(err => {
        this.log = err
        console.error(err);
      })
    }


    return { cartridge, manifest, items, resources, log, upload }

  }
}

</script>

<style>
</style>
