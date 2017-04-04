import Vue from 'vue'
import Router from 'vue-router'
import Hello from '@/components/Hello'
import Dashboard from '@/components/Dashboard'
import Items from '@/components/Items'

Vue.use(Router)

export default new Router({
    routes: [
        {
            path: '/',
            name: 'Hello',
            component: Hello
        },

        {
            path: '/dashboard',
            name: 'Dashboard',
            component: Dashboard
        },

        {
            path: '/app/:app/items',
            name: 'items',
            component: Items
        }
    ]
})
