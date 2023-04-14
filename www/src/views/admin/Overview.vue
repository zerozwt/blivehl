<template>
<n-layout>
    <n-layout-header bordered>
        <div class="clear-float">
            <div style="float:left; padding-top: 1px;">
                <n-breadcrumb separator=">">
                    <n-breadcrumb-item><router-link to="/">首页</router-link></n-breadcrumb-item>
                    <n-breadcrumb-item>数据管理</n-breadcrumb-item>
                </n-breadcrumb>
            </div>
            <div style="float: right; padding-bottom: 4px;"><LoginStatus @login="AfterLogin"/></div>
        </div>
    </n-layout-header>
    <n-layout-content>
        <div class="clear-float overview-title">
            <div style="float:left"><p>主播列表</p></div>
            <div style="float:right">
                <n-input type="text" v-model:value="RoomID" :allow-input="(value) => !value || /^[0-9]+$/.test(value)" placeholder="直播间号码">
                    <template #suffix>
                        <n-button size="tiny" @click="OnBtnGo" :disabled="RoomID == ''">GO!</n-button>
                    </template>
                </n-input>
            </div>
        </div>
        <n-grid x-gap="8" y-gap="8" :cols="8">
            <n-gi v-for="item in Rooms" :key="item.room_id">
                <router-link :to="'/admin/room/'+item.room_id">
                    <n-card :title="item.name" hoverable>
                        <template #cover>
                            <n-image width="256" :src="item.icon" preview-disabled/>
                        </template>
                    </n-card>
                </router-link>
            </n-gi>
        </n-grid>
    </n-layout-content>
</n-layout>
</template>

<script setup>
import {RouterLink} from 'vue-router'
import LoginStatus from '@/components/LoginStatus.vue'
import { reactive, ref } from 'vue';
import router from "@/router"
import createAPICaller from '@/utils'
import { useMessage } from 'naive-ui';

const API = createAPICaller(router)
const message = useMessage()

const RoomID = ref("")
const Rooms = reactive([])

const AfterLogin = () => {
    API.get("/api/admin/rooms").then(rsp => {
        let data = rsp.data;
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        Rooms.length = 0
        data.data.list.forEach(item => {
            Rooms.push(item)
        })
    }).catch(err => message.error(JSON.stringify(err)))
}

const OnBtnGo = () => {
    if (!RoomID.value) return
    router.push('/admin/room/'+RoomID.value)
}
</script>

<style>
.overview-title {
    margin: 8px;
}

.overview-title p {
    padding-top: 2px;
}
</style>