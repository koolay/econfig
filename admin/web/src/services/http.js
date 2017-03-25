import Vue from 'vue'

/**
 * Responsible for all HTTP requests.
 */
/* global alert */
export default {
    request (method, url, data, successCb = null, errorCb = null) {
        return Vue.http[method](url, data).then(response => {
            successCb(response.body)
        }, error => {
            if (errorCb) {
                errorCb(error)
            } else {
                alert(JSON.stringify(error))
            }
        })
    },

    get (url, successCb = null, errorCb = null) {
        return this.request('get', url, {}, successCb, errorCb)
    },

    post (url, data, successCb = null, errorCb = null) {
        return this.request('post', url, data, successCb, errorCb)
    },

    put (url, data, successCb = null, errorCb = null) {
        return this.request('put', url, data, successCb, errorCb)
    },

    delete (url, successCb = null, errorCb = null) {
        return this.request('delete', url, null, successCb, errorCb)
    },

    /**
     * A shortcut method to ping and check if the user session is still valid.
     */
    ping () {
        return this.get('/')
    }
}
