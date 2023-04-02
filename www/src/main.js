import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import naive from 'naive-ui'

import './assets/main.css'

if (import.meta.env.DEV) {
    console.log("DEV_ENV");
    await import('./mocks');
}

const app = createApp(App)

app.use(router).use(naive)

app.mount('#app')
