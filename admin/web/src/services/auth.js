
import http from './http'

export default {
    login (username, password, cb) {
        cb = arguments[arguments.length - 1]
        if (localStorage.token) {
            if (cb) cb(true)
            this.onChange(true)
            return
        }
        pretendRequest(username, password, (res) => {
            if (res.authenticated) {
                localStorage.token = res.token
                if (cb) cb(true)
                this.onChange(true)
            } else {
                if (cb) cb(false)
                this.onChange(false)
            }
        })
    },

    getToken () {
        return localStorage.token
    },

    logout (cb) {
        delete localStorage.token
        if (cb) cb()
        this.onChange(false)
    },

    loggedIn () {
        return !!localStorage.token
    },

    onChange () {}
}

function pretendRequest (username, password, cb) {
    http.post('/api/login', {
        'username': username,
        'password': password
    }, function (response) {
        if (response.data.result) {
            cb({
                authenticated: true,
                token: response.data.token
            })
        } else {
            cb({ authenticated: false })
        }
    })
}
