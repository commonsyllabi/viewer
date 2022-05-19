<template>
  <div class="modal-dialog modal-dialog-centered" id="modal-dialog">
    <div class="modal-content p-3">
      <form id="submit-form" action="/syllabi/" method="POST">
        <div>
          <label for="title">Title of the course</label>
          <input type="text" name="title" id="title" v-model="title" />
        </div>

        <div>
          <label for="description">Description of the course</label>
          <input type="text" name="description" id="description" v-model="description" />
        </div>

        <div>
          <label for="email">Email</label>
          <input type="email" name="email" id="email" v-model="email" />
        </div>

        <div>
          <label for="email">Confirm email</label>
          <input type="email" name="email-conf" id="email-conf" />
        </div>


        <button id="course-submit" class="btn btn-primary mb-4" @click="submit()">submit</button>
        <button id="button-submit" class="btn btn-primary mb-4" @click="$emit('close')">close</button>
        <div class="log">{{ log }}</div>
      </form>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

//props
defineProps<{
  title: String,
  description: String
}>();

const email = ref("")

const HOST = import.meta.env.DEV ? "http://localhost:3046" : ""
const log = ref("")

let submit = () => {

  if (isInvalidEmail()) {
    log.value = "please make sure that the emails are matching!"
    return
  }

  const pformElem = document.getElementById("upload-form") as HTMLFormElement;
  const pformData = new FormData(pformElem);

  const formElem = document.getElementById("submit-form") as HTMLFormElement;
  const formData = new FormData(formElem);


  formData.set('attachments[]', pformData.get('cartridge') as FormDataEntryValue)
  formData.forEach((v, k) => {
    console.log(k, v)
  })

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

let isInvalidEmail = (): boolean => {
  const e1 = document.getElementById("email") as HTMLInputElement
  const e2 = document.getElementById("email-conf") as HTMLInputElement
  if (!e1 || !e2)
    return true
  if (e1.value != e2.value)
    return true

  return false
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
</script>

<style lang="scss" scoped>
</style>