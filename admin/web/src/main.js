// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import Resource from 'vue-resource'
import iView from 'iview'
import NProgress from 'nprogress'
import ls from './services/ls'
// import config from './config'

import 'iview/dist/styles/iview.css'
import 'nprogress/nprogress.css'

Vue.config.productionTip = false
Vue.use(Resource)
Vue.use(iView)

Vue.http.interceptors.push((request, next) => {
    const token = ls.get('token')
    if (token) {
        Vue.http.headers.common['Authorization'] = `Bearer ${token}`
    }
    next((response) => {
        NProgress.done()
    })
})

Vue.http.options.root = '/'
/* eslint-disable no-new */
new Vue({
    el: '#app',
    router,
    template: '<App/>',
    components: { App }
})
