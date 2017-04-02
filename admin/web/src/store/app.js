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

    get: function (projectId, successCb, errorCb) {
        const apiUrl = `/api/app/${projectId}`
        http.get(apiUrl, data => {
            successCb(data)
        })
    }

}
