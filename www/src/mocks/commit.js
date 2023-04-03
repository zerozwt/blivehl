const Commit = {
    code: 0,
    msg: "",
    data: {}
}

export default {
    'post|^/api/commit$': opt => Commit
}