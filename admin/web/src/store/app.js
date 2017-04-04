/**
 * project mangement
 */

import http from './../services/http'

export default {

    list: function (cb) {
        const apiUrl = `/api/apps`
        http.get(apiUrl, data => {
            cb(data)
        })
    },

    get: function (appName, successCb, errorCb) {
        const apiUrl = `/api/app/${appName}`
        http.get(apiUrl, data => {
            successCb(data)
        })
    }

}
