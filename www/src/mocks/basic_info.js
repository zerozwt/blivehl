const BasicInfo = {
    code: 0,
    msg: "",
    data: {
        streamer: {
            uid: 12345,
            name: "主播昵称",
            icon: "/icon.png",
        }
    },
}

export default {
    'get|^/api/room/basic\\?room_id=[0-9]+$': opt => BasicInfo,
};