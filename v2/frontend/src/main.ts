import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import { globalLoader } from 'vue-global-loader'
import { createNotivue } from 'notivue'
import 'notivue/notification.css'

import { useAuthStore } from '@/stores/auth'

const app = createApp(App)

const pinia = createPinia()
app.use(pinia)

app.use(router)
app.use(globalLoader)

app.use(
  createNotivue({
    position: 'top-right',
    limit: 3,
  }),
)

// IMPORTANT: init auth BEFORE mount
const auth = useAuthStore(pinia)

auth.initAuth().finally(() => {
  app.mount('#app')
})
