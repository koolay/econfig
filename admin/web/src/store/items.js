/**
 * configuration items
 */

import http from './../services/http'

export default {

    list: function (appName, cb) {
        const apiUrl = `/api/app/${appName}/items`
        http.get(apiUrl, data => {
            cb(data)
        })
    },

    listOfTmp: function (appName, successCb, errorCb) {
        const apiUrl = `/api/app/${appName}/tmp-items`
        http.get(apiUrl, data => {
            successCb(data)
        })
    },

    listOfStore: function (appName, successCb, errorCb) {
        const apiUrl = `/api/app/${appName}/store-items`
        http.get(apiUrl, data => {
            successCb(data)
        })
    },

    saveItem: function (appName, key, value, successCb, errorCb) {
        const apiUrl = `/api/app/${appName}/item`
        http.put(apiUrl, { app: appName, key, value }, data => {
            successCb(data)
        })
    }

}
