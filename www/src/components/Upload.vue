<template>
  <div class="modal-dialog modal-dialog-centered small modal-lg rounded" id="modal-dialog">
    <div class="modal-content p-5">
      <h5>Upload this cartridge to the public repo</h5>
      <p>This cartridge will be viewable to others if you upload it to the public repo. Please only upload course
        materials that you have authored or have the rights to.</p>
      <p class="important">
        Important
      </p>
      <p>You may delete your shared cartridge at any time <b>via the link that will be sent to your email</b> after the
        upload.</p>
      <p>
        <b>Please verify that the email address that you enter here is correct, and that you have access to it.</b>
      </p>
      <form id="submit-form" action="/syllabi/" method="POST" v-show="!isSubmitted">
        <div>
          <label class="w-100 visually-hidden" for="title">Title of the course</label>
          <input class="text-input mb-2 w-50 visually-hidden" type="text" name="title" id="title" v-model="title" />
        </div>

        <div>
          <label class="w-100 visually-hidden" for="description">Description of the course</label>
          <input class="text-input mb-2 w-50 visually-hidden" type="text" name="description" id="description"
            v-model="description" />
        </div>

        <div>
          <label class="w-100" for="email">Email</label>
          <input class="text-input mb-2 w-50" type="email" name="email" id="email" v-model="email" />
        </div>

        <div>
          <label class="w-100" for="email">Confirm email</label>
          <input class="text-input mb-2 w-50" type="email" name="email-conf" id="email-conf" />
        </div>

        <p>If you have any questions, you can email us at <a
            href="mailto:admin@commonsyllabi.org">admin@commonsyllabi.org</a>.</p>

        <div class="buttons">
          <button id="upload-close" text="close" class="close-btn" @click="$emit('close')">close</button>
          <button id="course-submit" text="submit" class="btn btn-primary mb-4 cc-btn"
            @click.prevent="submit()">upload</button>
        </div>
      </form>
      <div class="log" id="submit-log">{{ log }}</div>
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

const isSubmitted = ref(false)
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
      isSubmitted.value = true
      log.value = "Thank you! Your cartridge was submitted successfully."
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
a {
  color: black;
}

.important {
  color: red;
  font-weight: bold;
}

.text-input {
  border: 1px solid grey;
}

.buttons {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cc-btn {
  border-radius: 25px;
  border-color: black;
  background-color: white;
  color: black;
}

.cc-btn:hover {
  background-color: white;
  color: black;
  font-weight: bold;
}

.close-btn {
  border: none;
  background-color: white;
  color: black;
  text-decoration: underline;
}
</style>