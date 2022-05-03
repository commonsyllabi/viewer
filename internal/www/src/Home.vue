<template>
    <h1>Home</h1>
    <div>{{ msg }}</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

const msg = ref("Hello!")
const HOST = import.meta.env.DEV ? "http://localhost:3046": ""

onMounted(() => {
    
    fetch(`${HOST}/syllabi/`,
    {
        method: 'GET'
    })
    .then(res => {
        return res.json()
    })
    .then(data => {
        console.log(data);
        if (data.length > 0) 
            msg.value = "No syllabi :("
        else
            msg.value = `There are ${data.length} syllabi.`
        
    })
    .catch(err => {
        console.error(err)
    })
})
</script>

<style>

</style>