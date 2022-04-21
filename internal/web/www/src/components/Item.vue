<script setup>

const props = defineProps({
    item: {}
})

</script>

<template>
        <h3>{{ props.item.Item.Identifier }} - {{ props.item.Item.Title }}</h3>

        <div v-if="props.item.children">
          <h4>Children</h4>
          <div class="sub-item" v-for="child in props.item.Children">
            <div>{{ child.Identifier }} - {{ child.Title }}</div>
          </div>
        </div>

        <h4>Resources</h4>
        <ul class="sub-resource" v-for="res in props.item.Resources">
          <li>{{ res.XMLName.Local }} - {{ res.Identifier }}</li>
        </ul>

        <hr>
</template>