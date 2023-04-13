const RoomList = {
    code: 0,
    msg: "",
    data: {
        list: [
            {room_id: 1, name: "主播1", icon: "/icon.png"},
            {room_id: 2, name: "主播2", icon: "/icon.png"},
            {room_id: 3, name: "主播3", icon: "/icon.png"},
            {room_id: 4, name: "主播4", icon: "/icon.png"},
            {room_id: 5, name: "主播5", icon: "/icon.png"},
        ],
    },
};

export default {
    'get|^/apiroom/list$': opt => RoomList,
};