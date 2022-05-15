<template>
    <div class="container p-3">
        <h1>Home</h1>
        <div>{{ msg }}</div>
        <div class="syllabi">
            <ul>
                <li v-for="syllabus in syllabi">
                    <div>
                        <a :href="'/syllabi/' + syllabus.ID">{{ syllabus.title }}</a>
                    </div>
                    <div>{{ syllabus.description }}</div>
                </li>
            </ul>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import { SyllabusType } from './js/types';

const msg = ref("")
const HOST = import.meta.env.DEV ? "http://localhost:3046" : ""
const syllabi = new Array<SyllabusType>()

onMounted(() => {
    fetch(`${HOST}/syllabi/`,
        {
            method: 'GET'
        })
        .then(res => {
            return res.json()
        })
        .then(data => {
            Object.assign(syllabi, JSON.parse(data))
            console.log(syllabi);
            if (syllabi.length == 0)
                msg.value = "No syllabi :("
            else
                msg.value = `There are ${syllabi.length} syllabi.`

        })
        .catch(err => {
            console.error(err)
            msg.value = "Network error :|"
        })
})
</script>

<style lang="scss">
@import "./css/global-vars.scss";
</style>