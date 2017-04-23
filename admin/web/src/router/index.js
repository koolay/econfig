import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import Dashboard from '@/components/Dashboard'
import Items from '@/components/Items'
import Login from '@/components/Login'
import auth from '@/services/auth'

Vue.use(Router)

function requireAuth (to, from, next) {
    if (!auth.loggedIn()) {
        next({
            path: '/login',
            query: { redirect: to.fullPath }
        })
    } else {
        next()
    }
}

export default new Router({
    routes: [
        {
            path: '/',
            name: 'Home',
            component: Home
        },

        {
            path: '/login',
            name: 'Login',
            component: Login
        },
        {
            path: '/logout',
            beforeEnter (to, from, next) {
                auth.logout()
                next('/')
            }
        },
        {
            path: '/dashboard',
            name: 'Dashboard',
            component: Dashboard,
            beforeEnter: requireAuth
        },

        {
            path: '/app/:app/items',
            name: 'items',
            component: Items,
            beforeEnter: requireAuth
        }
    ]
})
