<template>
    <n-layout>
        <n-layout-header bordered>
            <n-breadcrumb>
                <n-breadcrumb-item><router-link to="/">首页</router-link></n-breadcrumb-item>
                <n-breadcrumb-item>{{ BasicInfoLoaded ? BasicInfo.name : $route.params.roomid }}</n-breadcrumb-item>
            </n-breadcrumb>
        </n-layout-header>
        <n-layout-content v-if="BasicInfoLoaded">
            <n-card title="主播基本信息" class="basic-info">
                <n-space>
                    <n-image width="256" :src="BasicInfo.icon" preview-disabled/>
                    <n-space vertical>
                        <div class="strong-word" >{{ BasicInfo.name }}</div>
                        <span>UID: {{ BasicInfo.uid }}</span>
                        <span>直播间: {{ BasicInfo.room_id }}</span>
                        <n-button type="primary" :disabled="LightUpPrepare" :loading="LightUpPrepare" @click="PrpareLight">
                            <n-space><n-icon><LightBulbIcon/></n-icon>现在点灯！</n-space>
                        </n-button>
                    </n-space>
                </n-space>
            </n-card>
        </n-layout-content>
    </n-layout>
    <n-modal v-model:show="LightUpModalShow" :mask-closable="false">
        <n-card title="点灯" closable @close="event => {LightUpModalShow = false}" class="lightup-modal">
            <n-layout has-sider>
                <n-layout-sider>
                    <n-image width="256" :src="CurrentLiveInfo.cover" preview-disabled/>
                </n-layout-sider>
                <n-layout-content>
                    <n-space vertical>
                        <div class="strong-word">{{ CurrentLiveInfo.title }}</div>
                        <div>开播时间: {{ CurrentLiveInfo.live_start_time }}</div>
                        <span></span>
                        <span>点灯时刻: {{ LightUpTime }}</span>
                        <n-input class="lightup-message" type="text" placeholder="发生什么事了？发生什么事了？发生什么事了？" v-model:value="LightUpComment"/>
                        <n-button type="primary">
                            <n-space><n-icon><LightBulbIcon/></n-icon>点灯！</n-space>
                        </n-button>
                    </n-space>
                </n-layout-content>
            </n-layout>
        </n-card>
    </n-modal>
</template>

<script setup>
import {RouterLink, useRoute, useRouter} from 'vue-router'
import {computed, onMounted, reactive, ref} from 'vue'
import axios from 'axios';
import { useMessage } from 'naive-ui';
import LightBulbIcon from '../components/icons/LightBulbIcon.vue';

const BasicInfoLoaded = ref(false)
const BasicInfo = reactive({
    room_id: 0,
    uid: 0,
    name: "",
    icon: "",
})

const router = useRouter()
const route = useRoute()
const message = useMessage()

onMounted(() => {
    router.isReady().then(() => {
        axios.get("/api/basic_info", {params: {room_id: route.params.roomid}}).then(rsp => {
            let data = rsp.data
            if (data.code != 0) {
                message.error(`[${data.code}]请求失败: ${data.msg}`)
                return
            }
            let streamerInfo = data.data.streamer
            BasicInfo.room_id = route.params.roomid
            BasicInfo.uid = streamerInfo.uid
            BasicInfo.name = streamerInfo.name
            BasicInfo.icon = streamerInfo.icon
            BasicInfoLoaded.value = true
        }).catch(err => message.error(JSON.stringify(err)))
    })
})

const LightUpPrepare = ref(false)
const LightUpModalShow = ref(false)
const CurrentLiveInfo = reactive({
    title: "",
    live_start_time: "",
    cover: "",
    live_id: 0,
    light_ts: 0,
})
const LightUpTime = computed(() => {
    let timeNumber = value => value < 10 ? '0'+value : value
    let ret = `${timeNumber(Math.floor((CurrentLiveInfo.light_ts%3600)/60))}分${timeNumber(CurrentLiveInfo.light_ts%60)}秒`
    return CurrentLiveInfo.light_ts > 3600 ? `${Math.floor(CurrentLiveInfo.light_ts/3600)}小时`+ret : ret
})
const LightUpComment = ref("")

const PrpareLight = () => {
    LightUpPrepare.value = true
    axios.get("/api/prepare", {params: {room_id: BasicInfo.room_id}}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        let liveData = data.data
        if (!liveData.live_status) {
            message.error(`主播未开播或已经下播，不能点灯。`)
            return
        }
        CurrentLiveInfo.title = liveData.title
        CurrentLiveInfo.live_start_time = liveData.live_start_time
        CurrentLiveInfo.cover = liveData.cover
        CurrentLiveInfo.live_id = liveData.live_id
        CurrentLiveInfo.light_ts = liveData.light_ts
        LightUpModalShow.value = true
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {LightUpPrepare.value = false})
}

</script>

<style>
.basic-info {
    margin: 8px 0;
}

.strong-word {
    font-size: 24px;
    font-weight: bold;
}

.lightup-modal {
    max-width: 960px;
}

.lightup-message {
    max-width: 512px;
}
</style>