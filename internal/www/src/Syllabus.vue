<template>
    <div>
        <p>
            <b>Is this your syllabus?</b>
        </p>
        <p>Enter the email address that you used in order to receive an edit link.</p>
        <form action="/api/request-email" method="POST" id="request-email">
            <input name="email" type="email" v-model="email" autocomplete="email" />
            <input type="number" name="id" v-model="id" hidden />
            <button @click="submit()" type="button">request</button>
        </form>
        <div class="log">{{ log }}</div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const email = ref("")
const log = ref("")
const syll = document.getElementById("syllabus") as HTMLElement
const id = syll.getAttribute("syll_id") as string

let submit = () => {
    const formElem = document.getElementById("request-email") as HTMLFormElement;
    const formData = new FormData(formElem);

    if (!syll || !id)
        return

    log.value = `Sending email to: ${email.value}`
    fetch("/api/request-email", {
        body: formData,
        method: "POST"
    })
        .then(res => {
            if (res.ok)
                log.value = `Sent email to: ${email.value}. Check your inbox!`
        })
        .catch(err => {
            log.value = `There was an error verifying your email.`
            console.error(err)
        })

}
</script>

<style lang="scss">
@import "./css/global-vars.scss";

.log {
    font-size: 0.8em;
    margin-top: 5px;
}
</style>