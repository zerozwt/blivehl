import axios from 'axios'

const jumpToLoginOrResolve = (router, rsp, resolve) => {
    let data = rsp.data
    if (data.code == 114514) {
        router.push("/login")
        return
    }
    if (data.code == 1919810) {
        router.push("/")
        return
    }
    resolve(rsp)
}

const createAPICaller = (router) => {
    return {
        get: (path, config) => new Promise((resolve, reject) => {
            axios.get(path, config).then(rsp => {
                jumpToLoginOrResolve(router, rsp, resolve)
            }).catch(err => reject(err))
        }),
        post: (path, data) => new Promise((resolve, reject) => {
            axios.post(path, data).then(rsp => {
                jumpToLoginOrResolve(router, rsp, resolve)
            }).catch(err => reject(err))
        }),
    }
}

export default createAPICaller