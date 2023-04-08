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
            <LiveBox v-for="item in LiveList" :key="item.live_id" :live-data="item" :roomid="$route.params.roomid" @copyTimeline="onCopytimeline"/>
            <div v-if="!LiveListEnded">
                <n-button type="tertiary" class="more-button" @click="AppendLiveList" :disabled="Appending">点击加载更多记录</n-button>
            </div>
        </n-layout-content>
    </n-layout>
    <n-modal v-model:show="LightUpModalShow" :mask-closable="false">
        <n-card :title="CurrentLiveInfo.modal_title" closable @close="event => {LightUpModalShow = false}" class="lightup-modal">
            <n-layout has-sider>
                <n-layout-sider>
                    <n-image width="256" :src="CurrentLiveInfo.cover" preview-disabled/>
                </n-layout-sider>
                <n-layout-content>
                    <n-space vertical>
                        <div class="strong-word">{{ CurrentLiveInfo.title }}</div>
                        <div>开播时间: {{ CurrentLiveInfo.live_start_time }}</div>
                        <span></span><span></span>
                        <span>点灯时刻: {{ RenderTimestamp(CurrentLiveInfo.light_ts) }}</span>
                        <n-input type="text" placeholder="发生什么事了？发生什么事了？发生什么事了？" v-model:value="LightUpComment"/>
                        <n-button type="primary" :disabled="HightLightCommitting" :loading="HightLightCommitting" @click="CommitHighLight">
                            <n-space><n-icon><LightBulbIcon/></n-icon>点灯！</n-space>
                        </n-button>
                    </n-space>
                </n-layout-content>
            </n-layout>
        </n-card>
    </n-modal>
    <n-modal v-model:show="ShowCopyTimeline">
        <n-card title="复制高能时间轴" closable @close="event => {ShowCopyTimeline = false}" class="lightup-modal">
            <n-input type="textarea" :value="TimelineData" :allow-input="value => false" rows="10"/>
        </n-card>
    </n-modal>
</template>

<script setup>
import {RouterLink, useRoute, useRouter} from 'vue-router'
import {onMounted, provide, reactive, ref} from 'vue'
import axios from 'axios';
import { useMessage } from 'naive-ui';
import LightBulbIcon from '../components/icons/LightBulbIcon.vue';
import LiveBox from '../components/LiveBox.vue';

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

const LiveList = reactive([])
const LiveListEnded = ref(false)
const InitLiveList = () => {
    axios.get("/api/live_list", {params:{room_id: route.params.roomid, until: 0}}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        let list = data.data.list
        while(LiveList.length>0) LiveList.pop()
        list.forEach((item) => {LiveList.push(reactive({
            title: item.title,
            live_start_time: item.live_start_time,
            cover: item.cover,
            live_id: item.live_id,
        }))})
        LiveListEnded.value = data.data.ended
    }).catch(err => message.error(JSON.stringify(err)))
}

const Appending = ref(false)
const AppendLiveList = () => {
    if (!LiveList.length) return
    Appending.value = true
    axios.get("/api/live_list", {params:{room_id: route.params.roomid, until: LiveList[LiveList.length-1].live_id-1}}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        data.data.list.forEach((item) => {LiveList.push(reactive({
            title: item.title,
            live_start_time: item.live_start_time,
            cover: item.cover,
            live_id: item.live_id,
        }))})
        LiveListEnded.value = data.data.ended
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {Appending.value = false})
}

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
            InitLiveList()
        }).catch(err => message.error(JSON.stringify(err)))
    })
})

const LightUpPrepare = ref(false)
const LightUpModalShow = ref(false)
const CurrentLiveInfo = reactive({
    modal_title: "",
    title: "",
    live_start_time: "",
    cover: "",
    live_id: 0,
    light_ts: 0,
    commit_cb: null,
})

const RenderTimestamp = (value) => {
    let timeNumber = value => value < 10 ? '0'+value : value
    let ret = `${timeNumber(Math.floor((value%3600)/60))}分${timeNumber(value%60)}秒`
    return value > 3600 ? `${Math.floor(value/3600)}小时`+ret : ret
}
provide('RenderTimestamp', RenderTimestamp)

const LightUpComment = ref("")

const OpenModal = (liveData, header, comment, commit_cb) => {
    LightUpComment.value = comment
    CurrentLiveInfo.modal_title = header
    CurrentLiveInfo.title = liveData.title
    CurrentLiveInfo.live_start_time = liveData.live_start_time
    CurrentLiveInfo.cover = liveData.cover
    CurrentLiveInfo.live_id = liveData.live_id
    CurrentLiveInfo.light_ts = liveData.light_ts
    CurrentLiveInfo.commit_cb = commit_cb
    LightUpModalShow.value = true
}

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
        OpenModal(liveData, "点灯", "", InitLiveList)
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {LightUpPrepare.value = false})
}

const HightLightCommitting = ref(false)
const CommitHighLight = () => {
    let comment = LightUpComment.value
    if (!comment) comment = "(暂未填写描述)"
    HightLightCommitting.value = true
    axios.post("/api/commit", {
        room_id: Number(BasicInfo.room_id),
        live_id: CurrentLiveInfo.live_id,
        ts: CurrentLiveInfo.light_ts,
        comment: comment,
    }).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        CurrentLiveInfo.commit_cb()
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {
        HightLightCommitting.value = false
        LightUpModalShow.value = false
    })
}

const ShowCopyTimeline = ref(false)
const TimelineData = ref("")
const onCopytimeline = (copyData) => {
    TimelineData.value = copyData
    ShowCopyTimeline.value = true
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
    max-width: 800px;
}

.more-button {
    margin: 0 auto;
    display: block;
}
</style>