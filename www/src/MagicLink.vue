<template>
  <Header></Header>
  <main>
    <div class="container p-3">
      <h1 class="pb-3">Deleting syllabus</h1>
      <p>
        Click on the button below to delete the syllabus
        <b>{{ id }} - {{ title }}</b>
      </p>
      <button
        id="btn"
        class="btn cc-btn mb-3"
        @click="requestDelete()"
        :disabled="isDeleting"
      >
        Delete
      </button>
      <div id="log" v-html="log"></div>
    </div>
  </main>
  <Footer></Footer>
</template>

<script setup>
  import { onMounted, ref } from 'vue'
  import Footer from './components/Footer.vue'
  import Header from './components/Header.vue'

  const id = ref('')
  const title = ref('')
  const log = ref('')
  const isDeleting = ref(false)

  let requestDelete = () => {
    isDeleting.value = true

    fetch(`/syllabi/${id.value}`, {
      method: 'DELETE',
    })
      .then((res) => {
        if (res.ok) {
          log.value =
            'Successfully deleted syllabus. Redirecting you in a couple of seconds...'
          setTimeout(() => {
            window.location = '/'
          }, 2000)
        } else {
          log.value =
            'Error deleting syllabus. Write us at <a href="mailto:contact@commonsyllabi.org">contact@commonsyllabi.org</a>'
        }
      })
      .catch((err) => {
        log.value =
          'Error deleting syllabus. Write us at <a href="mailto:contact@commonsyllabi.org">contact@commonsyllabi.org</a>'
        console.error(err)
      })
  }

  onMounted(() => {
    let app = document.getElementById('app')
    id.value = app.dataset.id
    title.value = app.dataset.title
  })
</script>

<style lang="scss">
  @import './css/global-vars.scss';

  .cc-btn,
  .cc-btn:hover {
    border-radius: 25px;
    background-color: white;
    color: black;
    border-color: black;
  }

  .cc-btn:hover {
    text-decoration: underline;
  }

  .cc-btn:disabled {
    opacity: 0.6;
  }
</style>
