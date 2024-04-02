import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import PrimeVue from 'primevue/config'
import router from './router/router';
import Toast from 'primevue/toast';


import './assets/style.css';


import Tooltip from 'primevue/tooltip';
import ToastService from 'primevue/toastservice';


//in main.js

const pinia = createPinia()
const app =  createApp(App)
app.use(PrimeVue, { ripple: true })
app.use(router)
app.use(pinia)

app.component('Toast', Toast)

// service
app.use(ToastService)


app.directive('tooltip', Tooltip);


app.mount('#app')
