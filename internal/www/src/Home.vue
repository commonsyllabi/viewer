<template>
    <h1>Home</h1>
    <div>{{ msg }}</div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';

const msg = "Hello!"
onMounted(() => {
    fetch('http://localhost:3046/syllabi/',
    {
        method: 'GET'
    })
    .then(res => {
        return res.json()
    })
    .then(data => {
        console.log(data);
        
    })
    .catch(err => {
        console.error(err)
    })
})
</script>

<style>

</style>