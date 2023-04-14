const Login = {
    code: 0,
    msg: "",
    data: {},
}

const UserInfo = {
    code: 0,
    msg: "",
    data: {
        name: "这是用户名",
        admin: true,
    },
}

export default {
    'post|^/api/user/login$': opt => Login,
    'get|^/api/user/logout$': opt => Login,
    'get|^/api/user/info$': opt => UserInfo,
    'post|^/api/user/pass$': opt => Login,
}