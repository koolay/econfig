
import http from './http'
import ls from './ls'

export default {
    login (account, password, cb) {
        http.post('/api/login', {
            'account': account,
            'password': password,
            'captcha': ''
        }, response => {
            if (response.code === 0) {
                ls.set('token', response.data.token)
                this.onChange(true)
                return cb({
                    code: response.code,
                    token: response.data.token
                })
            } else {
                return cb({ code: response.code, msg: response.msg })
            }
        })
    },

    getToken () {
        return ls.get('token')
    },

    logout (cb) {
        ls.remove('token')
        this.onChange(false)
        if (cb) cb()
    },

    loggedIn () {
        return ls.get('token')
    },

    onChange (loggedIn) {}

}
