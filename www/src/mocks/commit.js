const Commit = {
    code: 0,
    msg: "",
    data: {}
}

export default {
    'post|^/api/highlight/commit$': opt => Commit
}