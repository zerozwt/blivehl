<template>
    <n-space>
        <div style="padding-top: 2px;">{{ UserName ? UserName : "-" }}</div>
        <div>
            <n-button quaternary size="small" @click="ChangePassword">
                <template #icon><n-icon><PasswordIcon/></n-icon></template>
                修改密码
            </n-button>
        </div>
        <div>
            <n-button quaternary size="small" @click="Logout">
                <template #icon><n-icon><ShutdownIcon/></n-icon></template>
                登出
            </n-button>
        </div>
    </n-space>
    <n-modal v-model:show="ShowChangePass">
        <n-card title="修改密码" closable @close="event => {ShowChangePass = false}" class="changepw-modal">
            <n-space vertical>
                <p>旧密码:</p>
                <n-input type="password" v-model:value="OldPass" placeholder="你现在的密码是啥？"/>
                <p>新密码:</p>
                <n-input type="password" v-model:value="NewPass" placeholder="新密码在这里输入一遍"/>
                <p>重复新密码:</p>
                <n-input type="password" v-model:value="NewPass2" placeholder="防止打错所以需要再输入一遍"/>
            </n-space>
            <template #action>
                <div class="change-pass-action">
                    <n-button type="primary" :disabled="ChangeButtonDisabled" :loading="ChangeButtonLoading" style="float: right" @click="Change">修改密码</n-button>
                </div>
            </template>
        </n-card>
    </n-modal>
</template>

<script setup>
import ShutdownIcon from './icons/ShutdownIcon.vue'
import PasswordIcon from './icons/PasswordIcon.vue'
import { computed, onMounted, ref } from 'vue';
import router from "@/router"
import createAPICaller from '@/utils'
import { useMessage } from 'naive-ui';

const API = createAPICaller(router)
const message = useMessage()
const emit = defineEmits(['login'])

const UserName = ref("")
const ShowChangePass = ref(false)
const OldPass = ref("")
const NewPass = ref("")
const NewPass2 = ref("")

const ChangeButtonLoading = ref(false)
const ChangeButtonDisabled = computed(() => {
    return ChangeButtonLoading.value || OldPass.value == "" || NewPass.value == '' || NewPass.value != NewPass2.value
})

const ChangePassword = () => {
    OldPass.value = ""
    NewPass.value = ""
    NewPass2.value = ""
    ChangeButtonLoading.value = false
    ShowChangePass.value = true
}

const Logout = () => {
    API.get("/api/user/logout").then(rsp => {
        router.push("/login")
    }).catch(err => message.error(JSON.stringify(err)))
}

const Change = () => {
    ChangeButtonLoading.value = true
    API.post("/api/user/pass", {
        old_pass: OldPass.value,
        new_pass: NewPass.value,
    }).then(rsp => {
        let data = rsp.data;
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        message.success("密码修改成功")
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {
        ChangeButtonLoading.value = false
        ShowChangePass.value = false
    })
}

onMounted(() => {
    API.get("/api/user/info").then(rsp => {
        let data = rsp.data;
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        UserName.value = data.data.name
        emit("login")
    }).catch(err => message.error(JSON.stringify(err)))
})
</script>

<style>
.changepw-modal {
    max-width: 400px;
}
.change-pass-action:after {
    content: "";
    display: block;
    clear: both;
}
</style>