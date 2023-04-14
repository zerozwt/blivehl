<template>
    <n-layout>
        <n-layout-header bordered>
            <div class="clear-float">
                <div style="float:left; padding-top: 1px;">
                    <n-breadcrumb separator=">">
                        <n-breadcrumb-item><router-link to="/">首页</router-link></n-breadcrumb-item>
                        <n-breadcrumb-item><router-link to="/admin/overview">数据管理</router-link></n-breadcrumb-item>
                        <n-breadcrumb-item>{{ BasicInfoLoaded ? BasicInfo.name : $route.params.roomid }}</n-breadcrumb-item>
                    </n-breadcrumb>
                </div>
                <div style="float: right; padding-bottom: 4px;"><LoginStatus @login="AfterLogin"/></div>
            </div>
        </n-layout-header>
        <n-layout-content>
            <div style="padding-top: 8px;">
                <n-data-table
                remote
                ref="table"
                :columns="TableColumns"
                :data="TableData"
                :loading="TableLoading"
                :pagination="TablePage"
                :row-key="(row) => row.live_id"
                @update:page="OnPageChange"
                />
            </div>
        </n-layout-content>
    </n-layout>
    <n-modal v-model:show="ShowDetail" transform-origin="center">
        <n-card closable @close="event => {ShowDetail = false}" title="直播高能数据" class="admin-detail">
            <n-layout has-sider>
                <n-layout-sider bordered>
                    <n-image width="256" :src="Detail.cover" preview-disabled/>
                    <p>{{ BasicInfo.name }}</p>
                    <p>{{ Detail.title }}</p>
                    <p>{{ Detail.live_start_time }}</p>
                    <n-button class="dl-button" type="primary" @click="() => {OnTableDownload(Detail.live_id)}">下载数据</n-button>
                </n-layout-sider>
                <n-layout-content style="padding-left: 12px;">
                    <n-scrollbar style="max-height: 480px;">
                        <n-timeline>
                            <n-timeline-item v-for="item in DetailData" type="info" :title="item.comment" :time="RenderTimestamp(item.time)" :content="item.author"/>
                        </n-timeline>
                    </n-scrollbar>
                </n-layout-content>
            </n-layout>
        </n-card>
    </n-modal>
</template>

<script setup>
import {RouterLink, useRoute, useRouter} from 'vue-router'
import LoginStatus from '@/components/LoginStatus.vue'
import { h, reactive, ref } from 'vue';
import router from "@/router"
import createAPICaller from '@/utils'
import { NButton, NSpace, useMessage } from 'naive-ui';

const API = createAPICaller(router)
const message = useMessage()
const routerObj = useRouter()
const route = useRoute()

const BasicInfoLoaded = ref(false)
const BasicInfo = reactive({
    room_id: 0,
    uid: 0,
    name: "",
    icon: "",
})

const TableLoading = ref(false)
const TableColumns = [
    {
        title: "直播标题",
        key: "title",
    },
    {
        title: "直播时间",
        key: "live_start_time",
    },
    {
        title: "操作",
        key: "live_id",
        render(row) {
            return h(NSpace, {}, () => [
                h(NButton, {size: "tiny", type: "info", onClick: () => {OnDisplayDetail(row)}}, () => "详情"),
                h(NButton, {size: "tiny", type: "primary", onClick: () => {OnTableDownload(row.live_id)}}, () => "下载"),
            ])
        }
    }
]
const TableData = ref([])
const TablePage = reactive({
    page: 1,
    pageCount: 1,
    pageSize: 10,
})

const LoadTableData = (page, pageSize) => {
    TableLoading.value = true
    API.get("/api/admin/lives", {params: {
        room_id: route.params.roomid,
        page: page,
        size: pageSize,
    }}).then(rsp => {
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        let apiData = data.data
        TableData.value = apiData.list
        TablePage.page = page
        TablePage.pageCount = Math.floor(apiData.total / TablePage.pageSize)
        TablePage.itemCount = apiData.total
    }).catch(err => message.error(JSON.stringify(err))).finally(() => {TableLoading.value = false})
}

const OnPageChange = (currentPage) => {
    LoadTableData(currentPage, TablePage.pageSize)
}

const OnTableDownload = (live_id) => {
    window.open('/api/admin/download?room_id='+route.params.roomid+'&live_id='+live_id, '_blank')
}

const AfterLogin = () => {
    routerObj.isReady().then(() => {
        API.get("/api/room/basic", {params: {room_id: route.params.roomid, save_recent: false}}).then(rsp => {
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
            LoadTableData(1, TablePage.pageSize)
        }).catch(err => message.error(JSON.stringify(err)))
    })
}

const ShowDetail = ref(false)
const Detail = reactive({
    title: "",
    live_start_time: "",
    cover: "",
    live_id: 0,
})
const DetailData = reactive([])

const OnDisplayDetail = (row) => {
    Detail.title = row.title
    Detail.live_start_time = row.live_start_time
    Detail.cover = row.cover
    Detail.live_id = row.live_id
    ShowDetail.value = true
    DetailData.length = 0
    LoadDetail(row.live_id)
}

const LoadDetail = (live_id) => {
    API.get("/api/admin/timeline", {params: {room_id: route.params.roomid, live_id: live_id}}).then(rsp => {
        if (live_id != Detail.live_id || !ShowDetail.value) {
            return
        }
        let data = rsp.data
        if (data.code != 0) {
            message.error(`[${data.code}]请求失败: ${data.msg}`)
            return
        }
        let apiData = data.data
        apiData.timeline.forEach((item) => {DetailData.push(item)})
    }).catch(err => message.error(JSON.stringify(err)))
}

const RenderTimestamp = (value) => {
    let timeNumber = value => value < 10 ? '0'+value : value
    let ret = `${timeNumber(Math.floor((value%3600)/60))}分${timeNumber(value%60)}秒`
    return value >= 3600 ? `${Math.floor(value/3600)}小时`+ret : ret
}
</script>

<style>
.admin-detail {
    max-width: 720px;
}
.dl-button {
    margin-top: 8px;
}
</style>