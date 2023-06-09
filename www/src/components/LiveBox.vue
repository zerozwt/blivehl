<template>
    <n-card :title="props.liveData.title" class="live-box">
        <n-layout has-sider>
            <n-layout-sider>
                <n-image width="256" :src="props.liveData.cover" preview-disabled/>
            </n-layout-sider>
            <n-space vertical>
                <span>开播时间: {{ props.liveData.live_start_time }}</span>
                <span></span>
                <n-space>
                    <n-button circle @click="emitCopyTimeline" v-if="timeline.length>0"><n-icon><CopyIcon/></n-icon></n-button>
                    <n-button circle @click="onDownload"><n-icon><DownloadIcon/></n-icon></n-button>
                    <n-button circle :disabled="loadingTimeline" @click="refreshTimeline"><n-icon><RfreshIcon/></n-icon></n-button>
                    <div style="padding-top: 6px;">高能瞬间:</div>
                </n-space>
                <n-timeline>
                    <n-timeline-item v-for="item in timeline" :key="item.time" type="info" :time="RenderTimestamp(item.time)">
                        <TimelineEntry :comment="item.comment" :time="item.time" @update="onEntryUpdate"/>
                    </n-timeline-item>
                </n-timeline>
            </n-space>
        </n-layout>
    </n-card>
</template>

<script setup>
import axios from 'axios';
import { useMessage } from 'naive-ui';
import { inject, onMounted, onUpdated, reactive, ref } from 'vue';
import TimelineEntry from './TimelineEntry.vue';
import RfreshIcon from '../components/icons/RfreshIcon.vue'
import DownloadIcon from '../components/icons/DownloadIcon.vue'
import CopyIcon from '../components/icons/CopyIcon.vue'

const props = defineProps(['liveData', 'roomid'])
const emit = defineEmits(['copyTimeline'])
const timeline = reactive([])
const loadingTimeline = ref(true)
const message = useMessage()

const RenderTimestamp = inject('RenderTimestamp')

const refreshTimeline = () => {
    loadingTimeline.value = true
    axios.get("/api/highlight/timeline", {params: {room_id: props.roomid, live_id: props.liveData.live_id}}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        while(timeline.length>0) timeline.pop()
        data.data.timeline.forEach((item) => {
            timeline.push(item)
        })
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {loadingTimeline.value = false})
}

onMounted(refreshTimeline)
onUpdated(refreshTimeline)

const onDownload = () => {
    window.open('/api/highlight/download?room_id='+props.roomid+'&live_id='+props.liveData.live_id, '_blank')
}

const onEntryUpdate = (ts, comment) => {
    axios.post("/api/highlight/commit", {
        room_id: Number(props.roomid),
        live_id: Number(props.liveData.live_id),
        ts: Number(ts),
        comment: !comment ? "(暂未填写描述)" : comment,
    }).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        refreshTimeline()
    }).catch(err => message.error(JSON.stringify(err)))
}

const emitCopyTimeline = () => {
    if (timeline.length <= 0) return
    let data = []
    timeline.forEach((item) => {
        data.push(RenderTimestamp(item.time)+': '+item.comment)
    })
    let reverse = []
    for (let i = data.length-1; i >= 0; i--) {
        reverse.push(data[i])
    }
    emit('copyTimeline', reverse.join('\n'))
}
</script>

<style>
.live-box {
    margin: 16px 0;
}
</style>