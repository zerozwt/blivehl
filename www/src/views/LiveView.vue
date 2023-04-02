<template>
    <n-layout>
        <n-layout-header bordered>
            <n-breadcrumb>
                <n-breadcrumb-item><router-link to="/">首页</router-link>></n-breadcrumb-item>
                <n-breadcrumb-item>{{ $route.params.roomid }}</n-breadcrumb-item>
            </n-breadcrumb>
        </n-layout-header>
    </n-layout>
</template>

<script setup>
import {RouterLink} from 'vue-router'
</script>

<style>
</style>