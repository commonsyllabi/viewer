<template>
    <Header></Header>
    <main>
        <div class="container p-3 " id="syllabus" v-if="syllabus.title != ''">
            <div class="syllabus">
                <h1>{{ syllabus.title }}</h1>
                <p>{{ syllabus.description }}</p>
                <h2>Attachments</h2>
                <div>
                    <ol>

                        <li v-for="att in syllabus.attachments">
                            <a :href="'/attachments/' + att.id + '?type=file'" target="_blank">{{ att.name }}</a>
                        </li>

                    </ol>
                </div>
            </div>


            <div class="request-link">
                <p>
                    <b>Is this your syllabus?</b>
                </p>
                <p>Enter the email address that you used in order to receive an edit link.</p>
                <form action="/api/magic-link" method="POST" id="request-email">
                    <input name="email" type="email" v-model="email" autocomplete="email" />
                    <input type="number" name="id" v-model="syllabus.id" hidden />
                    <button @click="submit()" type="button" class="btn cc-btn">Request</button>
                </form>
                <div class="log">{{ log }}</div>
            </div>
        </div>
        <div class="container p-3" v-else>
            <h2>Sorry!</h2>
            <p>We couldn't load data from the desired syllabus.</p>
            <p>Back to <a href="/">the listing</a>.</p>
        </div>
    </main>
    <Footer></Footer>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'

import Footer from './components/Footer.vue';
import Header from './components/Header.vue';

const HOST = import.meta.env.DEV ? "http://localhost:3046" : ""
const email = ref("")
const log = ref("")

const syllabus = reactive({
    id: "",
    title: "",
    description: "",
    attachments: new Array<{
        id: "",
        name: ""
    }>()
})

let submit = () => {
    const formElem = document.getElementById("request-email") as HTMLFormElement;
    const formData = new FormData(formElem);

    if (!syllabus || !syllabus.id)
        return

    log.value = `Sending email to: ${email.value}`
    fetch("/api/magic-link", {
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


onMounted(() => {
    let id = parseInt(window.location.pathname.split('/')[2])
    if (isNaN(id)) {
        console.error("URL parameter is not a number:", id);
        return
    }

    fetch(`${HOST}/syllabi/${id}`, {
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(res => {
            return res.json()
        })
        .then(data => {
            let d = JSON.parse(data)
            syllabus.id = d.ID
            syllabus.title = d.title
            syllabus.description = d.description
            for (let att of d.Attachments)
                syllabus.attachments.push({
                    id: att.ID,
                    name: att.Name
                })
        })
})
</script>

<style lang="scss">
@import "./css/global-vars.scss";

.syllabus {
    margin-bottom: 5rem;
}

.request-link {
    width: fit-content;

    form {
        display: flex;
    }
}

.log {
    font-size: 0.8em;
    margin-top: 5px;
}

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
</style>