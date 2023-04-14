import { createRouter, createWebHistory } from 'vue-router'
import IndexView from '../views/IndexView.vue'
import LiveView from '../views/LiveView.vue'
import LoginView from '../views/LoginView.vue'

import Overview from '../views/admin/Overview.vue'
import LiveListView from '../views/admin/LiveListView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: IndexView
    },
    {
      path: '/room/:roomid',
      name: 'live',
      component: LiveView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/admin/overview',
      name: "admin",
      component: Overview
    },
    {
      path: '/admin/room/:roomid',
      name: "admin_live_list",
      component: LiveListView
    }
  ]
})

export default router
