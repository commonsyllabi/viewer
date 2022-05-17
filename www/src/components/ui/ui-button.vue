<template>
    <button :id="id" :class="'govuk-button ' + type + ' ' + classes" data-module="govuk-button"
        @click.prevent="$emit('click', $event)" data-prevent-double-click="true">
        {{ text }}
    </button>
</template>

<script setup>

defineProps({
    id: {
        type: String,
        default: ""
    },
    text: {
        type: String,
        default: "Click here"
    },
    type: {
        type: String,
        default: "default"
        },
    disabled: false,
    classes: {
        type: String,
        default: ""
        },
});

</script>

<style lang="scss" scoped>

button{
    margin: 5px;
}
.default {
    background-color: green;
    color: white;
}

.secondary {
    background-color: lightgrey;
    color: black;
    font-size: 0.9em;
}

.warning {
    background-color: red;
    color: white;
}

.disabled {
    opacity: 0.7;
}
</style>