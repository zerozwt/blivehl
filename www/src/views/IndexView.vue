<template>
    <n-layout>
        <n-layout-header bordered>
            <div class="navbar">
                <div style="float:left; padding-top: 1px;">
                    <n-breadcrumb separator=">">
                        <n-breadcrumb-item>首页</n-breadcrumb-item>
                    </n-breadcrumb>
                </div>
                <div style="float: right; padding-bottom: 4px;"><LoginStatus @login="AfterLogin"/></div>
            </div>
        </n-layout-header>
        <div class="rooms-list">
            <n-layout>
                <n-layout-header bordered class="room-bar">
                    <n-input type="text" v-model:value="roomid" :allow-input="onlyAllowNumber" round size="large" placeholder="直播间号码">
                        <template #suffix>
                            <n-button type="primary" @click="jumpto" round><n-icon><LightBulbIcon/></n-icon>开始点灯</n-button>
                        </template>
                    </n-input>
                </n-layout-header>
                <n-layout-content>
                    <p style="margin-bottom: 8px">最近点灯:</p>
                    <n-grid x-gap="8" y-gap="8" :cols="4">
                        <n-gi v-for="item in rooms" :key="item.room_id">
                            <router-link :to="'/room/'+item.room_id">
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
        </div>
    </n-layout>
</template>

<script setup>
import {ref, reactive, onMounted} from 'vue'
import LightBulbIcon from '../components/icons/LightBulbIcon.vue'
import router from "@/router"
import createAPICaller from '@/utils'
import {useMessage} from 'naive-ui'
import {RouterLink} from 'vue-router'
import LoginStatus from '../components/LoginStatus.vue'

const API = createAPICaller(router)

const rooms = reactive([])
const onlyAllowNumber = (value) => !value || /^[0-9]+$/.test(value)
const roomid = ref("")

const jumpto = () => {
    if (!roomid.value) {return}
    router.push("/room/" + roomid.value)
}

const message = useMessage()

const AfterLogin = ()=> {
    API.get("/api/room/list").then(rsp => {
        let data = rsp.data;
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        rooms.length = 0
        data.data.list.forEach(item => {
            rooms.push(item)
        })
    }).catch(err => message.error(JSON.stringify(err)))
}
</script>

<style>
.rooms-list {
    margin: 256px auto 0 auto;
}

.room-bar {
    padding: 8px 0;
    margin-bottom: 8px;
}

.navbar:after {
    content: "";
    display: block;
    clear: both;
}
</style>