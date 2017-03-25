
import ls from 'local-storage'

export default {
    get (key) {
        return ls(key)
    },

    set (key, val) {
        return ls(key, val)
    },

    remove (key) {
        return ls.remove(key)
    }
}
