<template>
    <n-space>
        <div>
            <div v-if="!Editing">
                <n-button circle size="tiny" @click="onEdit"><n-icon><EditIcon/></n-icon></n-button>
            </div>
            <div v-else>
                <n-space size="small">
                    <n-button circle size="tiny" @click="onSave"><n-icon><SaveIcon/></n-icon></n-button>
                    <n-button circle size="tiny" @click="onCancel"><n-icon><CloseIcon/></n-icon></n-button>
                </n-space>
            </div>
        </div>
        <div>
            <div v-if="!Editing" style="padding-top: 2px;">{{ props.comment }}</div>
            <div v-else><n-input type="text" size="tiny" v-model:value="EditText" placeholder="发生什么事了？发生什么事了？发生什么事了？"/></div>
        </div>
    </n-space>
</template>

<script setup>
import { onMounted, onUpdated, ref } from 'vue';
import EditIcon from '../components/icons/EditIcon.vue'
import SaveIcon from '../components/icons/SaveIcon.vue'
import CloseIcon from '../components/icons/CloseIcon.vue'

const props = defineProps(['comment', 'time'])
const emit = defineEmits(['update'])

const Editing = ref(false)
const EditText = ref("")

onMounted(() => {
    EditText.value = props.comment
})
onUpdated(() => {
    EditText.value = props.comment
})

const onEdit = () => {
    Editing.value = true
}

const onCancel = () => {
    Editing.value = false
    EditText.value = props.comment
}

const onSave = () => {
    Editing.value = false
    emit("update", props.time, EditText.value)
    EditText.value = props.comment
}
</script>

<style>
</style>