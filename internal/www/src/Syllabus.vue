<template>
    <div>
        <p>
            <b>Is this your syllabus?</b>
        </p>
        <p>Enter the email address that you used in order to receive an edit link.</p>
        <form action="/api/request-email" method="POST" id="request-email">
            <input type="email" v-model="email" />
            <button @click="submit()" type="button">request</button>
        </form>
        <div class="log">{{ log }}</div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const email = ref("")
const log = ref("")

let submit = () => {
    const formElem = document.getElementById("request-email") as HTMLFormElement;
    const formData = new FormData(formElem);

    console.log("requesting for email:", email.value);
    log.value = `sending email to: ${email.value}`

}
</script>

<style lang="scss">
@import "./css/global-vars.scss";

.log {
    font-size: 0.8em;
    margin-top: 5px;
}
</style>