<template>
    <n-layout>
        <n-layout-header bordered>
            <n-breadcrumb separator=">">
                <n-breadcrumb-item>首页</n-breadcrumb-item>
                <n-breadcrumb-item>登录</n-breadcrumb-item>
            </n-breadcrumb>
        </n-layout-header>
        <n-layout-content>
            <div class="login-card">
                <n-card title="登录">
                    <n-space vertical>
                        <p>用户名:</p>
                        <n-input type="text" v-model:value="UserName" placeholder="你是谁？"/>
                        <p>密码:</p>
                        <n-input type="password" v-model:value="UserPass" placeholder="你的密码是？" />
                    </n-space>
                    <template #action>
                        <div class="login-action-bar">
                            <n-button style="float: right" type="primary" :loading="Loading" :disabled="Disabled" @click="Login">
                                登录
                            </n-button>
                        </div>
                    </template>
                </n-card>
            </div>
        </n-layout-content>
    </n-layout>
</template>

<script setup>
import axios from 'axios';
import { useMessage } from 'naive-ui';
import { computed, ref } from 'vue';
import router from "@/router"

const message = useMessage()

const UserName = ref("")
const UserPass = ref("")
const Loading = ref(false)
const Disabled = computed(() => {
    return Loading.value || UserName.value == "" || UserPass.value == ""
})

const Login = () => {
    Loading.value = true
    axios.post("/api/user/login", {
        user: UserName.value,
        pass: UserPass.value,
    }).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        router.push("/")
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {
        Loading.value = false
    })
}
</script>

<style>
.login-card {
    margin: 240px auto 0 auto;
    max-width: 480px;
}
.login-action-bar:after {
    content: "";
    display: block;
    clear: both;
}
</style>